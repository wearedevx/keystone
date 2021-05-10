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

func CanUserHasRightsEnvironment(Repo repo.IRepo, user *User, project *Project, environment *Environment, right string) (bool, error) {
	projectMember := ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}

	err := Repo.GetProjectMember(&projectMember).Err()

	if err != nil {
		fmt.Println("Error:", err)
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
	case Invite:
		return rolesEnvironmentType.Invite, nil
	default:
		return false, fmt.Errorf("unknown role %v on env %v", projectMember.Role, environment)
	}
}

func CanUserReadEnvironment(Repo repo.IRepo, user *User, project *Project, environment *Environment) (bool, error) {
	return CanUserHasRightsEnvironment(Repo, user, project, environment, "read")
}

func CanUserWriteOnEnvironment(Repo repo.IRepo, user *User, project *Project, environment *Environment) (bool, error) {
	return CanUserHasRightsEnvironment(Repo, user, project, environment, "write")
}

func CanUserInviteOnEnvironment(Repo repo.IRepo, user *User, project *Project, environment *Environment) (bool, error) {
	return CanUserHasRightsEnvironment(Repo, user, project, environment, "invite")
}

// devops can invite on:
// - dev

// Retrieve all role with  invite=true where
func CanUserInviteRole(Repo repo.IRepo, user *User, project *Project, roleToInvite *Role) (bool, error) {
	projectMember := ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}

	Repo.GetProjectMember(&projectMember)

	if err := Repo.Err(); err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	roles := make([]Role, 0)

	Repo.GetInvitableRoles(projectMember.Role, &roles)
	fmt.Println("keystone ~ rights.go ~ roles", roles)

	if Repo.Err() != nil {
		fmt.Println("Error when retriving invite roles", Repo.Err())
	}
	return false, nil

}
