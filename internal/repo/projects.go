// Package repo provides ...
package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
)

func (r *Repo) createProject(project *Project) *Repo {
	r.err = r.db.Create(project).Error

	return r
}

func (r *Repo) getProjectByUUID(uuid string) (Project, bool) {
	var foundProject Project

	r.err = r.db.Where("uuid = ?", uuid).First(&foundProject).Error

	return foundProject, r.err == nil
}

func (r *Repo) getUserProjectWithName(user User, name string) (Project, bool) {
	var foundProject Project

	r.err = r.db.Model(&Project{}).Joins("join project_permissions pp on pp.project_id = id").Joins("join user u on pp.user_id = u.id").Where("u.id = ? and name = ?", user.ID, name).First(&fofoundProject).Error

	return foundProject, r.err == nil
}

func (r *Repo) GetOrCreateProject(project *Project, user User) {
	var foundProject Project

	if foundProject, ok := r.getUserProjectWithName(user, project.Name); ok == true {
		*project = foundProject
		return
	}

	r.createProject(&project)
}

func (r *Repo) ProjectSetRoleForUser(project Project, user User, role UserRole) {
	perm := ProjectPermissions{
		UserID:    user.ID,
		ProjectID: project.ID,
		role:      role,
	}

	r.err = r.db.Create(perm).Error
}
