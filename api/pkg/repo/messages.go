package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) GetMessagesForUserOnEnvironment(user User, environment Environment, messages *[]Message) error {

	repo.err = repo.GetDb().Model(&Message{}).Where("recipient_id = ? AND environment_id = ?", user.ID, environment.EnvironmentID).Find(&messages).Error

	return repo.err
}
