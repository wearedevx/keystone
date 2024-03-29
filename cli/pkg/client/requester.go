package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wearedevx/keystone/api/pkg/apierrors"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
)

type methodType string

const (
	GET    methodType = "GET"
	POST   methodType = "POST"
	PUT    methodType = "PUT"
	DELETE methodType = "DELETE"
)

type requester struct {
	log      *log.Logger
	userID   string
	jwtToken string
}

func newRequester(userID string, token string) requester {
	return requester{
		userID:   userID,
		jwtToken: token,
		log:      log.New(log.Writer(), "[requester] ", 0),
	}
}

func (r *requester) request(
	method methodType,
	expectedStatusCode int,
	path string,
	data interface{},
	result interface{},
	params map[string]string,
) error {
	requestPayload := make([]byte, 0)
	buf := bytes.NewBuffer(requestPayload)

	if data != nil {
		if err := json.NewEncoder(buf).Encode(&data); err != nil {
			return err
		}
	}

	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Set(key, value)
		// if err := json.NewEncoder(buf).Encode(&data); err != nil {
		// 	return err
		// }
	}

	// fmt.Println(ApiURL + path)
	Url, err := url.Parse(ApiURL + path)
	if err != nil {
		return err
	}
	Url.RawQuery = queryParams.Encode()

	r.log.Printf("http: %s %s", string(method), Url.String())

	req, err := http.NewRequest(string(method), Url.String(), buf)
	if err != nil {
		return err
	}

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

	if resp == nil {
		return auth.ErrorServiceNotAvailable
	}

	if err != nil {
		return err
	}

	sbuf := new(strings.Builder)
	_, err = io.Copy(sbuf, resp.Body)
	if err != nil {
		return err
	}

	bodyBytes := []byte(sbuf.String())

	// minimum length for json response 2 bytes: {} or []
	if result != nil && len(bodyBytes) >= 2 {
		err := json.Unmarshal(bodyBytes, result)
		if err != nil {
			return apierrors.FromString(string(bodyBytes))
		}
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return auth.ErrorUnauthorized
	}

	if resp.StatusCode != expectedStatusCode &&
		resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	return nil
}

func (r *requester) get(
	path string,
	result interface{},
	params map[string]string,
) error {
	return r.request(GET, http.StatusOK, path, nil, result, params)
}

func (r *requester) post(
	path string,
	data interface{},
	result interface{},
	params map[string]string,
) error {
	return r.request(POST, http.StatusCreated, path, data, result, params)
}

func (r *requester) put(
	path string,
	data interface{},
	result interface{},
	params map[string]string,
) error {
	return r.request(PUT, http.StatusOK, path, data, result, params)
}

func (r *requester) del(
	path string,
	data interface{},
	result interface{},
	params map[string]string,
) error {
	return r.request(DELETE, http.StatusNoContent, path, data, result, params)
}
