package display

import (
	"github.com/wearedevx/keystone/cli/ui"
)

func HookCommand(hook string) {
	ui.Print(hook)
}

func HookPathDoesNotExist(path string) {
	ui.PrintError("%s does not exist", path)
}

func HookAddedSuccessfully() {
	ui.PrintSuccess("Hook added successfully")
}

func ThereIsNoHookYet() {
	ui.Print("You have not registered a hook yet. To add one, try `ks hook add <path-to-a-script>`")
}

func ExecutingHook(hook string) {
	ui.PrintDim("Executing hook '%s'", hook)
}
