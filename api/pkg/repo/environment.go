package repo

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/wearedevx/keystone/api/pkg/models"
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

	repo.err = repo.GetDb().
		Preload("EnvironmentType").
		Where(*environment).
		First(&environment).
		Error

	return repo
}

func (repo *Repo) GetOrCreateEnvironment(environment *Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Preload("EnvironmentType").
		Where(*environment).
		FirstOrCreate(environment).
		Error

	return repo
}

func (repo *Repo) GetEnvironmentsByProjectUUID(projectUUID string, foundEnvironments *[]Environment) IRepo {

	var project Project
	repo.err = repo.GetDb().Model(&Project{}).Where("uuid = ?", projectUUID).First(&project).Error

	repo.err = repo.GetDb().Model(&Environment{}).Where("project_id = ?", project.ID).Find(&foundEnvironments).Error

	return repo
}

func (repo *Repo) SetNewVersionID(environment *Environment) error {
	newVersionID := uuid.NewV4().String()
	environment.VersionID = newVersionID
	repo.err = repo.GetDb().Model(&Environment{}).Update("version_id", newVersionID).Error
	return repo.Err()
}

func (repo *Repo) GetEnvironmentPublicKeys(environmentID string, publicKeys *PublicKeys) IRepo {
	rows, err := repo.GetDb().Raw(`select u.public_key as PublicKey, u.id as UserID
	from environments as e
	inner join environment_types as et on et.id = e.environment_type_id
	inner join roles_environment_types as ret on ret.environment_type_id = et.id
	inner join roles as r on ret.role_id = r.id
	inner join project_members as pm on r.id = pm.role_id
	inner join users as u on u.id = pm.user_id
	where e.id = ?
	and ret.read = true
	`, environmentID).Rows()

	repo.err = err

	var PublicKey []byte
	var UserID uint

	for rows.Next() {
		rows.Scan(&PublicKey, &UserID)

		publicKeys.Keys = append(publicKeys.Keys, models.UserPublicKey{
			UserID:    fmt.Sprint(UserID),
			PublicKey: PublicKey,
		})
	}

	return repo
}
