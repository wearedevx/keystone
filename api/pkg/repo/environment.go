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
		Preload("Project.Organization").
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

	repo.err = repo.GetDb().
		Model(&models.Environment{}).
		Joins("inner join projects p on p.id = project_id").
		Where("p.uuid = ?", projectUUID).
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
		// all the not deleted devices that are related to a user who can read
		// the given environment (who has a role with `read=true` for the environment type
		// of the given environment
		Raw(`SELECT d.id, d.uid, d.name, d.public_key, u.user_id, u.id AS UserID
	FROM environments AS e
	INNER JOIN environment_types AS et ON et.id = e.environment_type_id
	INNER JOIN roles_environment_types AS ret ON ret.environment_type_id = et.id
	INNER JOIN roles AS r ON ret.role_id = r.id
	INNER JOIN project_members AS pm ON r.id = pm.role_id AND pm.project_id = e.project_id
	INNER JOIN users AS u ON u.id = pm.user_id
	INNER JOIN user_devices AS ud ON u.id = ud.user_id
	INNER JOIN devices AS d ON ud.device_id = d.id
	WHERE e.environment_id = ?
	AND ret.read = true
  AND d.deleted_at IS NULL
	`, environmentID).
		Rows()
	if err != nil {
		repo.err = err
		return repo
	}

	var PublicKey []byte
	var UserID, PublicKeyId uint
	var UserUID, DeviceUID, DeviceName string

	for rows.Next() {
		if err := rows.Scan(
			&PublicKeyId,
			&DeviceUID,
			&DeviceName,
			&PublicKey,
			&UserUID,
			&UserID,
		); err != nil {
			repo.err = err
			return repo
		}
		found := false

		for i, pk := range publicKeys.Keys {
			if pk.UserID == UserID {
				publicKeys.Keys[i].Devices = append(
					pk.Devices,
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
			publicKeys.Keys = append(publicKeys.Keys, models.UserDevices{
				UserID: UserID,
				Devices: []models.Device{
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
