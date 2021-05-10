package repo

import (
	uuid "github.com/satori/go.uuid"
	. "github.com/wearedevx/keystone/api/pkg/models"
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
func (repo *Repo) GetEnvironment(environmentID string) Environment {
	var foundEnvironment Environment

	if repo.err != nil {
		return foundEnvironment
	}

	repo.err = repo.GetDb().Model(&Environment{}).Where("environment_id = ?", environmentID).First(&foundEnvironment).Error

	return foundEnvironment
}

func (repo *Repo) GetEnvironmentsByProjectUUID(projectUUID string) []Environment {
	var foundEnvironments []Environment

	var project Project
	repo.err = repo.GetDb().Model(&Project{}).Where("uuid = ?", projectUUID).First(&project).Error

	repo.err = repo.GetDb().Model(&Environment{}).Where("project_id = ?", project.ID).Find(&foundEnvironments).Error

	return foundEnvironments
}

func (repo *Repo) SetNewVersionID(environment *Environment) error {
	newVersionID := uuid.NewV4().String()
	environment.VersionID = newVersionID
	repo.err = repo.GetDb().Model(&Environment{}).Update("version_id", newVersionID).Error
	return repo.Err()
}
