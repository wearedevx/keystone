// +build ignore

package main

import (
	. "fmt"
	. "io/ioutil"
	"os"
	"os/exec"
	"path"
	. "strings"

	"gopkg.in/yaml.v2"
)

type Param struct {
	Name string
	Typ  string `yaml:"type"`
}

type ErrorDef struct {
	Typ      string `yaml:"type"`
	Name     string
	Params   []Param
	Template string
}

type DefFile struct {
	Errors []ErrorDef
}

func open(p string) DefFile {
	var def DefFile

	f, err := ReadFile(p)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(f, &def)

	if err != nil {
		panic(err)
	}

	return def
}

func main() {
	wd, err := os.Getwd()
	filepath := path.Join(wd, "internal", "errors", "errors.yml")

	if err != nil {
		panic(err)
	}

	defs := open(filepath)

	var sb Builder
	sb.WriteString(`
package errors

`)

	kv_pairs := make([]string, 0)
	for _, def := range defs.Errors {
		key := def.Typ
		value := def.Template

		kv_pairs = append(kv_pairs, Sprintf("  \"%s\": `\n%s\n`", key, value))
	}

	var helpMap string
	if len(kv_pairs) == 0 {
		helpMap = "var helpTexts stristring]string = map[string]string {}\n\n"
	} else {
		helpMap = Sprintf("var helpTexts map[string]string = map[string]string {\n %s,\n }\n\n", Join(kv_pairs, ",\n"))
	}

	sb.WriteString(helpMap)

	for _, def := range defs.Errors {
		// Function signaturE
		sb.WriteString(Sprintf("func %s (", def.Typ))

		// Params
		for _, param := range def.Params {
			p := Sprintf("%s %s, ", ToLower(param.Name), param.Typ)

			sb.WriteString(p)
		}

		sb.WriteString("cause error)")

		// Return type
		sb.WriteString(" *Error {\n")

		kv_pairs = make([]string, 0)

		for _, param := range def.Params {
			typ := param.Typ
			name := param.Name
			argName := ToLower(name)

			kv_pairs = append(kv_pairs, Sprintf("  \"%s\": %s(%s)", name, typ, argName))
		}

		var metaMap string
		if len(kv_pairs) == 0 {
			metaMap = "  meta := map[string]string{}\n\n"
		} else {
			metaMap = Sprintf("meta := map[string]string{\n  %s,\n}\n", Join(kv_pairs, ", \n"))
		}

		sb.WriteString(metaMap)

		sb.WriteString(Sprintf("  return NewError(\"%s\", helpTexts[\"%s\"], meta, cause)\n", def.Name, def.Typ))

		// End function
		sb.WriteString("}\n\n")

	}

	outputPath := path.Join(wd, "internal", "errors", "generated_errors.go")
	WriteFile(outputPath, []byte(sb.String()), 0o644)

	c := exec.Command("gofmt", "-w", outputPath)
	c.Run()

}
