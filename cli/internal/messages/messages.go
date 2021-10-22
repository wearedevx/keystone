package messages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/apierrors"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/crypto"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

type messageService struct {
	err     *kserrors.Error
	ctx     *core.Context
	client  client.KeystoneClient
	printer ui.Printer
}

type MessageService interface {
	Err() *kserrors.Error
	GetMessages() core.ChangesByEnvironment
	SendEnvironments(environments []models.Environment) MessageService
	SendEnvironmentsToOneMember(environments []models.Environment, member string) MessageService
	DeleteMessages(messagesIds []uint) MessageService
}

func NewEchoMessageService(ctx *core.Context) (service MessageService) {
	return newMessageService(ctx, &ui.EchoPrinter{})
}

func NewMessageService(ctx *core.Context) (service MessageService) {
	return newMessageService(ctx, &ui.UiPrinter{})
}

func newMessageService(ctx *core.Context, Printer ui.Printer) (service MessageService) {
	client, ksErr := client.NewKeystoneClient()

	service = &messageService{
		err:     ksErr,
		ctx:     ctx,
		client:  client,
		printer: Printer,
	}

	return service
}

// Err returns the last keyssone error
func (s *messageService) Err() *kserrors.Error {
	return s.err
}

// Retrieve secure messages for user
func (s *messageService) GetMessages() core.ChangesByEnvironment {
	if s.err != nil {
		return core.ChangesByEnvironment{}
	}

	messagesByEnvironment := models.GetMessageByEnvironmentResponse{
		Environments: map[string]models.GetMessageResponse{},
	}

	s.fetchNewMessages(&messagesByEnvironment)
	if s.err != nil {
		return core.ChangesByEnvironment{}
	}

	changes := s.ctx.SaveMessages(messagesByEnvironment)

	if s.ctx.Err() != nil {
		s.err = s.ctx.Err()
		return core.ChangesByEnvironment{}
	}

	s.printChanges(changes, messagesByEnvironment)

	messagesIds := getMessagesIds(messagesByEnvironment)
	s.DeleteMessages(messagesIds)

	return changes
}

// DeleteMessages deletes messages
func (s *messageService) DeleteMessages(messagesIds []uint) MessageService {
	if s.err != nil {
		return s
	}

	for _, id := range messagesIds {
		_, err := s.client.Messages().DeleteMessage(id)
		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				s.err = kserrors.InvalidConnectionToken(err)
			} else {
				s.err = kserrors.UnkownError(err)
			}
			break
		}
	}

	return s
}

// getMessagesIds returns all messages’ ids
// from the API response
func getMessagesIds(messagesByEnvironment models.GetMessageByEnvironmentResponse) []uint {
	ids := []uint{}

	for _, msgEnv := range messagesByEnvironment.Environments {
		ids = append(ids, msgEnv.Message.ID)
	}

	return ids
}

// fetchNewMessages fetches Messages and dercrypts them
func (s *messageService) fetchNewMessages(result *models.GetMessageByEnvironmentResponse) {
	var err error

	if s.err != nil {
		return
	}

	sp := spinner.Spinner(" Syncing data...")
	sp.Start()

	projectID := s.ctx.GetProjectID()
	*result, err = s.client.Messages().GetMessages(projectID)

	if err != nil {
		if errors.Is(err, auth.ErrorUnauthorized) {
			s.err = kserrors.InvalidConnectionToken(err)
		}

		return
	}

	sp.Stop()

	s.err = s.decryptMessages(result)
}

// decryptMessages decrypts messages
// Since the payload in GetMessageByEnvironmentResponse is bytes anyway,
// the decryption is done in place.
func (s *messageService) decryptMessages(byEnvironment *models.GetMessageByEnvironmentResponse) (err *kserrors.Error) {
	privateKey, e := config.GetUserPrivateKey()
	if e != nil {
		return kserrors.CouldNotDecryptMessages("Failed to get the current user private key", e)
	}

	for environmentName, environment := range byEnvironment.Environments {
		msg := environment.Message
		if msg.Sender.UserID != "" && len(msg.Payload) > 0 {
			upks, e := s.client.Users().GetUserPublicKey(msg.Sender.UserID)
			if e != nil {
				return kserrors.CouldNotDecryptMessages(fmt.Sprintf("Failed to get the public key for user %s", msg.Sender.UserID), e)
			}

			if len(upks.PublicKeys) == 0 {
				return kserrors.CouldNotDecryptMessages(
					fmt.Sprintf("User %s has no public keys", msg.Sender.UserID),
					nil,
				)
			}

			var udevice models.Device
			for _, device := range upks.PublicKeys {
				if device.ID == msg.SenderDeviceID {
					udevice = device
				}
			}

			d, e := crypto.DecryptMessage(privateKey, udevice.PublicKey, msg.Payload)
			if e != nil {
				return kserrors.CouldNotDecryptMessages("Decryption failed", e)
			}

			environment.Message.Payload = d

			byEnvironment.Environments[environmentName] = environment
		}
	}

	return nil
}

