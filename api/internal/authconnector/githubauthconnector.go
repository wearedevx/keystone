// Package authconnector provides ...
package authconnector

import (
	"context"
	"strconv"

	"github.com/google/go-github/v32/github"
	"github.com/wearedevx/keystone/api/pkg/models"
	"golang.org/x/oauth2"
)

type githubAuthConnector struct {
	token string
}

func (ghac *githubAuthConnector) GetUserInfo(
	token *oauth2.Token,
) (user models.User, err error) {
	userEmail := ""

	ctx := context.Background()

	var client *github.Client
	var gUser *github.User
	var gEmails []*github.UserEmail

	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)

	gUser, _, err = client.Users.Get(ctx, "")
	if err != nil {
		return user, err
	}

	gEmails, _, err = client.Users.ListEmails(ctx, nil)
	if err != nil {
		return user, err
	}

	for _, email := range gEmails {
		if *email.Primary {
			userEmail = *email.Email
			break
		}
	}

	userName := *gUser.Login

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

	return user, err
}
