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

func (repo *Repo) GetMessagesForUserOnEnvironment(user User, environment Environment, message *Message) IRepo {
	repo.err = repo.GetDb().Model(&Message{}).Where("recipient_id = ? AND environment_id = ?", user.ID, environment.EnvironmentID).First(&message).Error
	return repo
}

func (repo *Repo) WriteMessage(user User, message Message) IRepo {
	message.SenderID = user.ID
	repo.err = repo.GetDb().Create(&message).Error
	return repo
}

func (repo *Repo) DeleteMessage(messageID int) error {
	repo.GetDb().Delete(&Message{}, messageID)
	return repo.Err()
}
