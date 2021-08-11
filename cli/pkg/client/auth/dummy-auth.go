// +build test

package auth

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type dummyAuthService struct {
	apiUrl       string
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *gitlab.Client
	token        *oauth2.Token
}

func (g dummyAuthService) Name() string { return "Gitlab" }

func DummyAuth(ctx context.Context, apiUrl string) AuthService {
	return &dummyAuthService{
		apiUrl: apiUrl,
		ctx:    ctx,
	}
}

func (g *dummyAuthService) Start() (string, error) {
	lr, _ := getLoginRequest(g.apiUrl)

	g.loginRequest = lr

	return "http://dummy-auth.com/" + lr.TemporaryCode + "/", nil
}

func fakeLoginSuccess(temporaryCode string) {
	timeout := time.Duration(20 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", "http://localhost:9001/auth-redirect/?state="+temporaryCode+"&code=youpicode", nil)

	if err == nil {
		_, err = client.Do(request)
	}

	if err != nil {
		panic(err)
	}

}

func (g *dummyAuthService) WaitForExternalLogin() error {
	c := make(chan pollResult)
	var result pollResult

	go pollLoginRequest(g.apiUrl, g.loginRequest.TemporaryCode, c)
	go fakeLoginSuccess(g.loginRequest.TemporaryCode)

	result = <-c

	if result.err != nil {
		return result.err
	}

	g.token = &oauth2.Token{
		AccessToken:  "access_token",
		TokenType:    "Bearer",
		RefreshToken: "refresh_token",
		Expiry:       time.Unix(0, 0),
	}

	return nil
}

func (g dummyAuthService) Finish(pk []byte) (models.User, string, error) {
	return completeLogin(g.apiUrl, models.GitlabAccountType, g.token, pk)
}

func (g dummyAuthService) CheckAccount(account map[string]string) (bool, error) {
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
