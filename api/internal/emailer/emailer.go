package emailer

import (
	"fmt"
	"html/template"
	"os"

	"github.com/keighl/mandrill"
)

var (
	mandrillKey  string
	baseTemplate *template.Template
)

func init() {
	mandrillKey = os.Getenv("MANDRILL_API_KEY")

	if mandrillKey == "" {
		mandrillKey = "SANDBOX_SUCCESS"
	}

	baseTemplate = template.Must(template.New("base").Parse(BASE_HTML))
}

func send(email *Email) (err error) {
	client := mandrill.ClientWithKey(mandrillKey)

	responses, err := client.MessagesSend(email.toMandrill())
	if err != nil {
		return fmt.Errorf(
			"error sending mail: %s | %w",
			err.Error(),
			ErrorEmailClientError,
		)
	}

	for _, response := range responses {
		if response.Status == "rejected" {
			fmt.Printf(
				"Email to %s was rejected because: %s\n",
				response.Email,
				response.RejectionReason,
			)
		}
		if response.Status == "invalid" {
			fmt.Printf(
				"Email to %s was deemed invalid by the service",
				response.Email,
			)
		}
	}

	return nil
}
