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

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Deploy changes website updates to host",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkSourcePath()
		jww.INFO.Println("Push: Source Dir Good: ", Source)
		checkDeployPath()
		jww.INFO.Println("Push: Deploy Record Dir Good: ", Deploy)

		var err error

		ftpDeployer = &deploy.FTPDeployer{}
		if err = ftpDeployer.Initialise(); err != nil {
			panic(err)
		}

		deployRecorder = &deploy.FileDeployer{Deploy}
		deployRecorder.Initialise()

		deploy.DeployChanges(Source, Deploy, !UnMinify, pushDeployCommandHandler, SkipFiles)

		ftpDeployer.Cleanup()
		deployRecorder.Cleanup()
	},
}

var ftpDeployer *deploy.FTPDeployer
var deployRecorder *deploy.FileDeployer

func pushDeployCommandHandler(cmd *deploy.DeployCommand) error {
	err := ftpDeployer.ApplyCommand(cmd)
	//err := deployRecorder.ApplyCommand(cmd)
	if err == nil {
		err = deployRecorder.ApplyCommand(cmd)
	}
	return err
}

func init() {
	RootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	pushCmd.PersistentFlags().String("ftppwd", "", "FTP Password. Avoids having to set it in config file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
