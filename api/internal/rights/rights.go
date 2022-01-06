package rights

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type UserRight string

const (
	Read   UserRight = "read"
	Write  UserRight = "write"
	Invite UserRight = "invite"
)

func doesUserHaveRightsOnEnvironment(
	Repo repo.IRepo,
	userID uint,
	projectID uint,
	environment *models.Environment,
	right string,
) (bool, error) {
	projectMember := models.ProjectMember{
		UserID:    userID,
		ProjectID: projectID,
	}

	err := Repo.GetProjectMember(&projectMember).Err()
	if err != nil {
		return false, err
	}

	rolesEnvironmentType := models.RolesEnvironmentType{
		EnvironmentTypeID: environment.EnvironmentTypeID,
		RoleID:            projectMember.RoleID,
	}

	if err = Repo.GetRolesEnvironmentType(&rolesEnvironmentType).Err(); err != nil {
		fmt.Println("Error getEnvironmentRoleRights", err)
		return false, err
	}

	userRight := UserRight(right)

	switch userRight {
	case Read:
		return rolesEnvironmentType.Read, nil
	case Write:
		return rolesEnvironmentType.Write, nil
	default:
		return false, fmt.Errorf(
			"unknown role %v on env %v",
			projectMember.Role,
			environment,
		)
	}
}

func CanUserReadEnvironment(
	Repo repo.IRepo,
	userID uint,
	projectID uint,
	environment *models.Environment,
) (bool, error) {
	return doesUserHaveRightsOnEnvironment(
		Repo,
		userID,
		projectID,
		environment,
		"read",
	)
}

func CanUserWriteOnEnvironment(
	Repo repo.IRepo,
	userID uint,
	projectID uint,
	environment *models.Environment,
) (bool, error) {
	return doesUserHaveRightsOnEnvironment(
		Repo,
		userID,
		projectID,
		environment,
		"write",
	)
}

// CanRoleAddRole tells if a user with a given role can add or set users
// with an other role
func CanRoleAddRole(
	Repo repo.IRepo,
	role models.Role,
	roleToInvite models.Role,
) (can bool, err error) {
	if role.CanAddMember {
		roles := make([]models.Role, 0)

		if role.ID == roleToInvite.ID {
			return true, nil
		}

		err = Repo.GetChildrenRoles(role, &roles).Err()
		if err != nil {
			fmt.Println("Error when retrieving invite roles", Repo.Err())
		} else {
			for _, childRole := range roles {
				if childRole.ID == roleToInvite.ID {
					can = true
					break
				}
			}
		}
	}

	return can && err == nil, err
}

// CanUserSetMemberRole checks if `user` can set the role of user `other` to `role` on `project`.
// `user` can set the `other`’s role to `role` if :
// - both users are members of `project`
// - `users`’s role has `CanAddMembers` set to `true`,
// - `users`’s role is a parent of `other`’s role.
// - users`’s role is a parent of the target `role`
func CanUserSetMemberRole(
	Repo repo.IRepo,
	user models.User,
	other models.User,
	role models.Role,
	project models.Project,
) (can bool, err error) {
	myMember := models.ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	otherMember := models.ProjectMember{
		UserID:    other.ID,
		ProjectID: project.ID,
	}

	if err = Repo.
		GetProjectMember(&myMember).
		GetProjectMember(&otherMember).
		Err(); err != nil {
		return false, err
	}

	canChangeOther, canChangeOtherErr := CanRoleAddRole(
		Repo,
		myMember.Role,
		otherMember.Role,
	)

	canSetTargetRole, canSetTargetRoleErr := CanRoleAddRole(
		Repo,
		myMember.Role,
		role,
	)

	// Owner of organization cannot have its role changed
	isMemberOwnerOfOrga, isMemberOwnerOfOrgaErr := IsUserOwnerOfOrga(
		Repo,
		other.ID,
		project,
	)

	if canChangeOtherErr != nil {
		return false, canChangeOtherErr
	}

	if canSetTargetRoleErr != nil {
		return false, canSetTargetRoleErr
	}

	if isMemberOwnerOfOrgaErr != nil {
		return false, isMemberOwnerOfOrgaErr
	}

	return canChangeOther && canSetTargetRole && !isMemberOwnerOfOrga, nil
}

// CanUserAddMemberWithRole checks if a given user can add members of
// a given role on the project
func CanUserAddMemberWithRole(
	Repo repo.IRepo,
	user models.User,
	role models.Role,
	project models.Project,
) (can bool, err error) {
	member := models.ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}

	if err = Repo.GetProjectMember(&member).Err(); err != nil {
		return false, err
	}

	can, err = CanRoleAddRole(Repo, member.Role, role)

	if err != nil {
		return false, err
	}

	return can, nil
}

func IsUserOwnerOfOrga(
	Repo repo.IRepo,
	userID uint,
	project models.Project,
) (bool, error) {
	orga := models.Organization{}
	if err := Repo.GetProjectsOrganization(project.UUID, &orga).Err(); err != nil {
		return false, err
	}

	if orga.UserID == userID {
		return true, nil
	}
	return false, nil
}

// HasOrganizationNotPaidAndHasNonAdmin function returns true if
// the project's organization is free (not paid)
// and there is at least one of the members is not admin
// which is invalid (although unlikely) state
func HasOrganizationNotPaidAndHasNonAdmin(
	Repo repo.IRepo,
	project models.Project,
) (has bool, err error) {
	var isPaid bool

	members := make([]models.ProjectMember, 0)
	isPaid, err = Repo.IsProjectOrganizationPaid(project.UUID)
	if err != nil {
		goto done
	}

	err = Repo.GetOrganizationMembers(project.OrganizationID, &members).Err()
	if err != nil {
		goto done
	}

	if !isPaid {
		for _, member := range members {
			if member.Role.Name != "admin" {
				has = true
				break
			}
		}
	}

done:
	return has, err
}
