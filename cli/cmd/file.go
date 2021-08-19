package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

const (
	FileFilterAll                  = "all"
	FileFilterAvailableOnly string = "available"
	FileFilterModifiedOnly  string = "modified"
)

var fileDisplayFilter string

type fileLine struct {
	required  bool
	available bool
	modified  bool
	path      string
}

func (f fileLine) requiredRune() (r rune) {
	r = ' '

	if f.required {
		r = '*'
	}

	return r
}

func (f fileLine) availableModifiedRune() (r rune) {
	r = ' '

	if f.available {
		r = 'A'
	}
	if f.modified {
		r = 'M'
	}

	return r
}

func (f fileLine) String(long bool) string {
	sb := strings.Builder{}

	if long {
		sb.WriteRune(f.requiredRune())
		sb.WriteRune(f.availableModifiedRune())
		sb.WriteRune(' ')
	}

	sb.WriteString(f.path)

	return sb.String()
}

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
Available files exist in the cache, but not in the keystone.yml file.
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
		var err *errors.Error

		ctx.MustHaveEnvironment(currentEnvironment)

		files := ctx.ListFiles()
		filesFromCache := ctx.ListFilesFromCache()
		linesFromCache := fromCacheToLines(ctx, filesFromCache, files)

		lines := make([]fileLine, len(files))
		for i, file := range files {
			lines[i] = fileToLine(file)
		}
		lines = append(lines, linesFromCache...)
		setModifiedFlags(ctx, &lines, currentEnvironment)

		lines = filterLines(lines, fileDisplayFilter)

		if err = ctx.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		if len(lines) == 0 {
			if !quietOutput {
				ui.Print(`No files are currently tracked as secret files.

To add files to secret files:
  $ ks file add <path-to-secret-file>
`)
			}
			os.Exit(0)
		}

		if quietOutput {
			for _, line := range lines {
				ui.Print(line.String(fileDisplayFilter == FileFilterAll))
			}
			return
		}

		ui.Print(ui.RenderTemplate("files list", `Files tracked as secret files:

{{ range . }}{{- (.String true) | indent 4 }}
{{ end }}
* Required; A Available; M Modified
`, lines))
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

	filesCmd.Flags().StringVarP(&fileDisplayFilter, "filter", "f", FileFilterAll, "Files to display")
}

func fileToLine(file keystonefile.FileKey) (line fileLine) {
	line.path = file.Path
	line.required = file.Strict

	return line
}

func fromCacheToLines(ctx *core.Context, filesFromCache []keystonefile.FileKey, files []keystonefile.FileKey) []fileLine {
	filesFromCacheToDisplay := make([]fileLine, 0)

	for _, fileFromCache := range filesFromCache {
		used := false

		for _, file := range files {
			fileAbs := filepath.Clean(filepath.Join(ctx.Wd, file.Path))
			cacheFileAbs := filepath.Clean(filepath.Join(ctx.Wd, fileFromCache.Path))

			if fileAbs == cacheFileAbs {
				used = true
				break
			}
		}
		if !used {
			line := fileToLine(fileFromCache)
			line.available = true

			filesFromCacheToDisplay = append(filesFromCacheToDisplay, line)
		}

	}
	return filesFromCacheToDisplay
}

func setModifiedFlags(ctx *core.Context, lines *[]fileLine, envname string) {
	for index, f := range *lines {
		(*lines)[index].modified = ctx.IsFileModified(f.path, envname)
	}
}

func filterLines(lines []fileLine, filter string) (filtered []fileLine) {
	filtered = make([]fileLine, 0, len(lines))

	if filter != FileFilterModifiedOnly && filter != FileFilterAvailableOnly {
		filter = FileFilterAll
	}

	if filter == FileFilterAll {
		return lines
	}

	switch filter {
	case FileFilterAvailableOnly:
		for _, line := range lines {
			if line.available {
				filtered = append(filtered, line)
			}
		}

	case FileFilterModifiedOnly:
		for _, line := range lines {
			if line.modified {
				filtered = append(filtered, line)
			}
		}

	default:
		panic(fmt.Errorf("Unknown filter: %s", filter))
	}

	return filtered
}
