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
	ID             uint            `json:"id" gorm:"primaryKey"`
	UUID           string          `json:"uuid" gorm:"not null;unique"`
	Name           string          `json:"name" gorm:"not null"`
	Members        []ProjectMember `json:"members"`
	UserID         uint            `json:"user_id"`
	User           User            `json:"user"`
	Environments   []Environment   `json:"environments"`
	OrganizationID uint            `json:"organization_id"`
	CreatedAt      time.Time       `json:"create_at"`
	UpdatedAt      time.Time       `json:"update_at"`
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

func (p *Project) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(p)

	*out = sb.String()

	return err
}

// API Types

type AddVariablePayload struct {
	VarName      string `json:"var_name"`
	UserEnvValue []struct {
		UserID      string `json:"user_id"`
		Environment string `json:"environment"`
		Value       []byte `json:"value"`
	} `json:"user_env_value"`
}

func (avp *AddVariablePayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(avp)
}

func (avp *AddVariablePayload) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(avp)

	*out = sb.String()

	return err
}

type SetVariablePayload struct {
	VarName   string `json:"var_name"`
	UserValue []struct {
		UserID string `json:"user_id"`
		Value  []byte `json:"value"`
	} `json:"user_env_value"`
}

func (svp *SetVariablePayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(svp)
}

func (svp *SetVariablePayload) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(svp)

	*out = sb.String()

	return err
}

type DestroyProjectPayload struct {
	ProjectId string `json:"project_id"`
}

func (svp *DestroyProjectPayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(svp)
}

func (svp *DestroyProjectPayload) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(svp)

	*out = sb.String()

	return err
}
