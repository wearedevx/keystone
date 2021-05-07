package login

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cli/cmd"
	"github.com/wearedevx/keystone/cli/tests/utils"
)

func TestMain(m *testing.M) {
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
	utils.WaitAPIStart()
	testscript.Run(t, testscript.Params{
		Dir:   ".",
		Setup: utils.SetupEnvVars,
	})
}
