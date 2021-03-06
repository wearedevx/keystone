package router

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	. "github.com/wearedevx/keystone/api/pkg/jwt"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type Route struct {
	method  func(*httprouter.Router) func(path string, handler httprouter.Handle)
	path    string
	handler httprouter.Handle
}

type Serde interface {
	Serialize(out *string) error
	Deserialize(in io.Reader) error
}

type Params struct {
	urlParams httprouter.Params
	urlQuery  url.Values
}

func newParams(req *http.Request, params httprouter.Params) Params {
	return Params{
		urlParams: params,
		urlQuery:  req.URL.Query(),
	}
}

func (p Params) Get(key string) interface{} {
	if v := p.urlParams.ByName(key); v != "" {
		return v
	}

	v := p.urlQuery.Get(key)

	if len(v) == 1 {
		return v[0]
	}

	return v
}

type Handler = func(params Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (Serde, int, error)

func AuthedHandler(handler Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		dw := newDoneWriter(w)

		wroteStatus := false
		// Get user and prepare the repo
		token := r.Header.Get("authorization")

		userID, err := VerifyToken(token)

		if err != nil {
			// 404 is returned purposefully, here to not reveal the existence of
			// resources for non authorized requesters
			http.Error(dw, "", http.StatusNotFound)
			return
		}

		err = repo.Transaction(func(Repo repo.IRepo) error {
			user := models.User{UserID: userID}

			Repo.GetUser(&user)

			if err = Repo.Err(); err != nil {
				status := http.StatusInternalServerError
				message := err.Error()

				if errors.Is(err, repo.ErrorNotFound) {
					status = http.StatusNotFound
					message = fmt.Sprintf("No user with id: %s", userID)
				}

				http.Error(dw, message, status)
				return err
			}

			p := newParams(r, params)
			// Actual call to the handler (i.e. Controller function)
			result, status, err := handler(p, r.Body, Repo, user)

			if err != nil {
				http.Error(dw, err.Error(), status)
			}

			// serialize the response for the user
			var serialized string

			if result != nil {
				err = result.Serialize(&serialized)
			}

			if err != nil {
				http.Error(dw, "Error Serializing results", http.StatusInternalServerError)
			}

			out := bytes.NewBufferString(serialized)

			// Write the response if any
			if out.Len() > 0 {
				dw.Header().Add("Content-Type", "application/json; charset=utf-8")
				dw.Header().Add("Content-Length", strconv.Itoa(out.Len()))
				dw.Write(out.Bytes())
			}

			if status != 200 && !wroteStatus {
				dw.WriteHeader(status)
			}

			// if Repo.Err() != nil {
			// 	fmt.Println(Repo.Err())
			// }

			return nil
		})
	}
}
