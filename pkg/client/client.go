package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/wearedevx/keystone/internal/crypto"
	. "github.com/wearedevx/keystone/internal/models"
	"golang.org/x/oauth2"
)

var ksauthURL string //= "http://localhost:9000"
var ksapiURL string  //= "http://localhost:9001"

type KeystoneClient interface {
	InitProject(name string) (Project, error)
}

func getLoginRequest() (LoginRequest, error) {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", ksauthURL+"/login-request", nil)
	request.Header.Set("Accept", "application/json; charset=utf-8")

	if err != nil {
		panic(err)
	}

	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var loginRequest LoginRequest

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&loginRequest)

		return loginRequest, err
	} else {
		message := fmt.Sprintf("Request Error: %s", resp.Status)
		fmt.Println(message)

		return loginRequest, errors.New(message)
	}
}

type pollResult struct {
	authCode string
	err      error
}

const MAX_ATTEMPTS int = 12

func pollLoginRequest(code string, c chan pollResult) {
	var done bool = false
	attemps := 0

	for !done {
		attemps = attemps + 1

		time.Sleep(5 * time.Second)

		timeout := time.Duration(20 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		u, _ := url.Parse(ksauthURL + "/login-request")
		q := u.Query()
		q.Add("code", code)

		request, err := http.NewRequest("GET", u.String()+"?"+q.Encode(), nil)
		request.Header.Set("Accept", "application/json; charset=utf-8")

		if err != nil {
			panic(err)
		}

		resp, err := client.Do(request)

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		var loginRequest LoginRequest

		if resp.StatusCode == http.StatusOK {
			err = json.NewDecoder(resp.Body).Decode(&loginRequest)

			if loginRequest.AuthCode != "" {
				r := pollResult{
					authCode: loginRequest.AuthCode,
				}

				c <- r

				done = true
			}

		} else {
			message := fmt.Sprintf("Request Error: %s", resp.Status)
			fmt.Println(message)

			r := pollResult{
				err: errors.New(message),
			}

			c <- r

			done = true
		}

		if attemps == MAX_ATTEMPTS {
			done = true
		}
	}

}

func completeLogin(accountType AccountType, tok *oauth2.Token, pk string) (User, error) {
	var user User
	payload := LoginPayload{
		AccountType: accountType,
		Token:       tok,
		PublicKey:   pk,
	}

	requestPayload := make([]byte, 0)
	buf := bytes.NewBuffer(requestPayload)
	json.NewEncoder(buf).Encode(&payload)

	req, err := http.NewRequest("POST", ksauthURL+"/complete", buf)
	req.Header.Add("Accept", "application/octet-stream")

	if err != nil {
		return user, err
	}

	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)

	if err != nil {
		return user, err
	}

	if resp.StatusCode != 200 {
		return user, fmt.Errorf("Failed to complete login: %s", resp.Status)
	}

	p, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return user, err
	}

	crypto.DecryptWithPublicKey(pk, p, &user)

	return user, nil
}