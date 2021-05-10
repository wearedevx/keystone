package client

import "github.com/wearedevx/keystone/api/pkg/models"

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

func (u *Users) GetPublicKeys(projectId string) ([]models.UserPublicKey, error) {
	var err error
	var result struct {
		keys []models.UserPublicKey
	}

	err = u.r.get("/projects/"+projectId+"/public-keys", &result, nil)

	return result.keys, err
}
