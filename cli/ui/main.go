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

type EchoPrinter struct{}

type UiPrinter struct{}

func (ep *EchoPrinter) Print(messageString string, args ...interface{}) {
	messageString = "echo " + "\"" + messageString + "\""

	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(formatted)
}

func (ep *EchoPrinter) PrintStdErr(messageString string, args ...interface{}) {
	messageString = "echo " + "\"" + messageString + "\""

	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(formatted)
}

func (up *UiPrinter) Print(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(formatted)
}

func (up *UiPrinter) PrintStdErr(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Fprintln(os.Stderr, formatted)
}

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

func PrintBox(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(aurora.Green(Box(formatted)))
}

func PrintDim(messageString string, args ...interface{}) {
	colored := aurora.Gray(11, messageString)

	formatted := colored.String()
	if len(args) > 0 {
		formatted = aurora.Sprintf(colored, args...)
	}

	fmt.Println(formatted)
}

func PrintStdErr(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Fprintln(os.Stderr, formatted)
}

func Print(messageString string, args ...interface{}) {
	formatted := messageString
	if len(args) > 0 {
		formatted = aurora.Sprintf(messageString, args...)
	}

	fmt.Println(formatted)
}
