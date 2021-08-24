package emailer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var templates = make(map[string]*template.Template)

const KEYSTONE_MAIL = "no-reply@keystone.sh"

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

	templates["added/text"] = template.Must(template.New("added/text").Parse(`
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
	templates["new_device_admin/html"] = template.Must(template.New("new_device_admin/html").Parse(`
<html>
<head>
</head>

<body>
<p>
	Hello!
</p>
<p>
{{.UserID}} has added a new device to its account.
</p>
<p>
You are admin in some of its project(s): {{.Projects}}
</p>
<p>
The new device name is: {{.DeviceName}}
</p>
<p>
If you think this new device is suspicious, feel free to contact {{.UserID}}.
</p>
<p>
Have a nice day!
</p>

<p>
The Keystone team
</p>

</body>
`))

	templates["new_device_admin/text"] = template.Must(template.New("new_device_admin/text").Parse(`
Hello!

{{.UserID}} has added a new device to its account.

You are admin in some of its project(s): {{.Projects}}

The new device name is: {{.DeviceName}}

If you think this new device is suspicious, feel free to contact {{.UserID}}.

Have a nice day!

The Keystone team
`))
	templates["new_device/html"] = template.Must(template.New("new_device/html").Parse(`
<html>
<head>
</head>

<body>
<p>
	Hello!
</p>
<p>
A new device have been added to your Keystone account {{.UserID}}.
</p>
<p>
The new device name is: {{.DeviceName}}
</p>
<p>
If you didn't connect with this new device, you can revoke its access using keystone app.
You should also change your access to the identity provider you chose to connect to Keystone.
</p>
<p>
Have a nice day!
</p>

<p>
The Keystone team
</p>

</body>
`))

	templates["new_device/text"] = template.Must(template.New("new_device/text").Parse(`
Hello!

A new device have been added to your Keystone account {{.UserID}}.

The new device name is: {{.DeviceName}}

If you didn't connect with this new device, you can revoke its access using keystone app.
You should also change your access to the identity provider you chose to connect to Keystone.

Have a nice day!

The Keystone team
`))
}

type inviteData struct {
	Inviter     string
	ProjectName string
}

type newDeviceAdminData struct {
	Projects   string
	UserID     string
	DeviceName string
}

type newDeviceData struct {
	DeviceName string
	UserID     string
}

func renderInviteTemplate(
	inviter string,
	projectName string,
) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["invite/html"].Execute(
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

	if err = templates["added/html"].Execute(
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

func renderNewDeviceAdmin(userID string, projects []string, deviceName string) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["new_device_admin/html"].Execute(
		htmlBuffer,
		newDeviceAdminData{DeviceName: deviceName, Projects: strings.Join(projects, ", "), UserID: userID},
	); err != nil {
		return "", "", err
	}

	if err = templates["new_device_admin/text"].Execute(
		textBuffer,
		newDeviceAdminData{DeviceName: deviceName, Projects: strings.Join(projects, ", "), UserID: userID},
	); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), textBuffer.String(), nil
}

func renderNewDevice(deviceName string, userID string) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["new_device/html"].Execute(
		htmlBuffer,
		newDeviceData{DeviceName: deviceName, UserID: userID},
	); err != nil {
		return "", "", err
	}

	if err = templates["new_device/text"].Execute(
		textBuffer,
		newDeviceData{DeviceName: deviceName},
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

func NewDeviceAdminMail(userID string, projects []string, deviceName string) (email *Email, err error) {
	html, text, err := renderNewDeviceAdmin(userID, projects, deviceName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     KEYSTONE_MAIL,
		Subject:  fmt.Sprintf("%s has registered a new device", userID),
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}

func NewDeviceMail(deviceName string, userID string) (email *Email, err error) {
	html, text, err := renderNewDevice(deviceName, userID)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     KEYSTONE_MAIL,
		Subject:  fmt.Sprintf("A new device has been registered"),
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}
