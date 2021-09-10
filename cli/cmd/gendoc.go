/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

// hels to determine output format
type documentationType string

const (
	Hugo documentationType = "hugo"
	Md                     = "md"
	Man                    = "man"
)

// documentation format
var doctype string

// output directory
var destination string

// init initializes the package
func init() {
	genCmd := newGenDocCommand(RootCmd)

	genCmd.Flags().StringVarP(&doctype, "type", "t", "md", "either 'hugo' or 'md'")
	genCmd.Flags().StringVarP(&destination, "destination", "d", "./doc", "target directory")

	RootCmd.AddCommand(genCmd)
}

// Generates a new Documentation Generation Command
func newGenDocCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "documentation",
		Hidden: true,
		Short:  "Generate keystone documentation",
		Long:   "Generate keystone documentation as markdown or man page",
		Example: `keystone documentation md
keystone documentation man`,
		Args: cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("Doc generation command")
			var err error

			switch documentationType(doctype) {
			case Hugo:
				err = doc.GenMarkdownTreeCustom(rootCmd, destination, hugoFrontMatterPrepender, hugoLinkHandler)

			case Md:
				err = doc.GenMarkdownTree(rootCmd, destination)

			case Man:
				err = doc.GenManTree(
					rootCmd,
					&doc.GenManHeader{
						Title:   "ks",
						Section: "1",
					},
					destination,
				)
			}

			if err != nil {
				panic(err)
			}
		},
	}
	return cmd
}

/*********************/
/* Hugo Link handler */
/*********************/

func hugoLinkHandler(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	return "/docs/cli/" + strings.ToLower(base) + "/"
}

/**************************/
/* Front Matter Prepender */
/**************************/

// fmTemplate is the Front Matter template used
// with Hugo static site generator
const fmTemplate = `---
date: %s
title: "%s"
description: |-
%s

slug: %s
url: %s
menu:
  docs:
    parent: "%s"
---
`

// hugoFrontMatterPrepender returns the Front Matter for the given markdown file
func hugoFrontMatterPrepender(filename string) string {
	now := time.Now().Format(time.RFC3339)
	name := filepath.Base(filename)
	description := getDescription(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	url := "/docs/cli/" + strings.ToLower(base) + "/"
	parent := "cli"

	return fmt.Sprintf(
		fmTemplate,
		now,
		strings.Replace(base, "_", " ", -1),
		description,
		base,
		url,
		parent,
	)
}

/*********************/
/* Hugo utilities */
/*********************/

// getDescription returns a descrition for a command based on the "Short"
// property of the command associated with the file `filename`
func getDescription(filename string) (description string) {
	source := findMatchingSource(filename)

	if utils.FileExists(source) {
		command := findCommand(source)
		if command == nil {
			log.Fatal(fmt.Errorf("Command not found in %s", source))
		}

		description = findStringValueAt(command, "Short")

		description = padMultiline(description)
	}

	return description
}

// findMatchingSource finds the matching source for a markdown file
// as given by the prepender callback.
// ALSO: getDescription()
func findMatchingSource(filename string) (source string) {
	filename = filepath.Base(filename)
	filename = strings.TrimPrefix(filename, "ks_")
	filename = strings.Replace(filename, ".md", ".go", -1)
	filename = strings.Replace(filename, "-", "_", -1)

	if filename == "ks.go" {
		filename = "root.go"
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	source = path.Join(wd, "cmd", filename)

	return source
}

// findCommand finds the Cobra command declared in the given source file
// and returns the associated AST node
// Exits with error if no command could be found
// ALSO: getDescription()
func findCommand(source string) (command *ast.CompositeLit) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		d, ok := n.(*ast.Ident)
		if ok && strings.HasSuffix(d.Name, "Cmd") && d.Obj != nil {
			a, ok := d.Obj.Decl.(*ast.ValueSpec)
			if ok {
				for _, v := range a.Values {
					e, _ := v.(*ast.UnaryExpr)
					command = e.X.(*ast.CompositeLit)
					return false
				}
			}
		}
		return true
	})

	return command
}

// findStringValueAt returns the values associated with `field` in the given
// composite literal
// NOTE: The implementation of the string fetching is quite naive and
// could break easily if the input gets too exotic.
// It currenty supports values that are plain string literals and
// concatenation of string literals; anything beyond thait is likely to break
// appart
func findStringValueAt(
	command *ast.CompositeLit,
	field string,
) (value string) {
	elmts := command.Elts

	for _, e := range elmts {
		kv := e.(*ast.KeyValueExpr)
		key, _ := kv.Key.(*ast.Ident)
		if key.Name == field {
			v, ok := kv.Value.(*ast.BinaryExpr)
			if ok {
				value = recomposeString(v)
				break
			}
			vb, ok := kv.Value.(*ast.BasicLit)
			if ok {
				value = vb.Value
				break
			}

		}
	}

	value = unQuote(value)

	return value
}

// unQuote removes starting and and ending quotes and back-ticks
// form any string
func unQuote(in string) string {
	in = strings.TrimPrefix(in, "`")
	in = strings.TrimSuffix(in, "`")
	in = strings.TrimPrefix(in, "\"")
	in = strings.TrimSuffix(in, "\"")

	return in
}

// padMultiline prepends a double space before every line
// while removing any leading tab
func padMultiline(in string) (out string) {
	parts := strings.Split(in, "\n")
	for i, part := range parts {
		parts[i] = "  " + strings.TrimPrefix(part, "\t")
	}

	out = strings.Join(parts, "\n")

	return out
}

// recomposeString makes a string out of binary Expressions
// such as string concatenation
func recomposeString(binaryExpression *ast.BinaryExpr) (result string) {
	var left string
	var right string

	l, ok := binaryExpression.X.(*ast.BasicLit)
	if !ok {
		lbin := binaryExpression.X.(*ast.BinaryExpr)
		left = recomposeString(lbin)
	} else {
		left = l.Value
	}

	r, ok := binaryExpression.Y.(*ast.BasicLit)
	if !ok {
		rbin := binaryExpression.Y.(*ast.BinaryExpr)
		right = recomposeString(rbin)
	} else {
		right = r.Value
	}

	result = strings.ReplaceAll(left+right, "`\"`", "`")
	return result
}
