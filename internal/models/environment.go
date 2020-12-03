package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Environment struct {
	gorm.Model
	Name      string   `json:"name" gorm:"not null"`
	ProjectID int      `json:"project_id" gorm:"not null;index"`
	Secrets   []Secret `json:"secrets" gorm:"many2many:project_environment_secrets;ForeignKey:ID;References:ID;"`
	Files     []File   `json:"files" gorm:"many2many:project_environment_files;ForeignKey:ID;References:ID;"`
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

type EnvironmentUserSecret struct {
	EnvironmentID int            `json:"environmentID" gorm:"primaryKey"`
	UserID        int            `json:"userID" gorm:"primaryKey"`
	SecretID      int            `json:"secretID" gorm:"primaryKey"`
	Value         string         `json:"value"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (eus *EnvironmentUserSecret) BeforeCreate(tx *gorm.DB) (err error) {
	eus.CreatedAt = time.Now()
	eus.UpdatedAt = time.Now()

	return nil
}

func (eus *EnvironmentUserSecret) BeforeUpdate(tx *gorm.DB) (err error) {
	eus.UpdatedAt = time.Now()

	return nil
}

func (eus *EnvironmentUserSecret) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(eus)
}

func (pes *EnvironmentUserSecret) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pes)

	*out = sb.String()

	return err
}
