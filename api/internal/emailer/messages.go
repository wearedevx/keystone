package emailer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/k3a/html2text"
	"github.com/vanng822/go-premailer/premailer"

	"github.com/wearedevx/keystone/api/pkg/models"
)

var templates = make(map[string]*template.Template)

const KEYSTONE_MAIL = "no-reply@keystone.sh"

func init() {
	templates["invite/html"] = template.Must(template.New("invite/html").Parse(`
<p>
	Hello!
</p>
<p>
{{.Inviter}} is inviting you to join a Keystone project!
</p>

<p>
To join the project <code>{{.ProjectName}}</code>, {{.Inviter}} needs your 
Keystone username. To get it :
</p>

<ol>
	<li>
		create, or login into your account: <code>ks login</code>;
	</li>
	<li>
		display your username: <code>ks whoami</code>.
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
`))

	templates["added/html"] = template.Must(template.New("added/html").Parse(`
<p>
	Hello!
</p>
<p>
{{.Inviter}} has added you to a Keystone project!
</p>

<p>
You now have access to <code>{{.ProjectName}}</code>.
</p>

<ol>
	<li>
		go in your project directory
	</li>
	<li>
		login into your account: <code>ks login</code>;
	</li>
	<li>
		use secret: <code>ks source<code>
	</li>
</ol>

<p>
Have a nice day!
</p>

<p>
The Keystone team
</p>
`))

	templates["new_device_admin/html"] = template.Must(
		template.New("new_device_admin/html").Parse(`
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
`),
	)

	templates["new_device/html"] = template.Must(
		template.New("new_device/html").Parse(`
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
<p>To revoke a device:</p>
<p style=" background-color:#f6f8fa ; border-radius: 6px; font-size: 95%; line-height: 1.45; overflow: auto; padding: 16px; "> $ ks device revoke {{.DeviceName}}</p>
<p>
Have a nice day!
</p>

<p>
The Keystone team
</p>
`),
	)

	templates["message_will_expire/html"] = template.Must(
		template.New("mesage_will_expire/html").Parse(`
		<p>
			Hello!
		</p>

		<p>
			Some messages you haven't read yet will expire in {{.NbDays}} days.
		</p>

		<p>
			Related projects:
		</p>

		{{ range $groupedProject := .GroupedProjects}}
			<p>
				- Project: {{$groupedProject.Project.Name}}
			</p>

			<p>
				Environments:
					<ul>
						{{range $environment := $groupedProject.Environments}}
							<li>- {{$environment.Name}}</li>
						{{ end }}
					</ul>
			</p>
		{{ end }}

		<p>
			Retrieve them before they expire with:
			<ol>
				<li>$ cd <my-project></li>
				<li>$ ks source</li>
			</ol>
		</p>
`),
	)
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

type GroupedMessageProject struct {
	Project      models.Project
	Environments map[string]models.Environment
}

type GroupedMessagesUser struct {
	Recipient models.User
	Projects  map[uint]GroupedMessageProject
}

type newGroupedMessageProjectData struct {
	NbDays          int
	GroupedProjects map[uint]GroupedMessageProject
}

func renderTemplate(
	templateName string,
	data interface{},
) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)

	// Render the content template
	if err = templates[templateName].Execute(
		htmlBuffer,
		data,
	); err != nil {
		return "", "", err
	}

	// Render Layout
	htmlString := htmlBuffer.String()
	finalBuffer := bytes.NewBufferString("")
	if err = baseTemplate.Execute(finalBuffer, map[string]template.HTML{
		"Preheader": "",
		"Content":   template.HTML(htmlString),
	}); err != nil {
		return "", "", err
	}

	// Transform to text
	finalHtml := finalBuffer.String()

	inlined, err := inlineHtml(finalHtml)
	if err != nil {
		return "", "", err
	}

	finalText := htmlToString(inlined)

	return inlined, finalText, nil
}

func htmlToString(html string) (text string) {
	text = html2text.HTML2Text(html)

	return text
}

