package emailer

import (
	"bytes"
	"text/template"
)

var templates map[string]*template.Template

func init() {
	templates["invite/html"] = template.Must(template.New("invite/html").Parse(`
<html>
<head>
</head>

<body>
<p>
	Hello!
</p>
<p>
{{.Inviter}} is inviting you to join a Keystone project!
</p>

<p>
To join the project <pre>{{.ProjectName}}</pre>, {{.Inviter}} needs your 
Keystone username. To get it :
</p>

<ol>
	<li>
		create, or login into your account: <pre>ks login</pre>;
	</li>
	<li>
		display your username: <pre>ks whoami</pre>.
	</li>
</ol>

<p>
The way you transmit your Keystone username to {{.Inviter}} is up to you.
</p>

<p>
Have a nice day!
</p>

<p>
The Keystone team
</p>

</body>
`))

	templates["invite/text"] = template.Must(template.New("invite/text").Parse(`
Hello!

{{.Inviter}} is inviting you to join a Keystone project!

To join the project '{{.ProjectName}}', {{.Inviter}} needs your Keystone username.
To get it:
	1. create, or login into your account: 'ks login';
    2. display your username: 'ks whoami'.

The way you transmit your Keystone username to {{.Inviter}} is up to you.

Have a nice day!

The Keystone team
`))
}

type inviteData struct {
	Inviter     string
	ProjectName string
}

func renderInviteTemplate(
	inviter string,
	projectName string,
) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["invite/text"].Execute(
		htmlBuffer,
		inviteData{Inviter: inviter, ProjectName: projectName},
	); err != nil {
		return "", "", err
	}

	if err = templates["invite/text"].Execute(
		textBuffer,
		inviteData{Inviter: inviter, ProjectName: projectName},
	); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), textBuffer.String(), nil
}

func InviteMail(inviter string, projectName string) (email *Email, err error) {
	html, text, err := renderInviteTemplate(inviter, projectName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     inviter,
		Subject:  "Your are invited to join a Keystone project",
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}
