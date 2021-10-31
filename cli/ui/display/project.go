package display

import (
	"strings"

	"github.com/wearedevx/keystone/cli/ui"
)

func DeletionSuccess(projectName string) {
	ui.Print(ui.RenderTemplate(
		"deletion ok",
		`{{ OK }} The project {{ . }} has successfully been destroyed.
Secrets and files are no longer accessible.
You may need to remove entries from your .gitignore file`,
		projectName,
	))
}

func ProjectInitSuccess() {
	ui.Print(ui.RenderTemplate("Init Success", `
{{ .Message | box | bright_green | indent 2 }}

{{ .Text | bright_black | indent 2 }}`, map[string]string{
		"Message": "All done!",
		"Text": `You can start adding environment variable with:
  $ ks secret add VARIABLE value

Load them with:
  $ eval $(ks source)

If you need help with anything:
  $ ks help [command]

`,
	}))
}

func InviteSuccess(usersUIDs []string, email string) {
	if len(usersUIDs) > 0 {
		ui.Print(ui.RenderTemplate("file add success", `
{{ OK }} {{ .Title | green }}

The email is associated with a Keystone account. They are registered as: {{ .Usernames | bright_green }}.

To add them to the project use "member add" command:
  $ ks member add <username>
`, map[string]string{
			"Title":     "User already on Keystone",
			"Usernames": strings.Join(usersUIDs, ", "),
		}))
	} else {
		ui.PrintSuccess("A email has been sent to %s, they will get back to you when their Keystone account will be created", email)
	}
}
