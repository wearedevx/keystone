package utils

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
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
	time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)

	return path.Join(os.TempDir(), "keystone_ksauth.pid")
}

func waitALittle() {
	min := 0
	max := 1000

	nbms := rand.Intn(max-min) + min

	time.Sleep(time.Duration(nbms) * time.Millisecond)
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

// func StartAuthCloudFunction() {
// 	gcloudPidFilePath := GetGcloudFuncAuthPidFilePath() // + time.Now().String()
// 	pid, _ := ioutil.ReadFile(gcloudPidFilePath)

// 	if len(pid) == 0 {
// 		pgid := startAuthCloudFuncProcess()

// 		pidString := []byte(strconv.Itoa(pgid))
// 		err := ioutil.WriteFile(gcloudPidFilePath, pidString, 0755)

// 		if err != nil {
// 			panic(err)
// 		}
// 	} else {
// 		// fmt.Println("PID DEJA PRESENT", pid)
// 	}
// }

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
	var done bool = false
	attemps := 0

	for !done {
		attemps = attemps + 1

		if attemps == maxAttempts {
			done = true
		}

		// Make a request to server
		request, _ := http.NewRequest("GET", serverUrl, nil)

		timeout := time.Duration(1 * time.Second)

		client := http.Client{
			Timeout: timeout,
		}

		_, err := client.Do(request)

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
	<-c

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

	// fmt.Println("keystone ~ setup-test.go ~ cmd.Process.Pid", cmd.Process.Pid)
	pgid, err := syscall.Getpgid(cmd.Process.Pid)

	if err != nil {
		// fmt.Println("Error gret gpid", err)
		waitALittle()
		return startCloudFunctionProcess(funcPath, serverUrl)
	}

	// fmt.Println("AVANT SLEEP")

	// fmt.Println("Ca wait server")
	waitForServerStarted(serverUrl)

	// time.Sleep(20000 * time.Millisecond)
	// fmt.Println("APRES SLEEP")

	return pgid
}

// func startAuthCloudFuncProcess() int {
// 	waitALittle()
// 	return startCloudFunctionProcess("../../api/ksauth/cmd/main.go", "http://127.0.0.1:9000")
// }

func startCloudApiFunc() int {
	waitALittle()
	return startCloudFunctionProcess("../../api/ksapi/cmd/main.go", "http://127.0.0.1:9001")
}

func CreateAndLogUser(env *testscript.Env) error {
	Repo := new(repo.Repo)

	ExtID := "56883564" + uuid.NewV4().String()
	Username := "LAbigael_" + uuid.NewV4().String()

	var user1 *User = &User{
		ExtID: ExtID,
		// UserID:      "00fb7666-de43-4559-b4e4-39b172117dd8",
		AccountType: "github",
		Username:    Username,
		Fullname:    "Test user1",
		Email:       "abigael.laldji@protonmail.com",
	}

	fmt.Println("keystone ~ functions.go ~ error MOU MOU")
	Repo.GetOrCreateUser(user1)

	token, _ := MakeToken(*user1)

	configDir := getConfigDir(env)

	pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	err := ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- Fullname: `+user1.Fullname+`
  account_type: "`+string(user1.AccountType)+`"
  email: `+user1.Email+`
  ext_id: "`+ExtID+`"
  fullname: `+user1.Fullname+`
  user_id: `+user1.UserID+`
  username: `+Username+`
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
		Username:    "Username owner " + uuid.NewV4().String(),
		Fullname:    "Fullname owner",
		Email:       "test+owner@example.com",
	}

	var devUser *User = &User{
		ExtID:       "my ext id",
		AccountType: "github",
		Username:    "Username dev " + uuid.NewV4().String(),
		Fullname:    "Fullname dev",
		Email:       "test+dev@example.com",
	}

	fmt.Println("keystone ~ functions.go ~ error DOUX DOUX")

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
