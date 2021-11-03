// +build test

package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/wearedevx/keystone/api/pkg/models"
	"golang.org/x/oauth2"
)

type dummyAuthService struct {
	apiUrl       string
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	token        *oauth2.Token
}

func (g dummyAuthService) Name() string { return "Dummy" }

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

func (g dummyAuthService) Finish(pk []byte, device string, deviceUID string) (models.User, string, error) {
	return completeLogin(g.apiUrl, models.GitlabAccountType, g.token, pk, device, deviceUID)
}

func (g dummyAuthService) CheckAccount(account map[string]string) (bool, error) {
	var err error

	if err != nil {
		return false, err
	}

	return true, nil
}
