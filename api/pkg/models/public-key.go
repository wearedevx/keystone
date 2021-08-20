package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PublicKey struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_public_keys_user_id"`
	Key       []byte    `json:"key"`
	Device    string    `json:"device"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (pm *PublicKey) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *PublicKey) BeforeUpdate(tx *gorm.DB) (err error) {
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *PublicKey) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *PublicKey) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type GetDevicesResponse struct {
	PublicKeys []PublicKey
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