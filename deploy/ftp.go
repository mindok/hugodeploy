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
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/dutchcoders/goftp"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"strings"
)

type FTPDeployer struct {
	HostID string
	Port   string
	UID    string
	PWD    string
	ftp    *goftp.FTP
}

func (f *FTPDeployer) GetName() string {
	return "FTP"
}

func (f *FTPDeployer) Initialise() error {
	serr := ""
	jww.INFO.Println("Getting FTP settings")
	// Gather together settings
	f.HostID = viper.GetString("ftp.host")
	f.Port = viper.GetString("ftp.port")
	f.UID = viper.GetString("ftp.user")
	f.PWD = viper.GetString("ftp.pwd")
	jww.INFO.Println("Got FTP settings: ", f.HostID, f.Port, f.UID)

	if f.HostID == "" {
		serr = serr + "HostID not found. Define ftp.host in config file. "
	}
	if f.Port == "" {
		serr = serr + "Port not found. Define ftp.port in config file. "
	}
	if f.UID == "" {
		serr = serr + "UID not found. Define ftp.user in config file. "
	}
	if f.PWD == "" {
		serr = serr + "PWD not found. Define ftp.pwd in config file. "
	}

	if serr != "" {
		return errors.New("Error initialising FTP Deployer. " + serr)
	}

	err := errors.New("") //Must be away to avoid this, but double function returns below barf

	//Create initial connection
	if f.ftp, err = goftp.Connect(f.HostID + ":" + f.Port); err != nil {
		jww.ERROR.Println("Failed initial FTP connection: ", err)
		return err
	}

	//Activate TLS
	config := tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequestClientCert,
	}

	if err = f.ftp.AuthTLS(config); err != nil {
		jww.ERROR.Println("Failed TLS Activation: ", err)
		return err
	}

	if err = f.ftp.Login(f.UID, f.PWD); err != nil {
		jww.ERROR.Println("Failed FTP Login: ", err)
		return err
	}

	//TODO: Assume everything is relative to root...
	if err = f.ftp.Cwd("/"); err != nil {
		jww.ERROR.Println("Failed to change to root directory: ", err)
		return err
	}

	jww.INFO.Println("Successfully connected to FTP")

	return nil

}

func (f *FTPDeployer) ApplyCommand(cmd *DeployCommand) error {
	switch cmd.Command {
	case COMMAND_FILE_ADD, COMMAND_FILE_UPD:
		return f.UploadFile(cmd.RelPath, cmd.Contents)

	case COMMAND_DIR_ADD:
		return f.MakeDirectory(cmd.RelPath)

	case COMMAND_DIR_DEL:
		return f.RemoveDirectory(cmd.RelPath)

	case COMMAND_FILE_DEL:
		return f.RemoveFile(cmd.RelPath)

	default:
		return errors.New("Not implemented")
	}
	//jww.WARN.Println("SFTP Cmds not implemented yet: ", cmd.RelPath)
	return nil
}

func (f *FTPDeployer) UploadFile(path string, data []byte) error {
	r := bytes.NewReader(data)

	if err := f.ftp.Stor(path, r); err != nil {
		jww.ERROR.Println("FTP Error uploading file: ", path, err)
		return err
	} else {
		jww.INFO.Println("Successfully FTP'd file: ", path)
	}
	return nil
}

func (f *FTPDeployer) RemoveDirectory(path string) error {
	jww.WARN.Println("FTP Directory delete requested, but not implemented")
	/*
		if err := f.ftp.Rmd(path); err != nil {
			jww.ERROR.Println("FTP Error deleting directory: ", path, err)
			return err
		} else {
			jww.INFO.Println("Successfully deleted directory: ", path)
		}
	*/
	return nil
}

func (f *FTPDeployer) RemoveFile(path string) error {
	if err := f.ftp.Dele(path); err != nil {
		if strings.Contains(err.Error(), "No such file") {
			jww.INFO.Println("Looks like FTP file already deleted: ", path)
			return nil
		} else {
			jww.ERROR.Println("FTP Error deleting file: ", path, err)
			return err
		}
	} else {
		jww.INFO.Println("Successfully deleted file: ", path)
	}
	return nil
}

func (f *FTPDeployer) MakeDirectory(path string) error {
	if err := f.ftp.Mkd(path); err != nil {
		jww.ERROR.Println("FTP Error creating directory: ", path, err)
		if strings.Contains(err.Error(), "File exists") {
			jww.INFO.Println("Looks like FTP directory already exists: ", path)
			return nil
		} else {
			return err
		}
	} else {
		jww.INFO.Println("Successfully created FTP directory: ", path)
	}
	return nil
}

func (f *FTPDeployer) Cleanup() error {

	f.ftp.Close()

	return nil
}
