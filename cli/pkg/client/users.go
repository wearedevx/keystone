package client

import (
	"github.com/wearedevx/keystone/api/controllers"
	"github.com/wearedevx/keystone/api/pkg/models"
)

type Users struct {
	r requester
}

func (u *Users) CheckUsersExist(userIds []string) (models.CheckMembersResponse, error) {
	var err error
	var result models.CheckMembersResponse

	payload := models.CheckMembersPayload{
		MemberIDs: userIds,
	}
	err = u.r.post("/users/exist", payload, &result, nil)

	return result, err
}

func (u *Users) GetEnvironmentPublicKeys(environmentId string) (models.PublicKeys, error) {
	var err error
	var result models.PublicKeys

	err = u.r.get("/environments/"+environmentId+"/public-keys", &result, nil)

	return result, err
}

func (u *Users) GetUserPublicKey(userID string) (result models.UserPublicKey, err error) {
	err = u.r.get("/users/"+userID+"/key", &result, nil)

	return result, err
}

func (u *Users) InviteUser(userEmail string) (result GenericResponse, err error) {
	payload := controllers.InvitePayload{
		Email:       userEmail,
		ProjectName: "keystone",
	}
	err = u.r.post("/users/invite", payload, &result, nil)

	return result, err
}
