package utils

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	uuid "github.com/satori/go.uuid"
	. "github.com/wearedevx/keystone/internal/jwt"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
)

func GetGcloudFuncAuthPidFilePath() string {
	return path.Join(os.TempDir(), "keystone_ksauth.pid")
}

func GetGcloudFuncApiPidFilePath() string {
	return path.Join(os.TempDir(), "keystone_ksapi.pid")
}

func listenCmdStartProcess(cmd *exec.Cmd, name string) {
	cmdReader, _ := cmd.StdoutPipe()

	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Println(name, "stdout:", scanner.Text())
		}
		done <- true
	}()

	cmdErrorReader, _ := cmd.StderrPipe()
	scannerError := bufio.NewScanner(cmdErrorReader)
	doneError := make(chan bool)

	go func() {
		for scannerError.Scan() {
			fmt.Println(name, "stderr:", scannerError.Text())
		}
		doneError <- true
	}()
}

func StartAuthCloudFunction() {
	gcloudPidFilePath := GetGcloudFuncAuthPidFilePath() // + time.Now().String()
	fmt.Println("keystone ~ setup-test.go ~ gcloudPidFilePath", string(gcloudPidFilePath))
	pid, _ := ioutil.ReadFile(gcloudPidFilePath)

	if len(pid) == 0 {
		pgid := startAuthCloudFuncProcess()

		pidString := []byte(strconv.Itoa(pgid))
		err := ioutil.WriteFile(gcloudPidFilePath, pidString, 0755)

		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("PID DEJA PRESENT", pid)
	}
}

func StartApiCloudFunction() {
	gcloudPidFilePath := GetGcloudFuncApiPidFilePath() // + time.Now().String()
	pid, _ := ioutil.ReadFile(gcloudPidFilePath)

	if len(pid) == 0 {
		pgid := startCloudApiFunc()

		pidString := []byte(strconv.Itoa(pgid))
		err := ioutil.WriteFile(gcloudPidFilePath, pidString, 0755)

		if err != nil {
			panic(err)
		}
	}
}

func pollServer(serverUrl string, c chan bool, maxAttempts int) {
	fmt.Println("keystone ~ setup-test.go ~ pollServer")
	var done bool = false
	attemps := 0

	for !done {
		fmt.Println("keystone ~ setup-test.go ~ done", done)
		attemps = attemps + 1

		if attemps == maxAttempts {
			done = true
		}

		fmt.Println("Start request ! 0")
		// Make a request to server
		request, _ := http.NewRequest("GET", serverUrl, nil)

		timeout := time.Duration(1 * time.Second)

		client := http.Client{
			Timeout: timeout,
		}

		fmt.Println("Start request ! 1")
		_, err := client.Do(request)

		fmt.Println("keystone ~ setup-test.go ~ err", err)
		// If it's started,

		if err == nil {
			done = true
			c <- true
		}

		time.Sleep(1 * time.Second)
	}
}

func waitForServerStarted(serverUrl string) {
	const max_attempts int = 20

	c := make(chan bool)

	go pollServer(serverUrl, c, max_attempts)

	// result := true
	result := <-c
	fmt.Println("keystone ~ setup-test.go ~ result", result)

}

func startCloudFunctionProcess(funcPath string, serverUrl string) int {

	// Start cloud functions
	ctx, _ := context.WithTimeout(context.Background(), 20000*time.Second)

	cmd := exec.CommandContext(ctx, "go", "run", "-tags", "test", funcPath)
	// cmd.Dir = funcPath
	cmd.Dir = "."
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	listenCmdStartProcess(cmd, funcPath)

	err := cmd.Start()

	if err != nil {
		panic(err)
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)

	if err != nil {
		panic(err)
	}

	fmt.Println("AVANT SLEEP")

	waitForServerStarted(serverUrl)

	// time.Sleep(20000 * time.Millisecond)
	// fmt.Println("APRES SLEEP")

	return pgid
}

func startAuthCloudFuncProcess() int {
	return startCloudFunctionProcess("../../functions/ksauth/cmd/main.go", "http://127.0.0.1:9000")
}

