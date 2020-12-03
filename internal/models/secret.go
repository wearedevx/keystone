package models

import (
	"encoding/json"
	"io"
	"strings"

	"gorm.io/gorm"
)

type SecretType string

const (
	SecretString SecretType = "string"
	SecretFile              = "file"
)

type Secret struct {
	gorm.Model
	Name string     `json:"name" gorm:"not null"`
	Type SecretType `json:"type" gorm:"not null"`
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
