// Package repo provides ...
package repo

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm/clause"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (r *Repo) createProject(project *Project) IRepo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()
	role := Role{
		Name: "admin",
	}

	r.err = db.Create(project).Error

	r.GetRole(&role)

	if r.err != nil {
		return r
	}

	projectMember := ProjectMember{
		Project: *project,
		Role:    role,
		UserID:  project.UserID,
	}

	r.err = db.Create(&projectMember).Error

	// Useless
	// @Kévin : Care to say why ?
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

		if r.err == nil {
			r.err = db.Preload("Environments").First(project, project.ID).Error
		}
	}

	return r
}

func (r *Repo) getAllEnvironmentTypes(environmentTypes *[]EnvironmentType) IRepo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	r.err = db.Model(EnvironmentType{}).Find(environmentTypes).Error

	return r
}

func (r *Repo) GetProjectByUUID(uuid string, project *Project) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().Where("uuid = ?", uuid).First(project).Error

	return r
}

func (r *Repo) GetProject(project *Project) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("Environments").
		Where(&project).
		First(project).
		Error

	return r
}

func (r *Repo) GetOrCreateProject(project *Project) IRepo {
	if r.Err() != nil {
		return r
	}

	if err := r.GetProject(project).Err(); err != nil {
		if errors.Is(err, ErrorNotFound) {
			r.err = nil
			return r.createProject(project)
		}
	}

	return r
}

// ProjectGetMembers returns all members of a project with
// their role
// TODO: implement paid role restrictions
func (r *Repo) ProjectGetMembers(project *Project, members *[]ProjectMember) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("User").
		Preload("Role").
		Where("project_id = ?", project.ID).
		Find(members).
		Error

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

	users := make(map[string]User)
	notFounds := make([]string, 0)

	r.FindUsers(userIDs, &users, &notFounds)

	return users, notFounds
}

func (r *Repo) ProjectAddMembers(project Project, memberRoles []MemberRole, currentUser models.User) IRepo {
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

	if r.err == nil {
		for _, memberRole := range memberRoles {
			userEmail := users[memberRole.MemberID].Email
			e, err := emailer.AddedMail(currentUser.Email, project.Name)
			if err != nil {
				r.err = err
				return r
			}

			if err = e.Send([]string{userEmail}); err != nil {
				r.err = err
				return r
			}
		}
	}

	return r
}

func (r *Repo) ProjectRemoveMembers(project Project, members []string) IRepo {
	if r.err != nil {
		return r
	}

	users := make(map[string]User)
	notFounds := make([]string, 0)

	r.FindUsers(members, &users, &notFounds)

	if r.err != nil {
		return r
	}

	if len(notFounds) != 0 {
		r.err = fmt.Errorf("Users not found: %s", strings.Join(notFounds, ", "))
		return r
	}

	db := r.GetDb()

	memberIDs := make([]uint, 0)
	for _, user := range users {
		memberIDs = append(memberIDs, user.ID)
	}

	r.err = db.
		Where("user_id IN (?)", memberIDs).
		Where("project_id = ?", project.ID).
		Delete(ProjectMember{}).
		Error

	return r
}

func (r *Repo) ProjectLoadUsers(project *Project) IRepo {
	if r.err != nil {
		return r
	}

	r.GetDb().Model(project).Association("Members.User")

	return r
}

func (r *Repo) ProjectSetRoleForUser(project Project, user User, role Role) IRepo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	pm := ProjectMember{
		Project: project,
		User:    user,
		Role:    role,
	}

	r.err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "project_id"}, {Name: "user_id"}},
		UpdateAll: true,
	}).
		Create(&pm).
		Error

	return r
}

func (r *Repo) CheckMembersAreInProject(project models.Project, members []string) (areInProjects []string, err error) {
	for _, member := range members {
		user := &models.User{UserID: member}

		if r.err = r.GetUser(user).Err(); r.err != nil {
			if errors.Is(r.err, ErrorNotFound) {
				r.err = nil
			}
			return areInProjects, r.err
		}

		projectMember := models.ProjectMember{
			UserID:    user.ID,
			ProjectID: project.ID,
		}

		if r.err = r.GetProjectMember(&projectMember).Err(); r.err == nil {
			areInProjects = append(areInProjects, member)
		} else {
			if errors.Is(r.err, ErrorNotFound) {
				r.err = nil
			}

		}
	}

	return areInProjects, r.err
}

func (r *Repo) DeleteProject(project *models.Project) IRepo {
	if r.Err() != nil {
		return r
	}

	db := r.GetDb()

	r.err = db.
		Delete(
			models.Project{},
			"uuid = ?",
			project.UUID,
		).
		Error

	return r
}
