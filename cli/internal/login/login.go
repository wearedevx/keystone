package login

import (
	"errors"
	"log"

	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

type LoginService struct {
	log *log.Logger
	err error
	c   auth.AuthService
}

// NewLoginService function returns an instance of LoginService
func NewLoginService(serviceName string) *LoginService {
	ls := new(LoginService)

	ls.log = log.New(log.Writer(), "[Login] ", 0)
	ls.c, ls.err = selectAuthService(serviceName)

	ls.log.Printf("Selected %v\n", ls.c.Name())

	return ls
}

// Err method returns the last error encountered
func (s *LoginService) Err() error {
	return s.err
}

// GetLoginLink method returns a link to the third party to initiate the oauth
// process
func (s *LoginService) GetLoginLink(url *string, name *string) *LoginService {
	if s.err != nil {
		return s
	}

	*name = s.c.Name()
	*url, s.err = s.c.Start()

	return s
}

// WaitForExternalLogin method polls the API regularly to check if the
// user has completed the login on the third party an authorized access
// for Keystone
func (s *LoginService) WaitForExternalLogin() *LoginService {
	s.err = s.c.WaitForExternalLogin()

	return s
}

// LogIntoExisitingAccount method logs an existing user in
func (s *LoginService) LogIntoExisitingAccount(accountIndex int) {
	config.SetCurrentAccount(accountIndex)

	s.log.Printf("Loging into existing account %d", accountIndex)

	publicKey, err := config.GetUserPublicKey()
	if errors.Is(err, config.ErrorNoPublicKey) {
		keyPair, keyErr := keys.New(keys.TypeEC)
		if keyErr != nil {
			s.err = keyErr
			return
		}
		config.SetUserPublicKey(keyPair.Public.Value)
		config.SetUserPrivateKey(keyPair.Private.Value)
		config.Write()

		publicKey = keyPair.Public.Value
	}

	_, jwtToken, refreshToken, err := s.c.Finish(
		publicKey,
		config.GetDeviceName(),
		config.GetDeviceUID(),
	)
	if err != nil {
		s.err = err
		return
	}

	config.SetAuthToken(jwtToken, refreshToken)
	config.Write()
}

// CreateAccountAndLogin method creates a new user account and logs them in
func (s *LoginService) CreateAccountAndLogin() {
	keyPair, err := keys.New(keys.TypeEC)
	if err != nil {
		s.err = err
		return
	}

	deviceName := config.GetDeviceName()
	deviceUID := config.GetDeviceUID()

	// Transfer credentials to the server
	// Create (or get) the user info
	user, jwtToken, refreshToken, err := s.c.Finish(
		keyPair.Public.Value,
		deviceName,
		deviceUID,
	)
	if err != nil {
		s.err = err
		return
	}

	s.log.Printf("Created account %s with device %s (%s), and public key %v\n",
		user.UserID,
		deviceName,
		deviceUID,
		keyPair.Public.Value,
	)

	// Save the user info in the local config
	accountIndex := config.AddAccount(
		map[string]string{
			"account_type": string(user.AccountType),
			"user_id":      user.UserID,
			"ext_id":       user.ExtID,
			"username":     user.Username,
			"fullname":     user.Fullname,
			"email":        user.Email,
		},
	)

	config.SetUserPublicKey(keyPair.Public.Value)
	config.SetUserPrivateKey(keyPair.Private.Value)

	config.SetCurrentAccount(accountIndex)
	config.SetAuthToken(jwtToken, refreshToken)
	config.Write()
}

// PerformLogin method logs a user in, creating the account if necessary
func (ls *LoginService) PerformLogin(
	currentAccount models.User,
	accountIndex int,
) bool {
	if ls.err != nil {
		return false
	}

	if accountIndex >= 0 {
		// Found an exiting matching account,
		// log into it
		ls.LogIntoExisitingAccount(accountIndex)

		return true
	}

	ls.CreateAccountAndLogin()

	return false
}

// finds an account matching `user` in the `account` slice
func (s *LoginService) FindAccount(
	user *models.User,
	current *int,
) *LoginService {
	if s.err != nil {
		return s
	}

	*current = -1

	for i, account := range config.GetAllAccounts() {
		isAccount, _ := s.c.CheckAccount(account)

		if isAccount {
			*current = i
			*user = config.UserFromAccount(account)

			s.log.Printf("Found account %d %s\n", i, user.UserID)
			break
		}
	}

	return s
}

// PromptDeviceName method asks the user for their device name
func (ls *LoginService) PromptDeviceName(skipPrompts bool) *LoginService {
	if ls.err != nil {
		return ls
	}

	existingName := config.GetDeviceName()
	deviceName := prompts.DeviceName(existingName, skipPrompts)
	viper.Set("device", deviceName)

	ls.log.Printf("Using device %s\n", deviceName)

	return ls
}

func selectAuthService(serviceName string) (auth.AuthService, error) {
	serviceName = prompts.SelectAuthService(serviceName)

	return auth.GetAuthService(serviceName, client.APIURL)
}
