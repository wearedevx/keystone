package repo

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/wearedevx/keystone/api/internal/emailer"
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

func (repo *Repo) GetMessagesForUserOnEnvironment(
	device models.Device,
	environment models.Environment,
	message *models.Message,
) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Model(&models.Message{}).
		Preload("Sender").
		Where("recipient_device_id = ? AND environment_id = ?", device.ID, environment.EnvironmentID).
		First(&message).
		Error

	if errors.Is(repo.Err(), ErrorNotFound) {
		repo.err = nil
	}

	return repo
}

func (repo *Repo) GetMessage(message *models.Message) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Where(message).
		First(message).
		Error

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
		Delete(&models.Message{}, messageID).Error

	return repo
}

// Deletes all messages older than a week
// TODO: avoid SQL dialect spicific features
func (repo *Repo) DeleteExpiredMessages() IRepo {
	if repo.err != nil {
		return repo
	}

	switch dialect {
	case DialectPostgres:
		// Per project message deletion
		// NOTE: ::interval is postgres specitfic
		repo.err = repo.GetDb().
			Exec(
				`DELETE
FROM messages m
WHERE 
	m.id IN (
		SELECT mm.id FROM messages mm
		JOIN environments e ON e.environment_id = mm.environment_id
		JOIN projects p ON p.id = e.project_id
		WHERE mm.created_at < now() - (p.ttl || ' days')::interval
	);`,
			).
			Error

	case DialectSQLite:
		// Per project message deletion
		// NOTE:  is postgres specitfic
		// This is for the tests only
		repo.err = repo.GetDb().
			Exec(
				`DELETE
FROM messages
WHERE 
	messages.id IN (
		SELECT mm.id FROM messages mm
		JOIN environments e ON e.environment_id = mm.environment_id
		JOIN projects p ON p.id = e.project_id
		WHERE mm.created_at < date('now', '-' || p.ttl || ' days')
	);`,
			).
			Error
	}

	return repo
}

// GetGroupedMessagesWillExpireByUser method  returns messages that are going
// to expire soon, group by user database id
// TODO: avoid using dialect specific features
func (repo *Repo) GetGroupedMessagesWillExpireByUser(
	groupedMessageUser *map[uint]emailer.GroupedMessagesUser,
) IRepo {
	if repo.Err() != nil {
		return repo
	}

	messages := []models.Message{}

	switch dialect {
	case DialectPostgres:
		repo.err = repo.
			GetDb().
			Model(&models.Message{}).
			Joins("inner join environments on messages.environment_id = environments.environment_id").
			Joins("inner join projects on environments.project_id = projects.id").
			Where("date_trunc('day', messages.created_at) < date_trunc('day', now() - (projects.ttl - projects.days_before_ttl_expiry - 1 || ' days')::interval)").
			Where("date_trunc('day', messages.created_at) > date_trunc('day', now() - (projects.ttl - projects.days_before_ttl_expiry + 1 || ' days')::interval)").
			Preload("Sender").
			Preload("Recipient").
			Preload("Environment").
			Preload("Environment.Project").
			Find(&messages).
			Error

	case DialectSQLite:
		repo.err = repo.
			GetDb().
			Model(&models.Message{}).
			Joins("inner join environments on messages.environment_id = environments.environment_id").
			Joins("inner join projects on environments.project_id = projects.id").
			Where("date(messages.created_at, 'start_of_day') < date('now', 'start_of_day', '-' || projects.ttl - projects.days_before_ttl_expiry - 1 || ' days')").
			Where("date(messages.created_at, 'start_of_day') > date('now', 'start_of_day', '-' || projects.ttl - projects.days_before_ttl_expiry + 1 || ' days')").
			Preload("Sender").
			Preload("Recipient").
			Preload("Environment").
			Preload("Environment.Project").
			Find(&messages).
			Error
	}

	if repo.Err() != nil {
		return repo
	}

	// Group messages by recipient, project and environment.
	for _, message := range messages {
		perUser, ok := (*groupedMessageUser)[message.RecipientID]
		if !ok {
			perUser = emailer.GroupedMessagesUser{
				Recipient: message.Recipient,
				Projects:  make(map[uint]emailer.GroupedMessageProject),
			}
		}

		perProject, ok := perUser.Projects[message.Environment.ProjectID]
		if !ok {
			perProject = emailer.GroupedMessageProject{
				Project:      message.Environment.Project,
				Environments: make(map[string]models.Environment),
			}
		}

		environment, ok := perProject.Environments[message.EnvironmentID]
		if !ok {
			environment = message.Environment
		}

		perProject.Environments[message.EnvironmentID] = environment
		perUser.Projects[message.Environment.ProjectID] = perProject
		(*groupedMessageUser)[message.RecipientID] = perUser

	}

	return repo
}

func (repo *Repo) RemoveOldMessageForRecipient(
	publicKeyID uint,
	environmentID string,
) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Model(&models.Message{}).
		Where("recipient_device_id = ?", publicKeyID).
		Where("environment_id = ?", environmentID).
		Delete(models.Message{}).Error

	return repo
}
