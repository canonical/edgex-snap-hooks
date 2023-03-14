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

package options

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/canonical/edgex-snap-hooks/v3/env"
	"github.com/canonical/edgex-snap-hooks/v3/log"
)

type configProcessor struct {
	appEnvVars            map[string]map[string]string
	envSegmentSeparator   string
	envHierarchySeparator string
	configHierarchy       bool
}

func newConfigProcessor(apps []string, hierarchy bool, hSep, sSep string) *configProcessor {
	cp := configProcessor{
		appEnvVars:            make(map[string]map[string]string),
		configHierarchy:       hierarchy,
		envHierarchySeparator: hSep,
		envSegmentSeparator:   sSep,
	}
	for _, app := range apps {
		cp.appEnvVars[app] = make(map[string]string)
	}
	return &cp
}

// add app's env var to memory
func (cp *configProcessor) addEnvVar(app, key, value string) error {
	envKey, err := cp.configKeyToEnvVar(key)
	if err != nil {
		return fmt.Errorf("error converting config key to environment variable key: %s", err)
	}
	log.Infof("Mapping %s to %s", key, envKey)
	cp.appEnvVars[app][envKey] = value
	return err
}

// convert snap option key to environment variable name
func (cp *configProcessor) configKeyToEnvVar(configKey string) (string, error) {
	if cp.configHierarchy {
		configKey = strings.ReplaceAll(configKey, ".", cp.envHierarchySeparator)
	} else if strings.Contains(configKey, ".") {
		return "", fmt.Errorf("config key must not contain dots: %s", configKey)
	}

	return strings.ToUpper(
		// replace the segment separator
		strings.ReplaceAll(configKey, "-", cp.envSegmentSeparator),
	), nil
}

// returns the suitable env file name for the service
func (cp *configProcessor) filename(service string) string {
	// The app-service-configurable snap is the one outlier snap that doesn't
	// include the service name in it's configuration path.
	var path string
	if env.SnapName == "edgex-app-service-configurable" {
		path = fmt.Sprintf("%s/config/overrides.env", env.SnapData)
	} else {
		path = fmt.Sprintf("%s/config/%s/overrides.env", env.SnapData, service)
	}
	return path
}

func (cp *configProcessor) writeEnvFiles() error {
	for app, envVars := range cp.appEnvVars {
		var buffer bytes.Buffer
		filename := cp.filename(app)

		// do not create a .env file if there are no snap options set for the app
		// remove .env file if exists
		if len(envVars) == 0 {
			if err := os.RemoveAll(filename); err != nil {
				return fmt.Errorf("failed to remove env file: %s", err)
			}
			continue
		}

		// Add comment on top of the file
		if _, err := fmt.Fprintln(&buffer, "# Sys-gen env vars from snap options:"); err != nil {
			return err
		}

		// add env vars to buffer
		for k, v := range envVars {
			if _, err := fmt.Fprintf(&buffer, "%s=\"%s\"\n", k, v); err != nil {
				return err
			}
		}

		log.Infof("Writing to env file %s: %s", filename, strings.ReplaceAll(buffer.String(), "\n", " "))

		dir := filepath.Dir(filename)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		tmp := filename + ".tmp"
		err = os.WriteFile(tmp, buffer.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("failed to write %s: %s", tmp, err)
		}

		err = os.Rename(tmp, filename)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %s", tmp, filename, err)
		}
	}

	return nil
}
