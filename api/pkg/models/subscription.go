package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CheckoutSessionStatus string

const (
	CheckoutSessionStatusPending  CheckoutSessionStatus = "pending"
	CheckoutSessionStatusSuccess  CheckoutSessionStatus = "success"
	CheckoutSessionStatusCanceled CheckoutSessionStatus = "canceled"
)

type CheckoutSession struct {
	ID        uint                  `json:"id"         gorm:"primaryKey" faker:"-"`
	SessionID string                `json:"session_id" gorm:"unique"`
	Status    CheckoutSessionStatus `json:"status"     gorm:"default:pending" default:"pending" faker:"-"`
	CreatedAt time.Time             `json:"created_at" faker:"-"`
	UpdatedAt time.Time             `json:"updated_at" faker:"-"`
}

func (pm *CheckoutSession) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *CheckoutSession) BeforeUpdate(tx *gorm.DB) (err error) {
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *CheckoutSession) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *CheckoutSession) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type StartSubscriptionResponse struct {
	SessionID string `json:"session_id"`
	Url       string `json:"url"`
}

func (e *StartSubscriptionResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *StartSubscriptionResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

type ManageSubscriptionResponse struct {
	Url string `json:"url"`
}

func (e *ManageSubscriptionResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *ManageSubscriptionResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
