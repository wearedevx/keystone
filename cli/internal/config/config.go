package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	. "github.com/wearedevx/keystone/cli/ui"
)

func castAccount(rawAccount map[interface{}]interface{}, account *map[string]string) {
	*account = make(map[string]string)

	for k, v := range rawAccount {
		(*account)[k.(string)] = v.(string)
	}
}

// Writes the global config to the disk
// Exits with 1 status code
func Write() {
	if err := viper.WriteConfig(); err != nil {
		Print(RenderTemplate("config write error", `
{{ ERROR }} {{ . | red }}

You have been successfully logged in, but the configuration file could not be written
`, err.Error()))
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

	if ty == "[]interface {}" {
		a := rawAccounts.([]interface{})
		accounts := make([]map[string]string, len(a))

		for i, r := range a {
			rawAccount := r.(map[interface{}]interface{})
			castAccount(rawAccount, &accounts[i])
		}

		return accounts
	} else if ty == "[]map[string]string" {
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
		if index = viper.Get("current").(int); index >= 0 && index < len(accounts) {
			user = userFromAccount(accounts[index])
		}
	}

	return user, index
}

func userFromAccount(account map[string]string) (user models.User) {
	devices := make([]models.Device, 0)
	devices = append(devices, models.Device{PublicKey: []byte(account["public_keys"])})
	user.AccountType = models.AccountType(account["account_type"])
	user.Email = account["email"]
	user.ExtID = account["ext_id"]
	user.Fullname = account["fullname"]
	user.Devices = devices
	user.UserID = account["user_id"]
	user.Username = account["username"]

	return user
}

func GetCurrentUserPrivateKey() (privateKey []byte, err error) {
	index := -1
	accounts := GetAllAccounts()

	if viper.IsSet("current") {
		index = viper.Get("current").(int)

		if index >= 0 && index < len(accounts) {
			account := accounts[index]

			privateKey = []byte(account["private_key"])
		}
	}

	return privateKey, err
}

func GetCurrentUserPublicKey() (publicKey []byte, err error) {
	index := -1
	accounts := GetAllAccounts()

	if viper.IsSet("current") {
		index = viper.Get("current").(int)

		if index >= 0 && index < len(accounts) {
			account := accounts[index]

			publicKey = []byte(account["public_key"])
		}
	}

	return publicKey, err
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

// finds an account matching `user` in the `account` slice
func FindAccount(c auth.AuthService) (user models.User, current int) {
	current = -1

	for i, account := range GetAllAccounts() {
		isAccount, _ := c.CheckAccount(account)

		if isAccount {
			current = i
			user = userFromAccount(account)
			break
		}
	}

	return user, current
}

// Create conf file if not exist
func createFileIfNotExist(filePath string) {
	// Check if need to create file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// path/to/whatever does not exist

		if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}

		err := ioutil.WriteFile(filePath, []byte(""), 0600)

		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
		}

		// Search config in home directory with name ".keystone" (without extension).
		viper.AddConfigPath(path.Join(home, ".config"))
		viper.SetConfigName("keystone")
		viper.SetConfigType("yaml")
		viper.SetConfigPermissions(0o600)

		createFileIfNotExist(path.Join(home, ".config", "keystone.yaml"))
	}

	defaultAccounts := make([]map[string]string, 0)

	viper.SetDefault("current", -1)
	viper.SetDefault("accounts", defaultAccounts)

	viper.SetDefault("device_uid", uuid.NewV4().String())

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
		err = viper.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}
