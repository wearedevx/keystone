package logintest

import (
	"os"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
	"github.com/wearedevx/keystone/tests/utils"
)

func TestMain(m *testing.M) {
	utils.StartAuthCloudFunction()

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
		Dir:   ".",
		Setup: utils.SetupEnvVars,
	})
}
