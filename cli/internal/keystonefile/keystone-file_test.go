package keystonefile

import (
	"testing"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

func TestKeystoneFile(t *testing.T) {
	t.Run("Creates a new structure", func(t *testing.T) {
		// Setup
		testDir, err := utils.CreateTestDir()
		if err != nil {
			t.Errorf("Error creating the test dir: %+v", err)
		}

		// Test
		file := NewKeystoneFile(testDir, models.Project{Name: "test_name"})

		t.Logf("Success: %+v\n", file)

		// TearDown
		utils.CleanTestDir(testDir)
	})

	t.Run("Saves a KeystoneFile", func(t *testing.T) {
		// Setup
		testDir, err := utils.CreateTestDir()
		if err != nil {
			t.Errorf("Error creating the test dir: %+v", err)
		}

		// Test
		file := NewKeystoneFile(testDir, models.Project{Name: "test_name"})

		err = file.Save().Err()

		if err != nil {
			t.Errorf("Error: %+v\n", err)
		} else {
			if !ExistsKeystoneFile(testDir) {
				t.Error("Error: Keystone file was not created")
			} else {
				t.Log("Success")
			}
		}

		// TearDown
		file.Remove()

		utils.CleanTestDir(testDir)
	})
}
