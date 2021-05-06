/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var memberId string
var role string

// memberSetRoleCmd represents the memberSetRole command
var memberSetRoleCmd = &cobra.Command{
	Use:   "set-role <member id> [role]",
	Short: "Sets the role for a member",
	Long: `Sets the role for a member.
If no role argument is provided, it will be prompted.

Roles determine access rights to environments.`,
	Example: `ks member set-role john@gitlab devops

ks member set-role sandra@github`,
	Args: func(cmd *cobra.Command, args []string) error {
		r := regexp.MustCompile("[\\w-_.]+@(gitlab|github)")
		argc := len(args)

		if argc == 0 || argc > 2 {
			return fmt.Errorf("invalid number of arguments. Expected 1 or 2, got %d", argc)
		}

		if argc >= 1 {
			memberId = args[0]
		}

		if argc == 2 {
			role = args[2]
		}

		if !r.Match([]byte(memberId)) {
			return fmt.Errorf("invalid member id: %s", memberId)
		}

		// TODO: check role is valid

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("memberSetRole called")
	},
}

func init() {
	memberCmd.AddCommand(memberSetRoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memberSetRoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memberSetRoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
