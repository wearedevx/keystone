package constants

import "strings"

var Version string

type EnvName string
type EnvNameList []EnvName

const (
	DEV     EnvName = "dev"
	STAGING EnvName = "staging"
	PROD    EnvName = "prod"
)

var EnvList EnvNameList = []EnvName{
	DEV,
	STAGING,
	PROD,
}

// String method formats the envlist for display
func (el EnvNameList) String() string {
	l := make([]string, 0, len(el))

	for i, e := range el {
		l[i] = string(e)
	}

	return strings.Join(l, ", ")
}
