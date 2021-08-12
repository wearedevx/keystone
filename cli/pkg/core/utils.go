package core

import (
	"fmt"
)

func PrintObject(object interface{}) {
	fmt.Printf("%+v\n", object)
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func Uniq(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}

	for e := range arr {
		if occured[arr[e]] != true {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result

}
