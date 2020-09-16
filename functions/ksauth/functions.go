// Package p contains an HTTP Cloud Function.
package ksauth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/internal/repo"
	"gorm.io/gorm"
)

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

	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	if err = Repo.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// Auth shows the code to copy paste into the cli
func Auth(w http.ResponseWriter, r *http.Request) {
	router := httprouter.New()

	router.POST("/login-request", postLoginRequest)
	router.GET("/login-request", getLoginRequest)
	router.GET("/auth-redirect/:code/", getAuthRedirect)

	router.ServeHTTP(w, r)
}
