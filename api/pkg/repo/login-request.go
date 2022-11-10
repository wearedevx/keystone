package repo

import (
	"errors"

	"gorm.io/gorm"

	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateLoginRequest() models.LoginRequest {
	lr, err := models.NewLoginRequest()
	if err != nil {
		repo.err = err
	}

	if repo.Err() == nil {
		repo.err = repo.GetDB().Create(&lr).Error
	}

	return lr
}

func (repo *Repo) GetLoginRequest(code string) (models.LoginRequest, bool) {
	lr := models.LoginRequest{}

	if repo.Err() == nil {
		repo.err = repo.
			GetDB().
			Where("temporary_code = ?", code).
			First(&lr).
			Error
	}

	return lr, repo.err == nil
}

func (repo *Repo) SetLoginRequestCode(
	code string,
	authCode string,
) models.LoginRequest {
	lr, err := models.NewLoginRequest()
	if err != nil {
		repo.err = err
	}

	if repo.Err() != nil {
		return lr
	}

	repo.err = repo.GetDB().Where(
		"temporary_code = ?",
		code,
	).First(&lr).Error

	if repo.err != nil || errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return lr
	}

	lr.TemporaryCode = code
	lr.AuthCode = authCode

	repo.err = repo.GetDB().Save(&lr).Error

	return lr
}

func (repo *Repo) DeleteLoginRequest(code string) bool {
	if repo.Err() == nil {
		return false
	}

	lr, err := models.NewLoginRequest()
	if err != nil {
		repo.err = err
		return false
	}

	repo.err = repo.GetDB().Where(
		"temporary_code = ?",
		code,
	).First(&lr).Error

	if errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return false
	}

	repo.err = repo.GetDB().Delete(&lr).Error

	return repo.Err() == nil
}
