package keystonefile

import (
	"testing"

	. "github.com/wearedevx/keystone/internal/models"
	. "github.com/wearedevx/keystone/internal/utils"
)

func TestKeystoneFile(t *testing.T) {
	t.Run("Creates a new structure", func(t *testing.T) {
		// Setup
		testDir, err := CreateTestDir()
		if err != nil {
			t.Errorf("Error creating the test dir: %+v", err)
		}

		// Test
		file := NewKeystoneFile(testDir, Project{Name: "test_name"})

		t.Logf("Success: %+v\n", file)

		// TearDown
		CleanTestDir(testDir)
	})

	t.Run("Saves a KeystoneFile", func(t *testing.T) {
		// Setup
		testDir, err := CreateTestDir()
		if err != nil {
			t.Errorf("Error creating the test dir: %+v", err)
		}

		// Test
		file := NewKeystoneFile(testDir, Project{Name: "test_name"})

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

		CleanTestDir(testDir)
	})
}
