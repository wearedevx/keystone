package env

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cli/cmd"
	"github.com/wearedevx/keystone/cli/tests/utils"
)

func TestMain(m *testing.M) {
	testscript.RunMain(m, map[string]func() int{
		"ks": cmd.Execute,
	})
}

func setupFunc(env *testscript.Env) error {
	utils.SetupEnvVars(env)
	utils.CreateAndLogUser(env)
	return nil
}

func TestCommands(t *testing.T) {
	if err := utils.WaitAPIStart(); err != nil {
		panic(err)
	}

	testscript.Run(t, testscript.Params{
		Dir:                  "./",
		Setup:                setupFunc,
		IgnoreMissedCoverage: true,
	})
}
