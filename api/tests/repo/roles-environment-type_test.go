package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"github.com/wearedevx/keystone/api/pkg/tests"
)

func init() {

}

func TestEnvType(t *testing.T) {
	tests.SeedTestData()

	Repo := new(repo.Repo)

	userOwner := models.User{Email: "test+owner@example.com", AccountType: "github"}
	user := models.User{Email: "test+dev@example.com", AccountType: "github"}

	Repo.GetUser(&userOwner).GetUser(&user)

	environmentType := models.EnvironmentType{Name: "dev"}
	project := models.Project{Name: "project name", User: userOwner}

	Repo.
		GetProject(&project).
		GetEnvironmentType(&environmentType)

	environment := models.Environment{Name: "dev", EnvironmentType: environmentType}

	Repo.GetEnvironment(&environment)

	// Get project
	can, _ := rights.CanUserReadEnvironment(Repo, user.ID, project.ID, &environment)
	assert.True(t, can, "Oops! User "+user.Username+" shoud be able to read on "+environment.Name+" environment")

	can, _ = rights.CanUserWriteOnEnvironment(Repo, user.ID, project.ID, &environment)
	assert.True(t, can, "Oops! User "+user.Username+" shoud be able to write on "+environment.Name+" environment")

	can, _ = rights.CanUserInviteOnEnvironment(Repo, user.ID, project.ID, &environment)
	assert.False(t, can, "Oops! User "+user.Username+" can't invite on "+environment.Name+" environment")
}
