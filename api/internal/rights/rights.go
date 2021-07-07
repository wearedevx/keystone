package rights

import (
	"errors"
	"fmt"

	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type UserRight string

const (
	Read   UserRight = "read"
	Write  UserRight = "write"
	Invite UserRight = "invite"
)

func DoesUserHaveRightsOnEnvironment(Repo repo.IRepo, userID uint, projectID uint, environment *Environment, right string) (bool, error) {
	projectMember := ProjectMember{
		UserID:    userID,
		ProjectID: projectID,
	}

	err := Repo.GetProjectMember(&projectMember).Err()

	if err != nil {
		return false, err
	}

	rolesEnvironmentType := RolesEnvironmentType{
		EnvironmentTypeID: environment.EnvironmentTypeID,
		RoleID:            projectMember.RoleID,
	}

	err = Repo.GetRolesEnvironmentType(&rolesEnvironmentType).Err()

	if errors.Is(err, repo.ErrorNotFound) {
		fmt.Println("No rolesEnvironmentType found")
		return false, err
	}

	if err != nil {
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
		return false, fmt.Errorf("unknown role %v on env %v", projectMember.Role, environment)
	}
}

func CanUserReadEnvironment(Repo repo.IRepo, userID uint, projectID uint, environment *Environment) (bool, error) {
	return DoesUserHaveRightsOnEnvironment(Repo, userID, projectID, environment, "read")
}

func CanUserWriteOnEnvironment(Repo repo.IRepo, userID uint, projectID uint, environment *Environment) (bool, error) {
	return DoesUserHaveRightsOnEnvironment(Repo, userID, projectID, environment, "write")
}

func CanUserInviteOnEnvironment(Repo repo.IRepo, userID uint, projectID uint, environment *Environment) (bool, error) {
	return DoesUserHaveRightsOnEnvironment(Repo, userID, projectID, environment, "invite")
}

// devops can invite on:
// - dev

// CanRoleAddRole tells if a user with a given role can add or set users
// with an other role
func CanRoleAddRole(Repo repo.IRepo, role Role, roleToInvite Role) (can bool, err error) {
	if role.CanAddMember {
		roles := make([]Role, 0)

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
// - `users`’s role is a parent of the target `role`
func CanUserSetMemberRole(Repo repo.IRepo, user User, other User, role Role, project Project) (can bool, err error) {
	myMember := ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	otherMember := ProjectMember{
		UserID:    other.ID,
		ProjectID: project.ID,
	}

	if err = Repo.
		GetProjectMember(&myMember).
		GetProjectMember(&otherMember).
		Err(); err != nil {
		return false, err
	}

	canChangeOther, canChangeOtherErr := CanRoleAddRole(Repo, myMember.Role, otherMember.Role)
	canSetTargetRole, canSetTargetRoleErr := CanRoleAddRole(Repo, myMember.Role, role)

	if canChangeOtherErr != nil {
		return false, canChangeOtherErr
	}

	if canSetTargetRoleErr != nil {
		return false, canSetTargetRoleErr
	}

	return canChangeOther && canSetTargetRole, nil
}

// CanUserAddMemberWithRole checks if a given user can add members of
// a given role on the project
func CanUserAddMemberWithRole(Repo repo.IRepo, user User, role Role, project Project) (can bool, err error) {
	member := ProjectMember{
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
