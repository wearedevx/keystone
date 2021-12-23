// +build test

package seed

import (
	"fmt"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) (err error) {
	fmt.Println("Seed")

	err = db.Exec(Sql).Error
	if err != nil {
		panic(err)
	}
	return nil
}

const Sql = `
insert or ignore into roles
(id, name,        can_add_member, parent_id, created_at,        updated_at)
values 
(1,  "admin",     true,           null,      current_timestamp, current_timestamp),
(2,  "devops",    true,           1,         current_timestamp, current_timestamp),
(3,  "lead-dev",  true,           2,         current_timestamp, current_timestamp),
(4,  "developer", false,          3,         current_timestamp, current_timestamp);

insert or ignore into environment_types
(id, name,      created_at,        updated_at)
values
(1,  "dev",     current_timestamp, current_timestamp),
(2,  "staging", current_timestamp, current_timestamp),
(3,  "prod",    current_timestamp, current_timestamp);

insert or ignore into roles_environment_types
(id, role_id, environment_type_id, name, read,  write, created_at,        updated_at)
values
-- admin
(1,  1,       1,                   "",   true,  true,  current_timestamp, current_timestamp),
(2,  1,       2,                   "",   true,  true,  current_timestamp, current_timestamp),
(3,  1,       3,                   "",   true,  true,  current_timestamp, current_timestamp),
-- devops
(4,  2,       1,                   "",   true,  true,  current_timestamp, current_timestamp),
(5,  2,       2,                   "",   true,  true,  current_timestamp, current_timestamp),
(6,  2,       3,                   "",   true,  true,  current_timestamp, current_timestamp),
-- lead-dev
(7,  3,       1,                   "",   true,  true,  current_timestamp, current_timestamp),
(8,  3,       2,                   "",   false, false, current_timestamp, current_timestamp),
(9,  3,       3,                   "",   false, false, current_timestamp, current_timestamp),
-- dev
(10, 4,       1,                   "",   true,  true,  current_timestamp, current_timestamp),
(11, 4,       2,                   "",   false, false, current_timestamp, current_timestamp),
(12, 4,       3,                   "",   false, false, current_timestamp, current_timestamp);
`

/* func seedRoles(Repo repo.IRepo) error {
	devRole := models.Role{Name: "developer"}
	leadDevRole := models.Role{Name: "lead-dev", CanAddMember: true}
	devopsRole := models.Role{Name: "devops", CanAddMember: true}
	adminRole := models.Role{Name: "admin", CanAddMember: true}

	if err := Repo.GetOrCreateRole(&adminRole).Err(); err != nil {
		return err
	}

	devopsRole.ParentID = adminRole.ID
	if err := Repo.GetOrCreateRole(&devopsRole).Err(); err != nil {
		return err
	}

	leadDevRole.ParentID = devopsRole.ID
	if err := Repo.GetOrCreateRole(&leadDevRole).Err(); err != nil {
		return err
	}

	devRole.ParentID = leadDevRole.ID
	if err := Repo.GetOrCreateRole(&devRole).Err(); err != nil {
		return err
	}

	devEnvironmentType := models.EnvironmentType{Name: "dev"}
	stagingEnvironmentType := models.EnvironmentType{Name: "staging"}
	prodEnvironmentType := models.EnvironmentType{Name: "prod"}

	if err := Repo.
		GetOrCreateEnvironmentType(&devEnvironmentType).
		GetOrCreateEnvironmentType(&stagingEnvironmentType).
		GetOrCreateEnvironmentType(&prodEnvironmentType).
		Err(); err != nil {
		return err
	}

	// developer
	devRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            devRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	devRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            devRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}
	devRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            devRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}

	// lead
	leadDevRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            leadDevRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	leadDevRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            leadDevRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}
	leadDevRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            leadDevRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}

	// devops
	devopsRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            devopsRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	devopsRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            devopsRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	devopsRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            devopsRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}

	// admin
	adminRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            adminRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	adminRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            adminRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	adminRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            adminRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}

	if err := Repo.
		GetOrCreateRoleEnvType(&devRoleDevEnvType).
		GetOrCreateRoleEnvType(&devRoleStagingEnvType).
		GetOrCreateRoleEnvType(&devRoleProdEnvType).
		GetOrCreateRoleEnvType(&leadDevRoleDevEnvType).
		GetOrCreateRoleEnvType(&leadDevRoleStagingEnvType).
		GetOrCreateRoleEnvType(&leadDevRoleProdEnvType).
		GetOrCreateRoleEnvType(&devopsRoleDevEnvType).
		GetOrCreateRoleEnvType(&devopsRoleStagingEnvType).
		GetOrCreateRoleEnvType(&devopsRoleProdEnvType).
		GetOrCreateRoleEnvType(&adminRoleDevEnvType).
		GetOrCreateRoleEnvType(&adminRoleStagingEnvType).
		GetOrCreateRoleEnvType(&adminRoleProdEnvType).
		Err(); err != nil {
		return err
	}

	fmt.Println("Things are seeded ?")

	return nil
} */
