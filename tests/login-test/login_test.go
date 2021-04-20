package logintest

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
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
	"github.com/wearedevx/keystone/tests/utils"
	"gopkg.in/h2non/gock.v1"
)

func startNock() {
	// defer gock.Off() // Flush pending mocks after test execution
	// defer fmt.Println("ca depile")

	gock.EnableNetworking()

	gock.New("https://github.com").
		Persist().
		Post("/login/oauth/access_token").
		// Reply(200).
		// JSON(map[string]string{"access_token": "access_token"})
		ReplyFunc(func(resp *gock.Response) {

			fmt.Println(" keystone ~ login_test.go ~ resp", resp)

			resp.BodyString("access_token=tutu&token_type=token_type&refresh_token=refresh_token&expires_in=100")
			resp.SetHeader("Authorization", "Bearer montoken")
			// return resp
		})
}

func TestMain(m *testing.M) {
	startNock()

	defer gock.Off()
	defer gock.DisableNetworking()

	fmt.Println(" keystone ~ login_test.go ~ TestMain LAUNCH", time.Now())

	// strr := "coucou " + time.Now().String()
	// fmt.Println(" keystone ~ login_test.go ~ strr", strr)
	// ioutil.WriteFile("gcloud-func.pid", []byte(strr), 0755)

	gcloudPidFilePath := getGcloudFuncAuthPidFilePath() // + time.Now().String()
	// gcloudPidFilePath := "gcloud-func.pid" // + time.Now().String()
	// fmt.Println(" keystone ~ login_test.go ~ gcloudPidFile", gcloudPidFilePath)
	// gcloudPidFile := "$WORK/keystone"

	pid, err := ioutil.ReadFile(gcloudPidFilePath)
	// fmt.Println(" keystone ~ login_test.go ~ content", pid)

	if err != nil {
		fmt.Println("No gcloud pid file", err)
	}

	// fmt.Println(" keystone ~ login_test.go ~ pid", pid)

	if len(pid) == 0 {
		fmt.Println(" keystone ~ login_test.go ~ ON START")
		pgid := startCloudAuthFunc(gcloudPidFilePath)

		// fmt.Println("Wait ?")
		// time.Sleep(5000 * time.Millisecond)
		// fmt.Println("Wait !")

		pidString := []byte(strconv.Itoa(pgid))
		err = ioutil.WriteFile(gcloudPidFilePath, pidString, 0755)
		fmt.Println(" keystone ~ login_test.go ~ pidString TO WRITE", pgid)
	}

	pid, err = ioutil.ReadFile(gcloudPidFilePath)
	fmt.Println(" keystone ~ login_test.go ~ contenttttt", pid)

	// startCloudAuthFunc()

	resRun := testscript.RunMain(m, map[string]func() int{
		"ks":                 cmd.Execute,
		"githubLoginSuccess": utils.GithubLoginSuccess,
	})

	fmt.Println(" keystone ~ login_test.go ~ cmd KILL")
	// fmt.Println(" keystone ~ login_test.go ~ cmd.Process", cmd.Process.Pid)

	// pidBytes, err := ioutil.ReadFile(gcloudPidFilePath)

	// if err != nil {
	// 	panic(err)
	// }

	// buf := bytes.NewBuffer(pidBytes) // b is []byte
	// pidToKill, err := binary.ReadVarint(buf)
	// pidToKill, err2 := strconv.Atoi(string(pidBytes))
	// fmt.Println(" keystone ~ login_test.go ~ pidToKill", pidToKill)

	// if err2 != nil {
	// 	panic(err2)
	// }

	// syscall.Kill(-pidToKill, 15)

	// e := os.Remove(gcloudPidFile)

	// if e != nil {
	// 	panic(e)
	// }

	os.Exit(resRun)
}

var ksAuthCmd int

func init() {
	ksAuthCmd = 0
}

func getGcloudFuncAuthPidFilePath() string {
	return path.Join(os.TempDir(), "keystone_ksauth.pid")
}

func setupInitFunc(env *testscript.Env) error {

	// fmt.Println("env.Getenv", env.Getenv("WORK"))
	// fmt.Println("env.WorkDir", env.WorkDir)
	// fmt.Println("env.KSAUTH_URL", os.Environ())

	homeDir := path.Join(env.Getenv("WORK"), "home")
	osTmpDir := os.TempDir()
	configDir := path.Join(homeDir, ".config")
	fmt.Println(" keystone ~ login_test.go ~ homeDir", homeDir)
	// pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	// Set home dir for test
	env.Setenv("HOME", homeDir)
	env.Setenv("TMPDIR", osTmpDir)
	env.Setenv("GCLOUDFILE", getGcloudFuncAuthPidFilePath())

	// log.Println("HOME ?", env.Getenv("HOME"))

	// Create config folder
	err := os.MkdirAll(configDir, 0777)

	if err != nil {
		panic(err)
	}

	// _ := new(repo.Repo)

	// // Migrate DB
	// repo.AutoMigrate()

	// env.Defer(func() {
	// 	fmt.Println("FINISHHHHH")
	// })

	return nil
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

func startCloudAuthFunc(gcloudPidFile string) int {
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

	// Write pid in file
	// pidString := []byte(strconv.Itoa(pgid))
	// err = ioutil.WriteFile("gcloud-func.pid", pidString, 0755)
	// fmt.Println(" keystone ~ login_test.go ~ pidString TO WRITE", pgid)

	if err != nil {
		panic(err)
	}

	return pgid
}

func TestLoginCommand(t *testing.T) {

	// pgid := startCloudAuthFunc()

	fmt.Println("Wait ?")
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("Wait !")

	//  t.Run("group", func(t *testing.T, testscript.Params{
	// 	Dir: ".",
	// 	// WorkdirRoot: "./",
	// 	Setup: setupInitFunc,
	// }) {
	//     t.Run("Test1", parallelTest1)
	//     t.Run("Test2", parallelTest2)
	//     t.Run("Test3", parallelTest3)
	// })

	testscript.Run(t, testscript.Params{
		Dir: ".",
		// WorkdirRoot: "./",
		Setup: setupInitFunc,
	})

	// cmd.Wait()

	// fmt.Println(" keystone ~ login_test.go ~ cmd KILL")
	// // fmt.Println(" keystone ~ login_test.go ~ cmd.Process", cmd.Process.Pid)

	// // Kill pgid
	// syscall.Kill(-pgid, 15)
}
