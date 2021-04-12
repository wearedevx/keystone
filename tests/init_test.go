package main

import (
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/wearedevx/keystone/cmd"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"toto": cmd.Execute,
	}))
}
func TestHelloWorld(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: ".",
	})
}
