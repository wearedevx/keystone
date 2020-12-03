package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
	"gorm.io/gorm/clause"
)

func (repo *Repo) CreateEnvironment(project Project, name string) Environment {
	env := Environment{
		Name:      name,
		ProjectID: project.ID,
	}

	if repo.Err() == nil {
		repo.err = repo.db.Create(&env).Error
	}

	return env
}

func (repo *Repo) EnvironmentSetUserRole(environment Environment, user User, role UserRole) *Repo {
	environmentPermissions := EnvironmentPermissions{
		UserID:        user.ID,
		EnvironmentID: environment.ID,
		Role:          role,
	}

	repo.err = repo.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "environment_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"role": role}),
	}).Create(&environment)

	return repo
}
