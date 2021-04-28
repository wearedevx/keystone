package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleReader UserRole = "read"
	RoleWriter          = "write"
	RoleOwner           = "owner"
)

type ProjectMember struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	Role          UserRole    `json:"role" gorm:"type:user_role"`
	ProjectOwner  bool        `json:"project_owner"`
	User          User        `json:"user"`
	UserID        uint        `json:"user_id"`
	Environment   Environment `json:"environment"`
	EnvironmentID uint        `json:"environment_id"`
	Project       Project     `json:"project"`
	ProjectID     uint        `json:"project_id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (pm *ProjectMember) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *ProjectMember) BeforeUpdate(tx *gorm.DB) (err error) {
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *ProjectMember) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *ProjectMember) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

// API Types
type MemberEnvironmentRole struct {
	ID          string
	Environment string
	Role        UserRole
}

type GetMembersResponse struct {
	Members []ProjectMember `json:"members"`
}

func (p *GetMembersResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(p)
}

func (p *GetMembersResponse) Serialize(out *string) error {
	var err error

	bout, err := json.Marshal(p)
	*out = string(bout)

	return err
}

type AddMembersPayload struct {
	Members []MemberEnvironmentRole
}

func (pm *AddMembersPayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *AddMembersPayload) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type AddMembersResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (pm *AddMembersResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *AddMembersResponse) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type RemoveMembersPayload struct {
	Members []string
}

func (pm *RemoveMembersPayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *RemoveMembersPayload) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type RemoveMembersResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (pm *RemoveMembersResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *RemoveMembersResponse) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}
