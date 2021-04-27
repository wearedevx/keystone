// Package repo provides ...
package repo

import (
	"fmt"

	. "github.com/wearedevx/keystone/internal/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) createProject(project *Project, user *User) *Repo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	r.err = db.Create(project).Error

	projectMember := ProjectMember{
		Project: *project,
		User:    *user,
	}

	r.err = db.Create(&projectMember).Error

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

func (r *Repo) getUserProjectWithName(user User, name string) (Project, bool) {
	var foundProject Project
	if r.err != nil {
		return foundProject, false
	}

	r.err = r.GetDb().Model(&Project{}).Joins("join users u on u.id = projects.user_id").Where("u.id = ? and projects.name = ?", user.ID, name).First(&foundProject).Error

	return foundProject, r.err == nil
}

func (r *Repo) GetOrCreateProject(project *Project, user User) *Repo {
	if r.err != nil {
		return r
	}

	if foundProject, ok := r.getUserProjectWithName(user, project.Name); ok == true {
		*project = foundProject
		return r
	}
	r.err = nil

	return r.createProject(project, &user)
}

func (r *Repo) ProjectGetMembers(project *Project, members *[]ProjectMember) *Repo {
	fmt.Println("project:", project)
	if r.err != nil {
		return r
	}

	r.GetDb().Preload("Environment").Preload("User").Where("project_id = ?", project.ID).Find(members)

	return r
}

func (r *Repo) ProjectAddMembers(project Project, mers []MemberEnvironmentRole) *Repo {
	if r.err != nil {
		return r
	}

	pms := make([]ProjectMember, 0)
	db := r.GetDb()

	for _, mer := range mers {
		if mer.Role != "" {
			var user User
			var environment Environment

			user.FromId(mer.ID)

			db.Where("name = ?", mer.Environment).First(&environment)
			db.Where("username = ? AND account_type = ?", user.Username, user.AccountType).First(&user)

			pms = append(pms, ProjectMember{
				UserID:        user.ID,
				EnvironmentID: environment.ID,
				ProjectID:     project.ID,
				ProjectOwner:  false,
				Role:          mer.Role,
			})
		}
	}

	r.err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "project_id"}, {Name: "environment_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role"}),
	}).Create(&pms).Error

	return r
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
