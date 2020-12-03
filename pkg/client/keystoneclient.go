package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wearedevx/keystone/internal/crypto"
	. "github.com/wearedevx/keystone/internal/models"
)

type SKeystoneClient struct {
	UserID string
	pk     string
}

func NewKeystoneClient(userID string, pk string) KeystoneClient {
	return &SKeystoneClient{
		UserID: userID,
		pk:     pk,
	}
}

func (client *SKeystoneClient) InitProject(name string) (Project, error) {
	var project Project

	payload := Project{
		Name: name,
	}

	requestPayload := make([]byte, 0)
	buf := bytes.NewBuffer(requestPayload)
	json.NewEncoder(buf).Encode(&payload)

	req, err := http.NewRequest("POST", ksapiURL+"/projects", buf)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-ks-user", client.UserID)

	if err != nil {
		return project, err
	}

	timeout := time.Duration(20 * time.Second)
	c := http.Client{
		Timeout: timeout,
	}

	resp, err := c.Do(req)

	if err != nil {
		return project, err
	}

	if resp.StatusCode != 200 {
		return project, fmt.Errorf("Failed to complete login: %s", resp.Status)
	}

	p, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return project, err
	}

	crypto.DecryptWithPublicKey(client.pk, p, &project)

	return project, nil
}
