package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/wearedevx/keystone/api/internal/activitylog"
	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	. "github.com/wearedevx/keystone/api/pkg/apierrors"
	"github.com/wearedevx/keystone/api/pkg/jwt"
	"github.com/wearedevx/keystone/api/pkg/notification"

	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/authconnector"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"github.com/wearedevx/keystone/api/templates"
)

// postUser Gets or Creates a user
// func PostUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	status := http.StatusOK
// 	var responseBody bytes.Buffer
// 	var err error
//
// 	user := &models.User{}
// 	var serializedUser string
//
// 	err = repo.Transaction(func(Repo repo.IRepo) error {
// 		alogger := activitylog.NewActivityLogger(Repo)
// 		log := models.ActivityLog{
// 			Action: "PostUser",
// 		}
// 		err = nil
// 		status = http.StatusOK
// 		msg := ""
//
// 		if err = user.Deserialize(r.Body); err != nil {
// 			status = http.StatusBadRequest
// 			err = apierrors.ErrorBadRequest(err)
// 			msg = err.Error()
//
// 			goto done
// 		}
//
// 		if err = Repo.GetOrCreateUser(user).Err(); err != nil {
// 			if errors.Is(err, repo.ErrorNotFound) {
// 				status = http.StatusNotFound
// 				err = repo.ErrorNotFound
// 			} else {
// 				status = http.StatusBadRequest
// 			}
// 			msg = err.Error()
//
// 			goto done
// 		}
//
// 		if err = user.Serialize(&serializedUser); err != nil {
// 			msg = "Internal Server Error"
// 			status = http.StatusInternalServerError
//
// 			goto done
// 		}
//
// 	done:
// 		alogger.Save(log.SetError(err))
//
// 		if err != nil {
// 			http.Error(w, msg, status)
// 		}
//
// 		return err
// 	})
//
// 	if err == nil {
// 		in := bytes.NewBufferString(serializedUser)
// 		responseBody = *in
//
// 		if responseBody.Len() > 0 {
// 			w.Header().Add("Content-Type", "application/octet-stream")
// 			w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
// 			_, err := w.Write(responseBody.Bytes())
// 			if err != nil {
// 				fmt.Printf("err: %+v\n", err)
// 				w.WriteHeader(500)
// 				return
// 			}
// 		}
//
// 		w.WriteHeader(status)
// 	}
// }

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
	_ router.Params,
	Repo repo.IRepo,
) (status int, err error) {
	status = http.StatusOK
	payload := models.LoginPayload{}

	var user models.User
	var serializedUser string
	var jwtToken string
	var responseBody bytes.Buffer
	var msg string

	var connector authconnector.AuthConnector

	var newDevices []models.Device

	alogger := activitylog.NewActivityLogger(Repo)
	log := models.ActivityLog{
		Action: "PostUserToken",
	}

	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		status = http.StatusBadRequest
		msg = "Bad Request"
		goto done
	}

	if len(payload.PublicKey) == 0 {
		status = http.StatusBadRequest
		msg = "Bad Request"
		err = errors.New("bad request")
		goto done
	}

	connector, err = authconnector.GetConnectorForAccountType(
		payload.AccountType,
	)
	if err != nil {
		status = http.StatusBadRequest
		msg = "Bad Request"
		goto done
	}

	user, err = connector.GetUserInfo(payload.Token)
	if err != nil {
		status = http.StatusBadRequest
		msg = "Bad Request"
		goto done
	}

	user.Devices = []models.Device{{
		PublicKey: payload.PublicKey,
		Name:      payload.Device,
		UID:       payload.DeviceUID,
	}}

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

	if err = Repo.GetNewlyCreatedDevices(&newDevices).Err(); err != nil {
		msg = "Internal Server Error"
		status = http.StatusInternalServerError
		goto done
	}

	for _, device := range newDevices {
		// Newly created devices only have one user
		user := device.Users[0]

		var adminProjectsMap map[string][]string
		if err = Repo.GetAdminsFromUserProjects(user.ID, &adminProjectsMap).Err(); err != nil {
			status = http.StatusInternalServerError
			msg = err.Error()
			goto done
		}
		if err = notification.SendEmailForNewDevices(device, adminProjectsMap, user); err != nil {
			status = http.StatusInternalServerError
			msg = err.Error()
			goto done
		}
		if err = Repo.SetNewlyCreatedDevice(false, device.ID, user.ID).Err(); err != nil {
			status = http.StatusInternalServerError
			msg = err.Error()
			goto done
		}

	}

	log.User = user
	jwtToken, err = jwt.MakeToken(user, payload.DeviceUID, time.Now())

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
		w.Write([]byte(msg))
	} else {

		serializedUserBytes := []byte(serializedUser)
		responseBody = *bytes.NewBuffer(serializedUserBytes)

		w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))

		if _, err := w.Write(responseBody.Bytes()); err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	}

	return status, err
}

