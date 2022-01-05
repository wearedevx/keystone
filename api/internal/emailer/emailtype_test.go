package emailer

import (
	"reflect"
	"testing"
)

func Test_filterOutEmptyMails(t *testing.T) {
	type args struct {
		recipients []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "filters out empty emails",
			args: args{
				recipients: []string{
					"test@mail.com",
					"invalidemailaddress",
					"",
					"another.test@example.com",
				},
			},
			want: []string{
				"test@mail.com",
				"invalidemailaddress",
				"another.test@example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterOutEmptyMails(tt.args.recipients); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterOutEmptyMails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_Send(t *testing.T) {
	type fields struct {
		FromEmail string
		FromName  string
		To        []string
		Subject   string
		HtmlBody  string
		TextBody  string
	}
	type args struct {
		recipients []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr bool
		key     string
	}{
		{
			name: "sends the email",
			fields: fields{
				FromEmail: "no-reply@keystone.sh",
				FromName:  "Keystone",
				To:        []string{},
				Subject:   "test email",
				HtmlBody:  "<html><body>test email</body></html>",
				TextBody:  "test email",
			},
			args: args{
				recipients: []string{
					"user@mail.com",
					"",
					"admin@tes.com",
					"invalidemail",
				},
			},
			want: fields{
				FromEmail: "no-reply@keystone.sh",
				FromName:  "Keystone",
				To: []string{
					"user@mail.com",
					"admin@tes.com",
					"invalidemail",
				},
				Subject:  "test email",
				HtmlBody: "<html><body>test email</body></html>",
				TextBody: "test email",
			},
			wantErr: false,
			key:     "SANDBOX_SUCCESS",
		},
		{
			name: "fails at seding the email",
			fields: fields{
				FromEmail: "no-reply@keystone.sh",
				FromName:  "Keystone",
				To:        []string{},
				Subject:   "test email",
				HtmlBody:  "<html><body>test email</body></html>",
				TextBody:  "test email",
			},
			args: args{
				recipients: []string{
					"user@mail.com",
					"",
					"admin@tes.com",
					"invalidemail",
				},
			},
			want: fields{
				FromEmail: "no-reply@keystone.sh",
				FromName:  "Keystone",
				To: []string{
					"user@mail.com",
					"admin@tes.com",
					"invalidemail",
				},
				Subject:  "test email",
				HtmlBody: "<html><body>test email</body></html>",
				TextBody: "test email",
			},
			wantErr: true,
			key:     "SANDBOX_ERROR",
		},
		{
			name: "error if no recipients",
			fields: fields{
				FromEmail: "no-reply@keystone.sh",
				FromName:  "Keystone",
				To:        []string{},
				Subject:   "test email",
				HtmlBody:  "<html><body>test email</body></html>",
				TextBody:  "test email",
			},
			args: args{
				recipients: []string{},
			},
			want: fields{
				FromEmail: "no-reply@keystone.sh",
				FromName:  "Keystone",
				To:        []string{},
				Subject:   "test email",
				HtmlBody:  "<html><body>test email</body></html>",
				TextBody:  "test email",
			},
			wantErr: true,
			key:     "SANDBOX_ERROR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Email{
				FromEmail: tt.fields.FromEmail,
				FromName:  tt.fields.FromName,
				To:        tt.fields.To,
				Subject:   tt.fields.Subject,
				HtmlBody:  tt.fields.HtmlBody,
				TextBody:  tt.fields.TextBody,
			}
			mandrillKey = tt.key

			if err := e.Send(tt.args.recipients); (err != nil) != tt.wantErr {
				t.Errorf("Email.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if e.FromEmail != tt.want.FromEmail {
				t.Errorf("Email.Send() got = %v, want %v", *e, tt.want)
				return
			}

			if e.FromName != tt.fields.FromName {
				t.Errorf("Email.Send() got = %v, want %v", *e, tt.want)
			}

			if !reflect.DeepEqual(e.To, tt.want.To) {
				t.Errorf("Email.Send() got.To = %v, want.To %v", e.To, tt.want.To)
			}

			if e.Subject != tt.fields.Subject {
				t.Errorf("Email.Send() got = %v, want %v", *e, tt.want)
			}

			if e.HtmlBody != tt.fields.HtmlBody {
				t.Errorf("Email.Send() got = %v, want %v", *e, tt.want)
			}

			if e.TextBody != tt.fields.TextBody {
				t.Errorf("Email.Send() got = %v, want %v", *e, tt.want)
			}

		})
	}
}
