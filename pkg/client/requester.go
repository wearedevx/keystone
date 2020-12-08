package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wearedevx/keystone/internal/crypto"
)

type methodType string

const (
	GET    methodType = "GET"
	POST              = "POST"
	PUT               = "PUT"
	DELETE            = "DELETE"
)

type requester struct {
	userID    string
	publicKey string
}

func newRequester(userID string, publicKey string) requester {
	return requester{userID, publicKey}
}

func (r *requester) request(method methodType, expectedStatusCode int, path string, data interface{}, result interface{}) error {
	requestPayload := make([]byte, 0)
	buf := bytes.NewBuffer(requestPayload)
	json.NewEncoder(buf).Encode(&data)

	req, err := http.NewRequest(string(method), ksapiURL+path, buf)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-ks-user", r.userID)

	if err != nil {
		return err
	}

	timeout := time.Duration(20 * time.Second)
	c := http.Client{
		Timeout: timeout,
	}

	resp, err := c.Do(req)

	if err != nil {
		return nil
	}

	if resp.StatusCode != expectedStatusCode && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status code %d", resp.StatusCode)
	}

	if result != nil {
		p, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil
		}

		crypto.DecryptWithPublicKey(r.publicKey, p, result)
	}

	return nil
}

func (r *requester) get(path string, data interface{}, result interface{}) error {
	return r.request(GET, http.StatusOK, path, data, result)
}

func (r *requester) post(path string, data interface{}, result interface{}) error {
	return r.request(POST, http.StatusCreated, path, data, result)
}

func (r *requester) put(path string, data interface{}, result interface{}) error {
	return r.request(PUT, http.StatusOK, path, data, result)
}

func (r *requester) del(path string, data interface{}) error {
	return r.request(DELETE, http.StatusNoContent, path, data, nil)
}
