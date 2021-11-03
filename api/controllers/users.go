package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/wearedevx/keystone/api/internal/activitylog"
	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	log "github.com/wearedevx/keystone/api/internal/utils/cloudlogger"
	. "github.com/wearedevx/keystone/api/pkg/apierrors"
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
	status := http.StatusOK
	var responseBody bytes.Buffer
	var err error

	user := &models.User{}
	var serializedUser string

	err = repo.Transaction(func(Repo repo.IRepo) error {
		alogger := activitylog.NewActivityLogger(Repo)
		log := models.ActivityLog{
			Action: "PostUser",
		}
		err = nil
		status = http.StatusOK
		msg := ""

		if err = user.Deserialize(r.Body); err != nil {
			status = http.StatusBadRequest
			err = apierrors.ErrorBadRequest(err)
			msg = err.Error()

			goto done
		}

		if err = Repo.GetOrCreateUser(user).Err(); err != nil {
			if errors.Is(err, repo.ErrorNotFound) {
				status = http.StatusNotFound
				err = repo.ErrorNotFound
			} else {
				status = http.StatusBadRequest
			}
			msg = err.Error()

			goto done
		}

		if err = user.Serialize(&serializedUser); err != nil {
			msg = "Internal Server Error"
			status = http.StatusInternalServerError

			goto done
		}

	done:
		alogger.Save(log.SetError(err))

		if err != nil {
			http.Error(w, msg, status)
		}

		return err
	})

	if err == nil {
		in := bytes.NewBufferString(serializedUser)
		responseBody = *in

		if responseBody.Len() > 0 {
			w.Header().Add("Content-Type", "application/octet-stream")
			w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
			_, err := w.Write(responseBody.Bytes())
			if err != nil {
				fmt.Printf("err: %+v\n", err)
				w.WriteHeader(500)
				return
			}
		}

		w.WriteHeader(status)
	}
}

// getUser gets a user
func GetUser(
	_ router.Params,
	_ io.ReadCloser,
	_ repo.IRepo,
	user models.User,
) (router.Serde, int, error) {
	return &user, http.StatusOK, nil
}

// Auth Complete route
func PostUserToken(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	var err error
	payload := models.LoginPayload{}

	var user models.User
	var serializedUser string
	var jwtToken string
	var responseBody bytes.Buffer

	var connector authconnector.AuthConnector

	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	connector, err = authconnector.GetConnectoForAccountType(
		payload.AccountType,
	)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err = connector.GetUserInfo(payload.Token)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user.Devices = []models.Device{{
		PublicKey: payload.PublicKey,
		Name:      payload.Device,
		UID:       payload.DeviceUID,
	}}

	if len(payload.PublicKey) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = repo.Transaction(func(Repo repo.IRepo) error {
		err = nil
		msg := ""
		status := http.StatusOK
		alogger := activitylog.NewActivityLogger(Repo)
		log := models.ActivityLog{
			Action: "PostUserToken",
		}

		if err = Repo.GetOrCreateUser(&user).Err(); err != nil {
			if errors.Is(err, ErrorBadDeviceName) {
				msg = err.Error()
				status = http.StatusConflict
			} else {
				msg = "Internal Server Error"
				status = http.StatusInternalServerError
			}

			goto done
		}

		log.User = user
		jwtToken, err = jwt.MakeToken(user, payload.DeviceUID)
		fmt.Printf("payload.DeviceUID: %+v\n", payload.DeviceUID)

		if err != nil {
			msg = "Internal Server Error"
			status = http.StatusInternalServerError

			goto done
		}

		if err = user.Serialize(&serializedUser); err != nil {
			msg = "Internal Server Error"
			status = http.StatusInternalServerError

			goto done
		}

	done:
		alogger.Save(log.SetError(err))

		if err != nil {
			http.Error(w, msg, status)
		}

		return err
	})

	if err == nil {
		serializedUserBytes := []byte(serializedUser)
		responseBody = *bytes.NewBuffer(serializedUserBytes)

		w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))

		if _, err := w.Write(responseBody.Bytes()); err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	}
}

