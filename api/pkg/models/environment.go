package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Environment struct {
	ID                uint            `json:"id"                  gorm:"primaryKey"`
	Name              string          `json:"name"                gorm:"not null"`
	EnvironmentTypeID uint            `json:"environment_type_id"`
	EnvironmentType   EnvironmentType `json:"environment_type"`
	ProjectID         uint            `json:"project_id"`
	Project           Project         `json:"project"`
	VersionID         string          `json:"version_id"`
	EnvironmentID     string          `json:"environment_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Environment) BeforeCreate(tx *gorm.DB) (err error) {
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	e.EnvironmentID = uuid.NewV4().String()
	e.VersionID = uuid.NewV4().String()

	return nil
}

func (e *Environment) BeforeUpdate(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now()

	return nil
}

func (e *Environment) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *Environment) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

type GetEnvironmentsResponse struct {
	Environments []Environment `json:"environments"`
}

func (e *GetEnvironmentsResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *GetEnvironmentsResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

type EnvironmentUserSecret struct {
	EnvironmentID uint      `json:"environmentID" gorm:"primaryKey"`
	UserID        uint      `json:"userID"        gorm:"primaryKey"`
	SecretID      uint      `json:"secretID"      gorm:"primaryKey"`
	Value         []byte    `json:"value"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
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

func (pes *EnvironmentUserSecret) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pes)

	*out = sb.String()

	return err
}
