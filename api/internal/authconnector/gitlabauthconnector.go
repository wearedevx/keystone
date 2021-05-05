package authconnector

import (
	"strconv"

	"github.com/wearedevx/keystone/internal/models"
	. "github.com/wearedevx/keystone/internal/utils"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type gitlabAuthConnector struct {
	token string
}

func (glac *gitlabAuthConnector) GetUserInfo(token *oauth2.Token) (models.User, error) {
	var err error
	var user models.User

	var client *gitlab.Client
	var gUser *gitlab.User

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			client, err = gitlab.NewOAuthClient(token.AccessToken)

			return err
		}),
		NewAction(func() error {
			gUser, _, err = client.Users.CurrentUser()

			return err
		}),
		NewAction(func() error {
			userName := "No name"

			if gUser.Name != "" {
				userName = gUser.Name
			}

			userID := gUser.Username + "@gitlab"

			user = models.User{
				ExtID:       strconv.Itoa(int(gUser.ID)),
				UserID:      userID,
				AccountType: models.AccountType("gitlab"),
				Username:    gUser.Username,
				Fullname:    userName,
				Email:       gUser.Email,
			}

			return nil
		}),
	})

	err = runner.Run().Error()

	return user, err
}
