package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SecretType string

const (
	SecretString SecretType = "string"
	SecretFile              = "file"
)

type Secret struct {
	ID        uint       `join:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"not null"`
	Type      SecretType `json:"type" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *Secret) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	return nil
}

func (s *Secret) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()

	return nil
}

func (u *Secret) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(u)
}

func (u *Secret) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
