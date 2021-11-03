package display

import (
	"fmt"

	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/ui"
)

type envListViewModel struct {
	Current      string
	Environments []string
}

func Environment(environment string) {
	ui.Print(environment)
}

func EnvironmentUsing(environmentName string) {
	ui.Print(ui.RenderTemplate("using env", `
{{ OK }} {{ .Message | bright_green }}

To load its variables:
  $ eval "$(ks source)"
`, map[string]string{
		"Message": fmt.Sprintf("Using the '%s' environment", environmentName),
		"EnvName": environmentName,
	}))
}

func EnvironmentList(environments []string, currentEnvironment string) {
	template := `{{ range  .Environments }}
{{ if eq . $.Current }} {{ "*" | blue }} {{ . | yellow }} {{ else }}   {{ . }} {{ end }} {{ end }}`

	ui.Print(ui.RenderTemplate("list env", template, envListViewModel{
		Current:      currentEnvironment,
		Environments: environments,
	}))
}

func EnvironmentSendSuccess() {
	ui.Print(
		ui.RenderTemplate(
			"send success",
			`{{ OK }} {{ "Environments sent successfully to members" | green }}`,
			nil,
		),
	)
}

// ———— PRIVATE UTILITIES ———— //

func pathList(files []keystonefile.FileKey) []string {
	r := make([]string, len(files))

	for index, file := range files {
		r[index] = file.Path
	}

	return r
}
