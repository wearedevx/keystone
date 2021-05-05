package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
)

func (repo *Repo) CreateEnvironment(project Project, environnementType EnvironmentType) Environment {
	if repo.err != nil {
		return Environment{}
	}

	env := Environment{
		EnvironmentTypeID: environnementType.ID,
		ProjectID:         project.ID,
	}

	if repo.Err() == nil {
		repo.err = repo.GetDb().Create(&env).Error
	}

	return env
}

func (repo *Repo) GetEnvironmentByProjectIDAndEnvType(project Project, environnementType EnvironmentType) (Environment, bool) {
	var foundEnvironment Environment

	if repo.err != nil {
		return foundEnvironment, false
	}

	repo.err = repo.GetDb().Preload("EnvironmentType").Where("project_id = ? and environment_type_id = ?", project.ID, environnementType.ID).First(&foundEnvironment).Error

	return foundEnvironment, repo.err == nil
}

func (repo *Repo) GetOrCreateEnvironment(project Project, environnementType EnvironmentType) Environment {
	if repo.err != nil {
		return Environment{}
	}

	if env, ok := repo.GetEnvironmentByProjectIDAndEnvType(project, environnementType); ok {
		return env
	}

	return repo.CreateEnvironment(project, environnementType)
}
