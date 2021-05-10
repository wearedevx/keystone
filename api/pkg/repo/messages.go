package repo

import (
	"encoding/json"
	"io"
	"strings"

	. "github.com/wearedevx/keystone/api/pkg/models"
)

type MessagesPayload struct {
	Messages []Message `json:"messages"`
}

func (gr *MessagesPayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(gr)
}

func (gr *MessagesPayload) Serialize(out *string) error {
	var sb strings.Builder

	err := json.NewEncoder(&sb).Encode(gr)

	*out = sb.String()

	return err
}

func (repo *Repo) GetMessagesForUserOnEnvironment(user User, environment Environment, message *Message) error {

	repo.GetDb().Model(&Message{}).Where("recipient_id = ? AND environment_id = ?", user.ID, environment.EnvironmentID).First(&message)

	return repo.err
}

func (repo *Repo) WriteMessage(user User, message Message) error {
	return nil
}
