package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateEnvironmentType(name string) (EnvironmentType, error) {
	envType := EnvironmentType{
		Name: name,
	}
	repo.err = repo.GetDb().Create(&envType).Error
	return envType, repo.err
}

func (repo *Repo) GetEnvironmentTypeByName(name string) (EnvironmentType, bool) {
	var foundEnvironmentType EnvironmentType

	// Why???
	// When before there's a error as Record not found, plouf
	// if repo.err != nil {
	// 	return foundEnvironmentType, false
	// }

	repo.err = repo.GetDb().Where("name = ?", name).First(&foundEnvironmentType).Error

	return foundEnvironmentType, repo.err == nil
}

func (repo *Repo) GetOrCreateEnvironmentType(name string) (EnvironmentType, error) {
	if repo.err != nil {
		return EnvironmentType{}, repo.err
	}

	if environmentType, ok := repo.GetEnvironmentTypeByName(name); ok {
		return environmentType, nil
	}

	return repo.CreateEnvironmentType(name)
}
