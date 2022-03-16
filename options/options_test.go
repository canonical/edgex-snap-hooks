// -*- Mode: Go; indent-tabs-mode: t -*-

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

package options

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/stretchr/testify/require"
)

const TEST_SERVICE = "configure-test"
const TEST_SERVICE2 = "configure-test2"

func TestOptionsSet(t *testing.T) {

	configDir := fmt.Sprintf("%s/%s/res/", env.SnapDataConf, TEST_SERVICE)
	configFile := path.Join(configDir, TEST_SERVICE+".env")
	os.MkdirAll(configDir, os.ModePerm)
	config2Dir := fmt.Sprintf("%s/%s/res/", env.SnapDataConf, TEST_SERVICE2)
	config2File := path.Join(config2Dir, TEST_SERVICE2+".env")
	os.MkdirAll(config2Dir, os.ModePerm)

	snapSet(t, "config.key.value", "value01")
	snapSet(t, "apps."+TEST_SERVICE+".config.app.key.value", "value02")
	snapSet(t, "apps."+TEST_SERVICE2+".config.app.key.value", "value03")

	require.NoError(t, ProcessOptions([]string{TEST_SERVICE, TEST_SERVICE2}), "Error setting options.")
	require.NoError(t, isInFile(configFile, "export KEY_VALUE=value01"), "Error validating .env file")
	require.NoError(t, isInFile(configFile, "export APP_KEY_VALUE=value02"), "Error validating .env file")

	require.NoError(t, isInFile(config2File, "export KEY_VALUE=value01"), "Error validating .env file")
	require.NoError(t, isInFile(config2File, "export APP_KEY_VALUE=value03"), "Error validating .env file")

}

// utility testing functions

func snapSet(t *testing.T, key, value string) {
	err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, value)).Run()
	require.NoError(t, err, "Error setting config value via snapctl.")
}

func isInFile(file string, line string) error {
	// read the whole file at once
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if strings.Contains(string(b), line) {
		return nil
	} else {
		return fmt.Errorf("Line %s not found in %s", line, file)
	}
}
