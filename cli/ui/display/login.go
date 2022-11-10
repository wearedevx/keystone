package display

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

// LoginLink function displays the link the user must follow to
// start the oauth process on a third party service
func LoginLink(name, url string) {
	ui.Print(
		ui.RenderTemplate(
			"login visit",
			`Visit the URL below to login with your {{ .Service }} account:

{{ .URL | indent 8 }}

Waiting for you to login with your {{ .Service }} Account...`,
			struct {
				Service string
				URL     string
			}{
				Service: name,
				URL:     url,
			},
		),
	)
}

// AlreadyLoggedIn function Message when user is logged in
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

// WelcomeBack message
func WelcomeBack(account models.User) {
	ui.Print(ui.RenderTemplate("login ok", `
{{ OK }} {{ . | bright_green }}
`, fmt.Sprintf("Welcome back, %s", account.Username)))
}

// LoginSucces message
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

// Logout message
func Logout() {
	ui.Print("User logged out")
}
