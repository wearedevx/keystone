package client

import (
	"log"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Users struct {
	log *log.Logger
	r   requester
}

// CheckUsersExist method check the existence of the users given in `userIds`
// These are Keystone userIds: `member@service`
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

// GetEnvironmentPublicKeys method get the public keys of all the devices
// that have a read access to the environment
func (u *Users) GetEnvironmentPublicKeys(
	environmentID string,
) (models.PublicKeys, error) {
	var err error
	var result models.PublicKeys

	err = u.r.get("/environments/"+environmentID+"/public-keys", &result, nil)

	return result, err
}

// GetUserKeys method returns the public key of a specific user
func (u *Users) GetUserKeys(
	userID string,
) (result models.UserDevices, err error) {
	err = u.r.get("/users/"+userID+"/key", &result, nil)

	return result, err
}

// InviteUser method sends a Keystone invitation mail
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
