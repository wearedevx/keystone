package tests

import (
	uuid "github.com/satori/go.uuid"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func SeedTestData() {
	Repo := new(repo.Repo)

	devEnvironmentType := models.EnvironmentType{Name: "dev"}
	stagingEnvironmentType := models.EnvironmentType{Name: "staging"}
	prodEnvironmentType := models.EnvironmentType{Name: "prod"}

	devRole := models.Role{Name: "dev"}
	devopsRole := models.Role{Name: "devops"}
	adminRole := models.Role{Name: "admin"}

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
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            devRole,
			EnvironmentType: devEnvironmentType,
			Read:            true,
			Write:           true,
		}).
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            devRole,
			EnvironmentType: stagingEnvironmentType,
			Read:            false,
			Write:           false,
		}).
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            devRole,
			EnvironmentType: prodEnvironmentType,
			Read:            false,
			Write:           false,
		}).

		// Staging
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            devopsRole,
			EnvironmentType: devEnvironmentType,
			Read:            true,
			Write:           true,
		}).
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            devopsRole,
			EnvironmentType: stagingEnvironmentType,
			Read:            true,
			Write:           true,
		}).
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            devopsRole,
			EnvironmentType: prodEnvironmentType,
			Read:            false,
			Write:           false,
		}).

		// ADMIN
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            adminRole,
			EnvironmentType: devEnvironmentType,
			Read:            true,
			Write:           true,
		}).
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            adminRole,
			EnvironmentType: stagingEnvironmentType,
			Read:            true,
			Write:           true,
		}).
		GetOrCreateRoleEnvType(&models.RolesEnvironmentType{
			Role:            adminRole,
			EnvironmentType: prodEnvironmentType,
			Read:            true,
			Write:           true,
		})

	var userProjectOwner = &models.User{
		ExtID:       "my iowner ext id",
		AccountType: "github",
		Username:    "Username owner " + uuid.NewV4().String(),
		Fullname:    "Fullname owner",
		Email:       "test+owner@example.com",
	}

	var devUser = &models.User{
		ExtID:       "my ext id",
		AccountType: "github",
		Username:    "Username dev " + uuid.NewV4().String(),
		Fullname:    "Fullname dev",
		Email:       "test+dev@example.com",
	}

	Repo.GetOrCreateUser(userProjectOwner)
	Repo.GetOrCreateUser(devUser)

	var project = &models.Project{
		Name: "project name",
		User: *userProjectOwner,
	}

	environmentType := models.EnvironmentType{Name: "dev"}

	Repo.
		GetOrCreateProject(project).
		GetOrCreateEnvironmentType(&environmentType)

	environment := models.Environment{
		Name:              environmentType.Name,
		ProjectID:         project.ID,
		EnvironmentTypeID: environmentType.ID,
	}

	Repo.GetOrCreateEnvironment(&environment)

	projectMember := models.ProjectMember{
		ProjectID: project.ID,
		UserID:    devUser.ID,
	}

	Repo.GetOrCreateProjectMember(&projectMember, "dev")
}
