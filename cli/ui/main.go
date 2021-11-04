package ui

import (
	"fmt"
	"os"

	aurora "github.com/logrusorgru/aurora/v3"
)

type Printer interface {
	Print(message string, args ...interface{})
	PrintStdErr(message string, args ...interface{})
}

type UiPrinter struct{}

// Print method print utility with aurora color support
// Mimics `fmt.Printf()`, but always adds a new line
func (up *UiPrinter) Print(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(formatted)
}

// PrintStdErr method same as `Print()` but outputs to `stderr`
func (up *UiPrinter) PrintStdErr(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Fprintln(os.Stderr, formatted)
}

// PrintError function pretty-prints an error to stderr
func PrintError(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	displayable := RenderTemplate("Error", `
{{ ERROR }} {{ .Message | red }}
`, map[string]string{
		"Message": formatted,
	})

	fmt.Fprintln(os.Stderr, displayable)
}

// PrintSuccess function pretty prints a success message
func PrintSuccess(messageString string, args ...interface{}) {

	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	displayable := RenderTemplate("Success", `
{{ OK }} {{ .Message | bright_green }}
`, map[string]string{
		"Message": formatted,
	})

	fmt.Println(displayable)
}

// PrintBox function prints in a box
func PrintBox(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(aurora.Green(Box(formatted)))
}

// PrintDim function prints dimmed
func PrintDim(messageString string, args ...interface{}) {
	colored := aurora.Gray(11, messageString)

	formatted := colored.String()
	if len(args) > 0 {
		formatted = aurora.Sprintf(colored, args...)
	}

	fmt.Println(formatted)
}

// PrintStdErr function prints to stderr
func PrintStdErr(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Fprintln(os.Stderr, formatted)
}

// Print function prints
func Print(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(formatted)
}
