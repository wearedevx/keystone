package prompts

import (
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/cli/ui"
)

// Base confirm prompt
// `message` is the question to confirm (e.g. "Continue").
// promptui will append a question mark at the end
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

// Ask the user to enter a free form input
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

// Ask the user to enter a free form input,
// but runs the user input against a validation function
func StringInputWithValidation(
	message string,
	defaultValue string,
	validation promptui.ValidateFunc,
) string {
	p := promptui.Prompt{
		Label:    message,
		Default:  defaultValue,
		Validate: validation,
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

// Default selection user interface.
// `message` is a message to give some context to the user
// `items` is a list of things to select from.
// Return the index of the item selected in the `items` list, and
// the string
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
