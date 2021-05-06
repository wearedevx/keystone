package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type gitlabAuthService struct {
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *gitlab.Client
	token        *oauth2.Token
}

func (g gitlabAuthService) Name() string { return "GitLab" }

func GitLabAuth(ctx context.Context) AuthService {
	return &gitlabAuthService{
		ctx: ctx,
	}
}

func (g *gitlabAuthService) Start() (string, error) {
	lr, err := getLoginRequest()

	g.loginRequest = lr

	g.conf = &oauth2.Config{
		// todo put the gitlab ones
		ClientID:     "d372c2f3eebd9c498b41886667609fbdcf149254bcb618ddc199047cbbc46b78",
		ClientSecret: "ffe9317fd42d32ea7db24c79f9ee25a3e30637b886f3bc99f951710c8cdc3650",
		Scopes:       []string{"read_user", "email"},
		RedirectURL:  ksapiURL + "/auth-redirect/",
		// RedirectURL:  ksauthURL + "/auth-redirect/" + lr.TemporaryCode,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://gitlab.com/oauth/authorize",
			TokenURL: "https://gitlab.com/oauth/token",
		},
	}

	return g.conf.AuthCodeURL(lr.TemporaryCode, oauth2.AccessTypeOffline), err
}

func (g *gitlabAuthService) WaitForExternalLogin() error {
	c := make(chan pollResult)
	var result pollResult

	go pollLoginRequest(g.loginRequest.TemporaryCode, c)

	result = <-c

	if result.err != nil {
		return result.err
	}

	token, err := g.conf.Exchange(g.ctx, result.authCode)

	if err != nil {
		return err
	}

	g.token = token
	g.client, err = gitlab.NewOAuthClient(token.AccessToken)

	return nil
}

func (g gitlabAuthService) Finish(pk []byte) (models.User, string, error) {
	return completeLogin(models.GitLabAccountType, g.token, pk)
}

func (g gitlabAuthService) CheckAccount(account map[string]string) (bool, error) {
	gUser, _, err := g.client.Users.CurrentUser()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	if account["account_type"] != string(models.GitLabAccountType) {
		return false, nil
	}

	if account["ext_id"] == strconv.Itoa(gUser.ID) {
		return true, nil
	}

	return false, nil
}
