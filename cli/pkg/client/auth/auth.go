package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/pkg/constants"
	"golang.org/x/oauth2"
)

var authRedirectURL string

func makeOAuthState(code string) (out string, err error) {
	state := models.AuthState{
		TemporaryCode: code,
		Version:       constants.Version,
	}

	out, err = state.Encode()

	if err != nil {
		return "", err
	}

	return out, nil
}

func getLoginRequest(apiUrl string) (loginRequest models.LoginRequest, err error) {
	var resp *http.Response
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", apiUrl+"/login-request", nil)
	request.Header.Set("Accept", "application/json; charset=utf-8")

	if err == nil {
		resp, err = client.Do(request)
	}

	if err != nil {
		errmsg := fmt.Sprintf("Failed to send login request (%s)", err.Error())
		println(errmsg)
		os.Exit(1)
		return loginRequest, err
	}

	defer resp.Body.Close()

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

func pollLoginRequest(apiUrl string, code string, c chan pollResult) {
	var done bool = false
	attemps := 0

	for !done {
		attemps = attemps + 1

		time.Sleep(5 * time.Second)

		timeout := time.Duration(20 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		u, _ := url.Parse(apiUrl + "/login-request")
		q := u.Query()
		q.Add("code", code)

		var resp *http.Response
		request, err := http.NewRequest("GET", u.String()+"?"+q.Encode(), nil)
		request.Header.Set("Accept", "application/json; charset=utf-8")

		if err == nil {
			resp, err = client.Do(request)
		}

		if err != nil {
			errmsg := fmt.Sprintf("Failed polling login request (%s)", err)
			println(errmsg)
			os.Exit(1)
			return
		}

		defer resp.Body.Close()

		var loginRequest models.LoginRequest

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

func completeLogin(apiUrl string, accountType models.AccountType, tok *oauth2.Token, pk []byte, device string, deviceUID string) (models.User, string, error) {
	var user models.User
	payload := models.LoginPayload{
		AccountType: accountType,
		Token:       tok,
		PublicKey:   pk,
		Device:      device,
		DeviceUID:   deviceUID,
	}

	requestPayload := make([]byte, 0)
	buf := bytes.NewBuffer(requestPayload)
	json.NewEncoder(buf).Encode(&payload)

	req, err := http.NewRequest("POST", apiUrl+"/complete", buf)
	req.Header.Add("Accept", "application/octet-stream")

	if err != nil {
		return user, "", err
	}

	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)

	if err != nil {
		return user, "", err
	}

	if resp.StatusCode != 200 {
		return user, "", fmt.Errorf("Failed to complete login: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	jwtToken := resp.Header.Get("Authorization")
	jwtToken = strings.Replace(jwtToken, "Bearer ", "", 1)

	if jwtToken == "" {
		err = fmt.Errorf("No token was returned")
	}

	return user, jwtToken, err
}
