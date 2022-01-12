package messages

import (
	"errors"
	"fmt"
	"log"
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
	log    *log.Logger
	err    *kserrors.Error
	ctx    *core.Context
	client client.KeystoneClient
}

type MessageService interface {
	Err() *kserrors.Error
	GetMessages() core.ChangesByEnvironment
	SendEnvironments(environments []models.Environment) MessageService
	SendEnvironmentsToOneMember(
		environments []models.Environment,
		member string,
	) MessageService
	DeleteMessages(messagesIds []uint) MessageService
}

// NewMessageService function returns a new instance of MessageService
func NewMessageService(ctx *core.Context) (service MessageService) {
	client, ksErr := client.NewKeystoneClient()

	service = &messageService{
		log:    log.New(log.Writer(), "[Messges] ", 0),
		err:    ksErr,
		ctx:    ctx,
		client: client,
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

	environmentsNeedingASend := changes.ChangedEnvironmentsWithoutPayload()
	if len(environmentsNeedingASend) > 0 {
		s.err = kserrors.EnvironmentsHaveChanged(
			strings.Join(environmentsNeedingASend, ", "),
			nil,
		)
	}

	messagesIds := getMessagesIds(messagesByEnvironment)
	s.DeleteMessages(messagesIds)

	if shouldRunHooks(changes) {
		s.ctx.RunHook()
	}

	return changes
}

func shouldRunHooks(changes core.ChangesByEnvironment) (shouldRun bool) {
	for _, c := range changes.Environments {
		if !c.IsEmpty() && !c.IsSingleVersionChange() {
			shouldRun = true
			break
		}
	}

	return shouldRun
}

// DeleteMessages deletes messages
func (s *messageService) DeleteMessages(messagesIds []uint) MessageService {
	if s.err != nil {
		return s
	}

	for _, id := range messagesIds {
		if id > 0 {
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
	}

	return s
}

// getMessagesIds returns all messages’ ids
// from the API response
func getMessagesIds(
	messagesByEnvironment models.GetMessageByEnvironmentResponse,
) []uint {
	ids := []uint{}

	for _, msgEnv := range messagesByEnvironment.Environments {
		ids = append(ids, msgEnv.Message.ID)
	}

	return ids
}

// fetchNewMessages fetches Messages and dercrypts them
func (s *messageService) fetchNewMessages(
	result *models.GetMessageByEnvironmentResponse,
) {
	var err error

	if s.err != nil {
		return
	}

	sp := spinner.Spinner("Syncing data...")
	sp.Start()

	projectID := s.ctx.GetProjectID()
	s.log.Printf("Fetching messages for project %s\n", projectID)
	*result, err = s.client.Messages().GetMessages(projectID)

	sp.Stop()

	if err != nil {
		if errors.Is(err, auth.ErrorUnauthorized) {
			s.err = kserrors.InvalidConnectionToken(err)
		}
		s.log.Printf("[WARNING] %v\n", err)

		return
	}

	for env, gmr := range result.Environments {
		s.log.Printf("Got message for environment %s: %s\n",
			env, gmr.Message.Uuid)
	}

	s.err = s.decryptMessages(result)
}

// decryptMessages decrypts messages
// Since the payload in GetMessageByEnvironmentResponse is bytes anyway,
// the decryption is done in place.
func (s *messageService) decryptMessages(
	byEnvironment *models.GetMessageByEnvironmentResponse,
) (err *kserrors.Error) {
	s.log.Printf("Decrypting Messages...\n")

	privateKey, e := config.GetUserPrivateKey()
	if e != nil {
		return kserrors.CouldNotDecryptMessages(
			"Failed to get the current user private key",
			e,
		)
	}

	if len(privateKey) == 0 {
		s.log.Printf("[Error] Invalid Private Key length: %d\n", len(privateKey))
	}

	for environmentName, environment := range byEnvironment.Environments {
		msg := environment.Message
		if msg.Sender.UserID != "" && len(msg.Payload) > 0 {
			udevices, e := s.client.Users().GetUserKeys(msg.Sender.UserID)
			if e != nil {
				return kserrors.CouldNotDecryptMessages(
					fmt.Sprintf(
						"Failed to get the public key for user %s",
						msg.Sender.UserID,
					),
					e,
				)
			}

			if len(udevices.Devices) == 0 {
				return kserrors.CouldNotDecryptMessages(
					fmt.Sprintf(
						"User %s has no public keys",
						msg.Sender.UserID,
					),
					nil,
				)
			}

			if len(udevices.Devices) == 0 {
				return kserrors.CouldNotDecryptMessages(
					fmt.Sprintf(
						"User %s has no public keys",
						msg.Sender.UserID,
					),
					nil,
				)
			}

			var udevice models.Device
			for _, device := range udevices.Devices {
				if device.ID == msg.SenderDeviceID {
					udevice = device
				}
			}

			if len(msg.Payload) > 0 {
				d, e := crypto.DecryptMessage(
					privateKey,
					udevice.PublicKey,
					msg.Payload,
				)
				if e != nil {
					return kserrors.CouldNotDecryptMessages("Decryption failed", e)
				}

				environment.Message.Payload = d

				byEnvironment.Environments[environmentName] = environment
				s.log.Printf("Message for %s decryted\n", environmentName)
			}
		}
	}

	return nil
}

// SendEnvironments sends environments to all members of the project
// The API providing the public keys, it should handle reading rights for
// each project member
func (s *messageService) SendEnvironments(
	environments []models.Environment,
) MessageService {
	s.log.Printf("Sending environmnts to everybody\n")
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
		messages, err := s.prepareMessages(
			currentUser,
			senderPrivateKey,
			environment,
		)
		if err != nil {
			s.err = err
			return s
		}
		messagesToWrite.Messages = append(messagesToWrite.Messages, messages...)
	}

	s.sendMessageAndUpdateEnvironment(messagesToWrite)

	if err := s.ctx.RunHook(); err != nil {
		ui.PrintError(err.Error())
	}

	return s
}

func (s *messageService) sendMessageAndUpdateEnvironment(
	messagesToWrite models.MessagesToWritePayload,
) *messageService {
	sp := spinner.Spinner("Sending secrets...")
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
func (s *messageService) SendEnvironmentsToOneMember(
	environments []models.Environment,
	member string,
) MessageService {
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

		userPublicKeys, err := s.client.Users().
			GetEnvironmentPublicKeys(environmentId)
		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				s.err = kserrors.InvalidConnectionToken(err)
			} else {
				s.err = kserrors.UnkownError(err)
			}
			return s
		}

		var recipientDevices models.UserDevices
		found := false

		for _, upk := range userPublicKeys.Keys {
			if upk.UserUID == member {
				recipientDevices.Devices = upk.Devices
				recipientDevices.UserID = upk.UserID
				found = true
			}
		}

		// If receiver has no access to environment, print error and continue to other environments
		if !found {
			kserrors.MemberHasNoAccessToEnv(fmt.Errorf("%s has no public key associated with the environment %s", member, environment.Name)).
				Print()
			sentEnvironmentCount -= 1
			continue
		}

		PayloadContent, err := s.ctx.PrepareMessagePayload(environment)
		if err != nil {
			s.err = kserrors.UnkownError(err)
			return s
		}

		for _, recipientPublicKey := range recipientDevices.Devices {
			message, err := s.prepareMessage(
				senderPrivateKey,
				environment,
				recipientPublicKey,
				recipientDevices.UserID,
				PayloadContent,
			)
			if err != nil {
				s.err = kserrors.UnkownError(err)
				return s
			}
			message.UpdateEnvironmentVersion = false

			messagesToWrite.Messages = append(messagesToWrite.Messages, message)
		}
	}

	s = s.sendMessageAndUpdateEnvironment(messagesToWrite)
	ui.PrintSuccess(
		"Secrets and files sent to user for %d environments.",
		len(environments),
	)

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

	userPublicKeys, err := s.client.Users().
		GetEnvironmentPublicKeys(environmentId)
	if err != nil {
		if errors.Is(err, auth.ErrorUnauthorized) {
			return messages, kserrors.PermissionDenied(environment.Name, err)
		}
		return messages, kserrors.CannotGetEnvironmentKeys(
			environment.Name,
			err,
		)
	}

	s.log.Printf("Will send environment %s to %d devices\n", environment.Name, len(userPublicKeys.Keys))

	PayloadContent, err := s.ctx.PrepareMessagePayload(environment)
	if err != nil {
		return messages, kserrors.PayloadErrors(err)
	}

	// Create one message per user
	for _, userDevices := range userPublicKeys.Keys {
		for _, device := range userDevices.Devices {
			// Do send to current device !!!
			// And all others also of course
			message, err := s.prepareMessage(
				senderPrivateKey,
				environment,
				device,
				userDevices.UserID,
				PayloadContent,
			)
			if err != nil {
				return messages, kserrors.CouldNotEncryptMessages(err)
			}

			messages = append(messages, message)
			s.log.Printf("Will send environment %s to device %s of user %d,\nUsing public key %d\n",
				environment.Name,
				device.UID,
				userDevices.UserID,
				device.PublicKey,
			)

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

	encryptedPayload, err := crypto.EncryptMessage(
		senderPrivateKey,
		userDevice.PublicKey,
		[]byte(payload),
	)
	if err != nil {
		return message, err
	}

	return models.MessageToWritePayload{
		Payload:                  encryptedPayload,
		RecipientID:              recipientID,
		EnvironmentID:            environment.EnvironmentID,
		RecipientDeviceID:        userDevice.ID,
		SenderDeviceUID:          config.GetDeviceUID(),
		UpdateEnvironmentVersion: true,
	}, nil
}
