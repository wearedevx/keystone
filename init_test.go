package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"ks": cmd.Execute,
	}))
}
func Test_ExecuteCommand(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "./tests",
		Setup: func(env *testscript.Env) error {
			fmt.Println(env.Getenv("WORK"))
			fmt.Println(env.WorkDir)

			env.Setenv("HOME", path.Join(env.Getenv("WORK"), "home"))

			return nil
		},
	})
}
