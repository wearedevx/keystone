package emailer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
)

// var templates = make(map[string]*template.Template)

const KEYSTONE_MAIL = "no-reply@keystone.sh"

type inviteData struct {
	Inviter     string
	ProjectName string
}

type newDeviceAdminData struct {
	Projects   string
	UserID     string
	DeviceName string
}

type newDeviceData struct {
	DeviceName string
	UserID     string
}

type GroupedMessageProject struct {
	Project      models.Project
	Environments map[string]models.Environment
}

type GroupedMessagesUser struct {
	Recipient models.User
	Projects  map[uint]GroupedMessageProject
}

type newGroupedMessageProjectData struct {
	NbDays          int
	GroupedProjects map[uint]GroupedMessageProject
}

func renderInviteTemplate(
	inviter string,
	projectName string,
) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["invite/html"].Execute(
		htmlBuffer,
		inviteData{Inviter: inviter, ProjectName: projectName},
	); err != nil {
		return "", "", err
	}

	if err = templates["invite/text"].Execute(
		textBuffer,
		inviteData{Inviter: inviter, ProjectName: projectName},
	); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), textBuffer.String(), nil
}

func renderAddedTemplate(
	inviter string,
	projectName string,
) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["added/html"].Execute(
		htmlBuffer,
		inviteData{Inviter: inviter, ProjectName: projectName},
	); err != nil {
		return "", "", err
	}

	if err = templates["added/text"].Execute(
		textBuffer,
		inviteData{Inviter: inviter, ProjectName: projectName},
	); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), textBuffer.String(), nil
}

func renderNewDeviceAdmin(userID string, projects []string, deviceName string) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["new_device_admin/html"].Execute(
		htmlBuffer,
		newDeviceAdminData{DeviceName: deviceName, Projects: strings.Join(projects, ", "), UserID: userID},
	); err != nil {
		return "", "", err
	}

	if err = templates["new_device_admin/text"].Execute(
		textBuffer,
		newDeviceAdminData{DeviceName: deviceName, Projects: strings.Join(projects, ", "), UserID: userID},
	); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), textBuffer.String(), nil
}

func renderNewDevice(deviceName string, userID string) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["new_device/html"].Execute(
		htmlBuffer,
		newDeviceData{DeviceName: deviceName, UserID: userID},
	); err != nil {
		return "", "", err
	}

	if err = templates["new_device/text"].Execute(
		textBuffer,
		newDeviceData{DeviceName: deviceName},
	); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), textBuffer.String(), nil
}

func renderMessageWillExpire(nbDays int, groupedProjects map[uint]GroupedMessageProject) (html string, text string, err error) {
	htmlBuffer := bytes.NewBufferString(html)
	textBuffer := bytes.NewBufferString(text)

	if err = templates["message_will_expire/html"].Execute(
		htmlBuffer,
		newGroupedMessageProjectData{NbDays: nbDays, GroupedProjects: groupedProjects},
	); err != nil {
		return "", "", err
	}

	if err = templates["message_will_expire/text"].Execute(
		textBuffer,
		newGroupedMessageProjectData{NbDays: nbDays, GroupedProjects: groupedProjects},
	); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), textBuffer.String(), nil
}

func InviteMail(inviter string, projectName string) (email *Email, err error) {
	html, text, err := renderInviteTemplate(inviter, projectName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     inviter,
		Subject:  "Your are invited to join a Keystone project",
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}

func AddedMail(inviter string, projectName string) (email *Email, err error) {
	html, text, err := renderInviteTemplate(inviter, projectName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     inviter,
		Subject:  "Your are added to join a Keystone project",
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}

func NewDeviceAdminMail(userID string, projects []string, deviceName string) (email *Email, err error) {
	html, text, err := renderNewDeviceAdmin(userID, projects, deviceName)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     KEYSTONE_MAIL,
		Subject:  fmt.Sprintf("%s has registered a new device", userID),
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}

func NewDeviceMail(deviceName string, userID string) (email *Email, err error) {
	html, text, err := renderNewDevice(deviceName, userID)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     KEYSTONE_MAIL,
		Subject:  fmt.Sprintf("A new device has been registered"),
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}

func MessageWillExpireMail(nbDays int, groupedProjects map[uint]GroupedMessageProject) (email *Email, err error) {
	html, text, err := renderMessageWillExpire(nbDays, groupedProjects)
	if err != nil {
		return nil, err
	}

	email = &Email{
		From:     KEYSTONE_MAIL,
		Subject:  "Some message will expire",
		HtmlBody: html,
		TextBody: text,
	}

	return email, nil
}
