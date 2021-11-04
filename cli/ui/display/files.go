package display

import (
	"fmt"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

const (
	FileFilterAll                  = "all"
	FileFilterAvailableOnly string = "available"
	FileFilterModifiedOnly  string = "modified"
)

// ————  ———— //

// Display the files, with filters
// and with a shorter version if quiet is on
func Files(
	files []core.FileDescriptor,
	filter string,
	quiet bool,
) {
	if len(files) == 0 {
		if !quiet {
			NoFilesTracked()
		}

		return
	}

	if quiet {
		for _, fd := range files {
			ui.Print(line(fd).String(filter == FileFilterAll))
		}
	}

	lines := filterLines(files, filter)

	fileList(lines)
}

// Message for when there are no files to display
func NoFilesTracked() {
	ui.Print(`No files are currently tracked as secret files.

To add files to secret files:
  $ ks file add <path-to-secret-file>
`)
}

// FileContentsForEnvironments function displays the contents of a file
// for each environment
func FileContentsForEnvironments(
	fileName string,
	environments []models.Environment,
	getContent func(string, string) ([]byte, error),
) {
	ui.Print(`The file already exist but is not used.`)
	for _, env := range environments {
		content, err := getContent(fileName, env.Name)

		ui.Print("\n")
		ui.Print("------------------" + env.Name + "------------------")
		ui.Print("\n")
		if err != nil {
			ui.Print("File not found for this environment")
		}

		ui.Print(string(content))
	}
}

// FileAddSuccess function Message when adding a file is successfull
func FileAddSuccess(filePath string, numberOfEnvironments int) {
	ui.Print(ui.RenderTemplate("file add success", `
{{ OK }} {{ .Title | green }}
The file has been added to {{ .NumberEnvironments }} environment(s).
It has also been gitignored.`, map[string]string{
		"Title":              fmt.Sprintf("Added '%s'", filePath),
		"NumberEnvironments": fmt.Sprintf("%d", numberOfEnvironments),
	}))
}

// FileAskForFileContentForEnvironment function Ask file content
func FileAskForFileContentForEnvironment(filePath, environmentName string) {
	ui.Print(
		fmt.Sprintf(
			"Enter content for file `%s` for the '%s' environment (Press any key to continue)",
			filePath,
			environmentName,
		),
	)
}

// FileFailUserInput function Message when input from $EDITOR failed
func FileFailUserInput(err error) {
	ui.PrintStdErr(
		fmt.Sprintf("Failed to read user input (%s)", err.Error()),
	)
}

// FileFailUserInput function Message when input from $EDITOR failed
func FileFailedGetContentFromEditor(err error) {
	ui.PrintStdErr(
		fmt.Sprintf("Failed to get content from editor (%s)", err.Error()),
	)
}

// FileIsNowOptional function Message when setting the file as optional
func FileIsNowOptional(filePath string) {
	ui.Print(ui.RenderTemplate(
		"set file optional",
		`File {{ . }} is now optional.`,
		filePath,
	))
}

// FileIsNowOptional function Message when setting the file as required
func FileIsNowRequired(filePath string) {
	ui.Print(ui.RenderTemplate(
		"set file optional",
		`File {{ . }} is now required.`,
		filePath,
	))
}

// FileNotManaged function Message when the file is not managed by Keystone
func FileNotManaged(filePath string) {
	ui.Print("File '" + filePath + "' is not managed by Keystone, ignoring")
}

// FileKept function Message informing the file is kept in cache after ks file rm
func FileKept() {
	ui.Print(
		`The file is kept in your keystone project for all the environments,
in case you need it again.
If you want to remove it from your device, use --purge`,
	)
}

// FileRemovedSuccess function Message when file removal happened successfully
func FileRemovedSuccess(filePath string) {
	ui.PrintSuccess("%s has been removed from the secret files.", filePath)
}

// FileSetSuccess function Message when file content updated successfully
func FileSetSuccess(filePath string) {
	ui.Print(ui.RenderTemplate("file set success", `
{{ OK }} {{ . | green }}
`,
		fmt.Sprintf("Modified '%s'", filePath),
	))
}

// ———— PRIVATE Utilities ———— //

// File list when no quiet
func fileList(lines []line) {
	ui.Print(ui.RenderTemplate("files list", `Files tracked as secret files:

{{ range . }}{{- (.String true) | indent 4 }}
{{ end }}
* Required; A Available; M Modified
`, lines))
}

// Add some display functions to the FileDescriptor type
// by aliasing the type
type line core.FileDescriptor

func (f line) requiredRune() (r rune) {
	r = ' '

	if f.Required {
		r = '*'
	}

	return r
}

func (f line) availableModifiedRune() (r rune) {
	r = ' '

	if f.Available {
		r = 'A'
	}
	if f.Modified {
		r = 'M'
	}

	return r
}

func (f line) String(long bool) string {
	sb := strings.Builder{}

	if long {
		sb.WriteRune(f.requiredRune())
		sb.WriteRune(f.availableModifiedRune())
		sb.WriteRune(' ')
	}

	sb.WriteString(f.Path)

	return sb.String()
}

// ————  ———— //

// Filters the files according to `filter`
// It returns the alias `line` type to enable the display functions
func filterLines(files []core.FileDescriptor, filter string) (filtered []line) {
	filtered = make([]line, 0, len(files))

	if filter != FileFilterModifiedOnly && filter != FileFilterAvailableOnly {
		filter = FileFilterAll
	}

	switch filter {
	case FileFilterAvailableOnly:
		for _, file := range files {
			if file.Available {
				filtered = append(filtered, line(file))
			}
		}

	case FileFilterModifiedOnly:
		for _, file := range files {
			if file.Modified {
				filtered = append(filtered, line(file))
			}
		}

	case FileFilterAll:
		for _, file := range files {
			filtered = append(filtered, line(file))
		}

	default:
		panic(fmt.Errorf("unknown filter: %s", filter))
	}

	return filtered
}
