package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	. "github.com/wearedevx/keystone/internal/jwt"
	"github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
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

type Handler = func(params Params, body io.ReadCloser, Repo repo.Repo, user models.User) (Serde, int, error)

func AuthedHandler(handler Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Get user and prepare the repo
		token := r.Header.Get("authorization")

		userID, err := VerifyToken(token)

		if err != nil {
			fmt.Println(err)
			// 404 is returned purposefully, here to not reveal the existence of
			// resources for non authorized requesters
			http.Error(w, "", http.StatusNotFound)
			return
		}

		Repo := new(repo.Repo)

		if user, ok := Repo.GetUser(userID); ok {
			p := newParams(r, params)
			// Actual call to the handler (i.e. Controller function)
			result, status, err := handler(p, r.Body, *Repo, user)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			// serialize the response for the user
			var serialized string

			err = result.Serialize(&serialized)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			out := bytes.NewBufferString(serialized)

			// Write the response if any
			if out.Len() > 0 {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				w.Header().Add("Content-Length", strconv.Itoa(out.Len()))
				w.WriteHeader(status)
				w.Write(out.Bytes())
			}

			if status != http.StatusOK {
				w.WriteHeader(status)
			}

		} else {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		if Repo.Err() != nil {
			fmt.Println(Repo.Err())
		}
	}
}
