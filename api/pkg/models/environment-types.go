package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EnvironmentType struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *EnvironmentType) BeforeCreate(tx *gorm.DB) (err error) {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	return nil
}

func (e *EnvironmentType) BeforeUpdate(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now()

	return nil
}

func (e *EnvironmentType) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *EnvironmentType) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
