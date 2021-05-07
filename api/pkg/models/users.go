package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type KeyRing struct {
	Cipher string `json:"cipher" gorm:"column:keys_cipher"`
	Sign   string `json:"sign" gorm:"column:keys_sign"`
}

type User struct {
	ID          uint        `json:"id" gorm:"primaryKey" faker:"-"`
	AccountType AccountType `json:"account_type" gorm:"default:custom" faker:"oneof: github, gitlab"`
	UserID      string      `json:"user_id" gorm:"uniqueIndex"`
	ExtID       string      `json:"ext_id" faker:"uuid_digit"`
	Username    string      `json:"username" faker:"username"`
	Fullname    string      `json:"fullname" gorm:"not null" faker:"name"`
	Email       string      `json:"email" gorm:"not null" faker:"email"`
	PublicKey   []byte      `json:"public_key"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()

	return nil
}

func (u *User) FromId(id string) {
	parts := strings.Split(id, "@")
	username := parts[0]
	accountType := AccountType(parts[1])

	u.Username = username
	u.AccountType = accountType
}

type LoginPayload struct {
	AccountType AccountType
	Token       *oauth2.Token
	PublicKey   []byte
}

type UserPublicKey struct {
	UserID    string `json:"user_id"`
	PublicKey []byte `json:"publick_key"`
}

type AccountType string

const (
	GitHubAccountType AccountType = "github"
	GitLabAccountType             = "gitlab"
	CustomAccountType             = "custom"
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

func (upk *UserPublicKey) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(upk)
}

func (upk UserPublicKey) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(upk)

	*out = sb.String()

	return err
}
