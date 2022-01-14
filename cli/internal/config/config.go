package config

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

var configFilePath string

var (
	ErrorNoPublicKey  error = errors.New("no public key")
	ErrorNoPrivateKey       = errors.New("no private key")
)

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

// UserFromAccount function returns a `models.User` instance
// from an account map
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

// GetUserPrivateKey function returns the currently logged in user's private key
func GetUserPrivateKey() (privateKey []byte, err error) {
	pk := viper.GetString("private_key")
	if pk == "" {
		return []byte{}, ErrorNoPrivateKey
	}

	privateKey, err = base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return []byte(pk), nil
	}

	return privateKey, nil
}

// GetUserPublicKey function returns the currently logged in users's public key
func GetUserPublicKey() (publicKey []byte, err error) {
	pk := viper.GetString("public_key")
	if pk == "" {
		return []byte{}, ErrorNoPublicKey
	}

	publicKey, err = base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return []byte(pk), nil
	}

	return publicKey, nil
}

// SetUserPrivateKey function sets the private key for the currenty logged in user
func SetUserPrivateKey(privateKey []byte) {
	encodedKey := base64.StdEncoding.EncodeToString(privateKey)
	viper.Set("private_key", encodedKey)
}

// SetUserPublicKey function sets the pulblic key for the currently logged in use
func SetUserPublicKey(publicKey []byte) {
	encodedKey := base64.StdEncoding.EncodeToString(publicKey)
	viper.Set("public_key", encodedKey)
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

// GetAuthToken function returns the latest auth token we obtained
func GetAuthToken() string {
	return viper.Get("auth_token").(string)
}

// GetDeviceName function returns the current device name
func GetDeviceName() string {
	if viper.Get("device") != nil {
		return viper.Get("device").(string)
	}
	return ""
}

// GetDeviceUID function returns the current device UID
func GetDeviceUID() string {
	if viper.Get("device_uid") != nil {
		return viper.Get("device_uid").(string)
	}
	return ""
}

// GetServiceApiKey function returns the API Key for the named CI service
func GetServiceApiKey(serviceName string) string {
	token := viper.Get(serviceName + "_auth_token")
	if token != nil {
		return token.(string)
	}
	return ""
}

// SetServiceApiKey function sets the API Key for the named CI service
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

// CheckExpiredTokenError function checks if the API error is a token exipiration
// and logsout the user if so
func CheckExpiredTokenError(err *kserrors.Error) {
	if err.Name() == "Invalid Connection Token" {
		Logout()
	}
}

// Logout function logs the user oud
func Logout() {
	SetCurrentAccount(-1)
	SetAuthToken("")
	Write()
}

func AddHook(command string) {
	viper.Set("hook", command)
}

func GetHook() (string, bool) {
	if command := viper.GetString("hook"); command != "" {
		return command, true
	}

	return "", false
}

func RemoveHook() {
	err := unset("hook")

	if err != nil {
		panic(err)
	}
}

// castAccount casts a map[interface{}]interface{}, which is returned by
// viper, into a more manageable map[string]string
func castAccount(
	rawAccount map[interface{}]interface{},
	account *map[string]string,
) {
	*account = make(map[string]string)

	for k, v := range rawAccount {
		(*account)[k.(string)] = v.(string)
	}
}

// deepCasts an entire object (as returned by viper.AllValues(), for
// instance) into an object where `map[interface{}]interface{}` have
// been cast to `map[string]interface{}` (recursively, at that), so that
// it can be used with json/yaml Marshal functions
func deepCast(input interface{}) (output interface{}) {
	t := fmt.Sprintf("%T", input)

	switch t {
	case "[]interface {}":
		v := input.([]interface{})
		if len(v) == 0 {
			output = input
		} else {
			first := v[0]
			t = fmt.Sprintf("%T", first)
			output = make([]interface{}, len(v))

			for i, r := range v {
				output.([]interface{})[i] = deepCast(r)
			}
		}

	case "map[interface {}]interface {}":
		m := input.(map[interface{}]interface{})
		output = make(map[string]interface{})

		for k, v := range m {
			key := k.(string)
			output.(map[string]interface{})[key] = deepCast(v)
		}

	case "map[string]interface {}":
		m := input.(map[string]interface{})
		output = make(map[string]interface{})

		for k, v := range m {
			output.(map[string]interface{})[k] = deepCast(v)
		}

	default:
		output = input
	}

	return output
}

// *Debug only* deepPrint recursively prints the type of a an
// interface{}, assuming it is map, as returned by viper.AllValues()
func deepPrint(input interface{}, level int) {
	t := fmt.Sprintf("%T", input)
	indent := strings.Repeat(" ", level)

	switch t {
	case "[]interface {}":
		v := input.([]interface{})
		if len(v) != 0 {
			first := v[0]
			t = fmt.Sprintf("%T", first)

			deepPrint(first, level+1)
		}

	case "map[interface {}]interface {}":
		m := input.(map[interface{}]interface{})

		for k, v := range m {
			key := k.(string)
			fmt.Printf("%s%v: %T\n", indent, key, v)
			deepPrint(v, level+1)
		}

	case "map[string]interface {}":
		m := input.(map[string]interface{})

		for k, v := range m {
			fmt.Printf("%s%v: %T\n", indent, k, v)
			deepPrint(v, level+1)
		}

	default:
	}
}

// unset removes a key from the configuration
func unset(vars ...string) error {
	cfg := viper.AllSettings()
	vals := cfg

	for _, v := range vars {
		parts := strings.Split(v, ".")
		for i, k := range parts {
			v, ok := vals[k]
			if !ok {
				// Doesn't exist no action needed
				break
			}

			switch len(parts) {
			case i + 1:
				// Last part so delete.
				delete(vals, k)

			default:
				m, ok := v.(map[string]interface{})
				if !ok {
					return fmt.Errorf("unsupported type: %T for %q", v, strings.Join(parts[0:i], "."))
				}
				vals = m
			}
		}
	}

	vals = deepCast(vals).(map[string]interface{})
	// deepPrint(vals, 0)

	b, err := json.MarshalIndent(vals, "", " ")
	if err != nil {
		return err
	}

	if err = viper.ReadConfig(bytes.NewReader(b)); err != nil {
		return err
	}

	Write()

	return nil
}
