package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type RolesEnvironmentType struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	RoleID            uint      `json:"role_id"`
	Role              Role      `json:"role"`
	EnvironmentTypeID uint      `json:"environment_type_id"`
	EnvironmentType   uint      `json:"environment_type"`
	Name              string    `json:"name" gorm:"not null"`
	Read              bool      `json:"read"`
	Write             bool      `json:"write"`
	Invite            bool      `json:"invite"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (e *RolesEnvironmentType) BeforeCreate(tx *gorm.DB) (err error) {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	return nil
}

func (e *RolesEnvironmentType) BeforeUpdate(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now()

	return nil
}

func (e *RolesEnvironmentType) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *RolesEnvironmentType) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
