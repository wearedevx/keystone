package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	log "github.com/wearedevx/keystone/api/internal/utils/cloudlogger"
	. "github.com/wearedevx/keystone/api/pkg/jwt"

	"github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"

	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/authconnector"
	"github.com/wearedevx/keystone/api/internal/router"
	. "github.com/wearedevx/keystone/api/internal/utils"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// postUser Gets or Creates a user
func PostUser(w http.ResponseWriter, r *http.Request, _params httprouter.Params) {
	var status int = http.StatusOK
	var responseBody bytes.Buffer
	var err error

	Repo := new(repo.Repo)
	var user *User = &User{}
	var serializedUser string

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			return user.Deserialize(r.Body)
		}).SetStatusError(http.StatusBadRequest),
		NewAction(func() error {
			Repo.GetOrCreateUser(user)

			return Repo.Err()
		}),
		NewAction(func() error {
			return user.Serialize(&serializedUser)
		}),
		NewAction(func() error {
			// 	// Server will not encrypt data.
			// 	// Crypto dependency only in cli side.
			// 	// Server is juste a mailbox.
			in := bytes.NewBufferString(serializedUser)
			// 	// _, e := crypto.EncryptForUser(user, in, &responseBody)
			responseBody = *in
			return nil
			// 	return e
		}),
	})

	if err = runner.Run().Error(); err != nil {
		log.Error(r, err.Error())
		http.Error(w, err.Error(), status)
		return
	}

	status = runner.Status()

	if responseBody.Len() > 0 {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
		w.Write(responseBody.Bytes())
	}

	w.WriteHeader(status)
}

// getUser gets a user
func GetUser(_ router.Params, _ io.ReadCloser, _ repo.Repo, user User) (router.Serde, int, error) {
	return &user, http.StatusOK, nil
}

// Auth Complete route
func PostUserToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	payload := models.LoginPayload{}

	Repo := new(repo.Repo)
	var user models.User
	var serializedUser string
	var jwtToken string
	var responseBody bytes.Buffer

	var connector authconnector.AuthConnector

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			return json.NewDecoder(r.Body).Decode(&payload)
		}),
		NewAction(func() error {
			connector, err = authconnector.GetConnectoForAccountType(payload.AccountType)

			return err
		}),
		NewAction(func() error {
			user, err = connector.GetUserInfo(payload.Token)

			return nil
		}),
		NewAction(func() error {
			user.PublicKey = payload.PublicKey
			Repo.GetOrCreateUser(&user)

			return Repo.Err()
		}),
		NewAction(func() error {
			jwtToken, err = MakeToken(user)

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
		// Use cloud logger?
		//log.Error(r, err.Error())
		log.Error(r, err.Error())
		http.Error(w, err.Error(), runner.Status())
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
	w.Write(responseBody.Bytes())
}

// Route uses a redirect URI for OAuth2 requests
func GetAuthRedirect(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// used to find the matching login request
	temporaryCode := r.URL.Query().Get("state")
	// code given by the third party
	code := r.URL.Query().Get("code")

	if len(temporaryCode) < 16 || len(code) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var response string
	var err error

	Repo := new(repo.Repo)

	Repo.SetLoginRequestCode(temporaryCode, code)

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

func PostLoginRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response string
	var err error
	Repo := new(repo.Repo)

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
func GetLoginRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response string
	var err error
	Repo := new(repo.Repo)

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
