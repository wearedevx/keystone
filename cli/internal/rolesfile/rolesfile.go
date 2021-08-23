package rolesfile

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func DefaultContent() string {
	return `roles:
developer:
    dev: write
    ci: ""
    staging: ""
    prod: ""
devops:
    dev: ""
    ci: write
    staging: write
    prod: write
admin:
    dev: owner
    ci: owner
    staging: owner
    prod: owner
`
}

type Roles struct {
	Roles map[string]map[string]string
}

func (r *Roles) Load(path string) error {
	/* #nosec */
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(contents, r)
}

func (r *Roles) List() []string {
	l := make([]string, 0)

	for roleName, _ := range r.Roles {
		l = append(l, roleName)
	}

	return l
}

func (r *Roles) GetRights(roleName string) (map[string]string, bool) {
	rights, ok := r.Roles[roleName]

	return rights, ok
}
