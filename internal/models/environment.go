package models

import (
	"encoding/json"
	"io"
	"strings"

	"gorm.io/gorm"
)

type Environment struct {
	gorm.Model
	ProjectID int
	Name      string   `json:"name" gorm:"not null"`
	Secrets   []Secret `json:"secrets" gorm:"many2many:project_environment_secrets;ForeignKey:id,project_id;References:secret_id;"`
	Files     []File   `json:"files" gorm:"many2many:project_environment_files;ForeignKey:id,project_id;References:file_id;"`
}

func (u *Environment) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(u)
}

func (u *Environment) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
