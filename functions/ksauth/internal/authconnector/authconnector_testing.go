// +build test

package authconnector

import (
	"github.com/wearedevx/keystone/internal/models"
	"golang.org/x/oauth2"
)

type AuthConnector interface {
	// Returns a ready-to-insert User instance
	// filled with data obtained from a third party.
	// The token must be retrieved by the cli
	GetUserInfo(token *oauth2.Token) (models.User, error)
}

// Factory method returning the appropriate
// AuthConnector for the given accontType
func GetConnectoForAccountType(accountType models.AccountType) (AuthConnector, error) {
	return new(dummyAuthConnector), nil
}
