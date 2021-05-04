// +build test

package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wearedevx/keystone/internal/models"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type dummyAuthService struct {
	ctx          context.Context
	conf         *oauth2.Config
	loginRequest models.LoginRequest
	client       *gitlab.Client
	token        *oauth2.Token
}

func (g dummyAuthService) Name() string { return "GitLab" }

func DummyAuth(ctx context.Context) AuthService {
	return &dummyAuthService{
		ctx: ctx,
	}
}

func (g *dummyAuthService) Start() (string, error) {
	lr, _ := getLoginRequest()

	g.loginRequest = lr

	return "http://dummy-auth.com/" + lr.TemporaryCode + "/", nil
}

func fakeLoginSuccess(temporaryCode string) {
	timeout := time.Duration(20 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", "http://localhost:9001/auth-redirect/?state="+temporaryCode+"&code=youpicode", nil)

	if err != nil {
		panic(err)
	}

	_, err = client.Do(request)

	if err != nil {
		panic(err)
	}

}

func (g *dummyAuthService) WaitForExternalLogin() error {
	c := make(chan pollResult)
	var result pollResult

	go pollLoginRequest(g.loginRequest.TemporaryCode, c)
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
	return completeLogin(models.GitLabAccountType, g.token, pk)
}

func (g dummyAuthService) CheckAccount(account map[string]string) (bool, error) {
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
