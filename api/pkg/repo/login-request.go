package repo

import (
	"errors"

	"gorm.io/gorm"

	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateLoginRequest() LoginRequest {
	lr := NewLoginRequest()

	if repo.Err() == nil {
		repo.err = repo.GetDb().Create(&lr).Error
	}

	return lr
}

func (repo *Repo) GetLoginRequest(code string) (LoginRequest, bool) {
	lr := LoginRequest{}

	if repo.Err() == nil {
		repo.err = repo.GetDb().Where(
			&LoginRequest{
				TemporaryCode: code,
			},
		).First(&lr).Error
	}

	return lr, repo.err == nil
}

func (repo *Repo) SetLoginRequestCode(code string, authCode string) LoginRequest {
	lr := NewLoginRequest()
	if repo.Err() != nil {
		return lr
	}

	repo.err = repo.GetDb().Where(
		"temporary_code = ?",
		code,
	).First(&lr).Error

	if repo.err != nil || errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return lr
	}

	lr.TemporaryCode = code
	lr.AuthCode = authCode

	repo.err = repo.GetDb().Save(&lr).Error

	return lr
}

func (repo *Repo) DeleteLoginRequest(code string) bool {
	if repo.Err() == nil {
		return false
	}

	lr := NewLoginRequest()

	repo.err = repo.GetDb().Where(
		"temporary_code = ?",
		code,
	).First(&lr).Error

	if errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return false
	}

	repo.err = repo.GetDb().Delete(&lr).Error

    return repo.Err() == nil
}
