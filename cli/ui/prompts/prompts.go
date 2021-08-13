package prompts

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

func PromptRole(memberId string, roles []models.Role) (models.Role, error) {

	templates := &promptui.SelectTemplates{
		Label: "Role for {{ . }}?",
		Active: fmt.Sprintf(
			"%s {{ .Name | underline }}",
			promptui.IconSelect,
		),
		Inactive: "  {{ .Name }}",
		Selected: fmt.Sprintf(
			`{{ "%s" | green }} {{ .Name | faint }}`,
			promptui.IconGood,
		),
		Details: `
--------- Role ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		role := roles[index]
		name := strings.Replace(strings.ToLower(role.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     memberId,
		Items:     roles,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()

	return roles[index], err
}

func Confirm(message string) bool {
	p := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	answer, err := p.Run()

	if err != nil {
		if err.Error() != "^C" && err.Error() != "" {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
	} else if strings.ToLower(answer) == "y" {
		return true
	}

	return false
}

func StringInput(message string, defaultValue string) string {
	p := promptui.Prompt{
		Label:   message,
		Default: defaultValue,
	}

	answer, err := p.Run()

	if err != nil {
		if err.Error() != "^C" && err.Error() != "" {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	return strings.Trim(answer, " ")
}

type SelectCIServiceItem struct {
	Name string
	Type string
}

func SelectCIService(items []SelectCIServiceItem) SelectCIServiceItem {
	prompt := promptui.Select{
		Label: "Select a CI service",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Active: fmt.Sprintf(
				"%s {{ .Name | underline }}",
				promptui.IconSelect,
			),
			Inactive: "  {{ .Name }}",
			Selected: fmt.Sprintf(
				`{{ "%s" | green }} {{ .Name | faint }}`,
				promptui.IconGood,
			),
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		if err.Error() != "^C" {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	return items[index]
}

func Select(message string, items []string) (index int, selected string) {
	prompt := promptui.Select{
		Label: message,
		Items: items,
	}

	index, selected, err := prompt.Run()

	if err != nil {
		if err.Error() != "^C" {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	return index, selected
}
