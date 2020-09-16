package repo

import (
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	. "github.com/wearedevx/keystone/internal/models"
	. "github.com/wearedevx/keystone/internal/utils"
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
	user := getEnv("DB_USER", "keystone-dev")
	password := getEnv("DB_PASSWORD", "keystone-dev")
	dbname := getEnv("DB_NAME", "keystone")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
}

// getPostrgres gets the postgres driver for GORM
func getPostgres() gorm.Dialector {
	config := postgres.Config{
		DSN: getDSN(),
	}

	if driver := os.Getenv("DB_DRIVERNAME"); driver != "" {
		config.DriverName = driver
	}

	return postgres.New(config)
}

func autoMigrate(db *gorm.DB) error {
	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			return db.AutoMigrate(&LoginRequest{})
		}),
		NewAction(func() error {
			return db.AutoMigrate(&User{})
		}),
	})

	return runner.Run().Error()

}

func (repo *Repo) Err() error {
	return repo.err
}

func (repo *Repo) Connect() {
	db, err := gorm.Open(getPostgres(), &gorm.Config{})

	if err == nil {
		repo.err = autoMigrate(db)
	}

	repo.db = db
	repo.err = err
}
