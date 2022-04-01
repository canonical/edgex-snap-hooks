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
	"encoding/json"
	"fmt"

	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

type snapOptions struct {
	Apps   map[string]map[string]map[string]interface{} `json:"apps"`
	Config map[string]interface{}                       `json:"config"`
}

func getConfigMap(config map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string)

	for env, value := range config {
		if err := flattenConfigJSON("", env, value, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Process the "config.<my.env.var>" configuration
//	 -> setting env variable for all apps (e.g. DEBUG=true, SERVICE_SERVERBINDADDRESS=0.0.0.0)
func processGlobalConfigOptions(services []string) error {
	var options snapOptions

	jsonString, err := snapctl.Get("config").Document().Run()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonString), &options)
	if err != nil {
		return err
	}

	if options.Config == nil {
		log.Debugf("No global configuration settings")
		return nil
	}

	configuration, err := getConfigMap(options.Config)
	if err != nil {
		return err
	}
	for _, service := range services {
		overrides := getEnvVarFile(service)
		for env, value := range configuration {
			overrides.setEnvVariable(env, value)
		}
		overrides.writeEnvFile(true)
	}
	return nil
}

// Process the "apps.<app>.config.<my.env.var>" configuration
//	-> setting env var MY_ENV_VAR for an app
func processAppConfigOptions(services []string) error {
	var options snapOptions

	// get the 'apps' json structure
	jsonString, err := snapctl.Get("apps").Document().Run()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonString), &options)
	if err != nil {
		return err
	}
	// iterate through the known services in this snap
	for _, service := range services {
		log.Debugf("Processing service:%s", service)

		// get the configuration specified for each service
		// and create the environment override file
		appConfig := options.Apps[service]
		log.Debugf("Processing appConfig:%v", appConfig)
		if appConfig != nil {
			config := appConfig["config"]
			log.Debugf("Processing config:%v", config)
			if config != nil {
				configuration, err := getConfigMap(config)

				log.Debugf("Processing configuration:%v", configuration)
				if err != nil {
					return err
				}
				overrides := getEnvVarFile(service)

				log.Debugf("Processing overrides:%v", overrides)
				for env, value := range configuration {
					log.Debugf("Processing overrides setEnvVariable:%v %v", env, value)
					overrides.setEnvVariable(env, value)
				}
				overrides.writeEnvFile(true)
			}
		}
	}
	return nil
}

// ProcessAppConfig processes snap configuration which can be used to override
// edgexfoundry configuration via environment variables sourced by the snap
// service wrapper script.
// A service specific file (named <service>.env) is created in  the
// $SNAP_DATA/config/res directory.
// The settings can either be app-specific or apply to all services/apps in the snap
// a) snap set edgex-snap-name apps.<app>.config.<my.env.var>
//	-> sets env var MY_ENV_VAR for an app
// b) snap set edgex-snap-name config.<my.env.var>
//	-> sets env variable for all apps (e.g. DEBUG=true, SERVICE_SERVERBINDADDRESS=0.0.0.0)
func ProcessAppConfig(services ...string) error {

	if len(services) == 0 {
		return fmt.Errorf("empty service list")
	}

	if err := processGlobalConfigOptions(services); err != nil {
		return err
	}

	if err := processAppConfigOptions(services); err != nil {
		return err
	}

	return nil

}
