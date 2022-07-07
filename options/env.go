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
	"strings"

	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/canonical/edgex-snap-hooks/v2/log"
)

type configProcessor struct {
	appEnvVars map[string]map[string]string
}

func newConfigProcessor(apps []string) *configProcessor {
	var cp configProcessor
	cp.appEnvVars = make(map[string]map[string]string)
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

// convert my-var to MY_VAR
func (cp *configProcessor) configKeyToEnvVar(configKey string) (string, error) {
	if strings.Contains(configKey, ".") {
		return "", fmt.Errorf("config key must not contain dots: %s", configKey)
	}
	return strings.ReplaceAll(strings.ToUpper(configKey), "-", "_"), nil
}

// returns the suitable env file name for the service
func (cp *configProcessor) filename(service string) string {
	// The app-service-configurable snap is the one outlier snap that doesn't
	// include the service name in it's configuration path.
	var path string
	if env.SnapName == "edgex-app-service-configurable" {
		path = fmt.Sprintf("%s/res/%s.env", env.SnapDataConf, service)
	} else {
		path = fmt.Sprintf("%s/%s/res/%s.env", env.SnapDataConf, service, service)
	}
	return path
}

func (cp *configProcessor) writeEnvFiles() error {
	for app, envVars := range cp.appEnvVars {
		var buffer bytes.Buffer
		filename := cp.filename(app)

		if len(envVars) == 0 {
			continue
		}

		// add env vars to buffer
		for k, v := range envVars {
			if _, err := fmt.Fprintf(&buffer, "%s=\"%s\"\n", k, v); err != nil {
				return err
			}
		}

		log.Infof("Writing to env file %s: %s", filename, strings.ReplaceAll(buffer.String(), "\n", " "))

		tmp := filename + ".tmp"
		err := os.WriteFile(tmp, buffer.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("failed to write %s  - %v", tmp, err)
		}

		err = os.Rename(tmp, filename)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s:%v", tmp, filename, err)
		}
	}

	return nil
}
