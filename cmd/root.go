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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hugodeploy",
	Short: "HugoDeploy deploys files to a webserver using SFTP",
	Long: `HugoDeploy tracks changes made to a local directory and
transfers those changes to a remote server. Currently the transfer
is done using SFTP.

HugoDeploy can generate a list of changed files using the preview
command and does the actual transfer using the push command.

By default HugoDeploy minifies files prior to transfer.

HugoDeploy requires a configuration file to nominate directories and
setup server connection details. By default the configuration file is
expected in the directory HugoDeploy is executed from.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//TODO: Probably should split out per command - init doesn't need everything listed
	initCoreCommonFlags(RootCmd)
}

var CfgFile, Source, Deploy string
var Verbose, Debug, UnMinify bool
var SkipFiles []string

func initCoreCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is path/hugodeploy.yaml|json|toml)")
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "debug trace output")
	//cmd.PersistentFlags().BoolVar(&Logging, "log", false, "Enable Logging")
	//cmd.PersistentFlags().StringVar(&LogFile, "logFile", "", "Log File path (if set, logging enabled automatically)")
	//cmd.PersistentFlags().BoolVar(&VerboseLog, "verboseLog", false, "verbose logging")

	cmd.PersistentFlags().StringVarP(&Source, "sourceDir", "s", "", "filesystem path to read files relative from")
	cmd.PersistentFlags().StringVarP(&Deploy, "deployRecordDir", "r", "", "filesystem path to keep a record of what has been deployed")
	cmd.PersistentFlags().BoolVarP(&UnMinify, "dontminify", "m", false, "disable minify")
}

func LoadDefaultSettings() {
	viper.SetDefault("sourceDir", "publish")
	viper.SetDefault("deployRecordDir", "deployed")
	viper.SetDefault("dontminify", false)
	viper.SetDefault("verbose", false)
	viper.SetDefault("debug", false)
	viper.SetDefault("skipfiles", []string{".git*", ".DS_Store"})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if CfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(CfgFile)
	}

	viper.SetConfigName("hugodeploy") // name of config file (without extension)
	viper.AddConfigPath(".")          // adding cwd directory as first search path
	//viper.AddConfigPath("/Users/johnjessop/Documents/Code/GoCode/src/github.com/mindok/hugodeploy")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Error: No valid config file found. ", err)
		//os.Exit(-1)
	}

	LoadDefaultSettings()

	if RootCmd.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("Verbose", Verbose)
	}
	if RootCmd.PersistentFlags().Lookup("debug").Changed {
		viper.Set("Debug", Debug)
	}
	if RootCmd.PersistentFlags().Lookup("dontminify").Changed {
		viper.Set("DontMinify", UnMinify)
	}
	if RootCmd.PersistentFlags().Lookup("sourceDir").Changed {
		viper.Set("sourceDir", Source)
	}
	if RootCmd.PersistentFlags().Lookup("deployRecordDir").Changed {
		viper.Set("deployRecordDir", Deploy)
	}

	if viper.GetBool("verbose") {
		//TODO: Maybe have additional flag for verbose v. debug tracing
		jww.SetStdoutThreshold(jww.LevelInfo)
	}
	if viper.GetBool("debug") {
		jww.SetStdoutThreshold(jww.LevelTrace)
	}

	SkipFiles = viper.GetStringSlice("skipfiles")

	jww.INFO.Println("Listing Config:")
	for _, x := range viper.AllKeys() {
		jww.INFO.Println(x, ":", viper.Get(x))
	}

}

func checkSourcePath() {
	//TODO: Need to do some fancy path fixing for relative paths etc
	Source = viper.GetString("sourceDir")
	jww.INFO.Println("Checking Source Dir exists: ", Source)
	b, err := exists(Source)
	if err != nil {
		er(err)
	}
	if !b {
		jww.CRITICAL.Println("Source Dir does not exist", Source)
		os.Exit(-1)
	}
}

// Only needs to be checked when we are not init-ing
func checkDeployPath() {
	Deploy = viper.GetString("deployRecordDir")
	jww.INFO.Println("Checking Deploy Dir exists: ", Deploy)
	b, err := exists(Deploy)
	if err != nil {
		er(err)
	}
	if !b {
		jww.CRITICAL.Println("Deploy Dir does not exist", Deploy)
		os.Exit(-1)
	}
}
