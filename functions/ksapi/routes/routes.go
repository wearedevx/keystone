package routes

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/internal/crypto"
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
		userID := r.Header.Get("x-ks-user")

		if userID == "" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		Repo := new(repo.Repo)
		Repo.Connect()

		if user, ok := Repo.GetUser(userID); ok {
			p := newParams(r, params)
			// Actual call to the handler (i.e. Controller function)
			result, status, err := handler(p, r.Body, repo, user)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			// encrypt the response for the user
			var serialized string

			err = result.Serialize(&serialized)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}

			in := bytes.NewBufferString(serialized)
			var out bytes.Buffer
			_, err = crypto.EncryptForPublicKey(user.Keys.Cipher, in, &out)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}

			// Write the response if any
			if out.Len() > 0 {
				w.Header().Add("Content-Type", "application/octet-stream")
				w.Header().Add("Content-Length", strconv.Itoa(out.Len()))
				w.Write(out.Bytes())
			}

			w.WriteHeader(status)

		} else {
			http.Error(w, "", http.StatusBadRequest)
		}
	}
}