// Route uses a redirect URI for OAuth2 requests
func GetAuthRedirect(
	w http.ResponseWriter,
	r *http.Request,
	_ router.Params,
	Repo repo.IRepo,
) (status int, err error) {
	var response string
	var temporaryCode string
	var code string
	status = http.StatusOK

	type tplData struct {
		Title   string
		Message string
	}

	tpl := "login-success"
	data := tplData{
		Title:   "You have been successfully authenticated",
		Message: `You may now return to your terminal and start using Keystone`,
	}

	// used to find the matching login request
	state := models.AuthState{}
	err = state.Decode(r.URL.Query().Get("state"))
	if err != nil {
		tpl = "login-fail"
		data.Title = "Bad Request"
		data.Message = "The link used is malformed"
		status = http.StatusBadRequest
		goto done
	}

	temporaryCode = state.TemporaryCode

	// code given by the third party
	code = r.URL.Query().Get("code")

	if len(temporaryCode) < 16 || len(code) == 0 {
		tpl = "login-fail"
		data.Title = "Bad Request"
		data.Message = "The provided code is invalid"
		status = http.StatusBadRequest
		goto done
	}

	Repo.SetLoginRequestCode(temporaryCode, code)
	if err = Repo.Err(); err != nil {
		tpl = "login-fail"
		status = http.StatusInternalServerError
		data.Title = "Internal Server Error"
		data.Message = "An unexpected error occurred while trying to log you in"

		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusBadRequest
			data.Title = "Bad Request"
			data.Message = "The link used is invalid or expired"
		}

		goto done
	}

done:
	response, err = templates.RenderTemplate(tpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(response)))
	if _, err = fmt.Fprint(w, response); err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	return status, err
}

func PostLoginRequest(
	w http.ResponseWriter,
	_ *http.Request,
	_ router.Params,
	Repo repo.IRepo,
) (status int, err error) {
	status = http.StatusOK
	var response string

	loginRequest := Repo.CreateLoginRequest()
	alogger := activitylog.NewActivityLogger(Repo)
	log := models.ActivityLog{
		Action: "PostLoginRequest",
	}
	var msg string

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

	if err == nil {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("Content-Length", strconv.Itoa(len(response)))
		if _, err := fmt.Fprint(w, response); err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	} else {
		w.Write([]byte(msg))
	}

	return status, err
}

// Route to poll to check wether the user has completed the login
func GetLoginRequest(
	w http.ResponseWriter,
	r *http.Request,
	_ router.Params,
	Repo repo.IRepo,
) (status int, err error) {
	var response string
	var loginRequest models.LoginRequest
	status = http.StatusOK

	alogger := activitylog.NewActivityLogger(Repo)
	alog := models.ActivityLog{
		Action: "GetLoginRequest",
	}

	temporaryCode := r.URL.Query().Get("code")

	if len(temporaryCode) < 16 {
		status = http.StatusBadRequest
		response = "Bad Request"
		err = errors.New("code too short")
		goto done
	}

	loginRequest, _ = Repo.GetLoginRequest(temporaryCode)

	if err = Repo.Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		response = err.Error()

		goto done
	}

	err = loginRequest.Serialize(&response)

	if err != nil {
		response = err.Error()
		status = http.StatusInternalServerError

		goto done
	}

done:
	alogger.Save(alog.SetError(err))

	w.Header().Add("Content-Length", strconv.Itoa(len(response)))

	if _, err := fmt.Fprint(w, response); err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	return status, err
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
	userDevices := models.UserDevices{
		UserID:  0,
		Devices: make([]models.Device, 0),
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
			err = apierrors.ErrorFailedToGetResource(err)
		}
	} else {
		userDevices.UserID = targetUser.ID
		userDevices.Devices = targetUser.Devices
	}

done:
	return &userDevices, status, log.SetError(err)
}
