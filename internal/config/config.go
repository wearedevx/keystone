package config

import (
	"os"

	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/pkg/client"
	. "github.com/wearedevx/keystone/ui"
)

func castAccount(rawAccount map[string]string, account *map[string]string) {
	for k, v := range rawAccount {
		*account = make(map[string]string)

		(*account)[k] = v
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
	rawAccounts := viper.Get("accounts").([]map[string]string)
	accounts := make([]map[string]string, len(rawAccounts))

	for i, r := range rawAccounts {
		rawAccount := r
		castAccount(rawAccount, &accounts[i])
	}

	return accounts
}

// Reads the current account from disk
// Returns the account as a map, and its index
// If the user is logged out, the map is empty,
// and the index is -1
func GetCurrentAccount() (map[string]string, int) {
	nullAccount := make(map[string]string)

	if viper.IsSet("current") {
		index := viper.Get("current").(int)
		accounts := GetAllAccounts()

		if index >= 0 && index < len(accounts) {
			account := accounts[index]

			return account, index
		}
	}

	return nullAccount, -1
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

// Returns `true` if the user is logged in
func IsLoggedIn() bool {
	_, index := GetCurrentAccount()

	return index >= 0
}

// finds an account matching `user` in the `account` slice
func FindAccount(c client.AuthService) (map[string]string, int) {
	current := -1
	a := make(map[string]string)

	for i, account := range GetAllAccounts() {
		isAccount, _ := c.CheckAccount(account)

		if isAccount {
			current = i
			a = account
			break
		}
	}

	return a, current
}
