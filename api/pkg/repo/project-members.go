package repo

import (
	"encoding/json"

	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) GetProjectMember(user *User, project *Project) (ProjectMember, error) {
	var projectMember ProjectMember

	repo.err = repo.GetDb().Preload("Role").Where("user_id = ? AND project_id = ?", user.ID, project.ID).First(&projectMember).Error

	return projectMember, repo.err
}

func (repo *Repo) CreateProjectMember(user *User, project *Project, role *Role) (ProjectMember, error) {
	envType := ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
		RoleID:    role.ID,
	}
	repo.err = repo.GetDb().Create(&envType).Error
	return envType, repo.err
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (repo *Repo) GetOrCreateProjectMember(project *Project, user *User, roleName string) (ProjectMember, error) {
	// if repo.err != nil {
	// 	fmt.Println("keystone ~ project-members.go ~ repo.errYOUPIIIIIIIIIIIIIIIIIIII", repo.err)
	// 	return ProjectMember{}, repo.err
	// }

	if projectMember, err := repo.GetProjectMember(user, project); err == nil {
		prettyPrint(projectMember)
		return projectMember, nil
	}

	role, _ := repo.GetRoleByName(roleName)

	return repo.CreateProjectMember(user, project, &role)
}
