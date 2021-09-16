package models

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
	"unsafe"

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
	Message       string      `json:"error" gorm:"column=error"`
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

func (pm ActivityLog) Error() string {
	return pm.Message
}

func (pm *ActivityLog) SetError(err error) *ActivityLog {
	if err != nil {
		pm.Message = err.Error()
		pm.Success = false
	} else {
		pm.Success = true
	}

	return pm
}

func (pm *ActivityLog) Ptr() unsafe.Pointer {
	return (unsafe.Pointer)(pm)
}

func ErrorIsActivityLog(err interface{}) bool {
	activityLogType := fmt.Sprintf("%T", ActivityLog{})
	errType := fmt.Sprintf("%T", err)

	return activityLogType == errType
}
