/*
 Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/wearedevx/keystone/api/pkg/models"
	kerrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"

	// . "github.com/wearedevx/keystone/cli/ui"

	"github.com/spf13/cobra"
)

var member string

// initCmd represents the init command
var sendEnvCmd = &cobra.Command{
	Use:   "send-env <member id>",
	Short: "Send current environment to member.",
	Long: `Send secrets and files from current environment to member.
If a member hasn't received secrets and files last time someone sent an update, it can be done again with this command.
`,
	Example: `ks member send-env john@gitlab`,

	Args: func(_ *cobra.Command, args []string) error {
		r := regexp.MustCompile(`[\w-_.]+@(gitlab|github)`)
		argc := len(args)

		if argc == 0 || argc > 1 {
			return fmt.Errorf("invalid number of arguments. Expected 1, got %d", argc)
		}

		member = args[0]

		if !r.Match([]byte(member)) {
			return fmt.Errorf("invalid member id: %s", member)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err *kerrors.Error

		ctx := core.New(core.CTX_RESOLVE)

		ctx.MustHaveEnvironment(currentEnvironment)

		currentEnvironment = ctx.CurrentEnvironment()

		messagesByEnvironment := &models.GetMessageByEnvironmentResponse{
			Environments: map[string]models.GetMessageResponse{},
		}

		fmt.Println("Syncing data...")
		fetchErr := ctx.FetchNewMessages(messagesByEnvironment)

		if fetchErr != nil {
			err.SetCause(fetchErr)
			err.Print()
		}

		_, writeErr := ctx.WriteNewMessages(*messagesByEnvironment)

		if writeErr != nil {
			writeErr.Print()
			// err.SetCause(writeErr)
			// err.Print()
			return
		}

		localEnvironment := ctx.LoadEnvironmentsFile().GetByName(currentEnvironment)
		environment := models.Environment{
			Name:          localEnvironment.Name,
			VersionID:     localEnvironment.VersionID,
			EnvironmentID: localEnvironment.EnvironmentID,
		}

		if pushErr := ctx.PushEnvForOneMember(environment, member); pushErr != nil {
			ui.PrintError(pushErr.Error())
			return
		}

		ui.PrintSuccess("Environment sent to user.")
		// Retrieve working directry

		// Print(RenderTemplate("Environment push", `
		// {{ .Message | box | bright_green | indent 2 }}

		// {{ .Text | bright_black | indent 2 }}`, map[string]string{
		// 	// "Message": "All done!",
		// 	// "Text": `You can start adding environment variable with:
		// 	//   $ ks secrets add VARIABLE value

		// 	// Load them with:
		// 	//   $ eval $(ks source)

		// 	// If you need help with anything:
		// 	//   $ ks help [command]

		// 	// `,
		// }))

	},
}

func init() {
	memberCmd.AddCommand(sendEnvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	sendEnvCmd.Flags().StringVar(&member, "all", "a", "Member to send env to.")
}
