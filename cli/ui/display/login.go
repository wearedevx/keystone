package display

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

func LoginLink(name, url string) {
	ui.Print(
		ui.RenderTemplate(
			"login visit",
			`Visit the URL below to login with your {{ .Service }} account:

{{ .Url | indent 8 }}

Waiting for you to login with your {{ .Service }} Account...`,
			struct {
				Service string
				Url     string
			}{
				Service: name,
				Url:     url,
			},
		),
	)
}

func AlreadyLoggedIn(account models.User) {
	username := account.Username
	if username == "" {
		username = account.Email
	}

	ui.Print(
		ui.RenderTemplate(
			"already logged in",
			`You are already logged in as {{ . }}`,
			username,
		),
	)
}

func WelcomeBack(account models.User) {
	ui.Print(ui.RenderTemplate("login ok", `
{{ OK }} {{ . | bright_green }}
`, fmt.Sprintf("Welcome back, %s", account.Username)))
}

func LoginSucces() {
	ui.Print(ui.RenderTemplate("login success", `
{{ OK }} {{ . | bright_green }}

Thank you for using Keystone!

To start managing secrets for a project:
  $ cd <path-to-your-project>
  $ ks init <your-project-name>

To invite collaborators:
  $ ks invite collaborator@email
`, "Thank you for using Keystone!"))
}

func Logout() {
	ui.Print("User logged out")
}
