// Copyright © 2015 Philosopher Businessman abp@philosopherbusinessman.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/mindok/hugodeploy/deploy"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// compareCmd represents the preview command
var compareCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview changes that will be pushed to your webhost",
	Long: `Preview allows you to view the changes that would be applied by push.
Preview uses the same comparison algorithms as push to determine what changes
need to be applied and lists those changes.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		checkSourcePath()
		jww.INFO.Println("Preview: Source Dir Good: ", Source)
		checkDeployPath()
		jww.INFO.Println("Preview: Deploy Record Dir Good: ", Deploy)

		deploy.DeployChanges(Source, Deploy, !UnMinify, previewDeployCommandHandler, SkipFiles)
	},
}

func previewDeployCommandHandler(cmd *deploy.DeployCommand) error {
	jww.FEEDBACK.Println("Command: ", cmd.GetCommandDesc(), " : ", cmd.RelPath)
	return nil
}

func init() {
	RootCmd.AddCommand(compareCmd)
}
