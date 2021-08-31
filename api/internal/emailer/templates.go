package emailer

import "text/template"

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

	templates["message_will_expire/html"] = template.Must(template.New("mesage_will_expire/html").Parse(`
<html>
	<head>
	</head>

	<body>
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
		
		<p>
			Have a nice day!
		</p>
		
		<p>
			The Keystone team
		</p>
	</body>
</html>
`))

	templates["message_will_expire/text"] = template.Must(template.New("new_device/text").Parse(`
Hello!

Some messages you haven't read yet will expire in {{.NbDays}} days.

Related projects:

{{ range $groupedProject := .GroupedProjects}}
- Project: {{$groupedProject.Project.Name}}

Environments:
	{{range $environment := $groupedProject.Environments}}
		- {{$environment.Name}}
	{{ end }}
{{ end }}

Retrieve them before they expire with:
$ cd <my-project>
$ ks source

Have a nice day!

The Keystone team
`))
}
