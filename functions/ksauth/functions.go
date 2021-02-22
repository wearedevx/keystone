// Package p contains an HTTP Cloud Function.A
package ksauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/google/go-github/v32/github"
	"github.com/julienschmidt/httprouter"
	. "github.com/weardevx/keystone/internal/jwt"
	log "github.com/wearedevx/keystone/internal/cloudlogger"
	"github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
	. "github.com/wearedevx/keystone/internal/utils"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// Route to post to to start a login sequence
func postLoginRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response string
	var err error
	Repo := new(repo.Repo)
	Repo.Connect()

	loginRequest := Repo.CreateLoginRequest()

	if err = Repo.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = loginRequest.Serialize(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Content-Length", strconv.Itoa(len(response)))
	fmt.Fprintf(w, response)
}

// Route to poll to check wether the user has completed the login
func getLoginRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response string
	var err error
	Repo := new(repo.Repo)
	Repo.Connect()

	temporaryCode := r.URL.Query().Get("code")

	if len(temporaryCode) < 16 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	loginRequest, found := Repo.GetLoginRequest(temporaryCode)

	if err = Repo.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		log.Error(r, "Login Request not found with: `%s`", temporaryCode)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	err = loginRequest.Serialize(&response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Content-Length", strconv.Itoa(len(response)))

	fmt.Fprintf(w, response)
}

// Route uses a redirect URI for OAuth2 requests
func getAuthRedirect(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	temporaryCode := params.ByName("code")

	if len(temporaryCode) < 16 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var response string
	var err error

	Repo := new(repo.Repo)
	Repo.Connect()

	Repo.SetLoginRequestCode(temporaryCode, r.URL.Query().Get("code"))

	if err = Repo.Err(); err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}

		http.Error(w, err.Error(), code)
		return
	}

	response = "OK"
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", strconv.Itoa(len(response)))
	fmt.Fprintf(w, response)
}

// Auth Complete route
func postUserToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	ctx := context.Background()
	payload := models.LoginPayload{}
	var gUser *github.User
	var gEmails []*github.UserEmail
	userEmail := ""
	userID, err := uuid.NewV4()
	Repo := new(repo.Repo)
	var user models.User
	var serializedUser string
	var jwtToken string
	var responseBody bytes.Buffer
	var client *github.Client

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			return json.NewDecoder(r.Body).Decode(&payload)
		}),
		NewAction(func() error {
			ts := oauth2.StaticTokenSource(payload.Token)
			tc := oauth2.NewClient(ctx, ts)

			client = github.NewClient(tc)

			return nil
		}),
		NewAction(func() error {
			gUser, _, err = client.Users.Get(ctx, "")

			return err
		}),
		NewAction(func() error {
			gEmails, _, err = client.Users.ListEmails(ctx, nil)

			return err
		}),
		NewAction(func() error {
			for _, email := range gEmails {
				if *email.Primary {
					userEmail = *email.Email
					break
				}
			}

			return nil
		}),
		NewAction(func() error {
			Repo.Connect()

			return Repo.Err()
		}),
		NewAction(func() error {
			user = models.User{
				ExtID:       strconv.Itoa(int(*gUser.ID)),
				UserID:      userID.String(),
				AccountType: models.AccountType(payload.AccountType),
				Username:    *gUser.Login,
				Fullname:    *gUser.Name,
				Email:       userEmail,
				PublicKey:   payload.PublicKey,
			}

			Repo.GetOrCreateUser(&user)

			return Repo.Err()
		}),
		NewAction(func() error {
			jwtToken, err := MakeToken(user)

			return err
		}),
		NewAction(func() error {
			return user.Serialize(&serializedUser)
		}),
		NewAction(func() error {
			serializedUserBytes := []byte(serializedUser)
			responseBody = *bytes.NewBuffer(serializedUserBytes)

			return nil
		}),
	})

	if err = runner.Run().Error(); err != nil {
		log.Error(r, err.Error())
		http.Error(w, err.Error(), runner.Status())
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
	w.Write(responseBody.Bytes())
}

// Auth
func Auth(w http.ResponseWriter, r *http.Request) {
	router := httprouter.New()

	router.POST("/login-request", postLoginRequest)
	router.GET("/login-request", getLoginRequest)
	router.GET("/auth-redirect/:code/", getAuthRedirect)
	router.POST("/complete", postUserToken)

	router.ServeHTTP(w, r)
}
