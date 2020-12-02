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
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/internal/errors"
	"github.com/wearedevx/keystone/pkg/core"
	. "github.com/wearedevx/keystone/ui"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Manage secret files",
	Long: `Manage secret files.

List tracked secret files:
  $ ks files
  Files tracked as secret files:

          config/wp-config.php
		  config/front.config.js
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error

		ctx := core.New(core.CTX_RESOLVE)
		ctx.MustHaveEnvironment(currentEnvironment)

		files := ctx.ListFiles()

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		if len(files) == 0 {
			if !quietOutput {
				Print(`No files are currently tracked as secret files.

To add files to secret files:
  $ ks files add <path-to-secret-file>
`)
			}
			return
		}

		if quietOutput {
			for _, file := range files {
				Print(file)
			}
			return
		}

		Print(RenderTemplate("files list", `Files tracked as secret files:

{{ range . }} {{- . | indent 8 }} {{ end }}
`, files))
	},
}

func init() {
	RootCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
