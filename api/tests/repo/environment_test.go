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
	tests.SeedTestData()

	Repo := new(repo.Repo)

	publicKeys := models.PublicKeys{
		Keys: make([]models.UserPublicKey, 0),
	}

	Repo.GetEnvironmentPublicKeys("1", &publicKeys)
	fmt.Printf("%+v\n", publicKeys)
}
