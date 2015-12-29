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
	jww "github.com/spf13/jwalterweatherman"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileDeployer struct {
	TargetDir string
}

func (f *FileDeployer) GetName() string {
	return "File"
}

func (f *FileDeployer) Initialise() error {
	if f.TargetDir == "" {
		panic(errors.New("TargetDir not set for FileDeployer - aborting before anything bad happens"))
		os.Exit(-1)
	}
	return nil
}

func (f *FileDeployer) ApplyCommand(cmd *DeployCommand) error {
	path := filepath.Join(f.TargetDir, cmd.RelPath)
	switch cmd.Command {
	case COMMAND_FILE_ADD, COMMAND_FILE_UPD:
		return f.UploadFile(path, cmd.Contents)

	case COMMAND_DIR_ADD:
		return f.MakeDirectory(path)

	case COMMAND_DIR_DEL:
		return f.RemoveDirectory(path)

	case COMMAND_FILE_DEL:
		return f.RemoveFile(path)

	default:
		return errors.New("Not implemented")
	}
	//jww.WARN.Println("SFTP Cmds not implemented yet: ", cmd.RelPath)
	return nil
}

func (f *FileDeployer) UploadFile(path string, data []byte) error {
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		jww.ERROR.Println("Error writing file: ", path, err)
		return err
	} else {
		jww.INFO.Println("Successfully wrote file: ", path)
	}
	return nil
}

func (f *FileDeployer) RemoveDirectory(path string) error {
	jww.WARN.Println("Removing directory: ", path)
	if err := os.RemoveAll(path); err != nil {
		jww.ERROR.Println("Error deleting dir: ", path, err)
		return err
	} else {
		jww.INFO.Println("Successfully deleted dir: ", path)
	}
	return nil
}

func (f *FileDeployer) RemoveFile(path string) error {
	jww.WARN.Println("Removing file: ", path)
	if err := os.Remove(path); err != nil {
		jww.ERROR.Println("Error deleting file: ", path, err)
		return err
	} else {
		jww.TRACE.Println("Successfully deleted file: ", path)
	}
	return nil
}

func (f *FileDeployer) MakeDirectory(path string) error {
	if err := os.Mkdir(path, 0777); err != nil {
		jww.ERROR.Println("Error creating directory: ", path, err)
		if strings.Contains(err.Error(), "File exists") {
			jww.INFO.Println("Looks like directory already exists: ", path)
			return nil
		} else {
			return err
		}
	} else {
		jww.INFO.Println("Successfully created directory: ", path)
	}
	return nil
}

func (f *FileDeployer) Cleanup() error {
	//Nothing to do
	return nil
}
