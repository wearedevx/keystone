package models

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LoginRequest struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TemporaryCode string    `json:"temporary_code" gorm:"not null"`
	AuthCode      string    `json:"auth_code"`
	Answered      bool      `json:"answered" gorm:"default:false"`
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

func randomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func NewLoginRequest() (LoginRequest, error) {
	temporaryCode, err := randomString(16)
	if err != nil {
		return LoginRequest{}, err
	}

	return LoginRequest{
		TemporaryCode: temporaryCode,
		Answered:      false,
	}, nil
}

func (lr *LoginRequest) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(lr)
}

func (lr *LoginRequest) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(lr)

	*out = sb.String()

	return err
}

// The OAuth state value. Is base64 encoded JSON data structure
// which allows us to transit cli/api version information, enforcing
// version matching between the two
type AuthState struct {
	TemporaryCode string `json:"temporary_code"`
	Version       string `json:"version"`
}

func (state *AuthState) Decode(input string) (err error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(decoded), state)

	if err != nil {
		return err
	}

	return nil
}

func (state AuthState) Encode() (out string, err error) {
	outb, err := json.Marshal(state)

	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(outb)

	return encoded, nil
}
