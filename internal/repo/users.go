package repo

import (
	"fmt"

	. "github.com/wearedevx/keystone/internal/models"

	uuid "github.com/satori/go.uuid"
)

func (r *Repo) GetUser(userID string) (User, bool) {
	var user User

	r.err = r.db.Where("user_id = ?", userID).First(&user).Error

	return user, r.err == nil
}

func (r *Repo) GetOrCreateUser(user *User) {
	var foundUser User
	// var err error

	r.err = r.db.Where(
		"account_type = ? AND ext_id = ?",
		user.AccountType,
		user.ExtID,
	).First(&foundUser).Error

	if r.err == nil {
		*user = foundUser
	} else { // if errors.Is(err, gorm.ErrRecordNotFound) {

		fmt.Println("ON CREE LE USER")
		user.UserID = uuid.NewV4().String()
		r.err = r.db.Create(&user).Error
	}

	return
}
