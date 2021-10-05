// +build !test

package repo

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
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

var dbHost string
var dbPort string
var dbUser string
var dbPassword string
var dbName string
var dbDialect string
var dbDriverName string

func getOrDefault(s string, d string) string {
	if s == "" {
		return d
	}

	return s
}

// getDSN builds the postgres DSN from environment variables
func getDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getOrDefault(dbHost, "127.0.0.1"),
		getOrDefault(dbPort, "5432"),
		getOrDefault(dbUser, "ks"),
		getOrDefault(dbPassword, "ks"),
		getOrDefault(dbName, "ks"),
	)
}

// getPostrgres gets the postgres driver for GORM
func getPostgres() gorm.Dialector {

	os.TempDir()
	config := postgres.Config{
		DSN: getDSN(),
	}

	if dbDriverName != "" {
		config.DriverName = dbDriverName
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
	return db
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
