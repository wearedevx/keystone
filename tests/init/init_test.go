package init

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
	"github.com/wearedevx/keystone/tests/utils"
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
	testscript.Run(t, testscript.Params{
		Dir:                  "./",
		Setup:                setupFunc,
		IgnoreMissedCoverage: true,
	})
}
