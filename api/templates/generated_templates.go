package templates

import (
	"fmt"
	"html/template"
	"strings"
)

var templates map[string]*template.Template = buildTemplates()

func buildTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template, 0)
	layout := `{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Keystone</title>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
		</style>
	</head>
	<body class="home">
		<header class="fixed-top navbar">
			<div class="container">
				<strong>Keystone</strong>
			</div>
		</header>	
		<div class="container" role="document">
			<div class="content">{{ template "content" .}}</div>
		</div>
	</body>
</html>	
{{ end }}
`
	templates["login-fail"] = template.Must(template.Must(
		template.
		New("login-fail").
		Parse(`{{ define "content" }}
<div class="error">
	<div class="title">{{ .Title }}</div>
	<div class="message">{{ .Message }}</div>
</div>
{{ end }}
`)).
		New("layout").
		Parse(layout))

	templates["login-success"] = template.Must(template.Must(
		template.
		New("login-success").
		Parse(`{{ define "content" }}
<div class="success">
	<div class="title">{{ .Title }}</div>
	<div class="message">{{ .Message }}</div>
</div>
{{ end }}
`)).
		New("layout").
		Parse(layout))

	return templates
}

func RenderTemplate(name string, data interface{}) (string, error) {
	sb := new(strings.Builder)

	if tpl, ok := templates[name]; ok {
		tpl.ExecuteTemplate(sb, "layout", data)
	} else {
		return "", fmt.Errorf("No such template: %s", name)
	}

	return sb.String(), nil
}
