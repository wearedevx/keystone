package emailer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/wearedevx/keystone/api/pkg/models"
)

func TestInviteMail(t *testing.T) {
	type args struct {
		inviter     models.User
		projectName string
	}
	tests := []struct {
		name      string
		args      args
		wantEmail *Email
		wantErr   bool
	}{
		{
			name: "ok",
			args: args{
				inviter: models.User{
					Email:    "ok@example.com",
					Username: "ok",
				},
				projectName: "okProject",
			},
			wantErr: false,
			wantEmail: &Email{
				FromEmail: KEYSTONE_MAIL,
				FromName:  "ok",
				Subject:   "You are invited to join a Keystone project",
				HtmlBody:  ``,
				TextBody: `Hello!

ok@example.com is inviting you to join a Keystone project!

To join the project okProject, ok@example.com needs your Keystone username. To get it :

create, or login into your account: ks login; display your username: ks whoami. The way you transmit your Keystone username to ok@example.com is up to you.

Have a nice day!

The Keystone team

DevX, 2 av Président Pierre Angot, 64000 Pau, Nouvelle-Aquitaine, France`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := InviteMail(tt.args.inviter, tt.args.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InviteMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEmail.FromName != tt.wantEmail.FromName {
				t.Errorf(
					"InviteMail() FromName = %v, want %v",
					gotEmail.FromName,
					tt.wantEmail.FromName,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.Subject != tt.wantEmail.Subject {
				t.Errorf(
					"InviteMail() Subject = %v, want %v",
					gotEmail.Subject,
					tt.wantEmail.Subject,
				)
			}

			got := trimNonBreaking(gotEmail.TextBody)
			want := emailText(tt.wantEmail.TextBody)

			if !charCmp(got, want) {
				t.Errorf("InviteMail() TextBody = %v, want %v", got, want)
			}
		})
	}
}

func TestAddedMail(t *testing.T) {
	type args struct {
		inviter     models.User
		projectName string
	}
	tests := []struct {
		name      string
		args      args
		wantEmail *Email
		wantErr   bool
	}{
		{
			name: "member-added-ok",
			args: args{
				inviter: models.User{
					Username: "ok@service",
					Email:    "ok@mail.co",
				},
				projectName: "okProject",
			},
			wantEmail: &Email{
				FromEmail: KEYSTONE_MAIL,
				FromName:  "ok@service",
				To:        []string{},
				Subject:   "You are added to a Keystone project",
				HtmlBody:  "",
				TextBody: `Hello! 

ok@mail.co has added you to a Keystone project! 

You now have access to okProject. 

go in your project directory login into your account: ks login; use secret: ks source Have a nice day! 

The Keystone team 

DevX, 2 av Président Pierre Angot, 64000 Pau, Nouvelle-Aquitaine, France`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := AddedMail(tt.args.inviter, tt.args.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddedMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotEmail.FromName != tt.wantEmail.FromName {
				t.Errorf(
					"InviteMail() FromName = %v, want %v",
					gotEmail.FromName,
					tt.wantEmail.FromName,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.Subject != tt.wantEmail.Subject {
				t.Errorf(
					"InviteMail() Subject = %v, want %v",
					gotEmail.Subject,
					tt.wantEmail.Subject,
				)
			}

			got := trimNonBreaking(gotEmail.TextBody)
			want := emailText(tt.wantEmail.TextBody)

			if !charCmp(got, want) {
				t.Errorf("InviteMail() TextBody = %v, want %v", got, want)
			}
		})
	}
}

func TestNewDeviceAdminMail(t *testing.T) {
	type args struct {
		userID     string
		projects   []string
		deviceName string
	}
	tests := []struct {
		name      string
		args      args
		wantEmail *Email
		wantErr   bool
	}{
		{
			name: "new-device-to-admin-ok",
			args: args{
				userID: "memberx@github",
				projects: []string{
					"this-one",
					"this-other-one",
				},
				deviceName: "The-Computer",
			},
			wantEmail: &Email{
				FromEmail: KEYSTONE_MAIL,
				FromName:  "Keystone",
				To:        []string{},
				Subject:   "memberx@github has registered a new device",
				HtmlBody:  "",
				TextBody: `Hello!

memberx@github has added a new device to its account.

You are admin in some of its project(s): this-one, this-other-one

The new device name is: The-Computer

If you think this new device is suspicious, feel free to contact memberx@github.

Have a nice day!

The Keystone team

DevX, 2 av Président Pierre Angot, 64000 Pau, Nouvelle-Aquitaine, France`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := NewDeviceAdminMail(
				tt.args.userID,
				tt.args.projects,
				tt.args.deviceName,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewDeviceAdminMail() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if gotEmail.FromName != tt.wantEmail.FromName {
				t.Errorf(
					"InviteMail() FromName = %v, want %v",
					gotEmail.FromName,
					tt.wantEmail.FromName,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.Subject != tt.wantEmail.Subject {
				t.Errorf(
					"InviteMail() Subject = %v, want %v",
					gotEmail.Subject,
					tt.wantEmail.Subject,
				)
			}

			got := trimNonBreaking(gotEmail.TextBody)
			want := emailText(tt.wantEmail.TextBody)

			if !charCmp(got, want) {
				t.Errorf("InviteMail() TextBody = %v, want %v", got, want)
			}
		})
	}
}

func TestNewDeviceMail(t *testing.T) {
	type args struct {
		deviceName string
		userID     string
	}
	tests := []struct {
		name      string
		args      args
		wantEmail *Email
		wantErr   bool
	}{
		{
			name: "new-device-ok",
			args: args{
				deviceName: "The-Computer",
				userID:     "memberx@github",
			},
			wantEmail: &Email{
				FromEmail: KEYSTONE_MAIL,
				FromName:  "Keystone",
				To:        []string{},
				Subject:   "A new device has been registered",
				HtmlBody:  "",
				TextBody: `Hello!

A new device have been added to your Keystone account memberx@github.

The new device name is: The-Computer

If you didn't connect with this new device, you can revoke its access using keystone app. You should also change your access to the identity provider you chose to connect to Keystone.

To revoke a device:

$ ks device revoke The-Computer

Have a nice day!

The Keystone team

DevX, 2 av Président Pierre Angot, 64000 Pau, Nouvelle-Aquitaine, France`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := NewDeviceMail(tt.args.deviceName, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewDeviceMail() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if gotEmail.FromName != tt.wantEmail.FromName {
				t.Errorf(
					"InviteMail() FromName = %v, want %v",
					gotEmail.FromName,
					tt.wantEmail.FromName,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.Subject != tt.wantEmail.Subject {
				t.Errorf(
					"InviteMail() Subject = %v, want %v",
					gotEmail.Subject,
					tt.wantEmail.Subject,
				)
			}

			got := emailText(trimNonBreaking(gotEmail.TextBody))
			want := emailText(tt.wantEmail.TextBody)

			if !charCmp(got, want) {
				t.Errorf("InviteMail() TextBody = %v, want %v", got, want)
			}
		})
	}
}

func TestMessageWillExpireMail(t *testing.T) {
	type args struct {
		nbDays          int
		groupedProjects map[uint]GroupedMessageProject
	}
	tests := []struct {
		name      string
		args      args
		wantEmail *Email
		wantErr   bool
	}{
		{
			name: "message-will-expire-ok",
			args: args{
				nbDays: 3,
				groupedProjects: map[uint]GroupedMessageProject{
					1: {
						Project: models.Project{
							Name: "that-one",
						},
						Environments: map[string]models.Environment{
							"dev": {
								Name: "dev",
							},
							"staging": {
								Name: "staging",
							},
							"prod": {
								Name: "prod",
							},
						},
					},
					2: {
						Project: models.Project{
							Name: "this-one",
						},
						Environments: map[string]models.Environment{
							"dev": {
								Name: "dev",
							},
							"staging": {
								Name: "staging",
							},
							"prod": {
								Name: "prod",
							},
						},
					},
					3: {
						Project: models.Project{
							Name: "that-one-too",
						},
						Environments: map[string]models.Environment{
							"dev": {
								Name: "dev",
							},
							"staging": {
								Name: "staging",
							},
							"prod": {
								Name: "prod",
							},
						},
					},
				},
			},
			wantEmail: &Email{
				FromEmail: KEYSTONE_MAIL,
				FromName:  "Keystone",
				To:        []string{},
				Subject:   "Some message will expire",
				HtmlBody:  "",
				TextBody: `Hello!

Some messages you haven't read yet will expire in 3 days.

Related projects:

- Project: that-one

Environments:

- dev - prod - staging


- Project: this-one

Environments:

- dev - prod - staging


- Project: that-one-too

Environments:

- dev - prod - staging


Retrieve them before they expire with:

$ cd $ ks source

DevX, 2 av Président Pierre Angot, 64000 Pau, Nouvelle-Aquitaine, France`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := MessageWillExpireMail(
				tt.args.nbDays,
				tt.args.groupedProjects,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"MessageWillExpireMail() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if gotEmail.FromName != tt.wantEmail.FromName {
				t.Errorf(
					"InviteMail() FromName = %v, want %v",
					gotEmail.FromName,
					tt.wantEmail.FromName,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.FromEmail != tt.wantEmail.FromEmail {
				t.Errorf(
					"InviteMail() FromEmail = %v, want %v",
					gotEmail.FromEmail,
					tt.wantEmail.FromEmail,
				)
			}
			if gotEmail.Subject != tt.wantEmail.Subject {
				t.Errorf(
					"InviteMail() Subject = %v, want %v",
					gotEmail.Subject,
					tt.wantEmail.Subject,
				)
			}

			got := emailText(trimNonBreaking(gotEmail.TextBody))
			want := emailText(tt.wantEmail.TextBody)

			if !charCmp(got, want) {
				t.Errorf("InviteMail() TextBody = %v, want %v", got, want)
			}
		})
	}
}

func trimNonBreaking(s string) string {
	return s[2 : len(s)-4]
}

func emailText(input string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	parts := strings.Split(input, "\n")

	for i, v := range parts {
		v = strings.TrimSpace(v)
		if len(v) > 1 {
			parts[i] = " " + v + " "
		} else {
			parts[i] = ""
		}
	}

	input = strings.Join(parts, "\r\n")

	return input
}

func charCmp(got, want string) bool {
	gotb := bytes.NewBuffer([]byte(got))
	wantb := bytes.NewBuffer([]byte(want))
	gotrd := bufio.NewReader(gotb)
	wantrd := bufio.NewReader(wantb)

	eq := true
	line := 1
	col := 1
	for {
		gotRune, _, aerr := gotrd.ReadRune()
		wantedRune, _, berr := wantrd.ReadRune()

		if aerr != nil || berr != nil {
			if errors.Is(aerr, io.EOF) {
				bf := wantrd.Buffered()
				if bf != 0 {
					fmt.Println("not the same length wanted more:", bf)
					eq = false
					break
				}
			}

			if errors.Is(berr, io.EOF) {
				af := gotrd.Buffered()
				if af != 0 {
					fmt.Println("not the same length got more:", af)
					eq = false
					break
				}
			}
			break
		}
		if gotRune != '\r' && gotRune != '\n' {
			col += 1
		}

		if gotRune != wantedRune {
			fmt.Printf(
				"gotRune 0x%x is not expected 0x%x @ %d:%d\n",
				gotRune,
				wantedRune,
				line,
				col,
			)
			eq = false
			break
		}

		if gotRune == '\n' {
			line += 1
			col = 1
		}

		// fmt.Printf("%s%s", string(gotRune), string(wantedRune))
	}

	return eq
}
