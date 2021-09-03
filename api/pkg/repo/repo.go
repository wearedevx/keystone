// +build !test

package repo

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	// . "github.com/wearedevx/keystone/internal/models"
	. "github.com/wearedevx/keystone/api/internal/utils"
	"github.com/wearedevx/keystone/api/pkg/message"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	err      error
	tx       *gorm.DB
	messages *message.MessageService
}

var db *gorm.DB

// getDSN builds the postgres DSN from environment variables
func getDSN() string {
	host := GetEnv("DB_HOST", "db")
	port := GetEnv("DB_PORT", "5432")
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
	fmt.Printf("config: %+v\n", config)

	if driver := os.Getenv("DB_DRIVERNAME"); driver != "" {
		config.DriverName = driver
	}

	return postgres.New(config)
}

func AutoMigrate() error {
	return nil
}

func Transaction(fn func(IRepo) error) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		repo := &Repo{
			err:      nil,
			tx:       tx,
			messages: message.NewMessageService(),
		}
		return fn(repo)

	})
	return err
}

func (repo *Repo) Err() error {
	if errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return ErrorNotFound
	}

	return repo.err
}

func (repo *Repo) GetDb() *gorm.DB {
	return repo.tx
}

func (repo *Repo) MessageService() *message.MessageService {
	return repo.messages
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

	db, err = gorm.Open(getPostgres(), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection established")
}
