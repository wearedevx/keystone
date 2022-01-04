package core

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/wearedevx/keystone/cli/ui"
)

type Hook struct {
	ctx     *Context
	Command string
}

var ErrorHookFailed = errors.New("hook failed")

// Run method executes the hook.
// Note: it displays the output of the command to stdout in casse of success
// and to stderr in case of failure
func (h *Hook) Run() (err error) {
	var output []byte
	projectPath := h.ctx.Wd
	projectId := h.ctx.GetProjectID()
	projectName := h.ctx.GetProjectName()

	ui.PrintDim("Executing hook '%s'", h.Command)
	output, err = exec.Command(h.Command, projectName, projectId, projectPath).Output()
	printer := ui.PrintDim

	if err != nil {
		ui.PrintStdErr("Error Executing hook:")

		output = err.(*exec.ExitError).Stderr
		err = ErrorHookFailed
		printer = ui.PrintStdErr
	}

	for _, line := range strings.Split(string(output), "\n") {
		printer("> %s", line)
	}

	return err
}
