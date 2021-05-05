// +build test

package authconnector

import (
	"github.com/bxcodec/faker/v3"
	"github.com/wearedevx/keystone/internal/models"
	"golang.org/x/oauth2"
)

type dummyAuthConnector struct{}

func (dac *dummyAuthConnector) GetUserInfo(token *oauth2.Token) (models.User, error) {
	user := models.User{}

	faker.FakeData(&user)
	user.Email = "email@example.com"

	return user, nil
}
