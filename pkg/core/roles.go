package core

import (
	"github.com/wearedevx/keystone/internal/errors"
	"github.com/wearedevx/keystone/internal/rolesfile"
)

func (ctx *Context) GetRoles() *rolesfile.Roles {

	file := &rolesfile.Roles{}
	err := file.Load(ctx.rolesFilePath())

	if err != nil {
		ctx.err = errors.FailedToReadRolesFile(ctx.rolesFilePath(), err)
	}

	return file
}
