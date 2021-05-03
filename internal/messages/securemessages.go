package securemessages

import (
	. "github.com/wearedevx/keystone/internal/models"
)

// Retrieve secure messages for user for environnement.
func GetMessages(user User, project Project, environment Environment) []Message {
	return nil
}

// Write secure messages from user to all environment members
func WriteMessages(user User, project Project, environment Environment, message string) error {
	// Retrieve members for environment

	// For each one, create secure message.

	// Save it

	return nil
}

// /environment/:env_id/messages?versionid=<versionid>
