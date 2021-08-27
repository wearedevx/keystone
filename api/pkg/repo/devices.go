package repo

import (
	"errors"
	"regexp"
	"strings"

	"github.com/wearedevx/keystone/api/internal/emailer"
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
		Where("recipient_device_id = ?", publicKey.ID).
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

func (r *Repo) AddNewDevice(device models.Device, userID string, userEmail string) IRepo {

	var result = models.GetDevicesResponse{
		Devices: []models.Device{},
	}

	r.GetPublicKeys(device.UserID, &result.Devices)

	for _, existingDevice := range result.Devices {
		if existingDevice.Name == device.Name {
			r.err = errors.New("Device name already registered for this account")
			return r
		}

	}

	db.Create(&device)

	// Get project on which user is present
	rows, err := r.GetDb().Raw(`
	SELECT u.email, array_agg(p.name) FROM users u
	LEFT join project_members pm on pm.user_id = u.id
	LEFT join roles r on r.id = pm.role_id
	LEFT join projects p on pm.project_id = p.id
	where r.name = 'admin' and p.id in (
	select pm.project_id from project_members pm where pm.user_id = ?) and u.user_id != ?
	group by u.user_id, u.email;`, device.UserID, userID).Rows()

	if err != nil {
		r.err = err
		return r
	}

	var adminEmail string
	var projects string
	for rows.Next() {
		rows.Scan(&adminEmail, &projects)
		re := regexp.MustCompile(`\{(.+)?\}`)
		res := re.FindStringSubmatch(projects)

		projects_list := strings.Split(res[1], ",")

		// Send mail to admins of user projects
		e, err := emailer.NewDeviceAdminMail(userID, projects_list, device.Name)

		if err != nil {
			r.err = err
			return r
		}

		if err = e.Send([]string{adminEmail}); err != nil {
			r.err = err
			return r
		}

	}

	// Send mail to user
	e, err := emailer.NewDeviceMail(device.Name, userID)

	if err != nil {
		r.err = err
		return r
	}

	if err = e.Send([]string{userEmail}); err != nil {

		r.err = err
		return r
	}

	return r
}
