package emailer

import "testing"

func Test_send(t *testing.T) {
	type args struct {
		email *Email
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		key     string
	}{
		{
			name: "it sends an email",
			args: args{
				email: &Email{
					FromEmail: "test@email.com",
					FromName:  "Test",
					To: []string{
						"recipient@mail.com",
					},
					Subject:  "Test",
					HtmlBody: "<html><head></head><body>Test</body></html>",
					TextBody: "Test",
				},
			},
			wantErr: false,
			key:     "SANDBOX_SUCCESS",
		},
		{
			name: "it fails at sending an email",
			args: args{
				email: &Email{
					FromEmail: "test@email.com",
					FromName:  "Test",
					To: []string{
						"recipient@mail.com",
					},
					Subject:  "Test",
					HtmlBody: "<html><head></head><body>Test</body></html>",
					TextBody: "Test",
				},
			},
			wantErr: true,
			key:     "SANDBOX_ERROR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mandrillKey = tt.key
			if err := send(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
			}
			mandrillKey = "SANDBOX_SUCCESS"
		})
	}
}
