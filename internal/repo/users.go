package repo

import (
	"errors"
	"fmt"

	. "github.com/wearedevx/keystone/internal/models"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (r *Repo) GetUser(userID string) (User, bool) {
	var user User

	r.err = r.db.Where("user_id = ?", userID).First(&user).Error

	return user, r.err == nil
}

func (r *Repo) GetOrCreateUser(user *User) {
	var foundUser User
	var err error

	r.err = r.db.Where(
		"account_type = ? AND ext_id = ?",
		user.AccountType,
		user.ExtID,
	).First(&foundUser).Error

	if r.err == nil {
		*user = foundUser
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		user.UserID = uuid.NewV4().String()
		r.err = r.db.Create(&user).Error
		fmt.Println(r.err)
	}

	return
}
