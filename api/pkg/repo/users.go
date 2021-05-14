package repo

import (
	"bytes"
	"errors"

	. "github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func (r *Repo) GetUser(user *User) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Where(user).First(user).Error

	return r
}

func (r *Repo) GetOrCreateUser(user *User) IRepo {
	if r.Err() != nil {
		return r
	}

	foundUser := User{
		AccountType: user.AccountType,
		ExtID:       user.ExtID,
	}

	r.err = r.GetDb().Where(*&foundUser).First(&foundUser).Error

	if r.err == nil {
		if bytes.Compare(foundUser.PublicKey, user.PublicKey) != 0 {
			foundUser.PublicKey = user.PublicKey
			db.Save(&foundUser)
		}

		*user = foundUser
	} else if errors.Is(r.err, gorm.ErrRecordNotFound) {
		user.UserID = user.Username + "@" + string(user.AccountType)
		r.err = r.GetDb().Create(&user).Error
	}

	return r
}

// From a slice of userIDs (<username>@<service>)
// fetchs the users.
// Returns the found users and a list of not found userIDs
func (r *Repo) FindUsers(userIDs []string, users *map[string]User, notFounds *[]string) IRepo {
	if r.err != nil {
		return r
	}

	userSlice := make([]User, 0)

	db := r.GetDb()

	r.err = db.Where("user_id IN ?", userIDs).Find(&userSlice).Error

	if r.err != nil {
		return r
	}

	for _, userID := range userIDs {
		found := false
		for _, user := range userSlice {
			if user.UserID == userID {
				found = true

				(*users)[userID] = user
				break
			}
		}

		if !found {
			*notFounds = append(*notFounds, userID)
		}
	}

	return r
}
