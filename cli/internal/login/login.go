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

func NewLoginService(serviceName string) *LoginService {
	ls := new(LoginService)

	ls.c, ls.err = selectAuthService(serviceName)

	return ls
}

func (s *LoginService) Err() error {
	return s.err
}

func (s *LoginService) GetLoginLink(url *string, name *string) *LoginService {
	if s.err != nil {
		return s
	}

	*name = s.c.Name()
	*url, s.err = s.c.Start()

	return s
}

func (s *LoginService) WaitForExternalLogin() *LoginService {
	s.err = s.c.WaitForExternalLogin()

	return s
}

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
		return err
	}

	config.SetAuthToken(jwtToken)
	config.Write()

	return nil
}

func (s *LoginService) CreateAccountAndLogin() error {
	keyPair, err := keys.New(keys.TypeEC)
	if err != nil {
		return err
	}

	// Transfer credentials to the server
	// Create (or get) the user info
	user, jwtToken, err := s.c.Finish(
		keyPair.Public.Value,
		config.GetDeviceName(),
		config.GetDeviceUID(),
	)
	if err != nil {
		return err
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

	return nil
}

func (ls *LoginService) PerformLogin(
	currentAccount models.User,
	accountIndex int,
) bool {
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
