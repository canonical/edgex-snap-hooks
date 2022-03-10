/*
 * Copyright (C) 2022 Canonical Ltd
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

package files

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/canonical/edgex-snap-hooks/v2/log"
)

// CopyFile copies a file within the snap
func copyFile(srcPath, destPath string) error {

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

func CopyFileFromSnapToSnapData(path string) error {
	var err error

	destFile := SnapData + path
	srcFile := Snap + path
	dir := filepath.Dir(destFile)

	// if configuration.toml already exists, it's been
	// provided by a content interface, so no need to
	// make the directory, which would cause any files
	// provided by the content interface to be deleted.
	if _, err = os.Stat(destFile); err == nil {
		return nil
	}

	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err = copyFile(srcFile, destFile); err != nil {
		return err
	}

	destFile = SnapData + path
	srcFile = Snap + path

	if err = os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
		return err
	}

	if err = copyFile(srcFile, destFile); err != nil {
		return err
	}

	log.Debugf("Copied %s to %s", srcFile, destFile)
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

func SnapDataDirectoryExists(dir string) bool {
	pathname := path.Join(SnapData, dir)

	if _, err := os.Stat(pathname); err == nil {
		return true
	}
	return false
}
