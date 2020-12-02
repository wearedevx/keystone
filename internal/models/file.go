package models

import (
	"encoding/json"
	"io"
	"strings"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Path string `json:"path" gorm:"not null"`
}

func (u *File) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(u)
}

func (u *File) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
