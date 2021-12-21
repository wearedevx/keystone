package ui

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	aurora "github.com/logrusorgru/aurora/v3"
)

var au aurora.Aurora

func init() {
	colors := os.Getenv("KSCOLORS") != "off"

	au = aurora.NewAurora(colors)
}

// RenderTemplate function uility to render template string
func RenderTemplate(name string, tpl string, viewData interface{}) string {
	var buf bytes.Buffer

	t := template.Must(template.New(name).Funcs(functions).Parse(tpl))

	err := t.Execute(&buf, viewData)
	if err != nil {
		println(
			fmt.Sprintf("Failed to render template %s (%s)", name, err.Error()),
		)
		os.Exit(1)
	}

	return buf.String()
}

var functions template.FuncMap = template.FuncMap{
	"ERROR": func() string {
		return au.Black(" ERROR ").BgRed().String()
	},
	"CAREFUL": func() string {
		return au.Black(" CAREFUL ").BgYellow().String()
	},
	"OK": func() string {
		return au.Black(" OK ").BgGreen().String()
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
	"add": func(a, b int) string {
		return strconv.Itoa(a + b)
	},
	"black": func(p string) string {
		return au.Black(p).String()
	},
	"red": func(p string) string {
		return au.Red(p).String()
	},
	"green": func(p string) string {
		return au.Green(p).String()
	},
	"yellow": func(p string) string {
		return au.Yellow(p).String()
	},
	"blue": func(p string) string {
		return au.Blue(p).String()
	},
	"magenta": func(p string) string {
		return au.Magenta(p).String()
	},
	"cyan": func(p string) string {
		return au.Cyan(p).String()
	},
	"white": func(p string) string {
		return au.White(p).String()
	},
	"bright_black": func(p string) string {
		return au.BrightBlack(p).String()
	},
	"bright_red": func(p string) string {
		return au.BrightRed(p).String()
	},
	"bright_green": func(p string) string {
		return au.BrightGreen(p).String()
	},
	"bright_yellow": func(p string) string {
		return au.BrightYellow(p).String()
	},
	"bright_blue": func(p string) string {
		return au.BrightBlue(p).String()
	},
	"bright_magenta": func(p string) string {
		return au.BrightMagenta(p).String()
	},
	"bright_cyan": func(p string) string {
		return au.BrightCyan(p).String()
	},
	"bright_white": func(p string) string {
		return au.BrightWhite(p).String()
	},
	"bg_black": func(p string) string {
		return au.BgBlack(p).String()
	},
	"bg_red": func(p string) string {
		return au.BgRed(p).String()
	},
	"bg_green": func(p string) string {
		return au.BgGreen(p).String()
	},
	"bg_yellow": func(p string) string {
		return au.BgYellow(p).String()
	},
	"bg_blue": func(p string) string {
		return au.BgBlue(p).String()
	},
	"bg_magenta": func(p string) string {
		return au.BgMagenta(p).String()
	},
	"bg_cyan": func(p string) string {
		return au.BgCyan(p).String()
	},
	"bg_white": func(p string) string {
		return au.BgWhite(p).String()
	},
	"bg_bright_black": func(p string) string {
		return au.BgBrightBlack(p).String()
	},
	"bg_bright_red": func(p string) string {
		return au.BgBrightRed(p).String()
	},
	"bg_bright_green": func(p string) string {
		return au.BgBrightGreen(p).String()
	},
	"bg_bright_yellow": func(p string) string {
		return au.BgBrightYellow(p).String()
	},
	"bg_bright_blue": func(p string) string {
		return au.BgBrightBlue(p).String()
	},
	"bg_bright_magenta": func(p string) string {
		return au.BgBrightMagenta(p).String()
	},
	"bg_bright_cyan": func(p string) string {
		return au.BgBrightCyan(p).String()
	},
	"bg_bright_white": func(p string) string {
		return au.BgBrightWhite(p).String()
	},
}
