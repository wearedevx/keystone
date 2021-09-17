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
	UserID        *uint       `json:"user_id"`
	User          User        `json:"user" gorm:"foreignKey:user_id"`
	ProjectID     *uint       `json:"project_id"`
	Project       Project     `json:"project" gorm:"foreignKey:project_id"`
	EnvironmentID *uint       `json:"environment_id"`
	Environment   Environment `json:"environment" gorm:"foreignKey:environment_id"`
	Action        string      `json:"action"`
	Success       bool        `json:"success"`
	Message       string      `json:"error"`
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

func ErrorIsActivityLog(err error) bool {
	activityLogPtrType := fmt.Sprintf("%T", &ActivityLog{})
	errType := fmt.Sprintf("%T", err)
	fmt.Printf("activityLogPtrType: %+v\n", activityLogPtrType)
	fmt.Printf("errType: %+v\n", errType)

	return activityLogPtrType == errType
}
