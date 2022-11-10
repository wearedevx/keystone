//go:build !test
// +build !test

package seed

import "gorm.io/gorm"

func Seed(_ *gorm.DB) error {
	return nil
}
