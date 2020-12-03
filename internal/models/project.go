package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	UUID         string        `json:"uuid" gorm:"not null;unique"`
	Name         string        `json:"name" gorm:"not null"`
	Users        []User        `json:"users" gorm:"many2many:project_permissions;"`
	Environments []Environment `json:"environments"`
	CreatedAt    time.Time     `json:"create_at"`
	UpdatedAt    time.Time     `json:"update_at"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	p.UUID = uuid.NewV4().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()

	return nil
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
