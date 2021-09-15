package models

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ActivityLog struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        uint        `json:"user_id"`
	User          User        `json:"user"`
	ProjectID     uint        `json:"project_id"`
	Project       Project     `json:"project"`
	EnvironmentID uint        `json:"environment_id"`
	Environment   Environment `json:"environment"`
	Action        string      `json:"action"`
	Success       bool        `json:"success"`
	Error         string      `json:"error"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (pm *ActivityLog) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *ActivityLog) BeforeUpdate(tx *gorm.DB) (err error) {
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *ActivityLog) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *ActivityLog) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

func (pm *ActivityLog) FromString(s string) bool {
	parts := strings.Split(s, " ")
	ok := false

	for _, p := range parts {
		key, value, err := splitKeyValue(p)
		if err != nil {
			switch key {
			case "user":
				pm.UserID = mustParseUint(value)
				ok = true
			case "project":
				pm.ProjectID = mustParseUint(value)
				ok = true
			case "environment":
				pm.EnvironmentID = mustParseUint(value)
				ok = true
			case "action":
				value = strings.TrimPrefix(value, "\"")
				value = strings.TrimSuffix(value, "\"")
				pm.Action = value
				ok = true
			case "success":
				pm.Success = value == "true"
				ok = true
			case "error":
				value = strings.TrimPrefix(value, "\"")
				value = strings.TrimSuffix(value, "\"")
				pm.Error = value
				ok = true
			}
		}

	}

	return ok
}

func mustParseUint(value string) uint {
	uid, err := strconv.ParseUint(value, 10, 64)

	if err != nil {
		panic(err)
	}

	return uint(uid)
}

func splitKeyValue(s string) (key, value string, err error) {
	r := strings.Split(s, "=")
	if len(r) != 2 {
		return "", "", errors.New("invalid key value pair")
	}

	return r[0], r[1], nil
}
