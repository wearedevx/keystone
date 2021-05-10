package prompts

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/api/pkg/models"
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
