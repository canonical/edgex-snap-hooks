// -*- Mode: Go; indent-tabs-mode: t -*-

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
	"os"

	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

/*
	https://warthogs.atlassian.net/browse/EDGEX-133
	This task is to create a package which implements the behaviour described in EDGEX-78: Decide on config option schema for setting environment variablesIN REVIEW .

 	The package would be responsible for reading the config options and writing the environment variables in env files. It would be a sub-package under GitHub - canonical/edgex-snap-hooks: Snap hooks library for EdgeX snaps.

	It may also need to start/stop services when corresponding options are changed.

	https://warthogs.atlassian.net/browse/EDGEX-78

	The schema agreed on is:
	apps.<app>.<ENV_KEY> -> setting env variable for an app
	apps.<app>.<option> -> setting another option for CLI executation or CLI arg override
	apps.<app>.auto-start (boolean) -> turn auto start on/off by seting to true/false
	<ENV_KEY> -> setting env variable for all apps (e.g. DEBUG=true, SERVICE_SERVERBINDADDRESS=0.0.0.0)

	Deprecation:
		<env>.<app>.<env-key> → can be done with env injection
		<app> → can be done with auto-start option
		startup-duration, startup-interval → can be done with env injection

*/

func processDeprecatedSettings(json string) error {

	log.Infof("Processing deprecated settings: %s", json)

	return nil
}

func ProcessOptions() error {

	log.Info("Processing options")

	appsJSON, err := snapctl.Get("apps").Run()
	if err != nil {
		log.Errorf("Reading config 'apps' failed: %v", err)
		os.Exit(1)
	}

	if err = processAppSettings(appsJSON); err != nil {
		log.Errorf("Processing apps config failed: %v", err)
		os.Exit(1)
	}

	/*	envJSON, err := snapctl.Get("env").Run()
		if err != nil {
			log.Error("Reading config 'env' failed: %v", err)
			os.Exit(1)
		}

		if err := processDeprecatedSettings(envJSON); err != nil {
			return err
		}
	*/
	return nil

}
