// Package repo provides ...
package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
)

func (r *Repo) createProject(project *Project, user *User) *Repo {
	if r.err != nil {
		return r
	}

	r.err = r.db.Create(project).Error

	if r.err == nil {
		envs := []string{"dev", "ci", "staging", "prod"}

		for _, env := range envs {
			environment := Environment{
				Name: env,
			}

			r.err = r.db.Create(&environment).Error

			if r.err != nil {
				break
			}

			projectMember := ProjectMember{
				Project:     *project,
				Environment: environment,
				User:        *user,
				Role:        RoleOwner,
			}

			r.err = r.db.Create(&projectMember).Error

			if r.err != nil {
				break
			}

			r.err = r.db.Preload("Members").First(project, project.ID).Error
		}
	}

	return r
}

func (r *Repo) GetProjectByUUID(uuid string, project *Project) *Repo {
	if r.err != nil {
		return r
	}

	r.err = r.db.Where("uuid = ?", uuid).First(project).Error

	return r
}

func (r *Repo) getUserProjectWithName(user User, name string) (Project, bool) {
	var foundProject Project
	if r.err != nil {
		return foundProject, false
	}

	r.err = r.db.Model(&Project{}).Joins("join project_members pm on pm.project_id = id").Joins("join users u on pp.user_id = u.id").Where("u.id = ? and name = ? and pm.project_owner = true", user.ID, name).First(&foundProject).Error

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

	return r.createProject(project)
}

func (r *Repo) ProjectLoadUsers(project *Project) *Repo {
	if r.err != nil {
		return r
	}

	r.db.Model(project).Association("Users")

	return r
}

func (r *Repo) ProjectSetRoleForUser(project Project, user User, role UserRole) *Repo {
	if r.err != nil {
		return r
	}

	// perm := ProjectPermissions{
	// 	UserID:    user.ID,
	// 	ProjectID: project.ID,
	// 	Role:      role,
	// }

	// r.err = r.db.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "user_id"}, {Name: "project_id"}},
	// 	DoUpdates: clause.Assignments(map[string]interface{}{"role": role}),
	// }).Create(&perm).Error

	return r
}
