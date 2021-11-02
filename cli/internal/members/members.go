package members

import (
	"io/ioutil"

	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
	"gopkg.in/yaml.v2"
)

func isMembersExist(c client.KeystoneClient, memberIDs []string) error {
	r, err := c.Users().CheckUsersExist(memberIDs)
	if err != nil {
		// The HTTP request must have failed
		return kserrors.UnkownError(err)
	}

	if r.Error != "" {
		return kserrors.UsersDontExist(r.Error, nil)
	}

	return nil
}

func GetMemberRolesFromFile(
	c client.KeystoneClient,
	filepath string,
	roles []models.Role,
) (map[string]models.Role, error) {
	var err error
	memberRoleNames := make(map[string]string)

	/* #nosec
	 * the file is going to be parsed, not executed in anyway
	 */
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dat, &memberRoleNames)
	if err != nil {
		return nil, err
	}

	memberIDs := make([]string, 0)
	for m := range memberRoleNames {
		memberIDs = append(memberIDs, m)
	}

	if err := isMembersExist(c, memberIDs); err != nil {
		return nil, err
	}

	display.WarningFreeOrga(roles)

	memberRoles, err := models.Roles(roles).
		MapWithMembersRoleNames(memberRoleNames)
	if err != nil {
		return nil, err
	}

	return memberRoles, nil
}

func GetMemberRolesFromArgs(
	c client.KeystoneClient,
	roleName string,
	memberIDs []string,
	roles []models.Role,
) (map[string]models.Role, error) {
	err := isMembersExist(c, memberIDs)
	if err != nil {
		return nil, err
	}

	foundRole := &models.Role{}

	display.WarningFreeOrga(roles)

	for _, role := range roles {
		if role.Name == roleName {
			*foundRole = role
		}
	}

	memberRoles := make(map[string]models.Role)

	for _, member := range memberIDs {
		memberRoles[member] = *foundRole
	}

	return memberRoles, nil
}

func GetMemberRolesFromPrompt(
	c client.KeystoneClient,
	memberIDs []string,
	roles []models.Role,
) (map[string]models.Role, error) {
	if err := isMembersExist(c, memberIDs); err != nil {
		return nil, err
	}

	display.WarningFreeOrga(roles)

	memberRole := make(map[string]models.Role)

	for _, memberId := range memberIDs {
		role, err := prompts.PromptRole(memberId, roles)
		if err != nil {
			return nil, err
		}

		memberRole[memberId] = role
	}

	return memberRole, nil
}

func SetMemberRole(
	c client.KeystoneClient,
	projectID, memberId, roleName string,
	roles []models.Role,
) error {
	var err error

	if _, ok := getRoleWithName(roleName, roles); ok {
		err = c.Project(projectID).SetMemberRole(memberId, roleName)
	} else {
		err = kserrors.RoleDoesNotExist(roleName, nil)
	}

	return err
}

// TODO: this should probably be inside the SetMemberRole function,
// or at least declared alongside it and called from there instead
// OR it – and most of the logic surrounding it - goes in a service internal
// package
func getRoleWithName(roleName string, roles []models.Role) (models.Role, bool) {
	found := false
	var role models.Role

	for _, existingRole := range roles {
		if existingRole.Name == roleName {
			found = true
			role = existingRole
			break
		}
	}

	return role, found
}
