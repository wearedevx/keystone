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
}

var db *gorm.DB

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
	fmt.Println("port:", port)
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

func AutoMigrate() error {
	return nil
}

func (repo *Repo) Err() error {
	return repo.err
}

func (repo *Repo) GetDb() *gorm.DB {
	return db
}

func init() {
	var err error

	db, err = gorm.Open(getPostgres(), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection established")
}
