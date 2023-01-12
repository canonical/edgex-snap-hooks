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

	"github.com/canonical/edgex-snap-hooks/v3/log"
	"github.com/canonical/edgex-snap-hooks/v3/snapctl"
)

type configOptions map[string]interface{}

type appOptions struct {
	Config    *configOptions `json:"config"`
	Autostart *bool          `json:"autostart"`
	// custom app options
	Proxy *proxyOptions `json:"proxy"`
}

type snapOptions struct {
	Apps   map[string]appOptions `json:"apps"`
	Config *configOptions        `json:"config"`
}

func getConfigMap(config configOptions) (map[string]string, error) {
	result := make(map[string]string)

	for env, value := range config {
		if err := flattenConfigJSON("", env, value, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Process the "config.<my.env.var>" configuration
//
//	-> setting env variable for all apps (e.g. DEBUG=true, SERVICE_SERVERBINDADDRESS=0.0.0.0)
func (cp *configProcessor) processGlobalConfigOptions(services []string) error {
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
		log.Debugf("No global config options")
		return nil
	}

	configuration, err := getConfigMap(*options.Config)
	if err != nil {
		return err
	}
	for _, service := range services {
		for env, value := range configuration {
			log.Debugf("Processing globally set env var for %s: %v=%v", service, env, value)
			if err := cp.addEnvVar(service, env, value); err != nil {
				return err
			}
		}
	}
	return nil
}

// Process the "apps.<app>.<custom.option>" where <custom.option> is not "config"
func ProcessAppCustomOptions(service string) error {
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

	log.Debugf("Processing custom options for service: %s", service)

	switch service {
	case "secrets-config":
		return processSecretsConfigOptions(options.Apps[service])
	}

	return nil
}

func validateAppConfigOptions(appConfigOptions map[string]appOptions, expectedServices []string) error {
	// make sure that set services in options are one of the expected services
	expected := make(map[string]bool)
	for _, s := range expectedServices {
		expected[s] = true
	}

	for setService, value := range appConfigOptions {
		if value.Config != nil && !expected[setService] {
			return fmt.Errorf("unsupported service in app config option: %s. Supported services are: %v",
				setService,
				expectedServices,
			)
		}
	}
	return nil
}

// Process the "apps.<app>.config.<my.env.var>" configuration
//
//	-> setting env var MY_ENV_VAR for an app
func (cp *configProcessor) processAppConfigOptions(services []string) error {
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

	err = validateAppConfigOptions(options.Apps, services)
	if err != nil {
		return err
	}

	// iterate through the known services in this snap
	for _, service := range services {
		log.Debugf("Processing service: %s", service)

		// get the configuration specified for each service
		// and create the environment override file
		appConfig := options.Apps[service]

		if appConfig.Config == nil {
			// no config options for this app
			continue
		}

		log.Debugf("Processing config: %v", appConfig.Config)
		configuration, err := getConfigMap(*appConfig.Config)

		log.Debugf("Processing flattened config: %v", configuration)
		if err != nil {
			return err
		}
		for env, value := range configuration {
			log.Debugf("Processing config option for %s: %v=%v", service, env, value)
			if err := cp.addEnvVar(service, env, value); err != nil {
				return err
			}
		}

	}
	return nil
}

var (
	// Snapd uses dots for hierarchy and hyphens as segment separators
	// These separators map to another character for environment variable names
	envSegmentSeparator   = "_"
	envHierarchySeparator = "_"
	configHierarchy       = false
)

// SetSegmentSeparator sets the separator used to replace hyphens in config.<x-y>
// Default is _
func SetSegmentSeparator(sep string) {
	envSegmentSeparator = sep
}

// SetHierarchySeparator sets the separator used to replace dots in config.<x.y>
// Default is _
func SetHierarchySeparator(sep string) {
	envHierarchySeparator = sep
}

// EnableConfigHierarchy is to allow config options such as config.<x.y> with
//
//	dots as the config key
func EnableConfigHierarchy() {
	configHierarchy = true
}

// ProcessConfig processes snap configuration which can be used to override
// app configuration via environment variables sourced by the snap
// service wrapper script.
// A service specific file (named <service>.env) is created in  the
// $SNAP_DATA/config/res directory.
// The settings can either be app-specific or apply to all services/apps in the snap
// a) snap set edgex-snap-name apps.<app>.config.<my.env.var>
//
//	-> sets env var MY_ENV_VAR for an app
//
// b) snap set edgex-snap-name config.<my.env.var>
//
//	-> sets env variable for all apps (e.g. DEBUG=true, SERVICE_SERVERBINDADDRESS=0.0.0.0)
func ProcessConfig(apps ...string) error {
	// uncomment to enable snap debugging
	// snapctl.Set("debug", "true")

	if len(apps) == 0 {
		return fmt.Errorf("empty apps list")
	}

	appOptionsStr, err := snapctl.Get("app-options").Run()
	if err != nil {
		return err
	}
	appOptions := (appOptionsStr == "true")

	log.Infof("Processing app options: %t", appOptions)

	isSet := func(v string) bool {
		return !(v == "" || v == "{}")
	}

	envOptions, err := snapctl.Get("env").Run()
	if err != nil {
		return err
	}
	if isSet(envOptions) {
		return fmt.Errorf("disallowed env options: %s", envOptions)
	}

	cp := newConfigProcessor(apps, configHierarchy, envHierarchySeparator, envSegmentSeparator)

	// process app-specific options
	if err := cp.processGlobalConfigOptions(apps); err != nil {
		return err
	}

	// process global options
	if err := cp.processAppConfigOptions(apps); err != nil {
		return err
	}

	if err := cp.writeEnvFiles(); err != nil {
		return err
	}

	return nil
}
