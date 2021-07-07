package authconnector

import (
	"strconv"

	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type gitlabAuthConnector struct {
	token string
}

func (glac *gitlabAuthConnector) GetUserInfo(token *oauth2.Token) (user User, err error) {
	var client *gitlab.Client
	var gUser *gitlab.User

	client, err = gitlab.NewOAuthClient(token.AccessToken)
	if err != nil {
		return user, err
	}

	gUser, _, err = client.Users.CurrentUser()
	if err != nil {
		return user, err
	}

	userName := "No name"

	if gUser.Name != "" {
		userName = gUser.Name
	}

	userID := gUser.Username + "@gitlab"

	user = User{
		ExtID:       strconv.Itoa(int(gUser.ID)),
		UserID:      userID,
		AccountType: AccountType("gitlab"),
		Username:    gUser.Username,
		Fullname:    userName,
		Email:       gUser.Email,
	}

	return user, err
}
