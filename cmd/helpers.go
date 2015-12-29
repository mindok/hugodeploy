// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Minor mods Philosopher Businessman abp@philosopherbusinessman.com Dec 2015
// to trim out unneeded stuff

package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var funcMap template.FuncMap
var projectPath = ""
var inputPath = ""
var projectBase = ""

// for testing only
var testWd = ""

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(-1)
}

// Check if a file or directory exists.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ProjectPath() string {
	if projectPath == "" {
		guessProjectPath()
	}

	return projectPath
}

// wrapper of the os package so we can test better
func getWd() (string, error) {
	if testWd == "" {
		return os.Getwd()
	}
	return testWd, nil
}

func guessProjectPath() {
	// if no path is provided... assume CWD.
	if inputPath == "" {
		x, err := getWd()
		if err != nil {
			er(err)
		}

		if projectPath == "" {
			projectPath = filepath.Clean(x)
			return
		}
	}

	// if provided, inspect for logical locations
	if strings.ContainsRune(inputPath, os.PathSeparator) {
		if filepath.IsAbs(inputPath) || filepath.HasPrefix(inputPath, string(os.PathSeparator)) {
			// if Absolute, use it
			projectPath = filepath.Clean(inputPath)
			return
		}
	} else {
		// hardest case.. just a word.
		if projectBase == "" {
			x, err := getWd()
			if err == nil {
				projectPath = filepath.Join(x, inputPath)
				return
			}

		}
	}
	projectPath = ""
	er("Can't create sensible project path")
}

// isDir checks if a given path is a directory.
func isDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// dirExists checks if a path exists and is a directory.
func dirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func writeStringToFile(path, file, text string) error {
	filename := filepath.Join(path, file)

	r := strings.NewReader(text)
	err := safeWriteToDisk(filename, r)

	if err != nil {
		return err
	}
	return nil
}

// Same as WriteToDisk but checks to see if file/directory already exists.
func safeWriteToDisk(inpath string, r io.Reader) (err error) {
	dir, _ := filepath.Split(inpath)
	ospath := filepath.FromSlash(dir)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			return
		}
	}

	ex, err := exists(inpath)
	if err != nil {
		return
	}
	if ex {
		return fmt.Errorf("%v already exists", inpath)
	}

	file, err := os.Create(inpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return
}
