package notification

import (
	"fmt"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// Send email to amdins and users for every new devices created
func SendEmailForNewDevices(r repo.IRepo) error {
	var newDevices []models.Device
	if err := r.GetNewlyCreatedDevices(&newDevices).Err(); err != nil {
		return err
	}

	for _, device := range newDevices {
		// Newly created devices only have one user
		user := device.Users[0]

		var adminProjectsMap map[string][]string
		if err := r.GetAdminsFromUserProjects(user.ID, &adminProjectsMap).Err(); err != nil {
			return err
		}
		for adminEmail, projectList := range adminProjectsMap {
			// Send mail to admins of user projects
			e, err := emailer.NewDeviceAdminMail(user.Username, projectList, device.Name)
			if err != nil {
				return err
			}

			if err = e.Send([]string{adminEmail}); err != nil {
				fmt.Printf("Add New Device Admin Mail err: %+v\n", err)
				return err
			}
		}

		// Send mail to user
		e, err := emailer.NewDeviceMail(device.Name, user.Username)
		if err != nil {
			return err
		}

		if user.Email != "" {
			if err = e.Send([]string{user.Email}); err != nil {
				fmt.Printf("Add New Device User Mail err: %+v\n", err)
				return err
			}
		}
		if err = r.GetDb().Model(&models.UserDevice{}).Where("user_id = ? and device_id = ?", user.ID, device.ID).Update("newly_created", false).Error; err != nil {
			return err
		}
	}
	return nil

}
