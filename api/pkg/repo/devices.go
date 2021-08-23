package repo

import (
	"errors"

	"github.com/wearedevx/keystone/api/pkg/models"
)

func (r *Repo) GetPublicKey(publicKey *models.Device) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Where(&publicKey).
		First(publicKey).
		Error

	return r
}

func (r *Repo) GetPublicKeys(userID uint, publicKeys *[]models.Device) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Where("user_id = ?", userID).
		Find(publicKeys).
		Error

	return r
}

func (r *Repo) RevokeDevice(userID uint, deviceName string) IRepo {
	if r.Err() != nil {
		return r
	}

	publicKey := models.Device{}
	r.err = r.GetDb().
		Where("user_id = ? and name = ?", userID, deviceName).
		Find(&publicKey).
		Error
	if r.err != nil {
		return r
	}

	if publicKey.ID == 0 {
		r.err = errors.New("No device found with this name")
		return r
	}

	// Remove message aimed for the device
	r.err = r.GetDb().
		Model(&models.Message{}).
		Where("public_key_id = ?", publicKey.ID).
		Delete(models.Message{}).Error

	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Model(&models.Device{}).
		Where("id = ?", publicKey.ID).
		Delete(models.Device{}).Error

	if r.err != nil {
		return r
	}

	return r
}
