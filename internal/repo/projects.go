// Package repo provides ...
package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) createProject(project *Project, user *User) *Repo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	project.UserID = user.ID
	r.err = db.Create(project).Error

	roleAdmin, _ := r.GetRoleByName("admin")

	projectMember := ProjectMember{
		Project: *project,
		RoleID:  roleAdmin.ID,
		// User:    *user,
		UserID: user.ID,
	}

	r.err = db.Create(&projectMember).Error

	// Useless
	if r.err == nil {
		envTypes := make([]EnvironmentType, 0)
		r.getAllEnvironmentTypes(&envTypes)

		if r.err != nil {
			return r
		}

		for _, envType := range envTypes {
			environment := Environment{
				Name:            envType.Name,
				EnvironmentType: envType,
				Project:         *project,
				VersionID:       "0",
			}

			r.err = db.Create(&environment).Error

			if r.err != nil {
				break
			}

		}

		r.err = db.Preload("Members").First(project, project.ID).Error
	}

	return r
}

func (r *Repo) getAllEnvironmentTypes(environmentTypes *[]EnvironmentType) *Repo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	r.err = db.Model(EnvironmentType{}).Find(environmentTypes).Error

	return r
}

func (r *Repo) GetProjectByUUID(uuid string, project *Project) *Repo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().Where("uuid = ?", uuid).First(project).Error

	return r
}

func (r *Repo) GetUserProjectWithName(user User, name string) (Project, bool) {
	var foundProject Project

	// Why?
	// if r.err != nil {
	// 	return foundProject, false
	// }

	// r.err = r.GetDb().Model(&Project{}).Joins("join users u on u.id = projects.user_id").Where("u.id = ? and projects.name = ?", user.ID, name).First(&foundProject).Error
	r.err = r.GetDb().Where("user_id = ? and name = ?", 1, name).First(&foundProject).Error

	return foundProject, r.err == nil
}

func (r *Repo) GetOrCreateProject(project *Project, user User) *Repo {
	// if r.err != nil {
	// 	fmt.Println("Error")
	// 	return r
	// }

	if foundProject, ok := r.GetUserProjectWithName(user, project.Name); ok == true {
		*project = foundProject
		return r
	}
	r.err = nil

	return r.createProject(project, &user)
}

func (r *Repo) ProjectLoadUsers(project *Project) *Repo {
	if r.err != nil {
		return r
	}

	r.GetDb().Model(project).Association("Members.User")

	return r
}

func (r *Repo) ProjectSetRoleForUser(project Project, user User, role Role) *Repo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	pm := ProjectMember{
		Project: project,
		User:    user,
		Role:    role,
	}

	r.err = db.Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&pm).
		Error

	return r
}
