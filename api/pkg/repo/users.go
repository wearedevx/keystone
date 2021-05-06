package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (r *Repo) GetUser(userID string) (User, bool) {
	var user User

	r.err = r.GetDb().Where("user_id = ?", userID).First(&user).Error

	return user, r.err == nil
}

func (r *Repo) GetUserByEmailAndAccountType(email string, accountType string) (User, bool) {
	var user User

	r.err = r.GetDb().Where("email = ? AND account_type = ?", email, accountType).First(&user).Error

	return user, r.err == nil
}

func (r *Repo) GetOrCreateUser(user *User) {
	var foundUser User
	// var err error

	r.err = r.GetDb().Where(
		"account_type = ? AND ext_id = ?",
		user.AccountType,
		user.ExtID,
	).First(&foundUser).Error

	if r.err == nil {
		*user = foundUser
	} else { // if errors.Is(err, gorm.ErrRecordNotFound) {
		user.UserID = user.Username + "@" + string(user.AccountType)
		r.err = r.GetDb().Create(&user).Error
	}
}

// From a slice of userIDs (<username>@<service>)
// fetchs the users.
// Returns the found users and a list of not found userIDs
func (r *Repo) FindUsers(userIDs []string) (map[string]User, []string) {
	users := make([]User, 0)
	userMap := make(map[string]User)
	notFounds := make([]string, 0)

	if r.err != nil {
		return userMap, notFounds
	}

	db := r.GetDb()

	r.err = db.Where("user_id IN ?", userIDs).Find(&users).Error

	if r.err != nil {
		return userMap, notFounds
	}

	for _, userID := range userIDs {
		found := false
		for _, user := range users {
			if user.UserID == userID {
				found = true

				userMap[userID] = user
				break
			}
		}

		if !found {
			notFounds = append(notFounds, userID)
		}
	}

	return userMap, notFounds
}
