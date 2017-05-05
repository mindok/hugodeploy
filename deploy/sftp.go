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

package deploy

import (
	"errors"
	"github.com/pkg/sftp"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"github.com/mindok/hugodeploy/cred"
)

/* NOTE: INCOMPLETE, UNTESTED CODE */

type SFTPDeployer struct {
	HostID     string
	Port       string
	UID        string
	PWD        string
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func (s *SFTPDeployer) GetName() string {
	return "SFTP"
}

func (s *SFTPDeployer) Initialise() error {
	serr := ""
	jww.INFO.Println("Getting SFTP settings")
	// Gather together settings
	s.HostID = viper.GetString("sftp.host")
	s.Port = viper.GetString("sftp.port")
	
	// get credentials
	credProviderName := viper.GetString("ftp.credentialProvider")
	if (credProviderName == "") {
		// when no credental provider is given
		// fall back to old style credentials
		credProviderName = "classic"
	} 
	
	jww.DEBUG.Println("CredentialProvider is ", credProviderName)
	
	credProvider := cred.GetCredentialProvider(credProviderName)
	if (credProvider == nil) {
		jww.ERROR.Println("Unknown or unsupported credential provider: ", credProviderName)
		return errors.New ( "Unknown or unsupported credential provider: " + credProviderName );
	}
	
	credentials, err := credProvider.GetCredentials("sftp")
	if (err != nil) {
		jww.ERROR.Println("Unable to retrieve credentials: ", err)
		return err;
	}			
			
	s.UID = credentials.UID
	s.PWD = credentials.PWD
	
	jww.INFO.Println("Got SFTP settings: ", s.HostID, s.Port, s.UID)

	if s.HostID == "" {
		serr = serr + "HostID not found. Define sftp.host in config file. "
	}
	if s.Port == "" {
		serr = serr + "Port not found. Define sftp.port in config file. "
	}	
	if serr != "" {
		return errors.New("Error initialising SFTP Deployer. " + serr)
	}

	err = errors.New("") //Must be away to avoid this, but double function returns below barf

	//Attempt to connect. First create the SSH client:
	config := &ssh.ClientConfig{
		User: s.UID,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.PWD),
		},
	}
	s.sshClient, err = ssh.Dial("tcp", s.HostID+":"+s.Port, config)
	if err != nil {
		jww.ERROR.Println("SSH subsystem failed to connect to ", s.HostID, " Error: ", err)
		return err
	}
	jww.INFO.Println("Successfully connected to SSH")

	//Now attempt to connect SFTP
	s.sftpClient, err = sftp.NewClient(s.sshClient)
	if err != nil {
		jww.ERROR.Println("SFTP failed to connect. Error: ", err)
		return err
	}
	jww.INFO.Println("Successfully connected to SFTP")

	return nil

}

func (s *SFTPDeployer) ApplyCommand(cmd *DeployCommand) error {
	return errors.New("SFTP Cms not implemented yet")
}

func (s *SFTPDeployer) Cleanup() error {

	s.sftpClient.Close()
	s.sshClient.Close()

	return nil
}
