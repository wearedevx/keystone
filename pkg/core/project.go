package core

import (
	. "github.com/wearedevx/keystone/internal/errors"
	. "github.com/wearedevx/keystone/internal/keystonefile"
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

	ctx.err = FailedToReadKeystoneFile(ksFile.Err())

	return ksFile.ProjectId
}
