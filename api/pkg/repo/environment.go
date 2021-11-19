package repo

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateEnvironment(environment *models.Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Create(&environment).Error

	return repo
}

func (repo *Repo) GetEnvironment(environment *models.Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Preload("EnvironmentType").
		Preload("Project").
		Where(*environment).
		First(&environment).
		Error

	return repo
}

func (repo *Repo) GetOrCreateEnvironment(
	environment *models.Environment,
) IRepo {
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

func (repo *Repo) GetEnvironmentsByProjectUUID(
	projectUUID string,
	foundEnvironments *[]models.Environment,
) IRepo {
	if repo.err != nil {
		return repo
	}

	var project models.Project
	repo.err = repo.GetDb().
		Model(&models.Project{}).
		Where("uuid = ?", projectUUID).
		First(&project).
		Error

	repo.err = repo.GetDb().
		Model(&models.Environment{}).
		Where("project_id = ?", project.ID).
		Find(&foundEnvironments).
		Error

	return repo
}

func (repo *Repo) SetNewVersionID(environment *models.Environment) error {
	newVersionID := uuid.NewV4().String()
	repo.err = repo.GetDb().
		Model(&models.Environment{}).
		Where(*environment).
		Update("version_id", newVersionID).
		Error
	environment.VersionID = newVersionID
	return repo.Err()
}

func (repo *Repo) GetEnvironmentPublicKeys(
	environmentID string,
	publicKeys *models.PublicKeys,
) IRepo {
	rows, err := repo.GetDb().
		Raw(`select d.id, d.uid, d.name, d.public_key, u.user_id, u.id as UserID
	from environments as e
	inner join environment_types as et on et.id = e.environment_type_id
	inner join roles_environment_types as ret on ret.environment_type_id = et.id
	inner join roles as r on ret.role_id = r.id
	inner join project_members as pm on r.id = pm.role_id and pm.project_id = e.project_id
	inner join users as u on u.id = pm.user_id
	inner join user_devices as ud on u.id = ud.user_id
	inner join devices as d on ud.device_id = d.id
	where e.environment_id = ?
	and ret.read = true
	`, environmentID).
		Rows()
	if err != nil {
		repo.err = err
		return repo
	}

	var PublicKey []byte
	var UserID uint
	var UserUID string
	var DeviceUID string
	var DeviceName string
	var PublicKeyId uint

	for rows.Next() {
		if err := rows.Scan(&PublicKeyId, &DeviceUID, &DeviceName, &PublicKey, &UserUID, &UserID); err != nil {
			repo.err = err
			return repo
		}
		found := false

		for i, pk := range publicKeys.Keys {
			if pk.UserID == UserID {
				publicKeys.Keys[i].PublicKeys = append(
					pk.PublicKeys,
					models.Device{
						PublicKey: PublicKey,
						Name:      DeviceName,
						UID:       DeviceUID,
						ID:        PublicKeyId,
					},
				)
				found = true
			}
		}

		if !found {
			publicKeys.Keys = append(publicKeys.Keys, models.UserPublicKeys{
				UserID: UserID,
				PublicKeys: []models.Device{
					{
						PublicKey: PublicKey,
						Name:      DeviceName,
						UID:       DeviceUID,
						ID:        PublicKeyId,
					},
				},
				UserUID: fmt.Sprint(UserUID),
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
