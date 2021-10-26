package prompts

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

// Prompts to select a role for a user
// `memberId` is a `username@service` userID
// `roles` is a list of roles to select from
// Returns the selected role
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

// Asks the user to select from a list of devices
func SelectDevice(devices []models.Device) models.Device {
	items := make([]map[string]string, 0)

	for _, device := range devices {
		var newItem = make(map[string]string, 0)
		newItem["Name"] = device.Name
		newItem["UID"] = device.UID
		newItem["CreatedAt"] = device.CreatedAt.Format("2006/01/02")

		if device.LastUsedAt.IsZero() {
			newItem["LastUsedAtString"] = "never used"
		} else {
			newItem["LastUsedAtString"] = fmt.Sprintf("last used on %s", device.CreatedAt.Format("2006/01/02"))
		}

		items = append(items, newItem)
	}

	prompt := promptui.Select{
		Label: "Select a Device to revoke",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Active: fmt.Sprintf(
				"%s {{ .Name | underline }}, {{ .LastUsedAtString  }}, added on {{ .CreatedAt  }}",
				promptui.IconSelect,
			),
			Inactive: "{{ .Name }}, {{ .LastUsedAtString  }}, added on {{ .CreatedAt  }}",
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

	if !Confirm(fmt.Sprintf("Sure you want to revoke %s", devices[index].Name)) {
		os.Exit(0)
	}

	return devices[index]
}

// Items for SelectCIService prompt
type SelectCIServiceItem struct {
	Name string
	Type string
}

// Asks the user to select from a list of CI services
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

// Asks the usre to select from a list of organizations
func OrganizationsSelect(organizations []models.Organization) models.Organization {
	templates := &promptui.SelectTemplates{
		Active: fmt.Sprintf(
			"%s {{ .Name | underline }}",
			promptui.IconSelect,
		),
		Inactive: "  {{ .Name }}",
		Selected: fmt.Sprintf(
			`{{ "%s" | green }} {{ .Name | faint }}`,
			promptui.IconGood,
		),
	}

	searcher := func(input string, index int) bool {
		orga := organizations[index]
		name := strings.Replace(strings.ToLower(orga.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Organization",
		Items:     organizations,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		if err.Error() != "^C" {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	return organizations[i]
}

func PasswordToEncrypt() string {
	return StringInput("Password to encrypt backup", "")
}

func PasswordToDecrypt() string {
	return StringInput("Password to decrypt backup", "")
}

func ConfirmDotKeystonDirRemoval() {
	ui.Print(ui.RenderTemplate("confirm files rm",
		`{{ CAREFUL }} You are about to remove the content of .keystone/ which contain all your local secrets and files.
This will override the changes you and other members made since the backup.
It will update other members secrets and files.`, map[string]string{}))
	if !Confirm("Continue") {
		os.Exit(0)
	}
}
