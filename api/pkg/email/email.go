package email

import (
	"errors"
	"fmt"
	"os"

	"github.com/mattbaird/gochimp"
)

type EmailService struct {
	Client *gochimp.MandrillAPI
	Error  error
}

func (c *EmailService) InitClient() {
	apiKey := os.Getenv("MANDRILL_KEY")
	mandrillApi, err := gochimp.NewMandrill(apiKey)

	if err != nil {
		fmt.Println("Error instantiating client")
		c.Error = errors.New("Error instantiating client")
	}

	c.Client = mandrillApi
}

func (c *EmailService) SendInviteEmail(email string) {
	templateName := "invite email"
	contentVar := gochimp.Var{"main", InviteEmail}
	content := []gochimp.Var{contentVar}

	_, err := c.Client.TemplateAdd(templateName, fmt.Sprintf("%s", contentVar.Content), true)

	if err != nil {
		fmt.Printf("Error adding template: %v", err)
		return
	}
	defer c.Client.TemplateDelete(templateName)
	renderedTemplate, err := c.Client.TemplateRender(templateName, content, nil)

	if err != nil {
		fmt.Printf("Error rendering template: %v", err)
		return
	}

	recipients := []gochimp.Recipient{
		gochimp.Recipient{Email: email},
	}

	message := gochimp.Message{
		Html:      renderedTemplate,
		Subject:   "You received an Invitation to join Keystone !",
		FromEmail: "noreply@keystone.sh",
		FromName:  "Keystone",
		To:        recipients,
	}

	_, err = c.Client.MessageSend(message, false)

	if err != nil {
		fmt.Println("Error sending message")
	}
}

const InviteEmail = `
You have been invited to keystone !!
`
