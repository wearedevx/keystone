package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v32/github"
	"github.com/wearedevx/keystone/internal/models"
	"golang.org/x/oauth2"
)

type PublicKey struct {
	Typ       string
	KeyID     string
	PublicKey string
}

type gitHubAuthService struct {
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *github.Client
	token        *oauth2.Token
}

func (g gitHubAuthService) Name() string { return "GitHub" }

func GitHubAuth(ctx context.Context) AuthService {
	return &gitHubAuthService{
		ctx: ctx,
	}
}

func (g *gitHubAuthService) Start() (string, error) {
	lr, err := getLoginRequest()

	g.loginRequest = lr

	g.conf = &oauth2.Config{
		ClientID:     "b073f661bc803aecee00",
		ClientSecret: "c2593f5b1e063625c7ed6e542c2757fdb050de2d",
		Scopes:       []string{"user", "user:email"},
		RedirectURL:  ksapiURL + "/auth-redirect/",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	return g.conf.AuthCodeURL(lr.TemporaryCode, oauth2.AccessTypeOffline), err

}

func (g *gitHubAuthService) WaitForExternalLogin() error {
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

	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(g.ctx, ts)

	g.token = token
	g.client = github.NewClient(tc)

	return nil
}

func (g gitHubAuthService) Finish(pk []byte) (models.User, string, error) {
	return completeLogin(models.GitHubAccountType, g.token, pk)
}

func (g gitHubAuthService) CheckAccount(account map[string]string) (bool, error) {
	gUser, _, err := g.client.Users.Get(g.ctx, "")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false, err
	}

	if account["account_type"] != string(models.GitHubAccountType) {
		return false, nil
	}

	if account["ext_id"] == strconv.Itoa(int(*gUser.ID)) {
		return true, nil
	}

	return false, nil
}