// printChanges displays changes for environments to the user
func (s *messageService) printChanges(changes core.ChangesByEnvironment, messagesByEnvironments models.GetMessageByEnvironmentResponse) {
	changedEnvironments := make([]string, 0)

	var envList []string = []string{"dev", "staging", "prod"}

	for _, environmentName := range envList {
		environment, ok := messagesByEnvironments.Environments[environmentName]
		if !ok {
			continue
		}

		messageID := environment.Message.ID

		if messageID != 0 {
			// IF changes detected
			if len(changes.Environments[environmentName]) > 0 {
				s.printer.PrintStdErr("Environment " + environmentName + ": " + strconv.Itoa(len(changes.Environments[environmentName])) + " secret(s) changed")

				for _, change := range changes.Environments[environmentName] {
					// No previous cotent => secret is new
					if len(change.From) == 0 {
						s.printer.PrintStdErr(ui.RenderTemplate("secret added", ` {{ "++" | green }} {{ .Secret }} : {{ .To }}`, map[string]string{
							"Secret": change.Name,
							"From":   change.From,
							"To":     change.To,
						}))
					} else if len(change.To) == 0 {
						s.printer.PrintStdErr(ui.RenderTemplate("secret deleted", ` {{ "--" | red }} {{ .Secret }} deleted.`, map[string]string{
							"Secret": change.Name,
						}))
					} else {
						s.printer.PrintStdErr("   " + change.Name + " : " + change.From + " ↦ " + change.To)
					}

				}
			} else {
				s.printer.PrintStdErr("Environment " + environmentName + " up to date ✔")
			}

			if err := s.ctx.Err(); err != nil {
				s.err = err
				return
			}
		} else {
			environmentChanged := s.ctx.EnvironmentVersionHasChanged(environmentName, environment.Environment.VersionID)

			if environmentChanged {
				s.printer.PrintStdErr("Environment " + environmentName + " has changed but no message available. Ask someone to push their secret ⨯")
				changedEnvironments = append(changedEnvironments, environmentName)
			} else {
				s.printer.PrintStdErr("Environment " + environmentName + " up to date ✔")
			}
		}
	}

	if len(changedEnvironments) > 0 {
		s.err = kserrors.EnvironmentsHaveChanged(strings.Join(changedEnvironments, ", "), nil)
	}
}

// SendEnvironments sends environments to all members of the project
// The API providing the public keys, it should handle reading rights for
// each project member
func (s *messageService) SendEnvironments(environments []models.Environment) MessageService {
	if s.err != nil {
		return s
	}

	messagesToWrite := models.MessagesToWritePayload{
		Messages: make([]models.MessageToWritePayload, 0),
	}

	currentUser, senderPrivateKey := s.getCurrentUserInformation()
	if s.err != nil {
		return s
	}

	for _, environment := range environments {
		messages, err := s.prepareMessages(currentUser, senderPrivateKey, environment)
		if err != nil {
			s.err = err
			return s
		}
		messagesToWrite.Messages = append(messagesToWrite.Messages, messages...)
	}

	return s.sendMessageAndUpdateEnvironment(messagesToWrite)
}

func (s *messageService) sendMessageAndUpdateEnvironment(messagesToWrite models.MessagesToWritePayload) *messageService {
	sp := spinner.Spinner(" Sending secrets...")
	sp.Start()

	result, err := s.client.Messages().SendMessages(messagesToWrite)
	sp.Stop()

	if err != nil {
		if errors.Is(err, auth.ErrorUnauthorized) {
			s.err = kserrors.InvalidConnectionToken(err)
			return s
		} else if errors.Is(err, apierrors.ErrorNeedsUpgrade) {
			s.err = kserrors.OrganizationNotPaid(err)
		} else {
			s.err = kserrors.UnkownError(err)
			return s
		}
	}

	for _, environment := range result.Environments {
		if err := s.ctx.UpdateEnvironment(environment).Err(); err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				s.err = kserrors.InvalidConnectionToken(err)
				return s
			} else {
				s.err = kserrors.UnkownError(err)
				return s
			}
		}
	}
	return s
}

