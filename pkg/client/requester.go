package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type methodType string

const (
	GET    methodType = "GET"
	POST              = "POST"
	PUT               = "PUT"
	DELETE            = "DELETE"
)

type requester struct {
	userID   string
	jwtToken string
}

func newRequester(userID string, token string) requester {
	return requester{userID: userID, jwtToken: token}
}

func (r *requester) request(method methodType, expectedStatusCode int, path string, data interface{}, result interface{}) error {
	requestPayload := make([]byte, 0)
	buf := bytes.NewBuffer(requestPayload)

	if data != nil {
		json.NewEncoder(buf).Encode(&data)
	}

	req, err := http.NewRequest(string(method), ksapiURL+path, buf)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.jwtToken))

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

	sbuf := new(strings.Builder)
	_, err = io.Copy(sbuf, resp.Body)
	// check errors

	if result != nil {
		err := json.Unmarshal([]byte(sbuf.String()), result)

		return err
	}

	return nil
}

func (r *requester) get(path string, result interface{}) error {
	return r.request(GET, http.StatusOK, path, nil, result)
}

func (r *requester) post(path string, data interface{}, result interface{}) error {
	return r.request(POST, http.StatusCreated, path, data, result)
}

func (r *requester) put(path string, data interface{}, result interface{}) error {
	return r.request(PUT, http.StatusOK, path, data, result)
}

func (r *requester) del(path string, data interface{}, result interface{}) error {
	return r.request(DELETE, http.StatusNoContent, path, data, result)
}
