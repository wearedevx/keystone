// +build test

package repo

import (
	"os"
	"path"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	. "github.com/wearedevx/keystone/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repo struct {
	err error
}

var db *gorm.DB

func autoMigrate() error {
	db.AutoMigrate(&LoginRequest{}, &Environment{}, &EnvironmentUserSecret{}, &Message{}, &Project{}, &ProjectMember{}, &Secret{}, &RolesEnvironmentType{})
	return nil
}

func (repo *Repo) Err() error {
	return repo.err
}

func (repo *Repo) GetDb() *gorm.DB {
	return db
}

func init() {
	dbFilePath := path.Join(os.TempDir(), "keystone_gorm.db")
	// fmt.Println("keystone ~ repo_testing.go ~ dbFilePath", dbFilePath)

	var err error
	db, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	autoMigrate()
}
