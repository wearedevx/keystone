package repo

import (
	"errors"
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func (r *Repo) GetUser(user *models.User) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("Devices").
		Where("user_id = ?", user.UserID).
		First(user).
		Error

	return r
}

func (r *Repo) findDeletedDevice(
	device *models.Device,
) (err error) {
	if r.Err() != nil {
		return r.Err()
	}

	err = r.GetDb().
		Unscoped().
		Where("uid = ?", device.UID).
		First(&device).
		Error

	return err
}

func (r *Repo) undeleteOrCreateDevices(
	user *models.User,
) *Repo {
	if r.err != nil {
		return r
	}

	for _, userDevice := range user.Devices {
		if err := r.findDeletedDevice(&userDevice); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := r.AddNewDevice(
					userDevice,
					user.ID,
					user.UserID,
					user.Email,
				).Err(); err != nil {
					r.err = err
					return r
				}
			} else {
				r.err = err
				break
			}
		} else {
			fmt.Printf("userDevice: %+v\n", userDevice)
			userDevice.DeletedAt = gorm.DeletedAt{}
			r.err = db.Save(&userDevice).Error
		}

		if r.err != nil {
			return r
		}
	}

	return r
}

func (r *Repo) GetOrCreateUser(user *models.User) IRepo {
	if r.Err() != nil {
		return r
	}

	db := r.GetDb()

	foundUser := models.User{
		AccountType: user.AccountType,
		Username:    user.Username,
		UserID:      fmt.Sprintf("%s@%s", user.Username, user.AccountType),
	}

	r.err = db.
		Where(foundUser).
		Preload("Devices").
		Preload("Organizations").
		First(&foundUser).
		Error

	if r.err == nil {
		// Undelete devices using data coming from the CLI
		// NOTE: Should we not update foundUser Device array here?
		r.undeleteOrCreateDevices(user)

		*user = foundUser
	} else if errors.Is(r.err, gorm.ErrRecordNotFound) {
		user.UserID = user.Username + "@" + string(user.AccountType)

		r.err = db.Omit("Devices").Create(&user).Error
		if r.err != nil {
			return r
		}

		// Devices
		for _, device := range user.Devices {
			if err := r.AddNewDevice(
				device,
				user.ID,
				user.UserID,
				user.Email).
				Err(); err != nil {
				r.err = err
				return r
			}
		}

		// Create default orga for user
		orga := models.Organization{
			UserID:  user.ID,
			Name:    user.UserID,
			Private: true,
		}

		if r.err = r.CreateOrganization(&orga).Err(); r.err != nil {
			return r
		}

		user.Organizations = append(user.Organizations, orga)
	}

	return r
}

// From a slice of userIDs (<username>@<service>)
// fetchs the users.
// Returns the found users and a list of not found userIDs
func (r *Repo) FindUsers(
	userIDs []string,
	users *map[string]models.User,
	notFounds *[]string,
) IRepo {
	if r.err != nil {
		return r
	}

	userSlice := make([]models.User, 0)

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

func (r *Repo) GetUserByEmail(email string, users *[]models.User) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Where("email = ?", email).Find(users).Error
	if len(*users) == 0 {
		r.err = ErrorNotFound
	}

	return r
}
