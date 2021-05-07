package repo

import (
	"encoding/json"
	"errors"

	. "github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func (repo *Repo) GetProjectMember(user *User, project *Project, projectMember *ProjectMember) *Repo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().Preload("Role").Where("user_id = ? AND project_id = ?", user.ID, project.ID).First(projectMember).Error

	return repo
}

func (repo *Repo) CreateProjectMember(user *User, project *Project, role *Role, projectMember *ProjectMember) *Repo {
	if repo != nil {
		return repo
	}

	projectMember.UserID = user.ID
	projectMember.ProjectID = project.ID
	projectMember.RoleID = role.ID

	repo.err = repo.GetDb().Create(&projectMember).Error

	return repo
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (repo *Repo) GetOrCreateProjectMember(project *Project, user *User, roleName string, projectMember *ProjectMember) *Repo {
	if repo.Err() != nil {
		return repo
	}

	repo.GetProjectMember(user, project, projectMember)

	if err := repo.Err(); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			repo.err = err
			return repo
		} else {
			// reset error to not block
			// the creation operation
			repo.err = nil
			role := Role{}

			repo.GetRoleByName(roleName, &role).
				CreateProjectMember(user, project, &role, projectMember)
		}
	}

	return repo
}
