package emailer

import "github.com/keighl/mandrill"

type Email struct {
	FromEmail string
	FromName  string
	To        []string
	Subject   string
	HtmlBody  string
	TextBody  string
}

func (e *Email) toMandrill() *mandrill.Message {
	m := &mandrill.Message{}

	for _, recipient := range e.To {
		m.AddRecipient(recipient, "", "to")
	}

	m.FromEmail = e.FromEmail
	m.FromName = e.FromName
	m.Subject = e.Subject
	m.HTML = e.HtmlBody
	m.Text = e.TextBody

	return m
}

func filterOutEmptyMails(recipients []string) []string {
	r := make([]string, 0)

	for _, email := range recipients {
		if len(email) > 0 {
			r = append(r, email)
		}
	}

	return r
}

func (e *Email) Send(recipients []string) (err error) {
	e.To = filterOutEmptyMails(recipients)

	if len(e.To) == 0 {
		return ErrorEmailNoRecipient
	}

	return send(e)
}
