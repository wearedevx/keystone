package prompts

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

// ———— MEMBERS PROMTS ———— //

// ConfirmRevokeAccessToMember asks the user to confirm they
// want to revoke the access to the given member,
// unless `forceYes` is true, in which case it returns true without
// asking the user.
func ConfirmRevokeAccessToMember(memberId string, forceYes bool) bool {
	if forceYes {
		return true
	}

	return Confirm("Revoke access to " + memberId)
}

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
		name := strings.ReplaceAll(strings.ToLower(role.Name), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

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

// ———— DEVICE PROMPTS ———— //

// Asks the user to select from a list of devices
func SelectDevice(devices []models.Device) models.Device {
	items := make([]map[string]string, 0)

	for _, device := range devices {
		newItem := make(map[string]string)
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

	if !Confirm(
		fmt.Sprintf("Sure you want to revoke %s", devices[index].Name),
	) {
		os.Exit(0)
	}

	return devices[index]
}

// DeviceName asks the user to enter a device name.
// If there is no existing name, it will use the device hostname as default
func DeviceName(existingName string, forceDefault bool) string {
	if existingName == "" {
		var defaultName string

		if hostname, err := os.Hostname(); err == nil {
			defaultName = hostname
		}

		if forceDefault {
			return defaultName
		}

		validate := func(input string) error {
			matched, err := regexp.MatchString(`^[a-zA-Z0-9\.\-\_]{1,}$`, input)
			if !matched {
				return errors.New(
					"incorrect device name. Device name must be alphanumeric with ., -, _",
				)
			}
			return err
		}

		deviceName := StringInputWithValidation(
			"Enter the name you want this device to have",
			defaultName,
			validate,
		)

		return deviceName
	}

	return existingName
}

// ———— CI SERVICE PROMPTS ———— //

func ServiceIntegrationName() string {
	return StringInput(
		"Enter a name for your integration",
		"",
	)
}

func ServiceConfigurationToRemove() string {
	return StringInput(
		"Enter the service name to remove",
		"",
	)
}

func ConfirmCiConfigurationRemoval(serviceName string) bool {
	ui.Print(ui.RenderTemplate("careful rm ci", `
{{ CAREFUL }} You are about to remove the {{ . }} CI service.

This cannot be undone.`,
		serviceName))

	return Confirm("Continue")
}

func ConfirmSendEnvironmentToCiService(environmentName string) bool {
	ui.Print(
		"You are about to send the '%s' environment to your CI services.",
		environmentName,
	)

	return Confirm("Continue")
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

// ———— PROJECT PROMPTS ———— //

func ConfirmProjectDestruction(projectName string) bool {
	ui.Print(ui.RenderTemplate("confirm project destroy",
		`{{ CAREFUL }} You are about to destroy the {{ .Project }} project.
Secrets and files managed by Keystone WILL BE LOST. Make sure you have backups.

Members of the project will no longer be able to get the latest updates,
or share secrets between them.

This is permanent, and cannot be undone.
`, map[string]string{
			"Project": projectName,
		}))

	result := StringInput(
		"Type the project name to confirm its destruction",
		"",
	)

	// expect result to be the project name
	return projectName == result
}

// ———— ORGANIZATION PROMPTS ————— //

// Asks the usre to select from a list of organizations
func OrganizationsSelect(
	organizations []models.Organization,
) models.Organization {
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
		name := strings.ReplaceAll(strings.ToLower(orga.Name), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

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

// ———— BACKUP/RESTORE PROMPTS ————— //

func PasswordToEncrypt() string {
	return StringInput("Password to encrypt backup", "")
}

func PasswordToDecrypt() string {
	return StringInput("Password to decrypt backup", "")
}

func ConfirmDotKeystonDirRemoval() bool {
	ui.Print(ui.RenderTemplate(
		"confirm files rm",
		`{{ CAREFUL }} You are about to remove the content of .keystone/ which contain all your local secrets and files.
This will override the changes you and other members made since the backup.
It will update other members secrets and files.`,
		map[string]string{},
	))
	return Confirm("Continue")
}

// ——— FILES PROMPTS ———— //

func ConfirmOverrideFileContents() bool {
	return Confirm("Do you want to overrid the contents")
}

func ConfirmFileReset(forceYes bool) bool {
	ui.Print(ui.RenderTemplate(
		"careful reset",
		`{{ CAREFUL }} {{ "Local changes will be lost" | yellow }}
The content of the files you are resetting will be replaced by their cached content.`,
		nil,
	))

	if forceYes {
		return true
	}

	return Confirm("Continue")
}

func ConfirmFileRemove(filePath, environmentName string, forceYes bool) bool {
	if forceYes {
		return true
	}

	ui.Print(ui.RenderTemplate(
		"confirm files rm",
		`{{ CAREFUL }} You are about to remove {{ .Path }} from the secret files.
Its current content will be kept locally.
Its content for other environments will be lost, it will no longer be gitignored.
This is permanent, and cannot be undone.`,
		map[string]string{
			"Path":        filePath,
			"Environment": environmentName,
		},
	))

	return Confirm("Continue")
}

// ——— SECRETS PROMPTS ——— //

func ConfirmOverrideSecretValue(forceYes bool) bool {
	if forceYes {
		return true
	}

	return Confirm("Do you want to overrid the value")
}

func ValueForEnvironment(
	secretName, environmentName, defaultValue string,
) string {
	ui.Print(
		"Enter the value of '%s' for the '%s' environment",
		secretName,
		environmentName,
	)

	return StringInput(secretName, defaultValue)
}

// ——— LOGIN PROMPTS ———— //

func SelectAuthService(serviceName string) string {
	if serviceName == "" {
		_, serviceName = Select(
			"Select an identity provider",
			[]string{
				"github",
				"gitlab",
			},
		)
	}

	return serviceName
}
