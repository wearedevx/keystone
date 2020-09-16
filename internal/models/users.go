package models

import (
	"encoding/json"
	"io"
	"strings"

	"gorm.io/gorm"
)

type KeyRing struct {
	Cipher string `json:"cipher" gorm:"not null"`
	Sign   string `json:"sign" gorm:"not null"`
}

type User struct {
	gorm.Model
	AccountType AccountType `json:"account_type" gorm:"default:custom"`
	UserID      string      `json:"user_id" gorm:"uniqueIndex"`
	ExtID       string      `json:"ext_id" gorm:"not null"`
	Username    string      `json:"username" gorm:"uniqueIndex"`
	Fullname    string      `json:"fullname" gorm:"not null"`
	Email       string      `json:"email" gorm:"not null"`
	Keys        KeyRing     `json:"keys" gorm:"embedded"`
}

type AccountType string

const (
	GitHub AccountType = "github"
	GitLab             = "gitlab"
	Custom             = "custom"
)

func (u *User) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(u)
}

func (u *User) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
