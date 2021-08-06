package emailer

import (
	"bytes"
	"text/template"
)

var templates = make(map[string]*template.Template)

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

	templates["added/html"] = template.Must(template.New("added/html").Parse(`
<html>
<head>
</head>

<body>
<p>
	Hello!
</p>
<p>
{{.Inviter}} has added you to a Keystone project!
</p>

<p>
You now have access to <pre>{{.ProjectName}}</pre>.
</p>

<ol>
	<li>
		go in your project directory
	</li>
	<li>
		login into your account: <pre>ks login</pre>;
	</li>
	<li>
		use secret: <pre>ks source<pre>
	</li>
</ol>

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

{{.Inviter}} has added you to a Keystone project!

You now have access to {{.ProjectName}}.

To get it:
  1. go in your project directory
  2. login into your account: <pre>ks login</pre>;
  3. use secret: <pre>eval "$(ks source)"<pre>

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

func renderAddedTemplate(
	inviter string,
	projectName string,
) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["added/text"].Execute(
		htmlBuffer,
		inviteData{Inviter: inviter, ProjectName: projectName},
	); err != nil {
		return "", "", err
	}

	if err = templates["added/text"].Execute(
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

func AddedMail(inviter string, projectName string) (email *Email, err error) {
	html, text, err := renderInviteTemplate(inviter, projectName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     inviter,
		Subject:  "Your are added to join a Keystone project",
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}
