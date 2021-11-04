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

// Environment function display environment name
func Environment(environment string) {
	ui.Print(environment)
}

// EnvironmentUsing function Message after switch
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

// EnvironmentList function displays a list of environment
// Emphasize the current one
func EnvironmentList(environments []string, currentEnvironment string) {
	template := `{{ range  .Environments }}
{{ if eq . $.Current }} {{ "*" | blue }} {{ . | yellow }} {{ else }}   {{ . }} {{ end }} {{ end }}`

	ui.Print(ui.RenderTemplate("list env", template, envListViewModel{
		Current:      currentEnvironment,
		Environments: environments,
	}))
}

// EnvironmentSendSuccess function Message when sharing envirionments is
// successfull
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
