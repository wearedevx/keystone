package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
	. "github.com/wearedevx/keystone/internal/jwt"
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

	var user1 *User = &User{
		ExtID:       "56883564",
		UserID:      "00fb7666-de43-4559-b4e4-39b172117dd8",
		AccountType: "github",
		Username:    "LAbigael",
		Fullname:    "Test user1",
		Email:       "abigael.laldji@protonmail.com",
	}

	Repo.GetOrCreateUser(user1)

	token, _ := MakeToken(*user1)

	homeDir := path.Join(env.Getenv("WORK"), "home")
	configDir := path.Join(homeDir, ".config")
	pathToKeystoneFile := path.Join(configDir, "keystone.yaml")

	// Set home dir for test
	env.Setenv("HOME", homeDir)

	env.Setenv("DB_PORT", os.Getenv("DB_PORT"))

	// Create config folder
	err := os.MkdirAll(configDir, 0777)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pathToKeystoneFile, []byte(`
accounts:
- Fullname: Test user1
  account_type: github
  email: abigael.laldji@protonmail.com
  ext_id: "56883564"
  fullname: Michel
  user_id: `+user1.UserID+`
  username: LAbigael
auth_token: `+token+`
current: 0
`), 0o777)

	if err != nil {
		return err
	}

	return nil
}

func TestCommands(t *testing.T) {

	fmt.Println("BEFOREEEEE")

	testscript.Run(t, testscript.Params{
		Dir:   "./init/",
		Setup: SetupFunc,
	})
	testscript.Run(t, testscript.Params{
		Dir:   "./env/",
		Setup: SetupFunc,
	})
	testscript.Run(t, testscript.Params{
		Dir:                  "./secrets/",
		Setup:                SetupFunc,
		IgnoreMissedCoverage: true,
	})
}
