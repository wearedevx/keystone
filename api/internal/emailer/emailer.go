package emailer

import (
	"fmt"

	"github.com/keighl/mandrill"
)

var mandrillKey string

func init() {
	if mandrillKey == "" {
		mandrillKey = "SANDBOX_SUCCESS"
	}
}

func send(email *Email) (err error) {
	client := mandrill.ClientWithKey(mandrillKey)

	responses, err := client.MessagesSend(email.toMandrill())
	fmt.Println(fmt.Printf("### EMAIL SENT TO %s ###", email.To))
	fmt.Println(email.TextBody)

	if err != nil {
		return fmt.Errorf("Error sending mail: %s | %w", err.Error(), EmailErrorClientError)
	}

	for _, response := range responses {
		if response.Status == "rejected" {
			fmt.Printf("Email to %s was rejected because: %s\n", response.Email, response.RejectionReason)
		}
		if response.Status == "invalid" {
			fmt.Printf("Email to %s was deemed invalid by the service", response.Email)
		}
	}

	return nil
}
