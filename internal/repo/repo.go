// +build !test

package repo

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	// . "github.com/wearedevx/keystone/internal/models"
	. "github.com/wearedevx/keystone/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	err error
}

var db *gorm.DB

// getDSN builds the postgres DSN from environment variables
func getDSN() string {
	host := GetEnv("DB_HOST", "db")
	port := GetEnv("DB_PORT", "5432")
	fmt.Println("port:", port)
	user := GetEnv("DB_USER", "ks")
	password := GetEnv("DB_PASSWORD", "ks")
	dbname := GetEnv("DB_NAME", "ks")

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

func AutoMigrate() error {
	return nil
}

func (repo *Repo) Err() error {
	return repo.err
}

func (repo *Repo) GetDb() *gorm.DB {
	return db
}

func (repo *Repo) notFoundAsBool(call func() error) (bool, error) {
	var err error
	found := false

	err = call()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	} else {
		found = true
	}

	return found, err
}

func init() {
	var err error

	db, err = gorm.Open(getPostgres(), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection established")
}
