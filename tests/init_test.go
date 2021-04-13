package main

import (
	"os"
	"path"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
	"github.com/wearedevx/keystone/internal/config"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"ks": cmd.Execute,
	}))
}

func SetupFunc(env *testscript.Env) error {
	Repo := new(repo.Repo)
	Repo.Connect()

	var user1 *User = &User{
		ExtID:       "56883564",
		UserID:      "00fb7666-de43-4559-b4e4-39b172117dd8",
		AccountType: "github",
		Username:    "LAbigael",
		Fullname:    "Test user1",
		Email:       "abigael.laldji@protonmail.com",
	}

	Repo.GetOrCreateUser(user1)

	user1Account := map[string]string{
		"account_type": "github",
		"email":        "abigael.laldji@protonmail.com",
		"ext_id":       "56883564",
		"fullname":     "Michel",
		"user_id":      "00fb7666-de43-4559-b4e4-39b172117dd8",
		"username":     "LAbigael",
		"Fullname":     "Test user1",
	}

	// Set home dir for test
	env.Setenv("HOME", path.Join(env.Getenv("WORK"), "home"))
	// Set home to test's dir to init config for the test
	os.Setenv("HOME", path.Join(env.Getenv("WORK"), "home"))

	// Create config folder
	os.MkdirAll(path.Join(os.Getenv("HOME"), ".config"), 0777)

	config.InitConfig("")

	config.AddAccount(user1Account)
	config.SetCurrentAccount(0)
	config.Write()

	return nil
}

func TestInitCommand(t *testing.T) {

	testscript.Run(t, testscript.Params{
		Dir:   "./init/",
		Setup: SetupFunc,
	})
}
