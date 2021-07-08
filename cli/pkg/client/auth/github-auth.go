package auth

import (
	"context"
	"strconv"

	"github.com/google/go-github/v32/github"
	"github.com/wearedevx/keystone/api/pkg/models"
	"golang.org/x/oauth2"
)

var githubClientId string
var githubClientSecret string

type PublicKey struct {
	Typ       string
	KeyID     string
	PublicKey string
}

type gitHubAuthService struct {
	apiUrl       string
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *github.Client
	token        *oauth2.Token
}

func (g gitHubAuthService) Name() string { return "GitHub" }

func GitHubAuth(ctx context.Context, apiUrl string) AuthService {
	return &gitHubAuthService{
		apiUrl: apiUrl,
		ctx:    ctx,
	}
}

func (g *gitHubAuthService) Start() (string, error) {
	lr, err := getLoginRequest(g.apiUrl)
	if err != nil {
		return "", err
	}

	g.loginRequest = lr

	g.conf = &oauth2.Config{
		ClientID:     githubClientId,
		ClientSecret: githubClientSecret,
		Scopes:       []string{"user", "user:email"},
		RedirectURL:  authRedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	state, err := makeOAuthState(lr.TemporaryCode)
	if err != nil {
		return "", err
	}

	return g.conf.AuthCodeURL(state, oauth2.AccessTypeOffline), err

}

func (g *gitHubAuthService) WaitForExternalLogin() error {
	c := make(chan pollResult)
	var result pollResult

	go pollLoginRequest(g.apiUrl, g.loginRequest.TemporaryCode, c)

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
	return completeLogin(g.apiUrl, models.GitHubAccountType, g.token, pk)
}

func (g gitHubAuthService) CheckAccount(account map[string]string) (bool, error) {
	gUser, _, err := g.client.Users.Get(g.ctx, "")

	if err != nil {
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
