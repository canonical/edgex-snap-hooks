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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/canonical/edgex-snap-hooks/v2/log"
)

type envVarOverrides struct {
	service  string
	filename string
	buffer   *bytes.Buffer
}

func getEnvVarFile(service string) *envVarOverrides {
	env := envVarOverrides{}
	env.service = service
	env.filename = env.getEnvFilename()
	env.buffer = &bytes.Buffer{}
	return &env
}

func (e *envVarOverrides) setEnvVariable(key string, value string) error {
	envKey, err := configKeyToEnvVar(key)
	if err != nil {
		return fmt.Errorf("error converting config key to environment variable key: %s", err)
	}
	log.Infof("Mapping %s to %s", key, envKey)
	_, err = fmt.Fprintf(e.buffer, "export %s=\"%s\"\n", envKey, value)
	return err
}

// convert my-var to MY_VAR
func configKeyToEnvVar(configKey string) (string, error) {
	if strings.Contains(configKey, ".") {
		return "", fmt.Errorf("config key must not contain dots: %s", configKey)
	}
	return strings.ReplaceAll(strings.ToUpper(configKey), "-", "_"), nil
}

func (e *envVarOverrides) getEnvFilename() string {

	// The app-service-configurable snap is the one outlier snap that doesn't
	// include the service name in it's configuration path.
	var path string
	if env.SnapName == "edgex-app-service-configurable" {
		path = fmt.Sprintf("%s/res/%s.env", env.SnapDataConf, e.service)
	} else {
		path = fmt.Sprintf("%s/%s/res/%s.env", env.SnapDataConf, e.service, e.service)
	}
	return path
}

func (e *envVarOverrides) writeEnvFile(append bool) error {
	buf := bytes.Buffer{}

	if append {
		current, err := ioutil.ReadFile(e.filename)
		if err == nil {
			buf.Write(current)
		}
	}
	buf.Write(e.buffer.Bytes())

	log.Infof("Writing settings to %s: %s", e.filename, strings.ReplaceAll(e.buffer.String(), "\n", " "))

	tmp := e.filename + ".tmp"
	err := ioutil.WriteFile(tmp, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s  - %v", tmp, err)
	}

	err = os.Rename(tmp, e.filename)
	if err != nil {
		return fmt.Errorf("failed to rename %s to %s:%v", tmp, e.filename, err)
	}

	return nil
}
