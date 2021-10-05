package paywall

import (
	"fmt"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/api/pkg/models"
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
	err := utils.MakeOrgaFree(env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.CreateFakeUserWithUsername("john.doe", models.GitlabAccountType, env)
	return nil
}

func TestCommands(t *testing.T) {
	utils.WaitAPIStart()

	testscript.Run(t, testscript.Params{
		Dir:                  "./",
		Setup:                setupFunc,
		IgnoreMissedCoverage: true,
	})
}
