package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
)

func (repo *Repo) CreateEnvironment(project Project, name string) Environment {
	if repo.err != nil {
		return Environment{}
	}

	env := Environment{
		Name: name,
	}

	if repo.Err() == nil {
		repo.err = repo.db.Create(&env).Error
	}

	return env
}

func (repo *Repo) GetEnvironmentByProjectIDAndName(project Project, name string) (Environment, bool) {
	var foundEnvironment Environment

	if repo.err != nil {
		return foundEnvironment, false
	}

	repo.err = repo.db.Where("project_id = ? and name = ?", project.ID, name).First(&foundEnvironment).Error

	return foundEnvironment, repo.err == nil
}

func (repo *Repo) GetOrCreateEnvironment(project Project, name string) Environment {
	if repo.err != nil {
		return Environment{}
	}

	if env, ok := repo.GetEnvironmentByProjectIDAndName(project, name); ok {
		return env
	}

	return repo.CreateEnvironment(project, name)
}

func (repo *Repo) EnvironmentSetUserRole(environment Environment, user User, role UserRole) *Repo {
	if repo.err != nil {
		return repo
	}

	// environmentPermissions := EnvironmentPermissions{
	// 	UserID:        user.ID,
	// 	EnvironmentID: environment.ID,
	// 	Role:          role,
	// }

	// repo.err = repo.db.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "user_id"}, {Name: "environment_id"}},
	// 	DoUpdates: clause.Assignments(map[string]interface{}{"role": role}),
	// }).Create(&environmentPermissions).Error

	return repo
}

func (repo *Repo) EnvironmentSetVariableForUser(environement Environment, secret Secret, user User, value []byte) {
	if repo != nil {
		return
	}

	repo.err = repo.db.Create(&EnvironmentUserSecret{
		EnvironmentID: environement.ID,
		UserID:        user.ID,
		SecretID:      secret.ID,
		Value:         value,
	}).Error
}
