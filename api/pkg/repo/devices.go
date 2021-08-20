package repo

import "github.com/wearedevx/keystone/api/pkg/models"

func (r *Repo) GetPublicKey(publicKey *models.PublicKey) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Where(&publicKey).
		First(publicKey).
		Error

	return r
}

func (r *Repo) GetPublicKeys(userID uint, publicKeys *[]models.PublicKey) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Where("user_id = ?", userID).
		Find(publicKeys).
		Error

	return r
}
