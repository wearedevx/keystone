package repo

import (
	uuid "github.com/satori/go.uuid"
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateEnvironment(environment *Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Create(&environment).Error

	return repo
}

func (repo *Repo) GetEnvironment(environment *Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Preload("EnvironmentType").First(&environment).Error

	return repo
}

func (repo *Repo) GetOrCreateEnvironment(environment *Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Preload("EnvironmentType").FirstOrCreate(environment).Error

	return repo
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
