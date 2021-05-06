// Package repo provides ...
package repo

import (
	"fmt"
	"strings"

	"gorm.io/gorm/clause"

	. "github.com/wearedevx/keystone/api/pkg/models"
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
		Role:    roleAdmin,
		// User:    *user,
		User: *user,
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
	found := false
	if r.err != nil {
		return foundProject, found
	}

	found, err := r.notFoundAsBool(func() error {
		return r.GetDb().
			Model(&Project{}).
			Where("projects.user_id = ? and projects.name = ?", user.ID, name).
			First(&foundProject).
			Error
	})

	r.err = err

	return foundProject, found
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
	// r.err = nil

	return r.createProject(project, &user)
}

func (r *Repo) ProjectGetMembers(project *Project, members *[]ProjectMember) *Repo {
	fmt.Println("project:", project)
	if r.err != nil {
		return r
	}

	r.GetDb().Preload("User").Preload("Role").Where("project_id = ?", project.ID).Find(members)

	return r
}

// From a list of MemberEnvironmentRole, fetches users from database
// Returns the found Users and a slice of not found userIDs
func (r *Repo) usersInMemberRoles(mers []MemberRole) (map[string]User, []string) {
	// Figure out members that do not exist in db
	userIDs := make([]string, 0)
	// only used so that userIDs are unique in the array
	uniqueMap := make(map[string]struct{})

	for _, mer := range mers {
		if _, ok := uniqueMap[mer.MemberID]; !ok {
			uniqueMap[mer.MemberID] = struct{}{}
			userIDs = append(userIDs, mer.MemberID)
		}
	}

	return r.FindUsers(userIDs)
}

func (r *Repo) ProjectAddMembers(project Project, memberRoles []MemberRole) *Repo {
	if r.err != nil {
		return r
	}

	pms := make([]ProjectMember, 0)
	db := r.GetDb()

	users, notFounds := r.usersInMemberRoles(memberRoles)

	if len(notFounds) != 0 {
		r.err = fmt.Errorf("Users not found: %s", strings.Join(notFounds, ", "))
		return r
	}

	for _, memberRole := range memberRoles {
		if memberRole.RoleID != 0 {
			user, hasUser := users[memberRole.MemberID]

			if hasUser {
				pms = append(pms, ProjectMember{
					UserID:    user.ID,
					ProjectID: project.ID,
					RoleID:    memberRole.RoleID,
				})
			}
		}
	}

	r.err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "project_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role_id"}),
	}).Create(&pms).Error

	return r
}

func (r *Repo) ProjectRemoveMembers(project Project, members []string) *Repo {
	if r.err != nil {
		return r
	}

	users, notFound := r.FindUsers(members)

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
