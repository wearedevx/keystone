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
	ID           uint          `json:"id" gorm:"primaryKey"`
	AccountType  AccountType   `json:"account_type" gorm:"default:custom"`
	UserID       string        `json:"user_id" gorm:"uniqueIndex"`
	ExtID        string        `json:"ext_id"`
	Username     string        `json:"username" gorm:"uniqueIndex"`
	Fullname     string        `json:"fullname" gorm:"not null"`
	Email        string        `json:"email" gorm:"not null"`
	Keys         KeyRing       `json:"keys" gorm:"embedded"`
	Projects     []Project     `json:"projects" gorm:"many2many:project_permissions;"`
	Environments []Environment `json:"environment" gorm:"many2many:environment_permissions;"`
	Secrets      []Secret      `json:"secrets" gorm:"many2many:environment_user_secrets;"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
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

type LoginPayload struct {
	AccountType AccountType
	Token       *oauth2.Token
	PublicKey   string
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

type UserRole string

const (
	RoleAdmin       UserRole = "admin"
	RoleContributor          = "contributor"
	RoleReader               = "reader"
)

type ProjectPermissions struct {
	UserID    uint           `json:"userID" gorm:"primaryKey"`
	ProjectID uint           `json:"projectID" gorm:"primaryKey"`
	Role      UserRole       `json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (perm *ProjectPermissions) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(perm)
}

func (perm *ProjectPermissions) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(perm)

	*out = sb.String()

	return err
}

type EnvironmentPermissions struct {
	UserID        uint           `json:"userID" gorm:"primaryKey"`
	EnvironmentID uint           `json:"environmentID" gorm:"primaryKey"`
	Role          UserRole       `json:"role"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (perm *EnvironmentPermissions) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(perm)
}

func (perm *EnvironmentPermissions) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(perm)

	*out = sb.String()

	return err
}
