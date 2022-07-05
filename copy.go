/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package hooks

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"
)

// CopyFile copies a file within the snap
func CopyFile(srcPath, destPath string) error {

	inFile, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// TODO: check file perm
	err = ioutil.WriteFile(destPath, inFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CopyDir copies a whole directory recursively
// snippet from https://blog.depa.do/post/copy-files-and-directories-in-go
func CopyDir(srcPath string, dstPath string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	srcinfo, err = os.Stat(srcPath)
	if err != nil {
		return err
	}

	// Temporarily change the process umask to allow creating directory with rwx permissions.
	oldMask := syscall.Umask(0)
	defer syscall.Umask(oldMask)

	err = os.MkdirAll(dstPath, srcinfo.Mode())
	if err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(srcPath); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(srcPath, fd.Name())
		dstfp := path.Join(dstPath, fd.Name())

		if fd.IsDir() {
			err = CopyDir(srcfp, dstfp)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcfp, dstfp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFileReplace copies a file within the snap and replaces strings using
// the string/replace values in the rStrings parameter.
func CopyFileReplace(srcPath, destPath string, rStrings map[string]string) error {

	inFile, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	rStr := string(inFile)
	for k, v := range rStrings {
		rStr = strings.Replace(rStr, k, v, 1)
	}

	// TODO: check file perm
	outBytes := []byte(rStr)
	err = ioutil.WriteFile(destPath, outBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
