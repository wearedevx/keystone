package core

import (
	"os"

	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/keystonefile"
)

func (ctx *Context) GetProjectName() string {
	if ctx.err != nil {
		return ""
	}

	ksFile := &KeystoneFile{}
	ksFile.Load(ctx.Wd)

	ctx.err = kserrors.FailedToReadKeystoneFile(ksFile.Err())

	return ksFile.ProjectName
}

func (ctx *Context) GetProjectID() string {
	if ctx.err != nil {
		return ""
	}

	ksFile := &KeystoneFile{}
	ksFile.Load(ctx.Wd)

	if ksFile.Err() != nil {
		ctx.err = kserrors.FailedToReadKeystoneFile(ksFile.Err())
		return ""
	}

	return ksFile.ProjectId
}

func (ctx *Context) MustHaveProject() {
	projectID := ctx.GetProjectID()

	if projectID == "" {
		ctx.err = kserrors.CannotFindProjectID(nil)
	}
}

// Removes the keystone.yaml, and the .keystone file
func (ctx *Context) Destroy() error {
	var err error

	if err = new(KeystoneFile).
		Load(ctx.Wd).
		Remove().
		Err(); err != nil {
		return err
	}

	if err = os.RemoveAll(ctx.dotKeystonePath()); err != nil {
		return err
	}

	return nil
}
