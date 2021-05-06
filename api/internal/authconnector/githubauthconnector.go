// Package authconnector provides ...
package authconnector

import (
	"context"
	"strconv"

	"github.com/google/go-github/v32/github"
	. "github.com/wearedevx/keystone/api/internal/utils"
	"github.com/wearedevx/keystone/api/pkg/models"
	"golang.org/x/oauth2"
)

type githubAuthConnector struct {
	token string
}

func (ghac *githubAuthConnector) GetUserInfo(token *oauth2.Token) (models.User, error) {
	var err error
	var user models.User
	userEmail := ""

	ctx := context.Background()

	var client *github.Client
	var gUser *github.User
	var gEmails []*github.UserEmail

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			ts := oauth2.StaticTokenSource(token)
			tc := oauth2.NewClient(ctx, ts)

			client = github.NewClient(tc)

			return nil
		}),
		NewAction(func() error {
			gUser, _, err = client.Users.Get(ctx, "")

			return err
		}),
		NewAction(func() error {
			gEmails, _, err = client.Users.ListEmails(ctx, nil)

			return err
		}),
		NewAction(func() error {
			for _, email := range gEmails {
				if *email.Primary {
					userEmail = *email.Email
					break
				}
			}

			return nil
		}),
		NewAction(func() error {
			userName := "No name"

			if gUser.Name != nil {
				userName = *gUser.Name
			}

			userID := *gUser.Login + "@github"

			user = models.User{
				ExtID:       strconv.Itoa(int(*gUser.ID)),
				UserID:      userID,
				AccountType: models.AccountType("github"),
				Username:    *gUser.Login,
				Fullname:    userName,
				Email:       userEmail,
			}

			return nil
		}),
	})

	err = runner.Run().Error()

	return user, err
}
