package display

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

// Organizations function displays a list of organizations
// withe ownership information
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

// OrganizationCreated function Message when organization is created
func OrganizationCreated(organization models.Organization) {
	ui.PrintSuccess("Organization %s has been created", organization.Name)
}

// ManageUrl function displays the link the user must follow to manage their
// subscritption
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

// UpgradeUrl function displays a link the user must follow to upgrade their
// organization
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

// OrganizationMembers function displays a list of orgnaiztion members
func OrganizationMembers(members []models.ProjectMember) {
	ui.Print(
		"%d members are in projects that belong to this organization:\n",
		len(members),
	)

	for _, member := range members {
		fmt.Printf("  - %s\n", member.User.UserID)
	}
}

// OrganizationStatusUpdate function displays the organizaiton private status
func OrganizationStatusUpdate(organization models.Organization) {
	if organization.Private {
		ui.PrintSuccess("Organization %s is now private", organization.Name)
	} else {
		ui.PrintSuccess("Organization %s is not private anymore", organization.Name)
	}
}

// OrganizationAccessibleProjects function displays a list of projects
// in an organization the user has access to
func OrganizationAccessibleProjects(projects []models.Project) {
	ui.Print(
		"You have access to %d project(s) in this organization :\n",
		len(projects),
	)

	for _, project := range projects {
		ui.Print("  - %s", project.Name)
	}
}

// OrganizationRenamed function Message when the organization was renamed
func OrganizationRenamed(from, to string) {
	ui.PrintSuccess("Organization %s has been renamed to %s", from, to)
}

// WarningFreeOrga function Message telling some things are not possible
// for free organizaitons
func WarningFreeOrga(roles []models.Role) {
	if len(roles) == 1 {
		ui.PrintStdErr(
			`WARNING: You are not allowed to set role other than admin for free organization
To learn more: https://keystone.sh
`)

	}
}
