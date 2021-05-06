package utils

import "os"

func GetEnv(varname string, fallback string) string {
	if value, ok := os.LookupEnv(varname); ok {
		return value
	}

	return fallback
}
