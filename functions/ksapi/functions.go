// Package p contains an HTTP Cloud Function.
package ksusers

import (
	"net/http"
	"strconv"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/julienschmidt/httprouter"

	"github.com/wearedevx/keystone/functions/ksapi/crypto"
	log "github.com/wearedevx/keystone/internal/cloudlogger"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
	. "github.com/wearedevx/keystone/internal/utils"
)

// postUser Gets or Creates a user
func postUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var status int = http.StatusOK
	var responseBody []byte
	var err error

	Repo := new(repo.Repo)
	var user User = User{}
	var serializedUser string

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			Repo.Connect()

			return Repo.Err()
		}),
		NewAction(func() error {
			return user.Deserialize(r.Body)
		}).SetStatusError(http.StatusBadRequest),
		NewAction(func() error {
			Repo.GetOrCreateUser(&user)

			return Repo.Err()
		}),
		NewAction(func() error {
			return user.Serialize(&serializedUser)
		}),
		NewAction(func() error {
			return crypto.EncryptForUser(&user, []byte(serializedUser), &responseBody)
		}),
	})

	if err = runner.Run().Error(); err != nil {
		log.Error(r, err.Error())
		http.Error(w, err.Error(), status)
		return
	}

	status = runner.Status()

	if len(responseBody) > 0 {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(len(responseBody)))
		w.Write(responseBody)
	}

	w.WriteHeader(status)
}

// getUser gets a user
func getUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var status int = http.StatusOK
	var responseBody []byte
	var err error

	if userID := r.Header.Get("x-ks-user"); userID != "" {
		Repo := new(repo.Repo)
		user := User{}
		var serializedUser string

		runner := NewRunner([]RunnerAction{
			NewAction(func() error {
				Repo.Connect()

				return Repo.Err()
			}),
			NewAction(func() error {
				user, _ = Repo.GetUser(userID)

				return Repo.Err()
			}),
			NewAction(func() error {
				return user.Serialize(&serializedUser)
			}),
			NewAction(func() error {
				return crypto.EncryptForUser(&user, []byte(serializedUser), &responseBody)
			}),
		})

		if err = runner.Run().Error(); err != nil {
			log.Error(r, err.Error())
			http.Error(w, err.Error(), status)
			return
		}

		status = runner.Status()
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if len(responseBody) > 0 {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(len(responseBody)))
		w.Write(responseBody)
	}

	w.WriteHeader(status)
}

// Auth shows the code to copy paste into the cli
func UserService(w http.ResponseWriter, r *http.Request) {
	router := httprouter.New()

	router.POST("/", postUser)
	router.GET("/", getUser)

	router.ServeHTTP(w, r)
}
