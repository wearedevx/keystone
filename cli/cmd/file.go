package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/ui"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "file",
	Short: "Manages secret files",
	Long: `Manages secret files.

List tracked secret files:
` + "```" + `
$ ks file
  Files tracked as secret files:
    config/wp-config.php
    config/front.config.js
` + "```" + `
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		var err *errors.Error

		ctx.MustHaveEnvironment(currentEnvironment)

		files := ctx.ListFiles()
		filesFromCache := ctx.ListFilesFromCache()
		filesFromCache = filterFilesFromCache(filesFromCache, files)

		files = append(files, filesFromCache...)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		if len(files) == 0 {
			if !quietOutput {
				ui.Print(`No files are currently tracked as secret files.

To add files to secret files:
  $ ks file add <path-to-secret-file>
`)
			}
			return
		}

		if quietOutput {
			for _, file := range files {
				ui.Print(file.Path)
			}
			return
		}

		filePaths := make([]string, len(files))
		for idx, file := range files {
			filePaths[idx] = file.Path
			if file.Strict {
				filePaths[idx] = filePaths[idx] + " *"
			}
		}

		ui.Print(ui.RenderTemplate("files list", `Files tracked as secret files:

{{ range . }}
{{- . | indent 4 }}
{{ end }}
* Required files; ° Unused files
`, filePaths))
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

func filterFilesFromCache(filesFromCache []keystonefile.FileKey, files []keystonefile.FileKey) []keystonefile.FileKey {
	filesFromCacheToDisplay := make([]keystonefile.FileKey, 0)

	for _, fileFromCache := range filesFromCache {
		found := false
		for _, file := range files {
			if file.Path == fileFromCache.Path {
				found = true
			}
		}
		if !found {
			fileFromCache.Path = fileFromCache.Path + " °"
			filesFromCacheToDisplay = append(filesFromCacheToDisplay, fileFromCache)
		}

	}
	return filesFromCacheToDisplay
}
