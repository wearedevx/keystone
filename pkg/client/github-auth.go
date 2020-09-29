package client

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/wearedevx/keystone/internal/models"
	"golang.org/x/oauth2"
)

type PublicKey struct {
	Typ       string
	KeyID     string
	PublicKey string
}

type AuthService interface {
	Name() string
	Start() (string, error)
	WaitForExternalLogin() error
	GetPublicKeys() ([]PublicKey, string, error)
	Finish(pkey PublicKey) (models.User, error)
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
		Scopes:       []string{"user", "read:public_key", "user:email"},
		RedirectURL:  ksauthURL + "/auth-redirect/" + lr.TemporaryCode + "/",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	return g.conf.AuthCodeURL("state", oauth2.AccessTypeOffline), err

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

func (g gitHubAuthService) GetPublicKeys() ([]PublicKey, string, error) {
	keys := make([]PublicKey, 0)
	userPublicKeys, _, err := g.client.Users.ListKeys(g.ctx, "", nil)

	if err != nil {
		return keys, "", err
	}

	for _, githubKey := range userPublicKeys {
		keys = append(keys, PublicKey{
			Typ:       "SSH",
			KeyID:     fmt.Sprintf("%s (ssh)", githubKey.GetTitle()),
			PublicKey: githubKey.GetKey(),
		})
	}

	if len(keys) == 0 {
		return keys, `
To associate the generated public key with your GitHub account,
login to your GitHub account, and navigate to
  https://github.com/settings/ssh/new
`, nil
	}

	return keys, ``, nil
}

func (g gitHubAuthService) Finish(pk PublicKey) (models.User, error) {
	return completeLogin(models.GitHubAccountType, g.token, pk.PublicKey)
}
