// build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	templatesDir := filepath.Join(wd, "templates")
	files, err := ioutil.ReadDir(templatesDir)

	sb := new(strings.Builder)

	sb.WriteString(`package templates

import (
	"fmt"
	"html/template"
	"strings"
)

var templates map[string]*template.Template = buildTemplates()

func buildTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template, 0)
`)

	bcontents, err := ioutil.ReadFile(filepath.Join(templatesDir, "layout.html"))
	if err != nil {
		panic(err)
	}

	sb.WriteString(fmt.Sprintf("\tlayout := `%s`", string(bcontents)))

	for _, file := range files {
		name := file.Name()
		if name != "layout.html" && strings.HasSuffix(name, ".html") {
			bcontents, err := ioutil.ReadFile(filepath.Join(templatesDir, name))
			if err != nil {
				panic(err)
			}

			tplName := strings.ReplaceAll(name, ".html", "")
			sb.WriteString(fmt.Sprintf(`
	templates["%s"] = template.Must(template.Must(
		template.
		New("%s").
		Parse(`+"`%s`"+`)).
		New("layout").
		Parse(layout))
`, tplName, tplName, string(bcontents)))
		}
	}

	sb.WriteString(`
	return templates
}`)

	sb.WriteString(`

func RenderTemplate(name string, data interface{}) (string, error) {
	sb := new(strings.Builder)

	if tpl, ok := templates[name]; ok {
		tpl.ExecuteTemplate(sb, "layout", data)
	} else {
		return "", fmt.Errorf("No such template: %s", name)
	}

	return sb.String(), nil
}
`)

	ioutil.WriteFile(filepath.Join(templatesDir, "generated_templates.go"), []byte(sb.String()), 0o644)
}
