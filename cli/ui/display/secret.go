package display

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

func SecretTable(secrets []core.Secret, environments []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Options.SeparateRows = true
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})

	topHeader := table.Row{"Secret name"}
	envHeader := table.Row{""}

	for _, environment := range environments {
		topHeader = append(topHeader, "Environments")
		envHeader = append(envHeader, environment)
	}

	t.AppendHeader(topHeader, table.RowConfig{AutoMerge: true})
	t.AppendHeader(envHeader)

	for _, secret := range secrets {
		name := secret.Name

		if secret.Required {
			name = name + " *"
		}
		if secret.FromCache {
			name = name + " A"
		}

		row := table.Row{name}

		for _, environment := range environments {
			value := secret.Values[core.EnvironmentName(environment)]

			if len(value) > 40 {
				value = value[:40]
				value += "..."
			}

			row = append(row, value)
		}

		t.AppendRow(row)
	}

	t.Render()
	fmt.Println(" * Required secrets; A Available secrets")
}

func SecretAlreadyExitsts(values map[core.EnvironmentName]core.SecretValue) {
	ui.Print(`The secret already exist. Values are:`)
	for env, value := range values {
		ui.Print(`%s: %s`, env, value)
	}
}

func SecretIsSetForEnvironment(secretName string, nbEnvironments int) {
	ui.PrintSuccess(
		"Secret '%s' is set for %d environment(s)",
		secretName,
		nbEnvironments,
	)
}

func SecretRemoved(secretName string) {
	ui.PrintSuccess("Secret '%s' removed", secretName)
}

func SecretUpdated(secretName, environmentName string) {
	ui.PrintSuccess(
		fmt.Sprintf(
			"Secret '%s' updated for the '%s' environment",
			secretName,
			environmentName,
		),
	)
}

const (
	OPTIONAL string = "optional"
	REQUIRED string = "required"
)

func SecretIsNow(secretName, prop string) {
	template := `Secret {{ .SecretName }} is now {{ .Prop }}.`

	ui.Print(
		ui.RenderTemplate(
			"set secret optional",
			template,
			struct {
				SecretName string
				Prop       string
			}{
				SecretName: secretName,
				Prop:       prop,
			},
		),
	)

	if prop == REQUIRED {
		ui.Print(`If you have setup a CI service, don’t forget to run:
  $ ks ci send
`)
	}
}
