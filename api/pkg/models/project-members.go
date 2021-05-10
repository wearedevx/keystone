package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProjectMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	User      User      `json:"user"`
	UserID    uint      `json:"user_id"`
	Project   Project   `json:"project"`
	ProjectID uint      `json:"project_id"`
	Role      Role      `json:"role"`
	RoleID    uint      `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
type MemberRole struct {
	MemberID string
	RoleID   uint
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
	Members []MemberRole
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

type CheckMembersResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (pm *CheckMembersResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *CheckMembersResponse) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type CheckMembersPayload struct {
	MemberIDs []string
}

func (pm *CheckMembersPayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *CheckMembersPayload) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type SetMemberRolePayload struct {
	MemberID string
	RoleName string
}

func (pm *SetMemberRolePayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *SetMemberRolePayload) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}
