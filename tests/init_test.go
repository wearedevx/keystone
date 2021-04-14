package main

import (
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
	Repo.Connect()

	var user1 *User = &User{
		ExtID:       "56883564",
		UserID:      "00fb7666-de43-4559-b4e4-39b172117dd8",
		AccountType: "github",
		Username:    "LAbigael",
		Fullname:    "Test user1",
		Email:       "abigael.laldji@protonmail.com",
	}

	token, _ := MakeToken(*user1)

	Repo.GetOrCreateUser(user1)

	homeDir := path.Join(env.Getenv("WORK"), "home")
	// Set home dir for test
	env.Setenv("HOME", homeDir)

	// Create config folder
	os.MkdirAll(path.Join(homeDir, ".config"), 0777)

	ioutil.WriteFile(path.Join(homeDir, ".config", "keystone.yaml"), []byte(`
accounts:
- Fullname: Test user1
  account_type: github
  email: abigael.laldji@protonmail.com
  ext_id: "56883564"
  fullname: Michel
  user_id: 00fb7666-de43-4559-b4e4-39b172117dd8
  username: LAbigael
auth_token: `+token+`
current: 0
`), 0o777)

	return nil
}

func TestInitCommand(t *testing.T) {

	testscript.Run(t, testscript.Params{
		Dir:   "./init/",
		Setup: SetupFunc,
	})
}
