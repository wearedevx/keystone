package tests

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func SeedTestData() {
	Repo := new(repo.Repo)

	devEnvironmentType, _ := Repo.GetOrCreateEnvironmentType("dev")
	stagingEnvironmentType, _ := Repo.GetOrCreateEnvironmentType("staging")
	prodEnvironmentType, _ := Repo.GetOrCreateEnvironmentType("prod")

	devRole := Repo.GetOrCreateRole("dev")
	devopsRole := Repo.GetOrCreateRole("devops")
	adminRole := Repo.GetOrCreateRole("admin")

	// DEV
	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devRole,
		EnvironmentType: devEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          false,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devRole,
		EnvironmentType: stagingEnvironmentType,
		Read:            false,
		Write:           false,
		Invite:          false,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devRole,
		EnvironmentType: prodEnvironmentType,
		Read:            false,
		Write:           false,
		Invite:          false,
	})

	// Staging
	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devopsRole,
		EnvironmentType: devEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devopsRole,
		EnvironmentType: stagingEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            devopsRole,
		EnvironmentType: prodEnvironmentType,
		Read:            false,
		Write:           false,
		Invite:          false,
	})

	// ADMIN
	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            adminRole,
		EnvironmentType: devEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
		Role:            adminRole,
		EnvironmentType: stagingEnvironmentType,
		Read:            true,
		Write:           true,
		Invite:          true,
	})

	Repo.GetOrCreateRoleEnvType(&RolesEnvironmentType{
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
	}

	Repo.GetOrCreateProject(project, *userProjectOwner)

	environmentType, _ := Repo.GetOrCreateEnvironmentType("dev")

	Repo.GetOrCreateEnvironment(*project, environmentType)

	Repo.GetOrCreateProjectMember(project, devUser, "dev")

}
