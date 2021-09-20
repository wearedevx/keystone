package repo

import (
	"errors"
	"regexp"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func (r *Repo) GetDevice(device *models.Device) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Where(&device).
		First(device).
		Error

	return r
}

func (r *Repo) GetDeviceByUserID(userID uint, device *models.Device) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Joins("left join user_devices on user_devices.device_id = devices.id").
		Joins("left join users on users.id = user_devices.user_id").
		Where(&device).
		Where("users.id = ?", userID).
		First(device).
		Error

	return r
}

func (r *Repo) GetDevices(userID uint, devices *[]models.Device) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Joins("left join user_devices on user_devices.device_id = devices.id").
		Joins("left join users on users.id = user_devices.user_id").
		Where("users.id = ?", userID).
		Find(devices).
		Error

	return r
}

func (r *Repo) RevokeDevice(userID uint, deviceName string) IRepo {
	if r.Err() != nil {
		return r
	}

	device := models.Device{}
	r.err = r.GetDb().
		Joins("left join user_devices on user_devices.device_id = devices.id").
		Joins("left join users on users.id = user_devices.user_id").
		Where("users.id = ? and devices.uid = ?", userID, deviceName).
		Find(&device).
		Error

	if r.err != nil {
		return r
	}

	if device.ID == 0 {
		r.err = errors.New("No device")
		return r
	}

	// Remove message aimed for the device
	r.err = r.GetDb().
		Model(&models.Message{}).
		Where("recipient_device_id = ? and sender_device_id = ?", device.ID, device.ID).
		Delete(models.Message{}).Error

	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Delete(&device).Error

	if r.err != nil {
		return r
	}

	return r
}

func (r *Repo) AddNewDevice(device models.Device, userID uint, userName string, userEmail string) IRepo {
	db := r.GetDb()

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\.\-\_]{1,}$`, device.Name)
	if !matched {
		r.err = errors.New("Incorrect device name. Device name must be alphanumeric with ., -, _")
		return r
	}

	if err := db.Where("uid = ?", device.UID).First(&device).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			r.err = db.Create(&device).Error
		} else {
			r.err = err
		}
	}

	if r.err != nil {
		return r
	}

	userDevice := models.UserDevice{UserID: userID, DeviceID: device.ID}

	err := db.SetupJoinTable(&models.User{}, "Devices", &models.UserDevice{})

	if err != nil {
		r.err = err
		return r
	}

	r.err = db.Create(&userDevice).Error

	if r.err != nil {
		return r
	}

	var projects_list []string
	var adminEmail string
	r.GetAdminsFromUserProjects(userID, userName, projects_list, &adminEmail)

	if r.err != nil {
		return r
	}

	// Send mail to admins of user projects
	e, err := emailer.NewDeviceAdminMail(userName, projects_list, device.Name)

	if err != nil {
		r.err = err
		return r
	}

	if err = e.Send([]string{adminEmail}); err != nil {
		r.err = err
		return r
	}

	// Send mail to user
	e, err = emailer.NewDeviceMail(device.Name, userName)

	if err != nil {
		r.err = err
		return r
	}

	if userEmail != "" {
		if err = e.Send([]string{userEmail}); err != nil {
			r.err = err
			return r
		}
	}

	return r
}
