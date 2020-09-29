// Package p contains an HTTP Cloud Function.
package ksusers

import (
	"bytes"
	"net/http"
	"strconv"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/julienschmidt/httprouter"

	log "github.com/wearedevx/keystone/internal/cloudlogger"
	"github.com/wearedevx/keystone/internal/crypto"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
	. "github.com/wearedevx/keystone/internal/utils"
)

// postUser Gets or Creates a user
func postUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var status int = http.StatusOK
	var responseBody bytes.Buffer
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
			in := bytes.NewBufferString(serializedUser)
			_, e := crypto.EncryptForUser(&user, in, &responseBody)

			return e
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
func getUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var status int = http.StatusOK
	var responseBody bytes.Buffer
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
				in := bytes.NewBufferString(serializedUser)
				_, e := crypto.EncryptForUser(&user, in, &responseBody)

				return e
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

	if responseBody.Len() > 0 {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
		w.Write(responseBody.Bytes())
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
