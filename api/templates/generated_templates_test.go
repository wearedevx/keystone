package templates

import (
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	type args struct {
		name string
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "renders login success",
			args: args{
				name: "login-success",
				data: struct {
					Title   string
					Message string
				}{
					Title:   "Success",
					Message: "We are a success",
				},
			},
			want: `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Keystone</title>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
		</style>
	</head>
	<body class="home">
		<header class="fixed-top navbar">
			<div class="container">
				<strong>Keystone</strong>
			</div>
		</header>	
		<div class="container" role="document">
			<div class="content">
<div class="success">
	<div class="title">Success</div>
	<div class="message">We are a success</div>
</div>
</div>
		</div>
	</body>
</html>
`,
			wantErr: false,
		},
		{
			name: "renders login fail",
			args: args{
				name: "login-fail",
				data: struct {
					Title   string
					Message string
				}{
					Title:   "Error",
					Message: "An error occurred while loging you in",
				},
			},
			want: `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Keystone</title>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
		</style>
	</head>
	<body class="home">
		<header class="fixed-top navbar">
			<div class="container">
				<strong>Keystone</strong>
			</div>
		</header>	
		<div class="container" role="document">
			<div class="content">
<div class="error">
	<div class="title">Error</div>
	<div class="message">An error occurred while loging you in</div>
</div>
</div>
		</div>
	</body>
</html>
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderTemplate(tt.args.name, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got = strings.TrimSpace(got)
			got = strings.ReplaceAll(got, "\t", "—")
			got = strings.ReplaceAll(got, " ", "·")
			tt.want = strings.TrimSpace(tt.want)
			tt.want = strings.ReplaceAll(tt.want, "\t", "—")
			tt.want = strings.ReplaceAll(tt.want, " ", "·")

			if got != tt.want {
				t.Errorf("RenderTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
