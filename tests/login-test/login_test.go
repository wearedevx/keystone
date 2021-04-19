package logintest

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
	"github.com/wearedevx/keystone/tests/utils"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"ks":                 cmd.Execute,
		"githubLoginSuccess": utils.GithubLoginSuccess,
	}))
}

func setupInitFunc(env *testscript.Env) error {
	// fmt.Println("env.Getenv", env.Getenv("WORK"))
	// fmt.Println("env.WorkDir", env.WorkDir)
	// fmt.Println("env.KSAUTH_URL", os.Environ())

	homeDir := path.Join(env.Getenv("WORK"), "home")

	env.Setenv("HOME", homeDir)

	// Repo := new(repo.Repo)
	// db := Repo.Connect()

	// // Migrate DB
	// repo.AutoMigrate(db)

	return nil
}

func startCloudAuthFunc() {
	fmt.Println("🚀 ~ file: init_test.go ~ line 49 ~ funcTestHelloWorld ~ cmd.Process??")

	// Start cloud functions
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	cmd := exec.CommandContext(ctx, "go", "run", "-tags", "test", "cmd/main.go")
	cmd.Dir = "../../functions/ksauth"

	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Printf(scanner.Text())
		}
		done <- true
	}()

	cmdErrorReader, _ := cmd.StderrPipe()
	scannerError := bufio.NewScanner(cmdErrorReader)
	doneError := make(chan bool)

	go func() {
		for scannerError.Scan() {
			fmt.Printf(scannerError.Text())
		}
		doneError <- true
	}()

	err := cmd.Start()
	fmt.Println(" keystone ~ init_test.go ~ err", err)

	// output, err := cmd.Start()
	// fmt.Println("🚀 ~ file: init_test.go ~ line 38 ~ funcTestHelloWorld ~ output", string(output))

	// defer fmt.Println("CA KILL 0")
	// defer fmt.Println("PROCESS", cmd.Process.Pid)
	// defer cmd.Process.Kill()
	// defer fmt.Println("CA KILL")
	// defer fmt.Println("CA KILL?")

	if err != nil {
		// Problemo
		fmt.Println("Ayayaye", err.Error())
	}

	fmt.Println("🚀 ~ file: init_test.go ~ line 49 ~ funcTestHelloWorld ~ cmd.Process", cmd.Process)

}

func TestLoginCommand(t *testing.T) {

	startCloudAuthFunc()
	fmt.Println("Wait ?")
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("Wait !")

	testscript.Run(t, testscript.Params{
		Dir: ".",
		// WorkdirRoot: "./",
		Setup: setupInitFunc,
	})

	// cmd.Wait()

	// cmd.Process.Kill()
}
