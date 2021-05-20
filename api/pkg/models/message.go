package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	Payload       []byte        `json:"payload"`
	Sender        ProjectMember `json:"sender"`
	SenderID      uint          `json:"sender_id"`
	Recipient     ProjectMember `json:"recipient"`
	RecipientID   uint          `json:"recipient_id"`
	EnvironmentID string        `json:"environment_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (msg *Message) BeforeCreate(tx *gorm.DB) (err error) {
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()

	return nil
}

func (msg *Message) BeforeUpdate(tx *gorm.DB) (err error) {
	msg.UpdatedAt = time.Now()

	return nil
}

func (msg *Message) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(msg)
}

func (msg *Message) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(msg)

	*out = sb.String()

	return err
}

type File struct {
	Path  string `json:"path"`
	Value string `json:"content"`
}

type SecretVal struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type MessagePayload struct {
	Files   []File      `json:"files"`
	Secrets []SecretVal `json:"secrets"`
}
type MessageToWritePayload struct {
	Payload []byte `json:"payload"`
	// SenderID // Set by server
	UserID        string `json:"userid"`
	RecipientID   uint   `json:"recipient_id"`
	EnvironmentID string `json:"environment_id"`
}

type MessagesToWritePayload struct {
	Messages []MessageToWritePayload `json:"messages"`
}

func (e *MessagesToWritePayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *MessagesToWritePayload) Serialize(out *string) (err error) {
	var sb strings.Builder
	err = json.NewEncoder(&sb).Encode(u)
	*out = sb.String()
	return err
}

func (e *MessagePayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *MessagePayload) Serialize(out *string) (err error) {
	var sb strings.Builder
	err = json.NewEncoder(&sb).Encode(u)
	*out = sb.String()
	return err
}

type GetMessageByEnvironmentResponse struct {
	Environments map[string]GetMessageResponse
}

type GetMessageResponse struct {
	Message   Message `json:"message"`
	VersionID string  `json:"versionid"`
}

func (e *GetMessageByEnvironmentResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(e)
}

func (u *GetMessageByEnvironmentResponse) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(u)

	*out = sb.String()

	return err
}
