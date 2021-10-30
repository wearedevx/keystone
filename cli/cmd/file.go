package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/ui/display"
)

var fileDisplayFilter string

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "file",
	Short: "Manages secret files",
	Long: `Manages secret files.

Without arguments, lists secret files:
` + "```" + `
$ ks file
  Files tracked as secret files:

    *  config/wp-config.php
     M config/front.config.js

  * Required; A Available; M Modified
` + "```" + `

Required files will stop ` + "`" + `ks source` + "`" + ` and ` + "`" + `ks ci send` + "`" + `
Available files exist in the cache, but not in the keystone.yaml file.
Modified files have different content from the cache.

For a machine parsable output, use the ` + "`" + `-q` + "`" + ` flag:
` + "```" + `
$ ks file -q
  *  config/wp-config.php
   M config/front.config.js
` + "```" + `

You can also filter output:
` + "```" + `
# Show only locally modified files
$ ks file -qf modified
   M config/front.config.js

# Show only available files
$ ks file -qf available
  A other-available.file
` + "```" + `
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		ctx.MustHaveEnvironment(currentEnvironment)
		shouldFetchMessages()

		files := ctx.ListAllFiles(currentEnvironment)
		exitIfErr(ctx.Err())

		display.Files(
			files,
			fileDisplayFilter,
			quietOutput,
		)
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

	filesCmd.Flags().StringVarP(&fileDisplayFilter, "filter", "f", display.FileFilterAll, "Files to display")
}
