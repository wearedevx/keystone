package repo

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func (r *Repo) GetUser(user *User) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Preload("Devices").Where("user_id = ?", user.UserID).First(user).Error

	return r
}

func (r *Repo) GetOrCreateUser(user *User) IRepo {
	if r.Err() != nil {
		return r
	}

	db := r.GetDb()

	foundUser := User{
		AccountType: user.AccountType,
		Username:    user.Username,
		UserID:      fmt.Sprintf("%s@%s", user.Username, user.AccountType),
	}

	db.SetupJoinTable(&User{}, "Devices", &UserDevice{})

	r.err = db.Where(&foundUser).
		Preload("Devices", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		First(&foundUser).Error

	if r.err == nil {
		for _, device := range user.Devices {
			found := false
			for _, fdevice := range foundUser.Devices {
				if fdevice.UID == device.UID {
					found = true
					fdevice.DeletedAt = gorm.DeletedAt{}
					if !bytes.Equal(device.PublicKey, fdevice.PublicKey) {
						fdevice.PublicKey = device.PublicKey
					}
					db.Save(&fdevice)
				}
			}
			if !found {
				userID := foundUser.ID
				if err := r.AddNewDevice(device, userID, foundUser.UserID, foundUser.Email).Err(); err != nil {
					r.err = err
					return r
				}
			}
		}

		*user = foundUser
	} else if errors.Is(r.err, gorm.ErrRecordNotFound) {
		user.UserID = user.Username + "@" + string(user.AccountType)

		// Devices
		r.err = db.Omit("Devices").Create(&user).Error
		if r.err != nil {
			return r
		}

		for _, device := range user.Devices {
			if err := r.AddNewDevice(device, user.ID, user.UserID, user.Email).Err(); err != nil {
				r.err = err
				return r
			}
		}

		// Create default orga for user
		orga := models.Organization{
			OwnerID: user.ID,
			Name:    user.UserID,
		}

		if err := r.CreateOrganization(&orga).Err(); err != nil {
			r.err = err
			return r
		}

		r.err = db.Create(&user).Error
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

func (r *Repo) GetUserByEmail(email string, users *[]User) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Where("email = ?", email).Find(users).Error
	if len(*users) == 0 {
		r.err = ErrorNotFound
	}

	return r
}
