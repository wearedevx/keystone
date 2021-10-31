// +build test

package repo

func (r *Repo) GetAdminsFromUserProjects(
	userID uint,
	adminProjectsMap *map[string][]string,
) IRepo {
	rows, err := r.GetDb().Raw(`
SELECT
	u.email,
	p.name
FROM
	users u
JOIN
	project_members pm
	ON pm.user_id = u.id
JOIN
	roles r
	ON r.id = pm.role_id
	AND r.name = 'admin'
JOIN
	projects p
	ON pm.project_id = p.id
WHERE
	p.id IN (
		SELECT
			p2.id
		FROM
			projects p2
		JOIN
			project_members pm2
			ON pm2.project_id = p2.id
			AND pm2.user_id = ?
	)
AND u.id <> ?`,
		userID,
		userID,
	).Rows()
	if err != nil {
		r.err = err
		return r
	}

	*adminProjectsMap = make(map[string][]string)
	var mail string
	var project string
	for rows.Next() {
		if err = rows.Scan(&mail, &project); err != nil {
			r.err = err
			return r
		}

		insertInMap(*adminProjectsMap, mail, project)
	}

	return r
}

func insertInMap(m map[string][]string, email, project string) {
	list, ok := m[email]
	if !ok {
		list = make([]string, 0)
	}

	m[email] = append(list, project)
}
