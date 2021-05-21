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
	"github.com/wearedevx/keystone/api/pkg/jwt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"

	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/authconnector"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// postUser Gets or Creates a user
func PostUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var status int = http.StatusOK
	var responseBody bytes.Buffer
	var err error

	Repo := new(repo.Repo)
	var user *models.User = &models.User{}
	var serializedUser string

	if err = user.Deserialize(r.Body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err = Repo.GetOrCreateUser(user).Err(); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err = user.Serialize(&serializedUser); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	in := bytes.NewBufferString(serializedUser)
	responseBody = *in

	if responseBody.Len() > 0 {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
		w.Write(responseBody.Bytes())
	}

	w.WriteHeader(status)
}

// getUser gets a user
func GetUser(_ router.Params, _ io.ReadCloser, _ repo.Repo, user models.User) (router.Serde, int, error) {
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

	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	connector, err = authconnector.GetConnectoForAccountType(payload.AccountType)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err = connector.GetUserInfo(payload.Token)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user.PublicKey = payload.PublicKey
	if len(user.PublicKey) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err = Repo.GetOrCreateUser(&user).Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jwtToken, err = jwt.MakeToken(user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err = user.Serialize(&serializedUser); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	serializedUserBytes := []byte(serializedUser)
	responseBody = *bytes.NewBuffer(serializedUserBytes)

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
	w.Write(responseBody.Bytes())
}

// Route uses a redirect URI for OAuth2 requests
func GetAuthRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	fmt.Fprint(w, response)
}

func PostLoginRequest(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
	fmt.Fprint(w, response)
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

	fmt.Fprint(w, response)
}