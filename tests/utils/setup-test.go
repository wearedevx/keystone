package utils

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	. "github.com/wearedevx/keystone/internal/jwt"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
)

func GetGcloudFuncAuthPidFilePath() string {
	return path.Join(os.TempDir(), "keystone_ksauth.pid")
}

func listenCmdStartProcess(cmd *exec.Cmd) {
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Println("stdout:", scanner.Text())
		}
		done <- true
	}()

	cmdErrorReader, _ := cmd.StderrPipe()
	scannerError := bufio.NewScanner(cmdErrorReader)
	doneError := make(chan bool)

	go func() {
		for scannerError.Scan() {
			fmt.Println("stderr:", scannerError.Text())
		}
		doneError <- true
	}()
}

func StartCloudAuth() {
	gcloudPidFilePath := GetGcloudFuncAuthPidFilePath() // + time.Now().String()
	pid, err := ioutil.ReadFile(gcloudPidFilePath)

	if err != nil {
		fmt.Println("No gcloud pid file", err)
	}

	if len(pid) == 0 {
		fmt.Println(" keystone ~ login_test.go ~ ON START")
		pgid := StartCloudAuthFunc()

		pidString := []byte(strconv.Itoa(pgid))
		err = ioutil.WriteFile(gcloudPidFilePath, pidString, 0755)

		if err != nil {
			panic(err)
		}
	}
}

func StartCloudAuthFunc() int {
	fmt.Println("🚀 ~ file: init_test.go ~ line 49 ~ funcTestHelloWorld ~ cmd.Process??")

	// Start cloud functions
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	cmd := exec.CommandContext(ctx, "go", "run", "-tags", "test", "cmd/main.go")
	cmd.Dir = "../../functions/ksauth"
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	listenCmdStartProcess(cmd)

	err := cmd.Start()

	if err != nil {
		// Problemo
		fmt.Println("Ayayaye 0", err.Error())
		panic(err)
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)

	fmt.Println(" keystone ~ init_test.go ~ err", err)

	if err != nil {
		// Problemo
		fmt.Println("Ayayaye 1", err.Error())
	}

	if err != nil {
		panic(err)
	}

	return pgid
}

func CreateAndLogUser(env *testscript.Env) error {
	Repo := new(repo.Repo)

	var user1 *User = &User{
		ExtID:       "56883564",
		UserID:      "00fb7666-de43-4559-b4e4-39b172117dd8",
		AccountType: "github",
		Username:    "LAbigael",
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
	env.Setenv("HOME", homeDir)
	env.Setenv("DB_PORT", os.Getenv("DB_PORT"))
	env.Setenv("TMPDIR", osTmpDir)
	env.Setenv("GCLOUDFILE", GetGcloudFuncAuthPidFilePath())

	// Create config folder
	err := os.MkdirAll(configDir, 0777)

	return err
}
