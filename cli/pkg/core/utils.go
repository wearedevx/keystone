package core

import (
	"fmt"
)

func PrintObject(object interface{}) {
	fmt.Printf("%+v\n", object)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
