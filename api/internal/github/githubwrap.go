// +build !test

package github

import (
	"context"

	"github.com/google/go-github/v32/github"
)

func GetUser(client *github.Client, ctx context.Context) (*github.User, *github.Response, error) {

	return client.Users.Get(ctx, "")
}

func ListEmails(client *github.Client, ctx context.Context) ([]*github.UserEmail, *github.Response, error) {
	return client.Users.ListEmails(ctx, nil)
}
