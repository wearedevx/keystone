// +build test

package repo

import (
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	. "github.com/wearedevx/keystone/internal/models"
	"gorm.io/driver/sqlite"
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
// func getPostgres() gorm.Dialector {
// 	config := postgres.Config{
// 		DSN: getDSN(),
// 	}

// 	if driver := os.Getenv("DB_DRIVERNAME"); driver != "" {
// 		config.DriverName = driver
// 	}

// 	return postgres.New(config)
// }

func AutoMigrate(db *gorm.DB) error {

	db.AutoMigrate(&LoginRequest{}, &Environment{}, &EnvironmentUserSecret{}, &Message{}, &Project{}, &ProjectMember{}, &Secret{})

	return nil
}

func (repo *Repo) Err() error {
	return repo.err
}

func (repo *Repo) GetDb() *gorm.DB {
	return repo.db
}

func (repo *Repo) Connect() {
	fmt.Println("REPO TESTTT")
	var err error
	// db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	db, err := gorm.Open(sqlite.Open("/Users/kevin/travail/devx/keystone/tests/gorm.db"), &gorm.Config{})

	// db, err := gorm.Open(sqlite.Open("tutu.db"), &gorm.Config{})

	repo.db = db
	repo.err = err
}
