package repo

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type MessagesPayload struct {
	Messages []models.Message `json:"messages"`
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

func (repo *Repo) GetMessagesForUserOnEnvironment(user models.User, environment models.Environment, message *models.Message) IRepo {
	if repo.err != nil {
		return repo
	}

	err := repo.GetDb().
		Model(&models.Message{}).
		Preload("Sender").
		Where("recipient_id = ? AND environment_id = ?", user.ID, environment.EnvironmentID).
		First(&message).
		Error

	if !errors.Is(err, ErrorNotFound) {
		// repo.err = err
	}

	return repo
}

func (repo *Repo) WriteMessage(user models.User, message models.Message) IRepo {
	if repo.err != nil {
		return repo
	}

	message.SenderID = user.ID
	repo.err = repo.GetDb().Create(&message).Error
	return repo
}

func (repo *Repo) DeleteMessage(messageID uint, userID uint) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Model(&models.Message{}).
		Where("recipient_id = ?", userID).
		Where("id = ?", messageID).
		Delete(messageID).Error

	return repo
}

func (repo *Repo) RemoveOldMessageForRecipient(userID uint, environmentID string) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Model(&models.Message{}).
		Where("recipient_id = ?", userID).
		Where("environment_id = ?", environmentID).
		Delete(models.Message{}).Error

	return repo
}
