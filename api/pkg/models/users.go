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

type InvitePayload struct {
	Email       string
	ProjectName string
}

func (pm *InvitePayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *InvitePayload) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type User struct {
	ID          uint        `json:"id" gorm:"primaryKey" faker:"-"`
	AccountType AccountType `json:"account_type" gorm:"default:custom" faker:"oneof: github, gitlab"`
	UserID      string      `json:"user_id" gorm:"uniqueIndex"`
	ExtID       string      `json:"ext_id" faker:"uuid_digit"`
	Username    string      `json:"username" faker:"username"`
	Fullname    string      `json:"fullname" gorm:"not null" faker:"name"`
	Email       string      `json:"email" gorm:"not null" faker:"email"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Devices     []Device    `json:"devices" faker:"-"`
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
	Device      string
	DeviceUID   string
}

type UserPublicKeys struct {
	UserID     uint     `json:"user_id"`
	UserUID    string   `json:"user_uid"` // UserID as string (e.g: toto@github)
	PublicKeys []Device `json:"publick_keys"`
}

type PublicKeys struct {
	Keys []UserPublicKeys `json:"keys"`
}

func (p *PublicKeys) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(p)
}

func (p *PublicKeys) Serialize(out *string) error {
	var sb strings.Builder
	err := json.NewEncoder(&sb).Encode(p)
	*out = sb.String()
	return err
}

type AccountType string

const (
	GitHubAccountType AccountType = "github"
	GitlabAccountType AccountType = "gitlab"
	CustomAccountType AccountType = "custom"
)

func (u *User) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(u)
}

func (u *User) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

func (upk *UserPublicKeys) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(upk)
}

func (upk UserPublicKeys) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(upk)

	*out = sb.String()

	return err
}

type GetInviteResponse struct {
	UserUIDs []string `json:"user_uids"`
}

func (e *GetInviteResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *GetInviteResponse) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
