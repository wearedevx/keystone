// +build test

package repo

import (
	"errors"
	"os"
	"path"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repo struct {
	err error
	tx  *gorm.DB
}

var db *gorm.DB

func autoMigrate() error {
	err := db.AutoMigrate(&User{}, &LoginRequest{}, &Environment{}, &EnvironmentUserSecret{}, &Message{}, &Project{}, &ProjectMember{}, &Secret{}, &RolesEnvironmentType{}, &Device{}, &Organization{})

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

func Transaction(fn func(IRepo) error) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		repo := &Repo{
			err: nil,
			tx:  tx,
		}
		return fn(repo)

	})
	return err
}

func init() {
	// dbFilePath := path.Join(os.TempDir(), "keystone_gorm-"+uuid.NewV4().String()+".db")
	dbFilePath := path.Join(os.TempDir(), "keystone_gorm.db")

	var err error
	db, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err)
	}

	autoMigrate()
}
