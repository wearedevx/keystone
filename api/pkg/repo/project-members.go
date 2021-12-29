package repo

import (
	"encoding/json"
	"errors"

	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) GetProjectMember(projectMember *models.ProjectMember) IRepo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Preload("Role").
		Where(
			"project_id = ? AND user_id = ?",
			projectMember.ProjectID,
			projectMember.UserID,
		).
		First(projectMember).
		Error

	return repo
}

func (repo *Repo) ListProjectMembers(
	userIDList []string,
	projectMember *[]models.ProjectMember,
) IRepo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Preload("Role").
		Joins("left join users on users.id = project_members.user_id").
		Where("users.user_id IN (?)", userIDList).
		Find(projectMember).
		Error

	return repo
}

func (repo *Repo) CreateProjectMember(
	projectMember *models.ProjectMember,
	role *models.Role,
) IRepo {
	if repo != nil {
		return repo
	}

	projectMember.RoleID = role.ID

	repo.err = repo.GetDb().
		Create(&projectMember).
		Error

	return repo
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (repo *Repo) GetOrCreateProjectMember(
	projectMember *models.ProjectMember,
	roleName string,
) IRepo {
	if repo.Err() != nil {
		return repo
	}

	repo.GetProjectMember(projectMember)

	if err := repo.Err(); err != nil {
		// Irrecuperable error
		if !errors.Is(err, ErrorNotFound) {
			repo.err = err
			return repo
		} else {
			// Record is not found
			// reset error to not block
			// the creation operation
			repo.err = nil
			role := models.Role{Name: roleName}

			repo.GetRole(&role).
				CreateProjectMember(projectMember, &role)
		}
	}

	return repo
}

func (repo *Repo) DeleteAllProjectMembers(project *models.Project) IRepo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.
		GetDb().
		Delete(models.ProjectMember{}, "project_id = ?", project.ID).
		Error

	return repo
}
