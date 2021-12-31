package notification

import (
	"fmt"

	apierrors "github.com/wearedevx/keystone/api/internal/errors"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/models"
)

// Send email to amdins and users for every new devices created
func SendEmailForNewDevices(device models.Device, adminProjectsMap map[string][]string, user models.User) error {
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
	return nil

}
func SendExpireMessageToUsers(
	groupedMessageUser map[uint]emailer.GroupedMessagesUser,
	errors *[]error,
) {
	// For each recipients, send message.
	for _, groupedMessagesUser := range groupedMessageUser {
		email, err := emailer.MessageWillExpireMail(
			5,
			groupedMessagesUser.Projects,
		)
		if err != nil {
			*errors = append(*errors, err)
		} else if err = email.Send([]string{groupedMessagesUser.Recipient.Email}); err != nil {
			fmt.Printf("Message will expire mail err: %+v\n", err)
			*errors = append(*errors, err)
		}
	}
}

func SendInvitationEmail(user models.User, payload models.InvitePayload) (err error) {
	var email *emailer.Email

	email, err = emailer.InviteMail(user, payload.ProjectName)
	if err != nil {
		return apierrors.ErrorFailedToCreateMailContent(err)
	}

	if err = email.Send([]string{payload.Email}); err != nil {
		fmt.Printf("Invite Mail err: %+v\n", err)
		return apierrors.ErrorFailedToSendMail(err)
	}
	return nil

}

func SendAddedMemberEmail(memberRoles []models.MemberRole, project models.Project, currentUser models.User, users map[string]models.User) error {
	for _, memberRole := range memberRoles {
		userEmail := users[memberRole.MemberID].Email
		e, err := emailer.AddedMail(currentUser, project.Name)
		if err != nil {
			return err
		}

		if err = e.Send([]string{userEmail}); err != nil {
			fmt.Printf("Project Add Member err: %+v\n", err)
			return err
		}
	}
	return nil
}
