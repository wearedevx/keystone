package display

import (
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

var envList = []string{
	"dev",
	"staging",
	"prod",
}

// TODO: should handle a `quiet` setting ?
// printChanges displays changes for environments to the user
func Changes(
	changes core.ChangesByEnvironment,
) {
	for _, environmentName := range envList {
		changesList, ok := changes.Environments[environmentName]
		if !ok { // means there are no changes, and versions are equal
			printEnvironmentUpToDate(environmentName)
			continue
		}

		if changesList.IsSingleVersionChange() {
			printChangesButNoMessage(environmentName)
			continue
		}

		if !changesList.IsEmpty() {
			printChangeList(environmentName, changesList)
		}

		printEnvironmentUpToDate(environmentName)
	}
}

func printChangesButNoMessage(environmentName string) {
	ui.PrintStdErr(
		"Environment %s has changed but no message available. Ask someone to push their secret ⨯",
		environmentName,
	)
}

func printEnvironmentUpToDate(environmentName string) {
	ui.PrintStdErr("Environment " + environmentName + " up to date ✔")
}

func printChangeList(environmentName string, changes core.Changes) {
	ui.PrintStdErr(
		"Environment %s: %d secret(s) changed",
		environmentName,
		len(changes),
	)

	for _, change := range changes {
		printChange(change)
	}
}

func printChange(change core.Change) {
	// No previous cotent => secret is new
	switch {
	case change.IsSecretAdd():
		ui.PrintStdErr(
			ui.RenderTemplate(
				"secret added",
				` {{ "++" | green }} {{ .Secret }} : {{ .To }}`,
				map[string]string{
					"Secret": change.Name,
					"From":   change.From,
					"To":     change.To,
				},
			),
		)

	case change.IsSecretDelete():
		ui.PrintStdErr(
			ui.RenderTemplate(
				"secret deleted",
				` {{ "--" | red }} {{ .Secret }} deleted.`,
				map[string]string{
					"Secret": change.Name,
				},
			),
		)

	case change.IsSecretChange():
		ui.PrintStdErr(
			"   " + change.Name + " : " + change.From + " ↦ " + change.To,
		)
	}
}
