package models

import (
	"encoding/json"
	"io"
	"math/rand"
	"strings"
	"time"
	"unsafe"

	"gorm.io/gorm"
)

type LoginRequest struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TemporaryCode string    `json:"temporary_code" gorm:"not null"`
	AuthCode      string    `json:"auth_code"`
	Answered      bool      `json:"answeredl" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (lr *LoginRequest) BeforeCreate(tx *gorm.DB) (err error) {
	lr.CreatedAt = time.Now()
	lr.UpdatedAt = time.Now()

	return nil
}

func (lr *LoginRequest) BeforeUpdate(tx *gorm.DB) (err error) {
	lr.UpdatedAt = time.Now()

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randomString(n int) string {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func NewLoginRequest() LoginRequest {
	return LoginRequest{
		TemporaryCode: randomString(16),
		Answered:      false,
	}
}

func (lr *LoginRequest) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(lr)
}

func (lr *LoginRequest) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(lr)

	*out = sb.String()

	return err
}
