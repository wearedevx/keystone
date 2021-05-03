package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wearedevx/keystone/internal/repo"
	"github.com/wearedevx/keystone/internal/rights"
	"github.com/wearedevx/keystone/tests/utils"
)

func init() {

}

func TestEnvType(t *testing.T) {
	utils.SeedTestData()

	Repo := new(repo.Repo)

	// Repo.GetRolesEnvironmentType(environment)

	user, _ := Repo.GetUserByEmailAndAccountType("test+dev@example.com", "github")
	project, _ := Repo.GetUserProjectWithName(user, "project name")
	environmentType, _ := Repo.GetEnvironmentTypeByName("dev")
	environment, _ := Repo.GetEnvironmentByProjectIDAndEnvType(project, environmentType)

	// Get project
	// repo RightsRepo, user *User, project *Project, environment *Environment
	can, _ := rights.CanUserReadEnvironment(Repo, &user, &project, &environment)
	assert.True(t, can, "Oops! User "+user.Username+" shoud be able to read on "+environment.Name+" environment")

	can, _ = rights.CanUserWriteOnEnvironment(Repo, &user, &project, &environment)
	assert.True(t, can, "Oops! User "+user.Username+" shoud be able to write on "+environment.Name+" environment")

	can, _ = rights.CanUserInviteOnEnvironment(Repo, &user, &project, &environment)
	assert.False(t, can, "Oops! User "+user.Username+" can't invite on "+environment.Name+" environment")
}
