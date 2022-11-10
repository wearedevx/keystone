//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type Error struct {
	Symbol   string `yaml:"symbol"`
	Value    string `yaml:"value"`
	HasCause bool   `yaml:"has_cause"`
}

type DefFile struct {
	Errors []Error `yaml:"errors"`
}

func open(p string) DefFile {
	var def DefFile

	f, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(f, &def)

	if err != nil {
		panic(err)
	}

	return def
}

func write(contents string, outpath string) error {
	fmt.Printf("Writing: %+v\n", outpath)
	err := ioutil.WriteFile(outpath, []byte(contents), 0o644)
	if err != nil {
		return err
	}

	c := exec.Command("gofmt", "-w", outpath)
	err = c.Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	wd, err := os.Getwd()
	filepath := path.Join(wd, "generators", "errors.yaml")

	if err != nil {
		panic(err)
	}

	defs := open(filepath)

	pkgApiErrorsPath := path.Join(wd, "pkg", "apierrors", "generated_errors.go")
	apierrors := createPkgApiErrors(defs.Errors)

	err = write(apierrors, pkgApiErrorsPath)
	if err != nil {
		panic(err)
	}

	internalErrorsPath := path.Join(
		wd,
		"internal",
		"errors",
		"generated_errors.go",
	)
	internalErrors := createInternalErrors(defs.Errors)

	err = write(internalErrors, internalErrorsPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}

func createPkgApiErrors(definitions []Error) string {
	var sb strings.Builder
	sb.WriteString(`
package apierrors

import "errors"

var (
`)

	for _, error := range definitions {
		sb.WriteString(fmt.Sprintf(
			"%s error = errors.New(\"%s\")\n",
			error.Symbol,
			error.Value,
		))
	}

	sb.WriteString(`)

func FromString(s string) error {
	switch s {
`)

	for _, error := range definitions {
		sb.WriteString(fmt.Sprintf(`case %s.Error():
		return %s
`, error.Symbol, error.Symbol))
	}

	sb.WriteString(`default:
		return errors.New(s)
	}
}`)

	return sb.String()
}

func createInternalErrors(definitions []Error) string {
	var sb strings.Builder
	sb.WriteString(`package errors

`)

	for _, error := range definitions {
		if error.HasCause {
			sb.WriteString(fmt.Sprintf(`func %s(cause error) error {
	return newError("%s", cause)
}

`, error.Symbol, error.Value))
		} else {
			sb.WriteString(fmt.Sprintf(`func %s() error {
	return newError("%s", nil)
}

`, error.Symbol, error.Value))
		}
	}

	return sb.String()
}
