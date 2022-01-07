// +build test

package repo

import (
	"errors"
	"os"
	"path"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/wearedevx/keystone/api/db/seed"
	"github.com/wearedevx/keystone/api/pkg/message"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repo struct {
	err      error
	tx       *gorm.DB
	messages message.MessageService
}

const (
	DialectPostgres = iota
	DialectSQLite
)

type Dialect int

var dialect Dialect = DialectSQLite

var db *gorm.DB

func autoMigrate() error {
	db.AutoMigrate(
		&LoginRequest{},
		&Environment{},
		&EnvironmentUserSecret{},
		&Message{},
		&Project{},
		&ProjectMember{},
		&Secret{},
		&Roles{},
		&EnvironmentType{},
		&RolesEnvironmentType{},
		&User{},
		&Device{},
		&UserDevice{},
		&Organization{},
		&ActivityLog{},
		&CheckoutSession{},
	)
	return nil
}

func (repo *Repo) Err() error {
	if errors.Is(repo.err, gorm.ErrRecordNotFound) {
		return ErrorNotFound
	}

	return repo.err
}

func (repo *Repo) ClearErr() IRepo {
	repo.err = nil

	return repo
}

func (repo *Repo) GetDb() *gorm.DB {
	return db
}

func (repo *Repo) GetDialect() Dialect {
	return dialect
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

func (repo *Repo) MessageService() message.MessageService {
	return repo.messages
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

func NewRepo() *Repo {
	return &Repo{
		err:      nil,
		tx:       db,
		messages: message.NewMessageService(),
	}
}

func init() {
	dbFilePath := path.Join(os.TempDir(), "keystone_gorm.db")

	var err error
	db, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err)
	}

	err = autoMigrate()
	if err != nil {
		// ignore... make the tests fail if there is an output
	}

	seed.Seed(db)
}
