package display

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

// Logs function displags all the logs
func Logs(logs []models.ActivityLogLite) {
	if len(logs) == 0 {
		ui.PrintStdErr("No logs to display")
		return
	}

	for _, log := range logs {
		printLog(log)
	}
}

func printLog(log models.ActivityLogLite) {
	fmt.Printf(
		"[%s] %s on %s%s: %s %s %s\n",
		log.CreatedAt.Format("2006-12-29 15:04:05"),
		log.UserID,
		log.ProjectName,
		formatEnvironmentForLog(log.EnvironmentName),
		formatSuccesForLog(log.Success),
		log.Action,
		log.ErrorMessage,
	)
}

func formatEnvironmentForLog(envName string) string {
	if envName == "" {
		return ""
	}

	return fmt.Sprintf(" (%s)", envName)
}

func formatSuccesForLog(s bool) string {
	if s {
		return ""
	} else {
		return "✘"
	}
}
