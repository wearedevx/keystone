// +build !test

package authconnector

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
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
func GetConnectorForAccountType(
	accountType models.AccountType,
) (AuthConnector, error) {
	switch accountType {
	case models.GitHubAccountType:
		return new(githubAuthConnector), nil

	case models.GitlabAccountType:
		return new(gitlabAuthConnector), nil

	default:
		return nil, fmt.Errorf("no connector for account type %s", accountType)
	}
}