func inlineHtml(html string) (out string, err error) {
	prem, err := premailer.NewPremailerFromString(html, premailer.NewOptions())
	if err != nil {
		return "", err
	}

	out, err = prem.Transform()
	if err != nil {
		return "", err
	}

	return out, nil
}

func renderInviteTemplate(
	inviter string,
	projectName string,
) (html string, text string, err error) {
	html, text, err = renderTemplate("invite/html", inviteData{
		Inviter:     inviter,
		ProjectName: projectName,
	})
	if err != nil {
		return "", "", err
	}

	return html, text, nil
}

func renderAddedTemplate(
	inviter string,
	projectName string,
) (html string, text string, err error) {
	html, text, err = renderTemplate("added/html", inviteData{
		Inviter:     inviter,
		ProjectName: projectName,
	})
	if err != nil {
		return "", "", err
	}

	return html, text, nil
}

func renderNewDeviceAdmin(
	userID string,
	projects []string,
	deviceName string,
) (html string, text string, err error) {
	html, text, err = renderTemplate(
		"new_device_admin/html",
		newDeviceAdminData{
			UserID:     userID,
			Projects:   strings.Join(projects, ", "),
			DeviceName: deviceName,
		},
	)
	if err != nil {
		return "", "", err
	}

	return html, text, nil
}

func renderNewDevice(
	deviceName string,
	userID string,
) (html string, text string, err error) {
	html, text, err = renderTemplate(
		"new_device/html",
		newDeviceData{
			DeviceName: deviceName,
			UserID:     userID,
		},
	)
	if err != nil {
		return "", "", err
	}

	return html, text, nil
}

func renderMessageWillExpire(
	nbDays int,
	groupedProjects map[uint]GroupedMessageProject,
) (html string, text string, err error) {
	html, text, err = renderTemplate(
		"message_will_expire/html",
		newGroupedMessageProjectData{
			NbDays:          nbDays,
			GroupedProjects: groupedProjects,
		},
	)
	if err != nil {
		return "", "", err
	}

	return html, text, nil
}

func InviteMail(
	inviter models.User,
	projectName string,
) (email *Email, err error) {
	html, text, err := renderInviteTemplate(inviter.Email, projectName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		FromEmail: KEYSTONE_MAIL,
		FromName:  inviter.GetName(),
		Subject:   "You are invited to join a Keystone project",
		HtmlBody:  html,
		TextBody:  text,
	}

	return email, nil
}

func AddedMail(
	inviter models.User,
	projectName string,
) (email *Email, err error) {
	html, text, err := renderInviteTemplate(inviter.Email, projectName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		FromEmail: KEYSTONE_MAIL,
		FromName:  inviter.GetName(),
		Subject:   "You are added to a Keystone project",
		HtmlBody:  html,
		TextBody:  text,
	}

	return email, nil
}

func NewDeviceAdminMail(
	userID string,
	projects []string,
	deviceName string,
) (email *Email, err error) {
	html, text, err := renderNewDeviceAdmin(userID, projects, deviceName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		FromEmail: KEYSTONE_MAIL,
		FromName:  "Keystone",
		Subject:   fmt.Sprintf("%s has registered a new device", userID),
		HtmlBody:  html,
		TextBody:  text,
	}

	return email, nil
}

func NewDeviceMail(deviceName string, userID string) (email *Email, err error) {
	html, text, err := renderNewDevice(deviceName, userID)
	if err != nil {
		return nil, err
	}

	email = &Email{
		FromEmail: KEYSTONE_MAIL,
		FromName:  "Keystone",
		Subject:   "A new device has been registered",
		HtmlBody:  html,
		TextBody:  text,
	}

	return email, nil
}

func MessageWillExpireMail(
	nbDays int,
	groupedProjects map[uint]GroupedMessageProject,
) (email *Email, err error) {
	html, text, err := renderMessageWillExpire(nbDays, groupedProjects)
	if err != nil {
		return nil, err
	}

	email = &Email{
		FromEmail: KEYSTONE_MAIL,
		FromName:  "Keystone",
		Subject:   "Some message will expire",
		HtmlBody:  html,
		TextBody:  text,
	}

	return email, nil
}
