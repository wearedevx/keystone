package repo

import "github.com/wearedevx/keystone/api/pkg/models"

func (r *Repo) CreateCheckoutSession(cs *models.CheckoutSession) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Create(cs).Error

	return r
}

func (r *Repo) GetCheckoutSession(sessionID string, cs *models.CheckoutSession) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Where("session_id = ?", sessionID).First(&cs).Error

	return r
}

func (r *Repo) UpdateCheckoutSession(cs *models.CheckoutSession) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Save(cs).Error

	return r
}

func (r *Repo) DeleteCheckoutSession(cs *models.CheckoutSession) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().Delete(cs).Error

	return r
}