// Route uses a redirect URI for OAuth2 requests
func GetAuthRedirect(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	var err error
	var response string

	// used to find the matching login request
	state := models.AuthState{}
	err = state.Decode(r.URL.Query().Get("state"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	temporaryCode := state.TemporaryCode

	// code given by the third party
	code := r.URL.Query().Get("code")

	if len(temporaryCode) < 16 || len(code) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = repo.Transaction(func(Repo repo.IRepo) error {
		Repo.SetLoginRequestCode(temporaryCode, code)
		if err = Repo.Err(); err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				code = http.StatusNotFound
			}

			http.Error(w, err.Error(), code)
			return err
		}
		return nil
	})

	if err == nil {
		response = `You have been successfully authenticated.
You may now return to your terminal and start using Keystone.

Thank you!`
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Length", strconv.Itoa(len(response)))
		if _, err := fmt.Fprint(w, response); err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	}
}

func PostLoginRequest(
	w http.ResponseWriter,
	_ *http.Request,
	_ httprouter.Params,
) {
	var response string
	err := repo.Transaction(func(Repo repo.IRepo) (err error) {
		loginRequest := Repo.CreateLoginRequest()
		alogger := activitylog.NewActivityLogger(Repo)
		log := models.ActivityLog{
			Action: "PostLoginRequest",
		}
		var msg string
		var status int

		if err = Repo.Err(); err != nil {
			msg = "Status Internal Server Error"
			status = http.StatusInternalServerError

			goto done
		}

		if err = loginRequest.Serialize(&response); err != nil {
			msg = "Status Internal Server Error"
			status = http.StatusInternalServerError

			goto done
		}

	done:
		alogger.Save(log.SetError(err))
		if err != nil {
			http.Error(w, msg, status)
		}

		return err
	})

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("Content-Length", strconv.Itoa(len(response)))
		if _, err := fmt.Fprint(w, response); err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	}
}

// Route to poll to check wether the user has completed the login
func GetLoginRequest(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	var response string
	var err error

	temporaryCode := r.URL.Query().Get("code")

	if len(temporaryCode) < 16 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = repo.Transaction(func(Repo repo.IRepo) (err error) {
		alogger := activitylog.NewActivityLogger(Repo)
		alog := models.ActivityLog{
			Action: "GetLoginRequest",
		}
		var status int
		var msg string

		loginRequest, found := Repo.GetLoginRequest(temporaryCode)

		if err = Repo.Err(); err != nil {
			msg = err.Error()
			status = http.StatusInternalServerError

			goto done
		}

		if !found {
			log.Error(r, "Login Request not found with: `%s`", temporaryCode)
			msg = "Not Found"
			status = http.StatusNotFound

			goto done
		}

		err = loginRequest.Serialize(&response)

		if err != nil {
			msg = err.Error()
			status = http.StatusInternalServerError

			goto done
		}

	done:
		alogger.Save(alog.SetError(err))
		if err != nil {
			http.Error(w, msg, status)
		}

		return nil
	})

	if err == nil {
		w.Header().Add("Content-Length", strconv.Itoa(len(response)))

		if _, err := fmt.Fprint(w, response); err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	}
}

func GetUserKeys(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	_ models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		Action: "GetUserKeys",
	}

	targetUser := models.User{}
	userPublicKeys := models.UserPublicKeys{
		UserID:     0,
		PublicKeys: make([]models.Device, 0),
	}

	userID := params.Get("userID")

	if userID == "" {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)

		goto done
	}

	targetUser.UserID = userID

	if err = Repo.GetUser(&targetUser).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToDeleteResource(err)
		}
	} else {
		userPublicKeys.UserID = targetUser.ID
		userPublicKeys.PublicKeys = targetUser.Devices
	}

done:
	return &userPublicKeys, status, log.SetError(err)
}
