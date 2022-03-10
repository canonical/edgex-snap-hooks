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

	"github.com/canonical/edgex-snap-hooks/v2/files"
	"github.com/canonical/edgex-snap-hooks/v2/log"
)

func getEnvVarName(setting string) string {
	result := strings.ToUpper(setting)
	result = strings.Replace(result, "-", "", -1)
	result = strings.Replace(result, ".", "_", -1)
	log.Infof("Mapping %s to %s", setting, result)
	return result
}

func setEnvVariable(b *bytes.Buffer, env string, value string) error {
	envvar := getEnvVarName(env)
	_, err := fmt.Fprintf(b, "export %s=%s\n", envvar, value)
	return err
}

func getEnvFilename(service string) string {

	// Handle security-* service naming. The service names in this
	// hook historically do not align with the actual binary commands.
	// As such, when handling configuration settings for them, we need
	// to translate the hook name to the actual binary name.
	if service == "security-proxy" {
		service = "security-proxy-setup"
	} else if service == "security-secret-store" {
		service = "security-secretstore-setup"
	}

	// The app-service-configurable snap is the one outlier snap that doesn't
	// include the service name in it's configuration path.
	var path string
	if files.SnapName == "edgex-app-service-configurable" {
		path = fmt.Sprintf("%s/res/%s.env", files.SnapDataConf, service)
	} else {
		path = fmt.Sprintf("%s/%s/res/%s.env", files.SnapDataConf, service, service)
	}
	return path
}

func writeEnvFile(b *bytes.Buffer, service string) error {

	filename := getEnvFilename(service)
	log.Infof("Writing settings to %s", filename)

	tmp := filename + ".tmp"
	err := ioutil.WriteFile(tmp, b.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s  - %v", tmp, err)
	}

	err = os.Rename(tmp, filename)
	if err != nil {
		return fmt.Errorf("failed to rename %s to %s", tmp, filename, err)
	}

	return nil
}

func setGlobalEnv(e string) error {
	log.Infof("Setting enviroment value %s", e)
	return nil
}
