package display

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

func Organizations(organizations []models.Organization, user models.User) {
	ui.Print("Organizations your are in:")
	ui.Print("---")

	for _, orga := range organizations {
		orgaString := orga.Name
		if orga.User.UserID == user.UserID {
			orgaString += " 👑"
		}
		if orga.Private {
			orgaString += " P"
		}
		ui.Print(orgaString)
	}

	ui.Print("")
	ui.Print(" 👑 : You own; P : private")
}

func OrganizationCreated(organization models.Organization) {
	ui.PrintSuccess("Organization %s has been created", organization.Name)
}

func ManageUrl(url string) {
	ui.Print(
		ui.RenderTemplate(
			"upgrade-url",
			`To manage your organization plan, visit the following link:

        {{. }}`,
			url,
		),
	)
}

func UpgradeUrl(url string) {
	ui.Print(
		ui.RenderTemplate(
			"upgrade-url",
			`To upgrade your organization plan, visit the following link:

        {{. }}`,
			url,
		),
	)
}

func OrganizationMembers(members []models.ProjectMember) {
	ui.Print(
		"%d members are in projects that belong to this organization:\n",
		len(members),
	)

	for _, member := range members {
		fmt.Printf("  - %s\n", member.User.UserID)
	}
}

func OrganizationStatusUpdate(organization models.Organization) {
	if organization.Private {
		ui.PrintSuccess("Organization %s is now private", organization.Name)
	} else {
		ui.PrintSuccess("Organization %s is not private anymore", organization.Name)
	}
}

func OrganizationAccessibleProjects(projects []models.Project) {
	ui.Print(
		"You have access to %d project(s) in this organization :\n",
		len(projects),
	)

	for _, project := range projects {
		ui.Print("  - %s", project.Name)
	}
}

func OrganizationRenamed(from, to string) {
	ui.PrintSuccess("Organization %s has been renamed to %s", from, to)
}

func WarningFreeOrga(roles []models.Role) {
	if len(roles) == 1 {
		ui.PrintStdErr(
			`WARNING: You are not allowed to set role other than admin for free organization
To learn more: https://keystone.sh
`)

	}
}
