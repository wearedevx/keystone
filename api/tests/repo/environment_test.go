package repo

import (
	"fmt"
	"testing"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"github.com/wearedevx/keystone/api/pkg/tests"
)

func init() {

}

func TestEnvironment(t *testing.T) {
	fmt.Println("api ~ environment_test.go ~ TestEnvironment")
	tests.SeedTestData()

	Repo := new(repo.Repo)

	publicKeys := models.PublicKeys{
		Keys: make([]models.UserPublicKey, 0),
	}

	Repo.GetEnvironmentPublicKeys("1", &publicKeys)
	// fmt.Println("api ~ environment_test.go ~ publicKeys", publicKeys)
	fmt.Printf("%+v\n", publicKeys)

	// fmt.Println("api ~ environment_test.go ~ publicKeys", publicKeys)

	// assert.True(t, publicKeys, "Oops! User "+user.Username+" shoud be able to write on "+environment.Name+" environment")
}
