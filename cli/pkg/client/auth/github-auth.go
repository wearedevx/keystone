package auth

import (
	"context"
	"strconv"

	"github.com/google/go-github/v32/github"
	"github.com/wearedevx/keystone/api/pkg/models"
	"golang.org/x/oauth2"
)

var (
	githubClientId     string
	githubClientSecret string
)

type PublicKey struct {
	Typ       string
	KeyID     string
	PublicKey string
}

type gitHubAuthService struct {
	apiURL       string
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *github.Client
	token        *oauth2.Token
}

// Name method returns the name of the service
func (g gitHubAuthService) Name() string { return "GitHub" }

// GitHubAuth function returns an instance of GitHubAuth service
func GitHubAuth(ctx context.Context, apiURL string) AuthService {
	return &gitHubAuthService{
		apiURL: apiURL,
		ctx:    ctx,
	}
}

// Start method initiate the oauth process by creating a login request
// on the Keystone server and requesting a login url to GitHub
func (g *gitHubAuthService) Start() (string, error) {
	lr, err := getLoginRequest(g.apiURL)
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

// WaitForExternalLogin method polls the API until the user has completed the
// login process and authorized the Keystone application
func (g *gitHubAuthService) WaitForExternalLogin() error {
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

	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(g.ctx, ts)

	g.token = token
	g.client = github.NewClient(tc)

	return nil
}

// Finish method finish the login process
func (g gitHubAuthService) Finish(
	pk []byte,
	device string,
	deviceUID string,
) (models.User, string, error) {
	return completeLogin(
		g.apiURL,
		models.GitHubAccountType,
		g.token,
		pk,
		device,
		deviceUID,
	)
}

// CheckAccount method returns true if the account is github one
func (g gitHubAuthService) CheckAccount(
	account map[string]string,
) (bool, error) {
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

func (g dummyAuthService) RefreshConnexion(refreshToken string) (string, string, error) {
	return getNewToken(g.apiURL, refreshToken)
}
