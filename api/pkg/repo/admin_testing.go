// +build test

package repo

import (
	"regexp"
	"strings"
)

func (r *Repo) GetAdminsFromUserProjects(userID uint, userName string, projects_list []string, adminEmail string) IRepo {
	// Get project on which user is present
	rows, err := r.GetDb().Raw(`
	SELECT u.email, group_concat(p.name) FROM users u
	LEFT join project_members pm on pm.user_id = u.id
	LEFT join roles r on r.id = pm.role_id
	LEFT join projects p on pm.project_id = p.id
	where r.name = 'admin' and p.id in (
	select pm.project_id from project_members pm where pm.user_id = ?) and u.user_id != ?
	group by u.user_id, u.email;`, userID, userName).Rows()

	if err != nil {
		r.err = err
		return r
	}

	var projects string
	for rows.Next() {
		rows.Scan(&adminEmail, &projects)
		re := regexp.MustCompile(`\{(.+)?\}`)
		res := re.FindStringSubmatch(projects)

		projects_list = strings.Split(res[1], ",")
	}
	return r
}
