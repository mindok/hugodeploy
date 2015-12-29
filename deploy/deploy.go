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

//TODO: Provide Deployer map to allow more flexible configuration

//Deployer interface
type Deployer interface {
	GetName() string //May not need this
	Initialise() error
	Cleanup() error
	ApplyCommand(cmd *DeployCommand) error
}

type CommandType int
type commandHandler func(cmd *DeployCommand) (err error)

const (
	COMMAND_FILE_ADD CommandType = 1 << iota
	COMMAND_DIR_ADD
	COMMAND_FILE_UPD
	COMMAND_FILE_DEL
	COMMAND_DIR_DEL
)

func (cmd *DeployCommand) GetCommandDesc() string {
	s := ""
	c := cmd.Command
	switch c {
	case COMMAND_DIR_ADD:
		s = "ADD DIR"
	case COMMAND_DIR_DEL:
		s = "DELETE DIR"
	case COMMAND_FILE_ADD:
		s = "ADD FILE"
	case COMMAND_FILE_DEL:
		s = "DELETE FILE"
	case COMMAND_FILE_UPD:
		s = "UPDATE FILE"
	default:
		s = ""
	}
	return s
}

type DeployCommand struct {
	RelPath  string
	Contents []byte
	Command  CommandType
}
