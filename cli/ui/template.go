package ui

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	. "github.com/logrusorgru/aurora/v3"
)

func RenderTemplate(name string, tpl string, viewData interface{}) string {
	var buf bytes.Buffer

	t := template.Must(template.New(name).Funcs(functions).Parse(tpl))

	err := t.Execute(&buf, viewData)

	if err != nil {
		println(fmt.Sprintf("Failed to render template %s (%s)", name, err.Error()))
		os.Exit(1)
	}

	return buf.String()
}

var functions template.FuncMap = template.FuncMap{
	"ERROR": func() string {
		return Black(" ERROR ").BgRed().String()
	},
	"CAREFUL": func() string {
		return Black(" CAREFUL ").BgYellow().String()
	},
	"OK": func() string {
		return Black(" OK ").BgGreen().String()
	},
	"box": func(p string) string {
		return Box(p)
	},
	"indent": func(amount int, p string) string {
		indent := ""

		for i := 0; i < amount; i++ {
			indent = indent + " "
		}

		lines := strings.Split(p, "\n")

		indented := make([]string, 0)
		for _, line := range lines {
			indented = append(indented, indent+line)
		}

		return strings.Join(indented, "\n")
	},
	"black": func(p string) string {
		return Black(p).String()
	},
	"red": func(p string) string {
		return Red(p).String()
	},
	"green": func(p string) string {
		return Green(p).String()
	},
	"yellow": func(p string) string {
		return Yellow(p).String()
	},
	"blue": func(p string) string {
		return Blue(p).String()
	},
	"magenta": func(p string) string {
		return Magenta(p).String()
	},
	"cyan": func(p string) string {
		return Cyan(p).String()
	},
	"white": func(p string) string {
		return White(p).String()
	},
	"bright_black": func(p string) string {
		return BrightBlack(p).String()
	},
	"bright_red": func(p string) string {
		return BrightRed(p).String()
	},
	"bright_green": func(p string) string {
		return BrightGreen(p).String()
	},
	"bright_yellow": func(p string) string {
		return BrightYellow(p).String()
	},
	"bright_blue": func(p string) string {
		return BrightBlue(p).String()
	},
	"bright_magenta": func(p string) string {
		return BrightMagenta(p).String()
	},
	"bright_cyan": func(p string) string {
		return BrightCyan(p).String()
	},
	"bright_white": func(p string) string {
		return BrightWhite(p).String()
	},
	"bg_black": func(p string) string {
		return BgBlack(p).String()
	},
	"bg_red": func(p string) string {
		return BgRed(p).String()
	},
	"bg_green": func(p string) string {
		return BgGreen(p).String()
	},
	"bg_yellow": func(p string) string {
		return BgYellow(p).String()
	},
	"bg_blue": func(p string) string {
		return BgBlue(p).String()
	},
	"bg_magenta": func(p string) string {
		return BgMagenta(p).String()
	},
	"bg_cyan": func(p string) string {
		return BgCyan(p).String()
	},
	"bg_white": func(p string) string {
		return BgWhite(p).String()
	},
	"bg_bright_black": func(p string) string {
		return BgBrightBlack(p).String()
	},
	"bg_bright_red": func(p string) string {
		return BgBrightRed(p).String()
	},
	"bg_bright_green": func(p string) string {
		return BgBrightGreen(p).String()
	},
	"bg_bright_yellow": func(p string) string {
		return BgBrightYellow(p).String()
	},
	"bg_bright_blue": func(p string) string {
		return BgBrightBlue(p).String()
	},
	"bg_bright_magenta": func(p string) string {
		return BgBrightMagenta(p).String()
	},
	"bg_bright_cyan": func(p string) string {
		return BgBrightCyan(p).String()
	},
	"bg_bright_white": func(p string) string {
		return BgBrightWhite(p).String()
	},
}
