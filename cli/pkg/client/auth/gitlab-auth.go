package auth

import (
	"context"
	"strconv"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

var gitlabClientId string
var gitlabClientSecret string

type gitlabAuthService struct {
	apiUrl       string
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *gitlab.Client
	token        *oauth2.Token
}

func (g gitlabAuthService) Name() string { return "GitLab" }

func GitLabAuth(ctx context.Context, apiUrl string) AuthService {
	return &gitlabAuthService{
		apiUrl: apiUrl,
		ctx:    ctx,
	}
}

func (g *gitlabAuthService) Start() (string, error) {
	lr, err := getLoginRequest(g.apiUrl)

	g.loginRequest = lr

	g.conf = &oauth2.Config{
		// todo put the gitlab ones
		ClientID:     gitlabClientId,
		ClientSecret: gitlabClientSecret,
		Scopes:       []string{"read_user", "email"},
		RedirectURL:  g.apiUrl + "/auth-redirect/",
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

	go pollLoginRequest(g.apiUrl, g.loginRequest.TemporaryCode, c)

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
	return completeLogin(g.apiUrl, models.GitLabAccountType, g.token, pk)
}

func (g gitlabAuthService) CheckAccount(account map[string]string) (bool, error) {
	gUser, _, err := g.client.Users.CurrentUser()

	if err != nil {
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
