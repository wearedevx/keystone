package core

import (
	. "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/keystonefile"
)

func (ctx *Context) GetProjectName() string {
	if ctx.err != nil {
		return ""
	}

	ksFile := &KeystoneFile{}
	ksFile.Load(ctx.Wd)

	ctx.err = FailedToReadKeystoneFile(ksFile.Err())

	return ksFile.ProjectName
}

func (ctx *Context) GetProjectID() string {
	if ctx.err != nil {
		return ""
	}

	ksFile := &KeystoneFile{}
	ksFile.Load(ctx.Wd)

	if ksFile.Err() != nil {
		ctx.err = FailedToReadKeystoneFile(ksFile.Err())
	}

	return ksFile.ProjectId
}

func (ctx *Context) MustHaveProject() {
	projectID := ctx.GetProjectID()

	if projectID == "" {
		ctx.err = CannotFindProjectID(nil)
	}
}
