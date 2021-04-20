// +build !test

package repo

import (
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	// . "github.com/wearedevx/keystone/internal/models"
	// . "github.com/wearedevx/keystone/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	err error
	db  *gorm.DB
}

func getEnv(varname string, fallback string) string {
	if value, ok := os.LookupEnv(varname); ok {
		return value
	}

	return fallback
}

// getDSN builds the postgres DSN from environment variables
func getDSN() string {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "keystone-dev")
	password := getEnv("DB_PASSWORD", "keystone-dev")
	dbname := getEnv("DB_NAME", "keystone")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

// getPostrgres gets the postgres driver for GORM
func getPostgres() gorm.Dialector {

	os.TempDir()
	config := postgres.Config{
		DSN: getDSN(),
	}

	if driver := os.Getenv("DB_DRIVERNAME"); driver != "" {
		config.DriverName = driver
	}

	return postgres.New(config)
}

func AutoMigrate(db *gorm.DB) error {
	// 	runner := NewRunner([]RunnerAction{
	// 		NewAction(func() error {
	// 			return db.AutoMigrate(&LoginRequest{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.AutoMigrate(&Project{}, &Environment{}, &User{}, &Secret{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.AutoMigrate(&EnvironmentPermissions{}, &ProjectPermissions{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.AutoMigrate(&EnvironmentUserSecret{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.SetupJoinTable(&User{}, "Projects", &ProjectPermissions{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.SetupJoinTable(&User{}, "Environments", &EnvironmentPermissions{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.SetupJoinTable(&Project{}, "Users", &ProjectPermissions{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.SetupJoinTable(&User{}, "EnvironmentsSecrets", &EnvironmentUserSecret{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.SetupJoinTable(&Environment{}, "UserSecrets", &EnvironmentUserSecret{})
	// 		}),
	// 		NewAction(func() error {
	// 			return db.SetupJoinTable(&Secret{}, "UserEnvironments", &EnvironmentUserSecret{})
	// 		}),
	return nil
}

// 	return runner.Run().Error()

// }

func (repo *Repo) Err() error {
	return repo.err
}

func (repo *Repo) GetDb() *gorm.DB {
	return repo.db
}

func (repo *Repo) Connect() {
	fmt.Println("CLOUD SQL")
	var err error
	db, err := gorm.Open(getPostgres(), &gorm.Config{})

	repo.db = db
	repo.err = err
}
