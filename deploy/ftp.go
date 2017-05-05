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
	"path"
	"strings"
	"os"
	
	"github.com/dutchcoders/goftp"
	"github.com/mindok/hugodeploy/cred"
	
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type FTPDeployer struct {
	HostID     string
	Port       string
	UID        string
	PWD        string
	RootDir    string
	DisableTLS bool
	ftp        *goftp.FTP
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
	
	credentials, err := credProvider.GetCredentials("ftp")
	if (err != nil) {
		jww.ERROR.Println("Unable to retrieve credentials: ", err)
		return err;
	}			
			
	f.UID = credentials.UID
	f.PWD = credentials.PWD
	
	f.RootDir = viper.GetString("ftp.rootdir")
	f.DisableTLS = false
	if viper.IsSet("ftp.disabletls") {
		f.DisableTLS = viper.GetBool("ftp.disabletls")
	}

	jww.INFO.Println("Got FTP settings: ", f.HostID, f.Port, f.UID, f.RootDir)

	if f.HostID == "" {
		serr = serr + "HostID not found. Define ftp.host in config file. "
	}
	if f.Port == "" {
		serr = serr + "Port not found. Define ftp.port in config file. "
	}		
	if f.RootDir == "" {
		f.RootDir = "/"
		jww.WARN.Println("FTP: Website root directory not found (ftp: rootdir in config). Defaulting to '/'")
	}

	if serr != "" {
		jww.ERROR.Println("Error initialising FTP Deployer: ", serr)
		panic(errors.New("Error initialising FTP Deployer. " + serr))
	}	

	jww.FEEDBACK.Println("Creating FTP connection... ")
	//Create initial connection
	if !viper.GetBool("debug") && !viper.GetBool("verbose") {
		if f.ftp, err = goftp.Connect(f.HostID + ":" + f.Port); err != nil {
			jww.ERROR.Println("Failed initial FTP connection: ", err)
			return err
		}
	} else {
		jww.INFO.Println("Connecting to ftp in debug mode")
		if f.ftp, err = goftp.ConnectDbg(f.HostID + ":" + f.Port); err != nil {
			jww.ERROR.Println("Failed initial FTP connection: ", err)
			return err
		}
	}

	//Activate TLS
	if !f.DisableTLS {
		config := tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.RequestClientCert,
		}

		if err = f.ftp.AuthTLS(&config); err != nil {
			jww.ERROR.Println("Failed TLS Activation: ", err)
			return err
		}
	} else {
		jww.WARN.Println("FTP TLS disabled - data will be transmitted in clear text")
	}

	if err = f.ftp.Login(f.UID, f.PWD); err != nil {
		jww.ERROR.Println("Failed FTP Login: ", err)
		return err
	}

	jww.FEEDBACK.Println("Successfully connected to FTP")

	return nil

}

func makeFtpPath(path string) string {
	fpath := path

	/**
	If the os's local path separator is not / (e.g. Windows) we have to adjust the path here
	because FTP only accepts forward slashes
	**/	
	if (os.PathSeparator != '/') {		
		fpath = strings.Replace(fpath, string(os.PathSeparator), string("/"), -1);
		// remove potentially doubled separators
		fpath = strings.Replace(fpath, string("//"), string("/"), -1);		
	}
	
	return fpath
}

func (f *FTPDeployer) ApplyCommand(cmd *DeployCommand) error {
	p := makeFtpPath(path.Join(f.RootDir, cmd.RelPath))		
	
	switch cmd.Command {
	case COMMAND_FILE_ADD, COMMAND_FILE_UPD:
		return f.UploadFile(p, cmd.Contents)

	case COMMAND_DIR_ADD:
		return f.MakeDirectory(p)

	case COMMAND_DIR_DEL:
		return f.RemoveDirectory(p)

	case COMMAND_FILE_DEL:
		return f.RemoveFile(p)

	default:
		return errors.New("Not implemented")
	}
	//jww.WARN.Println("SFTP Cmds not implemented yet: ", cmd.RelPath)
}

func (f *FTPDeployer) UploadFile(path string, data []byte) error {	
	r := bytes.NewReader(data)
	jww.FEEDBACK.Println("Sending file: ", path, "...")
	jww.DEBUG.Println("Data Size: ", len(data))	
	
	if err := f.ftp.Stor(path, r); err != nil {
		jww.ERROR.Println("FTP Error uploading file: ", path, err)
		jww.DEBUG.Println("Data that could not be sent: ", data)
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
	jww.FEEDBACK.Println("Deleting file: ", path, "...")

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
	jww.FEEDBACK.Println("Creating directory: ", path, "...")
	if err := f.ftp.Mkd(path); err != nil {
		if strings.Contains(err.Error(), "File exists") {
			jww.INFO.Println("Looks like FTP directory already exists: ", path)
			return nil
		} else {
			jww.ERROR.Println("FTP Error creating directory: ", path, err)
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
