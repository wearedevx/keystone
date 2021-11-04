package core

import (
	"fmt"
)

// PrintObject function prints an object. For debugging.
func PrintObject(object interface{}) {
	fmt.Printf("%+v\n", object)
}

// Contains function tells if a slice of string contains a string
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// Uniq function removes duplicatos form a slice of strings
func Uniq(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}

	for e := range arr {
		if !occured[arr[e]] {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result

}
