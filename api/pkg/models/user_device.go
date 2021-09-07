package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserDevice struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	User      User      `json:"user"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:user_device_user_id_device_id_key"`
	Device    Device    `json:"device"`
	DeviceID  uint      `json:"device_id" gorm:"uniqueIndex:user_device_user_id_device_id_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (pm *UserDevice) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *UserDevice) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *UserDevice) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}
