package rights

import (
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

type RightsRepo interface {
	GetRolesEnvironmentType(environment *Environment, role *Role) (*RolesEnvironmentType, error)
	GetProjectMember(user *User, project *Project) (ProjectMember, error)
	GetInvitableRoles(role Role, roles *[]Role) *repo.Repo
}

func CanUserHasRightsEnvironment(repo RightsRepo, user *User, project *Project, environment *Environment, right string) (bool, error) {
	projectMember, err := repo.GetProjectMember(user, project)

	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	rolesEnvironmentType, err := repo.GetRolesEnvironmentType(environment, &projectMember.Role)

	if err != nil {
		fmt.Println("Error getEnvironmentRoleRights", err)
		return false, err
	}

	if rolesEnvironmentType == nil {
		fmt.Println("No rolesEnvironmentType found", err)
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

func CanUserReadEnvironment(repo RightsRepo, user *User, project *Project, environment *Environment) (bool, error) {
	return CanUserHasRightsEnvironment(repo, user, project, environment, "read")
}

func CanUserWriteOnEnvironment(repo RightsRepo, user *User, project *Project, environment *Environment) (bool, error) {
	return CanUserHasRightsEnvironment(repo, user, project, environment, "write")
}

func CanUserInviteOnEnvironment(repo RightsRepo, user *User, project *Project, environment *Environment) (bool, error) {
	return CanUserHasRightsEnvironment(repo, user, project, environment, "invite")
}

// devops can invite on:
// - dev

// Retrieve all role with  invite=true where
func CanUserInviteRole(repo repo.Repo, user *User, project *Project, roleToInvite *Role) (bool, error) {
	projectMember, err := repo.GetProjectMember(user, project)

	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}

	var roles *[]Role

	repo.GetInvitableRoles(projectMember.Role, roles)
	fmt.Println("keystone ~ rights.go ~ roles", roles)

	if repo.Err() != nil {
		fmt.Println("Error when retriving invite roles", repo.Err())
	}
	return false, nil

}
