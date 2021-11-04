package display

import (
	"sort"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

// MembersByRole function displays list of project members, grouped by role
func MembersByRole(members []models.ProjectMember) {
	grouped := models.ProjectMembers(members).GroupByRole()

	for _, role := range getSortedRoles(grouped) {
		members := grouped[role]
		printRole(role, members)
	}

}

// MembersAdded function Message when members arr added
func MembersAdded() {
	ui.Print(ui.RenderTemplate("added members", `
{{ OK }} {{ "Members Added" | green }}

To send secrets and files to new member, use "member add" command.
  $ ks member send-env --all-env <member-id>
`, nil))
}

// RemovedMembers function Message when members are removed
func RemovedMembers() {
	ui.Print(ui.RenderTemplate("removed members", `
{{ OK }} {{ "Revoked Access To Members" | green }}
`, nil))
}

// SetRoleOk function Message when member role is set
func SetRoleOk() {
	ui.Print(ui.RenderTemplate("set role ok", `
{{ OK }} {{ "Role set" | green }}
`, nil))
}

// —————————————————
// PRIVATE UTILITIES

func getSortedRoles(m map[models.Role]models.ProjectMembers) []models.Role {
	roles := make([]models.Role, 0)
	for r := range m {
		roles = append(roles, r)
	}

	s := models.NewRoleSorter(roles)
	return s.Sort()
}

func printRole(role models.Role, members []models.ProjectMember) {
	ui.Print("%s: %s", role.Name, role.Description)
	ui.Print("---")

	memberIDs := make([]string, len(members))

	for idx, member := range members {
		memberIDs[idx] = member.User.UserID
	}

	sort.Strings(memberIDs)

	for _, member := range memberIDs {
		ui.Print("%s", member)
	}

	ui.Print("")
}
