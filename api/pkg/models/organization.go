package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID             uint      `json:"id"              gorm:"primaryKey" faker:"-"`
	Name           string    `json:"name"            gorm:"unique"`
	Paid           bool      `json:"paid" faker:"-"`
	Private        bool      `json:"private" faker:"-"`
	CustomerID     string    `json:"customer_id" faker:""`
	SubscriptionID string    `json:"subscription_id" faker:""`
	UserID         uint      `json:"user_id"`
	User           User      `json:"user"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (pm *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *Organization) BeforeUpdate(tx *gorm.DB) (err error) {
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *Organization) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *Organization) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

type GetOrganizationsResponse struct {
	Organizations []Organization
}

func (e *GetOrganizationsResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *GetOrganizationsResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}

type GetOrganizationByNameResponse struct {
	Organization Organization
}

func (e *GetOrganizationByNameResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *GetOrganizationByNameResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
