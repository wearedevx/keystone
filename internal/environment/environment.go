package environment

import (
	. "github.com/wearedevx/keystone/internal/models"
)

// Return environment by id
func GetEnvironnement(environmentID uint) *Environment {
	// TODO
	return nil
}

func GetMessagesByEnvironment(user *User, environmentID uint, userVersionId string) ([]Message, error) {
	// TODO

	// Get environment

	// Compare environment version id and user version id

	// If same, return ok

	// If not, get messages for user

	// - If there's no message, tell user to ask someone to write fresh message for him/her.

	// - Else return messages

	return nil, nil
}
