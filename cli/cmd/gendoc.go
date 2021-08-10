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
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// gendocCmd represents the gendoc command
// var gendocCmd = &cobra.Command{
// 	Use:   "gendoc",
// 	Short: "A brief description of your command",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:
//
// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("gendoc called")
//
// 		doc.GenMarkdownTree(cmd, "./doc")
// 	},
// }

type documentationType string

const (
	Hugo documentationType = "hugo"
	Md                     = "md"
	Man                    = "man"
)

var doctype string
var destination string

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
				err = doc.GenMarkdownTreeCustom(rootCmd, destination, prepender, linkHandler)

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

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
menu:
  docs:
    parent: "%s"
---
`

func prepender(filename string) string {
	now := time.Now().Format(time.RFC3339)
	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	url := "/docs/cli/" + strings.ToLower(base) + "/"
	parent := "cli"

	return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url, parent)
}

func linkHandler(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	return "/docs/cli/" + strings.ToLower(base) + "/"
}

func init() {
	genCmd := newGenDocCommand(RootCmd)

	genCmd.Flags().StringVarP(&doctype, "type", "t", "md", "either 'hugo' or 'md'")
	genCmd.Flags().StringVarP(&destination, "destination", "d", "./doc", "target directory")

	RootCmd.AddCommand(genCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
