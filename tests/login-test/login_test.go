package logintest

import (
	"fmt"
	"os"
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

	utils.StartCloudAuth()

	resRun := testscript.RunMain(m, map[string]func() int{
		"ks":                 cmd.Execute,
		"githubLoginSuccess": utils.GithubLoginSuccess,
	})

	os.Exit(resRun)
}

var ksAuthCmd int

func init() {
	ksAuthCmd = 0
}

func TestLoginCommand(t *testing.T) {

	time.Sleep(2000 * time.Millisecond)

	testscript.Run(t, testscript.Params{
		Dir:         ".",
		WorkdirRoot: "./",
		Setup:       utils.SetupEnvVars,
	})
}
