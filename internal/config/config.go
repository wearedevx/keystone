package config

import (
	"os"

	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/pkg/client"
	. "github.com/wearedevx/keystone/ui"
)

func castAccount(rawAccount map[interface{}]interface{}, account *map[string]string) {
	for k, v := range rawAccount {
		(*account)[k.(string)] = v.(string)
	}
}

func Write() {
	if err := viper.WriteConfig(); err != nil {
		Print(RenderTemplate("config write error", `
{{ ERROR }} {{ . | red }}

You have been successfully logged in, but the configuration file could not be written
`, err.Error()))
		os.Exit(1)
	}
}

func AddAccount(account map[string]string) int {
	accounts := GetAllAccounts()

	accounts = append(accounts, account)

	viper.Set("accounts", accounts)

	return len(accounts) - 1
}

func GetAllAccounts() []map[string]string {
	rawAccounts := viper.Get("accounts").([]interface{})
	accounts := make([]map[string]string, len(rawAccounts))

	for i, r := range rawAccounts {
		rawAccount := r.(map[interface{}]interface{})
		castAccount(rawAccount, &accounts[i])
	}

	return accounts
}

func GetCurrentAccount() (map[string]string, int) {
	nullAccount := make(map[string]string)

	if viper.IsSet("current") {
		index := viper.Get("current").(int)

		accounts := GetAllAccounts()

		if index < len(accounts) {
			account := accounts[index]

			return account, index
		}
	}

	return nullAccount, -1
}

func SetCurrentAccount(index int) {
	viper.Set("current", index)
}

func SetAuthToken(token string) {
	viper.Set("auth_token", token)
}

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
