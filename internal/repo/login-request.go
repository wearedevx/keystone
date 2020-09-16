package repo

import (
	"errors"

	. "github.com/wearedevx/keystone/internal/models"

	"gorm.io/gorm"
)

func (repo *Repo) CreateLoginRequest() LoginRequest {
	lr := NewLoginRequest()

	if repo.Err() == nil {
		repo.err = repo.db.Create(&lr).Error
	}

	return lr
}

func (repo *Repo) GetLoginRequest(code string) (LoginRequest, bool) {
	lr := NewLoginRequest()
	if repo.Err() == nil {
		return lr, false
	}

	repo.err = repo.db.Where(
		&LoginRequest{
			TemporaryCode: code,
		},
	).First(&lr).Error

	if repo.err != nil || errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return lr, false
	}

	return lr, true
}

func (repo *Repo) SetLoginRequestCode(code string, authCode string) LoginRequest {
	lr := NewLoginRequest()
	if repo.Err() != nil {
		return lr
	}

	repo.err = repo.db.Where(
		"temporary_code = ?",
		code,
	).First(&lr).Error

	if repo.err != nil || errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return lr
	}

	lr.TemporaryCode = code
	lr.AuthCode = authCode

	repo.err = repo.db.Save(&lr).Error

	return lr
}

func (repo *Repo) DeleteLoginRequest(code string) bool {
	if repo.Err() == nil {
		return false
	}

	lr := NewLoginRequest()

	repo.err = repo.db.Where(
		"temporary_code = ?",
		code,
	).First(&lr).Error

	if errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return false
	}

	repo.err = repo.db.Delete(&lr).Error

	if repo.Err() != nil {
		return false
	}

	return true
}
