// Package repo provides ...
package repo

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (r *Repo) createProject(project *models.Project) IRepo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()
	role := models.Role{
		Name: "admin",
	}

	r.err = db.Create(project).Error

	r.GetRole(&role)

	if r.err != nil {
		return r
	}

	projectMember := models.ProjectMember{
		Project: *project,
		Role:    role,
		UserID:  project.UserID,
	}

	r.err = db.Create(&projectMember).Error

	// Useless
	// @Kévin : Care to say why ?
	if r.err == nil {
		envTypes := make([]models.EnvironmentType, 0)
		r.getAllEnvironmentTypes(&envTypes)

		if r.err != nil {
			return r
		}

		for _, envType := range envTypes {
			environment := models.Environment{
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

func (r *Repo) getAllEnvironmentTypes(
	environmentTypes *[]models.EnvironmentType,
) IRepo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	r.err = db.Model(models.EnvironmentType{}).Find(environmentTypes).Error

	return r
}

func (r *Repo) GetProjectByUUID(uuid string, project *models.Project) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().Where("uuid = ?", uuid).First(project).Error

	return r
}

func (r *Repo) GetProject(project *models.Project) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("Environments").
		Preload("Organization").
		Where(&project).
		First(project).
		Error

	return r
}

func (r *Repo) GetOrCreateProject(project *models.Project) IRepo {
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
func (r *Repo) ProjectGetMembers(
	project *models.Project,
	members *[]models.ProjectMember,
) IRepo {
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

// ProjectGetAdmins returns all admins of a project
func (r *Repo) ProjectGetAdmins(
	project *models.Project,
	members *[]models.ProjectMember,
) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("User").
		Joins(
			"inner join roles as r on role_id = r.id and r.name = ?",
			"admin",
		).
		Where("project_id = ?", project.ID).
		Find(members).
		Error

	return r
}

func (r *Repo) ProjectIsMemberAdmin(
	project *models.Project,
	member *models.ProjectMember,
) bool {
	if r.err != nil {
		return false
	}

	err := r.GetDb().
		Joins(
			"inner join users as u on u.id = ?",
			member.UserID,
		).
		Joins(
			"inner join roles as r on role_id = r.id and r.name = ?",
			"admin",
		).
		Where("project_id = ?", project.ID).
		First(member).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.err = err
		}

		return false
	}

	return true
}

func (r *Repo) IsMemberOfProject(
	project *models.Project,
	member *models.ProjectMember,
) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Joins(
			"inner join users as u on u.id = ?",
			member.UserID,
		).
		Where("project_id = ?", project.ID).
		First(member).
		Error

	return r
}

// From a list of MemberEnvironmentRole, fetches users from database
// Returns the found Users and a slice of not found userIDs
func (r *Repo) usersInMemberRoles(
	mers []models.MemberRole,
) (map[string]models.User, []string) {
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

	users := make(map[string]models.User)
	notFounds := make([]string, 0)

	r.FindUsers(userIDs, &users, &notFounds)

	return users, notFounds
}

func (r *Repo) ProjectAddMembers(
	project models.Project,
	memberRoles []models.MemberRole,
	currentUser models.User,
) IRepo {
	if r.err != nil {
		return r
	}

	pms := make([]models.ProjectMember, 0)
	db := r.GetDb()

	users, notFounds := r.usersInMemberRoles(memberRoles)

	if len(notFounds) != 0 {
		r.err = fmt.Errorf("users not found: %s", strings.Join(notFounds, ", "))
		return r
	}

	for _, memberRole := range memberRoles {
		if memberRole.RoleID != 0 {
			user, hasUser := users[memberRole.MemberID]

			if hasUser {
				pms = append(pms, models.ProjectMember{
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
			e, err := emailer.AddedMail(currentUser, project.Name)
			if err != nil {
				r.err = err
				return r
			}

			if err = e.Send([]string{userEmail}); err != nil {
				fmt.Printf("Project Add Member err: %+v\n", err)
				r.err = err
				return r
			}
		}
	}

	return r
}

func (r *Repo) ProjectRemoveMembers(
	project models.Project,
	members []string,
) IRepo {
	if r.err != nil {
		return r
	}

	users := make(map[string]models.User)
	notFounds := make([]string, 0)

	r.FindUsers(members, &users, &notFounds)

	if r.err != nil {
		return r
	}

	if len(notFounds) != 0 {
		r.err = fmt.Errorf("users not found: %s", strings.Join(notFounds, ", "))
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
		Delete(models.ProjectMember{}).
		Error

	return r
}

func (r *Repo) ProjectLoadUsers(project *models.Project) IRepo {
	if r.err != nil {
		return r
	}

	r.GetDb().Model(project).Association("Members.User")

	return r
}

func (r *Repo) ProjectSetRoleForUser(
	project models.Project,
	user models.User,
	role models.Role,
) IRepo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()

	pm := models.ProjectMember{
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

func (r *Repo) CheckMembersAreInProject(
	project models.Project,
	members []string,
) (areInProjects []string, err error) {
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
			} else {
				break
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

func (r *Repo) GetUserProjects(userID uint, projects *[]models.Project) IRepo {
	if err := r.GetDb().Joins("left join project_members pm on projects.ID = pm.project_id").Where("pm.user_id = ?", userID).Find(&projects).Error; err != nil {
		r.err = err
	}

	return r
}

func (r *Repo) GetProjectsOrganization(
	projectID string,
	organization *models.Organization,
) IRepo {
	var project models.Project

	r.GetProjectByUUID(projectID, &project)
	if r.err != nil {
		return r
	}

	organization.ID = project.OrganizationID

	r.GetDb().First(&organization)

	return r
}

func (r *Repo) IsProjectOrganizationPaid(projectID string) (bool, error) {
	if projectID == "" {
		return false, nil
	}

	project := models.Project{UUID: projectID}

	if err := r.GetProject(&project).Err(); err != nil {
		return false, err
	}

	if project.Organization.Paid {
		return true, nil
	}

	return false, nil
}