// SendEnvironmentsToOneMember sends environments to only one specified
// member.
func (s *messageService) SendEnvironmentsToOneMember(environments []models.Environment, member string) MessageService {
	if s.err != nil {
		return s
	}

	messagesToWrite := models.MessagesToWritePayload{
		Messages: make([]models.MessageToWritePayload, 0),
	}

	_, senderPrivateKey := s.getCurrentUserInformation()

	sentEnvironmentCount := 3
	for _, environment := range environments {
		environmentId := environment.EnvironmentID

		userPublicKeys, err := s.client.Users().GetEnvironmentPublicKeys(environmentId)
		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				s.err = kserrors.InvalidConnectionToken(err)
			} else {
				s.err = kserrors.UnkownError(err)
			}
			return s
		}

		var recipientPublicKeys models.UserPublicKeys
		var found bool = false

		for _, upk := range userPublicKeys.Keys {
			if upk.UserUID == member {
				recipientPublicKeys.PublicKeys = upk.PublicKeys
				recipientPublicKeys.UserID = upk.UserID
				found = true
			}
		}

		// If receiver has no access to environment, print error and continue to other environments
		if !found {
			kserrors.MemberHasNoAccessToEnv(fmt.Errorf("%s has no public key associated with the environment %s", member, environment.Name)).Print()
			sentEnvironmentCount -= 1
			continue
		}

		PayloadContent, err := s.ctx.PrepareMessagePayload(environment)
		if err != nil {
			s.err = kserrors.UnkownError(err)
			return s
		}

		for _, recipientPublicKey := range recipientPublicKeys.PublicKeys {
			message, err := s.prepareMessage(senderPrivateKey, environment, recipientPublicKey, recipientPublicKeys.UserID, PayloadContent)
			if err != nil {
				s.err = kserrors.UnkownError(err)
				return s
			}
			message.UpdateEnvironmentVersion = false

			messagesToWrite.Messages = append(messagesToWrite.Messages, message)
		}
	}

	s = s.sendMessageAndUpdateEnvironment(messagesToWrite)
	ui.PrintSuccess("Secrets and files sent to user for %d environments.", len(environments))

	return s
}

// getCurrentUserInformation returns the currently logged in user
// and their private key.
func (s *messageService) getCurrentUserInformation() (models.User, []byte) {
	var currentUser models.User

	currentUser, index := config.GetCurrentAccount()
	if index < 0 {
		s.err = kserrors.MustBeLoggedIn(nil)
		return currentUser, []byte{}
	}

	senderPrivateKey, err := config.GetUserPrivateKey()
	if err != nil {
		s.err = kserrors.MustBeLoggedIn(nil)
		return currentUser, []byte{}
	}

	return currentUser, senderPrivateKey
}

// prepareMessages creates and encrypts messages
// for oll the user allowed to read the given environment
func (s *messageService) prepareMessages(
	currentUser models.User,
	senderPrivateKey []byte,
	environment models.Environment,
) ([]models.MessageToWritePayload, *kserrors.Error) {
	environmentId := environment.EnvironmentID
	messages := make([]models.MessageToWritePayload, 0)

	userPublicKeys, err := s.client.Users().GetEnvironmentPublicKeys(environmentId)
	if err != nil {
		if errors.Is(err, auth.ErrorUnauthorized) {
			return messages, kserrors.PermissionDenied(environment.Name, err)
		}
		return messages, kserrors.CannotGetEnvironmentKeys(environment.Name, err)
	}

	PayloadContent, err := s.ctx.PrepareMessagePayload(environment)
	if err != nil {
		return messages, kserrors.PayloadErrors(err)
	}

	// Create one message per user
	for _, userPublicKey := range userPublicKeys.Keys {
		for _, userDevice := range userPublicKey.PublicKeys {
			// Do send to current device !!!
			// And all others also of course
			message, err := s.prepareMessage(
				senderPrivateKey,
				environment,
				userDevice,
				userPublicKey.UserID,
				PayloadContent,
			)
			if err != nil {
				return messages, kserrors.CouldNotEncryptMessages(err)
			}

			messages = append(messages, message)
		}
	}

	return messages, nil
}

// prepareMessages creates and encryps one message
// for one environment and one project member.
// Read rights should have been checked beforehand
func (s *messageService) prepareMessage(
	senderPrivateKey []byte,
	environment models.Environment,
	userDevice models.Device,
	recipientID uint,
	payloadContent models.MessagePayload,
) (models.MessageToWritePayload, error) {
	message := models.MessageToWritePayload{}
	var payload string

	err := payloadContent.Serialize(&payload)
	if err != nil {
		return message, err
	}

	encryptedPayload, err := crypto.EncryptMessage(senderPrivateKey, userDevice.PublicKey, []byte(payload))
	if err != nil {
		return message, err
	}

	return models.MessageToWritePayload{
		// UserID:                   RecipientID,
		Payload:                  encryptedPayload,
		RecipientID:              recipientID,
		EnvironmentID:            environment.EnvironmentID,
		RecipientDeviceID:        userDevice.ID,
		SenderDeviceUID:          config.GetDeviceUID(),
		UpdateEnvironmentVersion: true,
	}, nil
}
