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

// getDSN builds the postgres DSN from environment variables
func getDSN() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

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

	return runner.Run().Err()
}

func (repo *Repo) Err() error {
	return repo.err
}

func (repo *Repo) Connect() {
	db, err := gorm.Open(getPostgres(), &gorm.Config{})

	if err == nil {
		autoMigrate(db)
	}

	repo.db = db
	repo.err = err
}
