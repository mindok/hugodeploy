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

//Package sync provides the routines that determine which files need to
//be transfered to the deployment target. The actual deployment actions
//are delegated to functions provided by the caller to allow different
//deployment targets. Note that we assume there are no huge files here - the
//file contents are passed around as byte arrays for comparisons etc
package deploy

import (
	"bytes"
	"errors"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var filetypeMime = map[string]string{
	"css":  "text/css",
	"htm":  "text/html",
	"html": "text/html",
	"js":   "text/javascript",
	"json": "application/json",
	"svg":  "image/svg+xml",
	"xml":  "text/xml",
}

type DeployScanner struct {
	minify     bool
	handleFunc commandHandler
	srcDir     string
	dstDir     string
	skipFiles  []string
	minifier   *minify.M
}

//DeployChanges recursively walks through srcDir and compares each file with the equivalent
//location in dstDir. If minify is set, each file in srcDir will be minified in memory.
//If the file in srcDir differs from that in dstDir handleFunc is called with an ADD or UPD
//command. DeployChanges then walks the dstDir to see if there are any files there which are not
//in srcDir, in which case handleFunc is called with a DEL command
func DeployChanges(srcDir string, dstDir string, minify bool, handleFunc commandHandler, skipFiles []string) error {
	deployer := &DeployScanner{minify, handleFunc, srcDir, dstDir, skipFiles, nil}
	deployer.initM()
	return deployer.Sync(dstDir, srcDir)
}

func (d *DeployScanner) initM() {
	d.minifier = minify.New()
	d.minifier.AddFunc("text/css", css.Minify)
	d.minifier.AddFunc("text/html", html.Minify)
	d.minifier.AddFunc("text/javascript", js.Minify)
	d.minifier.AddFunc("image/svg+xml", svg.Minify)
	d.minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	d.minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

// Sync copies files and directories inside src into dst.
func (d *DeployScanner) Sync(dst, src string) error {
	// make sure src & destination exist
	err := checkDirExists(src, "Source")
	if err != nil {
		return err
	}

	err = checkDirExists(dst, "Destination")
	if err != nil {
		return err
	}

	return d.syncRecover(dst, src)
}

// syncRecover handles errors and calls sync
func (d *DeployScanner) syncRecover(dst, src string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	d.sync(dst, src)
	return nil
}

func (d *DeployScanner) getRelativePath(src string) string {
	s := strings.Replace(src, d.srcDir, "", 1)
	sep := string(os.PathSeparator)
	if !strings.HasPrefix(s, sep) {
		s = sep + s
	}
	return s
}

//TODO: Improve implementation - this is rather naive
func (d *DeployScanner) shouldSkip(path string) bool {
	for _, s := range d.skipFiles {
		if strings.Contains(path, s) {
			return true
		}
	}
	return false
}

func (d *DeployScanner) makeCreateDirCmd(src string) *DeployCommand {
	return &DeployCommand{d.getRelativePath(src), nil, COMMAND_DIR_ADD}
}

func (d *DeployScanner) makeDeleteDirCmd(src string) *DeployCommand {
	return &DeployCommand{d.getRelativePath(src), nil, COMMAND_DIR_DEL}
}

func (d *DeployScanner) makeCreateFileCmd(src string, data []byte) *DeployCommand {
	//TODO: Unpack files, fix up paths, minify source
	return &DeployCommand{d.getRelativePath(src), data, COMMAND_FILE_ADD}
}

func (d *DeployScanner) makeDeleteFileCmd(src string) *DeployCommand {
	return &DeployCommand{d.getRelativePath(src), nil, COMMAND_FILE_DEL}
}

func (d *DeployScanner) makeUpdateFileCmd(src string, data []byte) *DeployCommand {
	//TODO: Unpack files, fix up paths, minify source
	return &DeployCommand{d.getRelativePath(src), data, COMMAND_FILE_UPD}
}

// sync updates dst to match with src, handling both files and directories.
func (d *DeployScanner) sync(dst, src string) {

	jww.FEEDBACK.Println("Comparing Dst: ", dst, " With Src: ", src)

	//Build a map of all files so we can figure out what to delete at the end
	srcFiles := make(map[string]os.FileInfo)
	//We need an ordered set of filesnames too...
	srcFileKeys := make([]string, 0)

	var scan = func(path string, fileInfo os.FileInfo, inpErr error) (err error) {
		jww.TRACE.Println("Scanning: ", path)
		if inpErr == nil {
			srcFiles[path] = fileInfo
			srcFileKeys = append(srcFileKeys, path)
		}
		return nil
	}

	err := filepath.Walk(src, scan)
	check(err)
	//jww.WARN.Println(srcFileKeys)

	//For each source file entry
	// If a directory
	//  Delete destination file of matching name if destination is a file
	//  Create destination directory if it was a file or if it is missing
	// If a file
	//  Check whether destination of matching name is a directory, and if so delete
	//  If destination was a directory, or is missing create destination file
	//  If destination file exists, but is different, update it
	for _, srcFile := range srcFileKeys {
		if d.shouldSkip(srcFile) {
			jww.FEEDBACK.Println("Skipping ", srcFile)
		} else {
			//TODO: May be a bit of a naive approach to figure out path to destination file
			dstFile := strings.Replace(srcFile, src, dst, 1)
			sstat := srcFiles[srcFile]

			jww.TRACE.Println("Checking source ", srcFile, " against destination ", dstFile)

			dstat, err := os.Stat(dstFile)
			if err != nil && !os.IsNotExist(err) {
				panic(err)
			}
			dExists := (err == nil)

			if sstat.IsDir() {
				jww.TRACE.Println("Src is a directory: ", srcFile)
				if dExists && dstat.IsDir() {
					jww.TRACE.Println("Dst is a dir - nothing to do: ", dstFile)
					jww.INFO.Println("Directories the same - skipping: ", dstFile)
				}
				if dExists && !dstat.IsDir() {
					jww.TRACE.Println("Dst is a file: ", dstFile)
					jww.INFO.Println("Replacing destination file with directory of same name: ", dstFile)
					check(d.handleFunc(d.makeDeleteFileCmd(srcFile)))
					check(d.handleFunc(d.makeCreateDirCmd(srcFile)))
				}
				if !dExists {
					jww.TRACE.Println("Dst doesn't exist: ", dstFile)
					jww.INFO.Println("Creating directory: ", dstFile)
					check(d.handleFunc(d.makeCreateDirCmd(srcFile)))
				}
			} else {
				jww.TRACE.Println("Src is a file: ", srcFile)
				data, err := d.getSourceData(srcFile)
				check(err)
				if dExists && dstat.IsDir() {
					jww.TRACE.Println("Dst is a dir: ", dstFile)
					jww.INFO.Println("Replacing directory with file of same name: ", dstFile)
					check(d.handleFunc(d.makeDeleteDirCmd(srcFile)))
					check(d.handleFunc(d.makeCreateFileCmd(srcFile, data)))
				}
				if !dExists {
					jww.TRACE.Println("Dst doesn't exist: ", dstFile)
					jww.INFO.Println("Creating file: ", dstFile)
					check(d.handleFunc(d.makeCreateFileCmd(srcFile, data)))
				}
				if dExists && !dstat.IsDir() {
					jww.TRACE.Println("Dst is a file: ", dstFile)
					if !d.filesEqual(srcFile, dstFile, data) {
						jww.TRACE.Println("Dst exists - updating")
						jww.INFO.Println("Updating file: ", dstFile)
						check(d.handleFunc(d.makeUpdateFileCmd(srcFile, data)))
					} else {
						jww.INFO.Println("Files the same - skipping: ", dstFile)
					}
				}
			}
		}
	}

	//TODO: Do in reverse order so that files get deleted before their parent directories
	dstDeleteFiles := make([]string, 0)
	dstDeleteDirs := make([]string, 0)
	var scanDeletes = func(path string, fileInfo os.FileInfo, inpErr error) (err error) {
		if inpErr == nil {
			srcFileExpected := strings.Replace(path, dst, src, 1)
			jww.TRACE.Println("Checking to deleted: ", path, ". Looking for: ", srcFileExpected)
			_, err := os.Stat(srcFileExpected)
			if err != nil && os.IsNotExist(err) && !d.shouldSkip(path) {
				if fileInfo.IsDir() {
					dstDeleteDirs = append(dstDeleteDirs, srcFileExpected)
				} else {
					dstDeleteFiles = append(dstDeleteFiles, srcFileExpected)

				}
			}
			if err != nil && !os.IsNotExist(err) {
				panic(err)
			}
		}
		return nil
	}

	err = filepath.Walk(dst, scanDeletes)
	check(err)

	for i := len(dstDeleteFiles) - 1; i >= 0; i-- {
		if d.shouldSkip(dstDeleteFiles[i]) {
			jww.FEEDBACK.Println("Skipping ", dstDeleteFiles[i])
		}
		check(d.handleFunc(d.makeDeleteFileCmd(dstDeleteFiles[i])))
	}
	for i := len(dstDeleteDirs) - 1; i >= 0; i-- {
		if d.shouldSkip(dstDeleteDirs[i]) {
			jww.FEEDBACK.Println("Skipping ", dstDeleteDirs[i])
		}
		check(d.handleFunc(d.makeDeleteDirCmd(dstDeleteDirs[i])))
	}
}

func (d *DeployScanner) getSourceData(src string) ([]byte, error) {
	contents, err := ioutil.ReadFile(src)
	jww.DEBUG.Println("getSourceData step 1: ", len(contents), " bytes read from: ", src)
	check(err)
	if d.minify {
		mediatype := getMediaType(src)
		jww.DEBUG.Println("Minifier media type ", mediatype, " for ", src)
		if mediatype != "" {
			contents, err = d.minifier.Bytes(mediatype, contents)
			jww.DEBUG.Println("getSourceData step 2: ", len(contents), " bytes when minified: ", src)
			check(err)
		}
	}
	return contents, err
}

func getMediaType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	jww.DEBUG.Println("Mediatype lookup - extension detected: ", ext)
	if mimetype, ok := filetypeMime[ext]; ok {
		return mimetype
	}
	return ""
}

// equal returns true if both files are equal, but takes into account
// that src may be minified.
func (d *DeployScanner) filesEqual(src, dst string, srcdata []byte) bool {
	// get file infos
	info1, err1 := os.Stat(src)
	info2, err2 := os.Stat(dst)
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		return false
	}
	check(err1)
	check(err2)

	// check sizes, but only if we aren't minifying as size will change
	if !d.minify {
		if info1.Size() != info2.Size() {
			return false
		}
	}

	// both have the same size, check the contents
	// Hopefully the files aren't too big as there is no chunking...
	contents, err := ioutil.ReadFile(dst)
	check(err)

	return bytes.Equal(contents, srcdata)
}

func checkDirExists(path, name string) error {
	sstat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !sstat.IsDir() {
		return errors.New(name + " must be a directory")
	}

	return nil
}

func check(err error) {
	if err != nil {
		jww.ERROR.Println("Error: ", err)
		panic(err)
	}
}
