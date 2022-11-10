package auth

import (
	"context"
	"strconv"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

var (
	gitlabClientId     string
	gitlabClientSecret string
)

type gitlabAuthService struct {
	apiURL       string
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *gitlab.Client
	token        *oauth2.Token
}

// Name method returns the name of the service
func (g gitlabAuthService) Name() string { return "Gitlab" }

// GitlabAuth function returns an instance of GitLabAuth service
func GitlabAuth(ctx context.Context, apiURL string) AuthService {
	return &gitlabAuthService{
		apiURL: apiURL,
		ctx:    ctx,
	}
}

// Start method initiate the oauth process by creating a login request
// on the Keystone server and requesting a login url to GitLab
func (g *gitlabAuthService) Start() (string, error) {
	lr, err := getLoginRequest(g.apiURL)
	if err != nil {
		return "", err
	}

	g.loginRequest = lr

	g.conf = &oauth2.Config{
		ClientID:     gitlabClientId,
		ClientSecret: gitlabClientSecret,
		Scopes:       []string{"read_user", "email"},
		RedirectURL:  authRedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://gitlab.com/oauth/authorize",
			TokenURL: "https://gitlab.com/oauth/token",
		},
	}

	state, err := makeOAuthState(lr.TemporaryCode)
	if err != nil {
		return "", err
	}

	return g.conf.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// WaitForExternalLogin method polls the API until the user has completed the
// login process and authorized the Keystone application
func (g *gitlabAuthService) WaitForExternalLogin() error {
	c := make(chan pollResult)
	var result pollResult

	go pollLoginRequest(g.apiURL, g.loginRequest.TemporaryCode, c)

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
	if err != nil {
		return err
	}

	return nil
}

// Finish method finish the login process
func (g gitlabAuthService) Finish(
	pk []byte,
	device string,
	deviceUID string,
) (models.User, string, error) {
	return completeLogin(
		g.apiURL,
		models.GitlabAccountType,
		g.token,
		pk,
		device,
		deviceUID,
	)
}

// CheckAccount method returns true if the account is gitlab one
func (g gitlabAuthService) CheckAccount(
	account map[string]string,
) (bool, error) {
	gUser, _, err := g.client.Users.CurrentUser()
	if err != nil {
		return false, err
	}

	if account["account_type"] != string(models.GitlabAccountType) {
		return false, nil
	}

	if account["ext_id"] == strconv.Itoa(gUser.ID) {
		return true, nil
	}

	return false, nil
}
