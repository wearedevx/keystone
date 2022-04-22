package core

import (
	"os"

	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
)

// GetProjectName method returns the project name from the keysotnefile
func (ctx *Context) GetProjectName() string {
	if ctx.err != nil {
		return ""
	}

	ksFile := &keystonefile.KeystoneFile{}
	ksFile.Load(ctx.Wd)

	ctx.err = kserrors.FailedToReadKeystoneFile(ksFile.Path, ksFile.Err())

	return ksFile.ProjectName
}

// GetProjectID method returns the project ID from the keystonefile
func (ctx *Context) GetProjectID() string {
	if ctx.err != nil {
		return ""
	}

	ksFile := &keystonefile.KeystoneFile{}
	ksFile.Load(ctx.Wd)

	if ksFile.Err() != nil {
		ctx.err = kserrors.FailedToReadKeystoneFile(ksFile.Path, ksFile.Err())
		return ""
	}

	return ksFile.ProjectId
}

// MustHaveProject method sets a context error if the keystonefile does
// not have a project id
func (ctx *Context) MustHaveProject() {
	projectID := ctx.GetProjectID()

	if projectID == "" {
		ctx.err = kserrors.CannotFindProjectID(nil)
	}
}

// Removes the keystone.yaml, and the .keystone file
func (ctx *Context) Destroy() error {
	var err error

	if err = new(keystonefile.KeystoneFile).
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
