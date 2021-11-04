package login

import (
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

type LoginService struct {
	err error
	c   auth.AuthService
}

// NewLoginService function returns an instance of LoginService
func NewLoginService(serviceName string) *LoginService {
	ls := new(LoginService)

	ls.c, ls.err = selectAuthService(serviceName)

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
func (s *LoginService) LogIntoExisitingAccount(
	accountIndex int,
) error {
	config.SetCurrentAccount(accountIndex)

	publicKey, _ := config.GetUserPublicKey()
	_, jwtToken, err := s.c.Finish(
		publicKey,
		config.GetDeviceName(),
		config.GetDeviceUID(),
	)
	if err != nil {
		s.err = err
		return
	}

	config.SetAuthToken(jwtToken)
	config.Write()
}

// CreateAccountAndLogin method creates a new user account and logs them in
func (s *LoginService) CreateAccountAndLogin() error {
	keyPair, err := keys.New(keys.TypeEC)
	if err != nil {
		s.err = err
		return
	}

	// Transfer credentials to the server
	// Create (or get) the user info
	user, jwtToken, err := s.c.Finish(
		keyPair.Public.Value,
		config.GetDeviceName(),
		config.GetDeviceUID(),
	)
	if err != nil {
		s.err = err
		return
	}

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

	config.SetUserPublicKey(string(keyPair.Public.Value))
	config.SetUserPrivateKey(string(keyPair.Private.Value))

	config.SetCurrentAccount(accountIndex)
	config.SetAuthToken(jwtToken)
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

	return ls
}

func selectAuthService(serviceName string) (auth.AuthService, error) {
	serviceName = prompts.SelectAuthService(serviceName)

	return auth.GetAuthService(serviceName, client.ApiURL)
}