func startCloudApiFunc() int {
	return startCloudFunctionProcess("../../functions/ksapi/cmd/main.go", "http://127.0.0.1:9001")
}

func CreateAndLogUser(env *testscript.Env) error {
	Repo := new(repo.Repo)

	var user1 *User = &User{
		ExtID: "56883564" + uuid.NewV4().String(),
		// UserID:      "00fb7666-de43-4559-b4e4-39b172117dd8",
		AccountType: "github",
		Username:    "LAbigael_" + uuid.NewV4().String(),
		Fullname:    "Test user1",
		Email:       "abigael.laldji@protonmail.com",
	}

	Repo.GetOrCreateUser(user1)

	token, _ := MakeToken(*user1)

	configDir := getConfigDir(env)

	pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	err := ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- Fullname: Test user1
  account_type: github
  email: abigael.laldji@protonmail.com
  ext_id: "56883564"
  fullname: Michel
  user_id: `+user1.UserID+`
  username: LAbigael
auth_token: `+token+`
current: 0
`), 0o777)

	return err
}

func getHomeDir(env *testscript.Env) string {
	return path.Join(env.Getenv("WORK"), "home")
}

func getConfigDir(env *testscript.Env) string {
	homeDir := getHomeDir(env)
	return path.Join(homeDir, ".config")
}

func SetupEnvVars(env *testscript.Env) error {
	homeDir := getHomeDir(env)
	configDir := getConfigDir(env)
	osTmpDir := os.TempDir()

	// Set home dir for test
	env.Setenv("GOPATH", "/DFDFDF")
	env.Setenv("HOME", homeDir)
	env.Setenv("DB_PORT", os.Getenv("DB_PORT"))
	env.Setenv("TMPDIR", osTmpDir)
	env.Setenv("GCLOUDFILE", GetGcloudFuncAuthPidFilePath())

	// Create config folder
	return os.MkdirAll(configDir, 0777)
}

func SeedTestData() {
	Repo := new(repo.Repo)

	devEnvironmentType, _ := Repo.GetOrCreateEnvironmentType("dev")
	stagingEnvironmentType, _ := Repo.GetOrCreateEnvironmentType("staging")
	prodEnvironmentType, _ := Repo.GetOrCreateEnvironmentType("prod")

	devRole := Repo.GetOrCreateRole("dev")
	devopsRole := Repo.GetOrCreateRole("devops")
	adminRole := Repo.GetOrCreateRole("admin")

	// DEV
	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devRole,
		EnvironmentType: devEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          false,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devRole,
		EnvironmentType: stagingEnvironmentType,
		Read:            false,
		Write:           false,
		Invite:          false,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devRole,
		EnvironmentType: prodEnvironmentType,
		Read:            false,
		Write:           false,
		Invite:          false,
	})

	// Staging
	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devopsRole,
		EnvironmentType: devEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devopsRole,
		EnvironmentType: stagingEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devopsRole,
		EnvironmentType: prodEnvironmentType,
		Read:            false,
		Write:           false,
		Invite:          false,
	})

	// ADMIN
	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            adminRole,
		EnvironmentType: devEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            adminRole,
		EnvironmentType: stagingEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            adminRole,
		EnvironmentType: prodEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	var userProjectOwner *User = &User{
		ExtID:       "my iowner ext id",
		AccountType: "github",
		Username:    "Username owner",
		Fullname:    "Fullname owner",
		Email:       "test+owner@example.com",
	}

	var devUser *User = &User{
		ExtID:       "my ext id",
		AccountType: "github",
		Username:    "Username dev",
		Fullname:    "Fullname dev",
		Email:       "test+dev@example.com",
	}

	Repo.GetOrCreateUser(userProjectOwner)
	Repo.GetOrCreateUser(devUser)

	var project *Project = &Project{
		Name: "project name",
	}

	Repo.GetOrCreateProject(project, *userProjectOwner)

	environmentType, _ := Repo.GetOrCreateEnvironmentType("dev")

	Repo.GetOrCreateEnvironment(*project, environmentType)

	Repo.GetOrCreateProjectMember(project, devUser, "dev")

}
