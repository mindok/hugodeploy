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
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises the setup",
	Long: `Initialises the setup by:
Creating a template config file if it doesn't exist
Creating the LastDeployed directory if it doesn't exist
Emptying it if it does.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !configFileExists() {
			createConfigFile()
		} else {

			if !deployDirExists() {
				createDeployDir()
			} else {
				emptyDir(Deploy)
			}
		}
	},
}

func deployDirExists() bool {
	Deploy = viper.GetString("deployRecordDir")
	jww.INFO.Println("Init: Checking Deploy Dir exists: ", Deploy)
	b, err := dirExists(Deploy)
	if err != nil {
		er(err)
	}

	if b {
		jww.INFO.Println("Init: Deploy Dir exists.")
	} else {
		jww.INFO.Println("Init: Deploy Dir does NOT exist.")
	}

	return b
}

func createDeployDir() {
	Deploy = viper.GetString("deployRecordDir")
	jww.INFO.Println("Init: Attempting to create Deploy Dir: ", Deploy)
	err := os.MkdirAll(Deploy, os.ModePerm)
	if err != nil {
		er(err)
	} else {
		jww.INFO.Println("Init: Created Deploy Dir: ")
	}
}

func askForConfirmation(match string) bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		jww.CRITICAL.Println("Fatal error waiting for confirmation", err)
		return false
	}
	return (response == match)
}

func emptyDir(path string) {

	fmt.Println("Type 'yes' to confirm emptying of ", path)
	if !askForConfirmation("yes") {
		jww.INFO.Println("Init: Cancelling at user request: ", path)
		return
	}

	jww.INFO.Println("Init: Emptying directory: ", path)

	fd, err := os.Open(path)
	if err != nil {
		jww.ERROR.Println("Init: Error reading path: ", path, err)
		if os.IsNotExist(err) {
			er(err)
			return
		}
	}

	err = nil
	for {
		names, err1 := fd.Readdirnames(100)
		for _, name := range names {
			fn := path + string(os.PathSeparator) + name
			jww.INFO.Println("Init: Removing: ", fn)
			err1 := os.RemoveAll(fn)
			if err == nil {
				err = err1
			}
		}
		if err1 == io.EOF {
			break
		}
		// If Readdirnames returned an error, use it.
		if err == nil {
			err = err1
		}
		if len(names) == 0 {
			break
		}
	}

	fd.Close()

	if err != nil {
		jww.ERROR.Println("Init: Error emptying directory: ", path, err)
	}
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func configFileExists() bool {
	jww.INFO.Println("Init: Checking config file exists")
	s := viper.ConfigFileUsed()
	jww.INFO.Println("Init: Reported file: ", s)
	b, err := exists(s)
	if err != nil {
		er(err)
	}
	if b {
		jww.INFO.Println("Init: Config File Found")
	} else {
		jww.INFO.Println("Init: Config File NOT Found")
	}
	return b
}

func createConfigFile() {
	path := ProjectPath()
	jww.INFO.Println("Init: Creating config file at path: ", path)
	filename := "hugodeploy.yaml"
	fullfilename := filepath.Join(path, filename)
	b, err := exists(fullfilename)
	if err != nil {
		er(err)
	}
	if b {
		jww.ERROR.Println("Config file exists. Strange that it wasn't found: ", fullfilename)
		os.Exit(-1)
	}
	template := `
# HugoDeploy Configuration File

# Connection settings for deployment target (FTP only)
ftp:
  host: <enter host id / ip address>
  port: <enter port - usually 21 for FTP over TLS>  
  rootdir: <enter root directory of website, e.g. /public_html/>
  disabletls: false  
  credentialProvider: classic
  user: <enter user id>
  pwd: <enter password>

# Connection settings for deployment target (SFTP only)
sftp:
  host: <enter host id / ip address>
  port: <enter port - usually 22 for SSH>
  user: <enter user id>
  pwd: <enter password>
  rootdir: <enter root directory of website, e.g. /public_html/>

# Location of files to publish. For hugo static sites this is PublishDir and defaults to public
sourcedir: published

# Skip files or directories which match the following patterns
skipfiles:
  - .DS_Store
  - .git
  - /tmp

# Location of directory used for tracking what has been deployed
deployRecordDir: deployed

# Want lots of messages? [Default false]
#verbose: true

# Disable minification? [Default false]
#DontMinify: true
`
	err = writeStringToFile(path, filename, template)
	if err != nil {
		jww.ERROR.Println("Config file exists. Strange that it wasn't found: ", fullfilename)
		os.Exit(-1)
	} else {
		jww.INFO.Println("Init: Successfully created config file: ", path)
	}
}
