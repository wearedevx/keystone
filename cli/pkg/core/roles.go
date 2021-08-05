package core

import (
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/rolesfile"
)

func (ctx *Context) GetRoles() *rolesfile.Roles {

	file := &rolesfile.Roles{}
	err := file.Load(ctx.rolesFilePath())

	if err != nil {
		ctx.err = kserrors.FailedToReadRolesFile(ctx.rolesFilePath(), err)
	}

	return file
}
