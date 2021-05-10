package tests

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func SeedTestData() {
	Repo := new(repo.Repo)

	devEnvironmentType := EnvironmentType{Name: "dev"}
	stagingEnvironmentType := EnvironmentType{Name: "staging"}
	prodEnvironmentType := EnvironmentType{Name: "prod"}

	devRole := Role{Name: "dev"}
	devopsRole := Role{Name: "devops"}
	adminRole := Role{Name: "admin"}

	Repo.
		// EnvTypes
		GetOrCreateEnvironmentType(&devEnvironmentType).
		GetOrCreateEnvironmentType(&stagingEnvironmentType).
		GetOrCreateEnvironmentType(&prodEnvironmentType).

		// Roles
		GetOrCreateRole(&devRole).
		GetOrCreateRole(&devopsRole).
		GetOrCreateRole(&adminRole).

		// DEV
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            devRole,
			EnvironmentType: devEnvironmentType,
			Read:            true,
			Write:           true,
			Invite:          false,
		}).
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            devRole,
			EnvironmentType: stagingEnvironmentType,
			Read:            false,
			Write:           false,
			Invite:          false,
		}).
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            devRole,
			EnvironmentType: prodEnvironmentType,
			Read:            false,
			Write:           false,
			Invite:          false,
		}).

		// Staging
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            devopsRole,
			EnvironmentType: devEnvironmentType,
			Read:            true,
			Write:           true,
			Invite:          true,
		}).
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            devopsRole,
			EnvironmentType: stagingEnvironmentType,
			Read:            true,
			Write:           true,
			Invite:          true,
		}).
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            devopsRole,
			EnvironmentType: prodEnvironmentType,
			Read:            false,
			Write:           false,
			Invite:          false,
		}).

		// ADMIN
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            adminRole,
			EnvironmentType: devEnvironmentType,
			Read:            true,
			Write:           true,
			Invite:          true,
		}).
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            adminRole,
			EnvironmentType: stagingEnvironmentType,
			Read:            true,
			Write:           true,
			Invite:          true,
		}).
		GetOrCreateRoleEnvType(&RolesEnvironmentType{
			Role:            adminRole,
			EnvironmentType: prodEnvironmentType,
			Read:            true,
			Write:           true,
			Invite:          true,
		})

	var userProjectOwner *User = &User{
		ExtID:       "my iowner ext id",
		AccountType: "github",
		Username:    "Username owner " + uuid.NewV4().String(),
		Fullname:    "Fullname owner",
		Email:       "test+owner@example.com",
	}

	var devUser *User = &User{
		ExtID:       "my ext id",
		AccountType: "github",
		Username:    "Username dev " + uuid.NewV4().String(),
		Fullname:    "Fullname dev",
		Email:       "test+dev@example.com",
	}

	fmt.Println("keystone ~ functions.go ~ error DOUX DOUX")

	Repo.GetOrCreateUser(userProjectOwner)
	Repo.GetOrCreateUser(devUser)

	var project *Project = &Project{
		Name: "project name",
		User: *userProjectOwner,
	}

	environmentType := EnvironmentType{Name: "dev"}

	Repo.
		GetOrCreateProject(project).
		GetOrCreateEnvironmentType(&environmentType)

	environment := Environment{
		Name:              environmentType.Name,
		ProjectID:         project.ID,
		EnvironmentTypeID: environmentType.ID,
	}

	Repo.GetOrCreateEnvironment(&environment)

	projectMember := ProjectMember{
		ProjectID: project.ID,
		UserID:    devUser.ID,
	}

	Repo.GetOrCreateProjectMember(&projectMember, "dev")
}
