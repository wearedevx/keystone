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

	"github.com/bxcodec/faker/v3"
	"github.com/rogpeppe/go-internal/testscript"
	. "github.com/wearedevx/keystone/api/pkg/jwt"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
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

	fmt.Println(os.Getwd())

	cmd := exec.CommandContext(ctx, "go", "run", "-tags", "test", funcPath)
	cmd.Dir = funcPath
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
	return startCloudFunctionProcess("../../../api/main.go", "http://127.0.0.1:9001")
}

func CreateAndLogUser(env *testscript.Env) error {
	Repo := new(repo.Repo)
	user := User{}

	faker.FakeData(&user)

	user.Email = "email@example.com"

	Repo.GetOrCreateUser(&user)

	if err := Repo.Err(); err != nil {
		fmt.Println("Get Or Create User", err)
		os.Exit(1)
	}

	token, _ := MakeToken(user)

	configDir := getConfigDir(env)

	pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	err := ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- Fullname: `+user.Fullname+`
  account_type: "`+string(user.AccountType)+`"
  email: `+user.Email+`
  ext_id: "`+user.ExtID+`"
  fullname: `+user.Fullname+`
  user_id: `+user.UserID+`
  username: `+user.Username+`
auth_token: `+token+`
current: 0
`), 0o777)

	if err != nil {
		fmt.Println("error writing accounts", err)
	}

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
