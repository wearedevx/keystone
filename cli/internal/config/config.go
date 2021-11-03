package config

import (
	"os"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

var configFilePath string

func castAccount(
	rawAccount map[interface{}]interface{},
	account *map[string]string,
) {
	*account = make(map[string]string)

	for k, v := range rawAccount {
		(*account)[k.(string)] = v.(string)
	}
}

// Writes the global config to the disk
// Exits with 1 status code
func Write() {
	utils.CreateFileIfNotExists(configFilePath, "")

	if err := viper.WriteConfigAs(configFilePath); err != nil {
		kserrors.CannotSaveConfig(err)
		os.Exit(1)
	}
}

// Adds an account to the config
// ! does not write to disk
func AddAccount(account map[string]string) int {
	accounts := GetAllAccounts()

	accounts = append(accounts, account)

	viper.Set("accounts", accounts)

	return len(accounts) - 1
}

// Reads all accounts from disk
func GetAllAccounts() []map[string]string {
	rawAccounts := viper.Get("accounts")
	ty := reflect.TypeOf(rawAccounts).String()

	switch ty {
	case "[]interface {}":
		a := rawAccounts.([]interface{})
		accounts := make([]map[string]string, len(a))

		for i, r := range a {
			rawAccount := r.(map[interface{}]interface{})
			castAccount(rawAccount, &accounts[i])
		}

		return accounts

	case "[]map[string]string":
		accounts := rawAccounts.([]map[string]string)

		return accounts
	}

	return make([]map[string]string, 0)
}

// Reads the current account from disk
// Returns the account as a map, and its index
// If the user is logged out, the map is empty,
// and the index is -1
func GetCurrentAccount() (user models.User, index int) {
	user = models.User{}
	index = -1
	accounts := GetAllAccounts()

	if viper.IsSet("current") {
		if index = viper.Get("current").(int); index >= 0 &&
			index < len(accounts) {
			user = UserFromAccount(accounts[index])
		}
	}

	return user, index
}

func UserFromAccount(account map[string]string) (user models.User) {
	devices := make([]models.Device, 0)
	devices = append(
		devices,
		models.Device{PublicKey: []byte(account["public_keys"])},
	)
	user.AccountType = models.AccountType(account["account_type"])
	user.Email = account["email"]
	user.ExtID = account["ext_id"]
	user.Fullname = account["fullname"]
	user.Devices = devices
	user.UserID = account["user_id"]
	user.Username = account["username"]

	return user
}

func GetUserPrivateKey() (privateKey []byte, err error) {
	privateKey = []byte(viper.Get("private_key").(string))
	return privateKey, err
}

func GetUserPublicKey() (publicKey []byte, err error) {
	publicKey = []byte(viper.Get("public_key").(string))
	return publicKey, err
}

func SetUserPrivateKey(privateKey string) {
	viper.Set("private_key", privateKey)
}

func SetUserPublicKey(publicKey string) {
	viper.Set("public_key", publicKey)
}

// Sets the current account as the index at `index`
// Means the user is logged in as the user of that account
func SetCurrentAccount(index int) {
	viper.Set("current", index)
}

// Saves the jwt token
func SetAuthToken(token string) {
	viper.Set("auth_token", token)
}

func GetAuthToken() string {
	return viper.Get("auth_token").(string)
}

func GetDeviceName() string {
	if viper.Get("device") != nil {
		return viper.Get("device").(string)
	}
	return ""
}

func GetDeviceUID() string {
	if viper.Get("device_uid") != nil {
		return viper.Get("device_uid").(string)
	}
	return ""
}

func GetServiceApiKey(serviceName string) string {
	token := viper.Get(serviceName + "_auth_token")
	if token != nil {
		return token.(string)
	}
	return ""
}

func SetServiceApiKey(serviceName string, token string) {
	viper.Set(serviceName+"_auth_token", token)
}

// Returns `true` if the user is logged in
func IsLoggedIn() bool {
	_, index := GetCurrentAccount()

	return index >= 0
}

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) (err error) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		configFilePath = cfgFile
	} else {
		configDir, err := ConfigDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigName("keystone")
		viper.SetConfigType("yaml")
		viper.SetConfigPermissions(0o600)

		configFilePath, err = ConfigPath()
		if err != nil {
			return err
		}
	}

	defaultAccounts := make([]map[string]string, 0)

	viper.SetDefault("current", -1)
	viper.SetDefault("accounts", defaultAccounts)

	viper.SetDefault("device_uid", uuid.NewV4().String())

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		err = viper.WriteConfig()
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func CheckExpiredTokenError(err *kserrors.Error) {
	if err.Name() == "Invalid Connection Token" {
		Logout()
	}
}

func Logout() {
	SetCurrentAccount(-1)
	SetAuthToken("")
	Write()
}
