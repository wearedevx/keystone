package models

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID           uint      `json:"id"             gorm:"primaryKey"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ParentID     uint      `json:"parent_id"`
	Parent       *Role     `json:"parent"         gorm:"references:ID;foreignKey:ParentID"`
	CanAddMember bool      `json:"can_add_member"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *Role) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return nil
}

func (u *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()

	return nil
}

func (e *Role) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *Role) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

type Roles []Role

func (rs Roles) MapWithMembersRoleNames(
	memberRoleNames map[string]string,
) (map[string]Role, error) {
	memberRoles := make(map[string]Role)

	for member, roleName := range memberRoleNames {
		var foundRole *Role
		for _, role := range rs {
			if role.Name == roleName {
				foundRole = &role

				break
			}
		}

		if foundRole != nil {
			memberRoles[member] = *foundRole
		} else {
			return nil, fmt.Errorf("role %s does not exist", roleName)
		}
	}

	return memberRoles, nil
}

// API Types
type GetRolesResponse struct {
	Roles []Role
}

func (e *GetRolesResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *GetRolesResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

// Sort interface
type RoleSorter struct {
	roles []Role
}

func NewRoleSorter(roles []Role) *RoleSorter {
	return &RoleSorter{
		roles: roles,
	}
}

func (rs *RoleSorter) Sort() []Role {
	sort.Sort(rs)

	return rs.roles
}

func (rs *RoleSorter) Len() int {
	return len(rs.roles)
}

func (rs *RoleSorter) Swap(i, j int) {
	rs.roles[i], rs.roles[j] = rs.roles[j], rs.roles[i]
}

func (rs *RoleSorter) Less(i, j int) bool {
	return rs.roles[i].Name < rs.roles[j].Name
}
