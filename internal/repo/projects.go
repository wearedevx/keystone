// Package repo provides ...
package repo

import (
	"fmt"
	"strings"

	. "github.com/wearedevx/keystone/internal/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) createProject(project *Project, user *User) *Repo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().Create(project).Error

	if r.err == nil {
		envs := []string{"dev", "ci", "staging", "prod"}

		for _, env := range envs {
			environment := Environment{
				Name: env,
			}

			r.err = r.GetDb().Create(&environment).Error

			if r.err != nil {
				break
			}

			projectMember := ProjectMember{
				Project:     *project,
				Environment: environment,
				User:        *user,
				Role:        RoleOwner,
			}

			r.err = r.GetDb().Create(&projectMember).Error

			if r.err != nil {
				break
			}

			r.err = r.GetDb().Preload("Members").First(project, project.ID).Error
		}
	}

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

	r.err = r.GetDb().Model(&Project{}).Joins("join project_members pm on pm.project_id = projects.id").Joins("join users u on pm.user_id = u.id").Where("u.id = ? and name = ? and pm.project_owner = true", user.ID, name).First(&foundProject).Error

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

// From a list of MemberEnvironmentRole, fetches users from database
// Returns the found Users and a slice of not found userIDs
func (r *Repo) usersInMemberEnvironmentsRole(mers []MemberEnvironmentRole) (map[string]User, []string) {
	// Figure out members that do not exist in db
	userIDs := make([]string, 0)
	// only used so that userIDs are unique in the array
	uniqueMap := make(map[string]struct{})

	for _, mer := range mers {
		if _, ok := uniqueMap[mer.ID]; !ok {
			uniqueMap[mer.ID] = struct{}{}
			userIDs = append(userIDs, mer.ID)
		}
	}

	return r.findUsers(userIDs)
}

func (r *Repo) environmentsInMemberEnvironmentsRole(mers []MemberEnvironmentRole) map[string]Environment {
	result := make(map[string]Environment)

	if r.err != nil {
		return result
	}

	db := r.GetDb()

	for _, mer := range mers {
		if _, ok := result[mer.Environment]; !ok {
			var environment Environment

			r.err = db.Where("name = ?", mer.Environment).First(&environment).Error

			if r.err == nil {
				result[mer.Environment] = environment
			}
		}
	}

	return result
}

func (r *Repo) ProjectAddMembers(project Project, mers []MemberEnvironmentRole) *Repo {
	if r.err != nil {
		return r
	}

	pms := make([]ProjectMember, 0)
	db := r.GetDb()

	users, notFounds := r.usersInMemberEnvironmentsRole(mers)

	if len(notFounds) != 0 {
		r.err = fmt.Errorf("Users not found: %s", strings.Join(notFounds, ", "))
		return r
	}

	envs := r.environmentsInMemberEnvironmentsRole(mers)

	for _, mer := range mers {
		if mer.Role != "" {
			user, hasUser := users[mer.ID]
			environment, hasEnv := envs[mer.Environment]

			if hasUser && hasEnv {
				pms = append(pms, ProjectMember{
					UserID:        user.ID,
					EnvironmentID: environment.ID,
					ProjectID:     project.ID,
					ProjectOwner:  false,
					Role:          mer.Role,
				})
			}
		}
	}

	r.err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "project_id"}, {Name: "environment_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role"}),
	}).Create(&pms).Error

	return r
}

func (r *Repo) ProjectRemoveMembers(project Project, members []string) *Repo {
	if r.err != nil {
		return r
	}

	users, notFound := r.findUsers(members)

	if len(notFound) != 0 {
		r.err = fmt.Errorf("Users not found: %s", strings.Join(notFound, ", "))
		return r
	}

	db := r.GetDb()

	memberIDs := make([]uint, 0)
	for _, user := range users {
		memberIDs = append(memberIDs, user.ID)
	}

	r.err = db.Where("user_id IN (?)", memberIDs).Delete(ProjectMember{}).Error

	return r
}

func (r *Repo) ProjectLoadUsers(project *Project) *Repo {
	if r.err != nil {
		return r
	}

	r.GetDb().Model(project).Association("Users")

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

	// r.err = r.GetDb().Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "user_id"}, {Name: "project_id"}},
	// 	DoUpdates: clause.Assignments(map[string]interface{}{"role": role}),
	// }).Create(&perm).Error

	return r
}
