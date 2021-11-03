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
	SecretFile   SecretType = "file"
)

type Secret struct {
	ID        uint       `gorm:"primaryKey" join:"id"`
	Name      string     `gorm:"not null"             json:"name"`
	Type      SecretType `gorm:"not null"             json:"type"`
	CreatedAt time.Time  `gorm:"not null"             json:"created_at"`
	UpdatedAt time.Time  `                            json:"updated_at"`
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

func (u *Secret) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
