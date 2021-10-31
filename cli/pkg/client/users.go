package client

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

type Users struct {
	r requester
}

func (u *Users) CheckUsersExist(
	userIds []string,
) (models.CheckMembersResponse, error) {
	var err error
	var result models.CheckMembersResponse

	payload := models.CheckMembersPayload{
		MemberIDs: userIds,
	}
	err = u.r.post("/users/exist", payload, &result, nil)

	return result, err
}

func (u *Users) GetEnvironmentPublicKeys(
	environmentId string,
) (models.PublicKeys, error) {
	var err error
	var result models.PublicKeys

	err = u.r.get("/environments/"+environmentId+"/public-keys", &result, nil)

	return result, err
}

func (u *Users) GetUserPublicKey(
	userID string,
) (result models.UserPublicKeys, err error) {
	err = u.r.get("/users/"+userID+"/key", &result, nil)

	return result, err
}

func (u *Users) InviteUser(
	userEmail string,
	projectName string,
) (result models.GetInviteResponse, err error) {
	payload := models.InvitePayload{
		Email:       userEmail,
		ProjectName: projectName,
	}
	err = u.r.post("/users/invite", payload, &result, nil)

	return result, err
}
