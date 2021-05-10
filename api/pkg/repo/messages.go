package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) GetMessagesForUserOnEnvironment(user User, environment Environment, message *Message) error {

	repo.GetDb().Model(&Message{}).Where("recipient_id = ? AND environment_id = ?", user.ID, environment.EnvironmentID).First(&message)

	return repo.err
}
