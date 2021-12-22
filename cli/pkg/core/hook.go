package core

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
)

type Hook struct {
	Command string
}

var ErrorHookFailed = errors.New("hook failed")

func (h *Hook) Run(ctx *Context) (err error) {
	var output []byte
	cacheDirPath := ctx.DotKeystonePath()
	projectId := ctx.GetProjectID()

	ui.PrintDim("Executing hook '%s'", h.Command)
	output, err = exec.Command(h.Command, projectId, cacheDirPath).Output()
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

func GetHook() (hook *Hook, ok bool) {
	var command string

	if command, ok = config.GetHook(); ok {
		hook = &Hook{Command: command}
	}

	return hook, ok
}

func AddHook(command string) {
	config.AddHook(command)
	config.Write()
}

func RunHookPostFetch(ctx *Context, changes ChangesByEnvironment) {
	if hook, ok := GetHook(); ok {
		shouldRun := false

		for _, c := range changes.Environments {
			if !c.IsEmpty() && !c.IsSingleVersionChange() {
				shouldRun = true
			}
		}

		if shouldRun {
			if utils.FileExists(hook.Command) {
				hook.Run(ctx)
			} else {
				ui.PrintError("Command \"%s\" not found", hook.Command)
			}
		}
	}
}

func RunHookPostSend(ctx *Context) {
	if hook, ok := GetHook(); ok {

		if utils.FileExists(hook.Command) {
			hook.Run(ctx)
		} else {
			ui.PrintError("Command \"%s\" not found", hook.Command)
		}
	}
}
