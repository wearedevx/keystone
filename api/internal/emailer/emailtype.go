package emailer

import "github.com/keighl/mandrill"

type Email struct {
	From     string
	To       []string
	Subject  string
	HtmlBody string
	TextBody string
}

func (e *Email) toMandrill() *mandrill.Message {
	m := &mandrill.Message{}

	for _, recipient := range e.To {
		m.AddRecipient(recipient, "", "to")
	}

	m.FromEmail = e.From
	m.Subject = e.Subject
	m.HTML = e.HtmlBody
	m.Text = e.TextBody

	return m
}

func (e *Email) Send(recipients []string) (err error) {
	e.To = recipients

	return send(e)
}
