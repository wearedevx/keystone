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
	repo.err = repo.GetDb().Model(&Environment{}).Where(environment).Update("version_id", newVersionID).Error
	environment.VersionID = newVersionID
	return repo.Err()
}

func (repo *Repo) GetEnvironmentPublicKeys(environmentID string, publicKeys *PublicKeys) IRepo {
	rows, err := repo.GetDb().Raw(`select pk.id as PublicKeyId ,pk.device as Device, pk.key as PublicKey, u.user_id as UserUID, u.id as UserID
	from environments as e
	inner join environment_types as et on et.id = e.environment_type_id
	inner join roles_environment_types as ret on ret.environment_type_id = et.id
	inner join roles as r on ret.role_id = r.id
	inner join project_members as pm on r.id = pm.role_id and pm.project_id = e.project_id
	inner join users as u on u.id = pm.user_id
	inner join public_keys as pk on u.id = pk.user_id
	where e.environment_id = ?
	and ret.read = true
	`, environmentID).Rows()

	repo.err = err

	var PublicKey []byte
	var UserID uint
	var UserUID string
	var Device string
	var PublicKeyId uint

	for rows.Next() {
		rows.Scan(&PublicKeyId, &Device, &PublicKey, &UserUID, &UserID)
		found := false

		for i, pk := range publicKeys.Keys {
			if pk.UserID == UserID {
				publicKeys.Keys[i].PublicKeys = append(pk.PublicKeys, models.PublicKey{Key: PublicKey, Device: Device, UserID: UserID, ID: PublicKeyId})
				found = true
			}
		}

		if !found {
			publicKeys.Keys = append(publicKeys.Keys, models.UserPublicKeys{
				UserID:     UserID,
				PublicKeys: []models.PublicKey{{Key: PublicKey, Device: Device, UserID: UserID, ID: PublicKeyId}},
				UserUID:    fmt.Sprint(UserUID),
			})
		}
	}
	return repo
}

func (repo *Repo) DeleteProjectsEnvironments(project *models.Project) IRepo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.
		GetDb().
		Delete(models.Environment{}, "project_id = ?", project.ID).
		Error

	return repo
}
