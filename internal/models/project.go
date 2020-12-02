package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name         string        `json:"name" gorm:"not null"`
	Users        []User        `json:"users" gorm:"many2many:project_permissions;"`
	Environments []Environment `json:"environments" gorm:"many2many:project_enviroment_secrets;"`
	Secrets      []Secret      `json:"secrets" gorm:"many2many:project_enviroment_secrets;"`
	Files        []File        `json:"files" gorm:"many2many:project_enviroment_files;"`
}

func (p *Project) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(p)
}

func (p *Project) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(p)

	*out = sb.String()

	return err
}

type ProjectEnvironmentSecret struct {
	ProjectID     int            `json:"projectID" gorm:"primaryKey"`
	EnvironmentID int            `json:"environmentID" gorm:"primaryKey"`
	SecretID      int            `json:"secretID" gorm:"primaryKey"`
	Value         string         `json:"value"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (pes *ProjectEnvironmentSecret) BeforeCreate(tx *gorm.DB) (err error) {
	pes.CreatedAt = time.Now()
	pes.UpdatedAt = time.Now()

	return nil
}

func (pes *ProjectEnvironmentSecret) BeforeUpdate(tx *gorm.DB) (err error) {
	pes.UpdatedAt = time.Now()

	return nil
}

func (pes *ProjectEnvironmentSecret) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pes)
}

func (pes *ProjectEnvironmentSecret) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pes)

	*out = sb.String()

	return err
}

type ProjectEnvironmentFile struct {
	ProjectID     int            `json:"projectID" gorm:"primaryKey"`
	EnvironmentID int            `json:"environmentID" gorm:"primaryKey"`
	FileID        int            `json:"secretID" gorm:"primaryKey"`
	Content       []byte         `json:"content"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (pef *ProjectEnvironmentFile) BeforeCreate(tx *gorm.DB) (err error) {
	pef.CreatedAt = time.Now()
	pef.UpdatedAt = time.Now()

	return nil
}

func (pef *ProjectEnvironmentFile) BeforeUpdate(tx *gorm.DB) (err error) {
	pef.UpdatedAt = time.Now()

	return nil
}

func (pef *ProjectEnvironmentFile) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pef)
}

func (pef *ProjectEnvironmentFile) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pef)

	*out = sb.String()

	return err
}
