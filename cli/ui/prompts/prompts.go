package prompts

import (
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

func PromptRole(memberId string, roles []models.Role) (models.Role, error) {

	templates := &promptui.SelectTemplates{
		Label:    "Role for {{ . }}?",
		Active:   " {{  .Name  }}",
		Inactive: " {{  .Name | faint }}",
		Selected: " {{ .Name }}",
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
		if err.Error() != "^C" {
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
		if err.Error() != "^C" {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	return strings.Trim(answer, " ")
}
