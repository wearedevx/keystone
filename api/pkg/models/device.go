package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID uint `json:"id" gorm:"primaryKey"`
	// UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_public_keys_user_id"`
	PublicKey  []byte    `json:"public_key" gorm:"type:bytea"`
	Name       string    `json:"name"`
	UID        string    `json:"uid"`
	Users      []User    `json:"users" gorm:"many2many:user_devices;"`
	LastUsedAt time.Time `json:last_used_at`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt
}

func (pm *Device) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()
	pm.LastUsedAt = time.Now()

	return nil
}

func (pm *Device) BeforeUpdate(tx *gorm.DB) (err error) {
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *Device) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *Device) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type GetDevicesResponse struct {
	Devices []Device
}

func (e *GetDevicesResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *GetDevicesResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

type RemoveDeviceResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (pm *RemoveDeviceResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *RemoveDeviceResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}
