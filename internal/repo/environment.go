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
		repo.err = repo.GetDb().Create(&env).Error
	}

	return env
}

func (repo *Repo) GetEnvironmentByProjectIDAndName(project Project, name string) (Environment, bool) {
	var foundEnvironment Environment

	if repo.err != nil {
		return foundEnvironment, false
	}

	repo.err = repo.GetDb().Where("project_id = ? and name = ?", project.ID, name).First(&foundEnvironment).Error

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
