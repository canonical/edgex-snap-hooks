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
	"strings"

	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

type configOptions map[string]interface{}

type snapOptions struct {
	Apps   map[string]map[string]configOptions `json:"apps"`
	Config configOptions                       `json:"config"`
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
		overrides.writeEnvFile(false)
	}
	return nil
}

func migrateLegacyInternalOptions() error {

	clear := []string{"env.security-bootstrapper", "env.security-secret-store"}

	/*
		"env": {
			"security-proxy": {
				"public-key": "",
				"user": "username,user_id,algorithm",
				"tls-certificate": "",
				"tls-private-key": "",
				"tls-sni": ""
			}
		}
	*/

	// apps.secrets-config.proxy.admin-public-key
	// apps.secrets-config.proxy.tls.key
	// apps.secrets-config.proxy.tls.cert
	// apps.secrets-config.proxy.tls.snis
	const proxyUserAttributes = "{user,id,algorithm}"
	namespaceMap := map[string]string{
		"env.security-secret-store.add-secretstore-tokens": "apps.security-secretstore-setup.config.add-secretstore-tokens",
		"env.security-secret-store.add-known-secrets":      "apps.security-secretstore-setup.config.add-known-secrets",
		"env.security-bootstrapper.add-registry-acl-roles": "apps.security-bootstrapper.config.add-registry-acl-roles",
		"env.security-proxy.user":                          "apps.secrets-config.proxy.admin." + proxyUserAttributes,
		"env.security-proxy.public-key":                    "apps.secrets-config.proxy.admin.public-key",
		"env.security-proxy.tls-certificate":               "apps.secrets-config.proxy.tls.cert",
		"env.security-proxy.tls-private-key":               "apps.secrets-config.proxy.tls.key",
		"env.security-proxy.tls-sni":                       "apps.secrets-config.proxy.tls.snis",
		// "env.security-proxy.public-key":                    "apps.security-proxy-setup.admin-public-key",
		// "env.security-proxy.tls-certificate":               "apps.security-proxy-setup.tls-certificate",
		// "env.security-proxy.tls-private-key":               "apps.security-proxy-setup.tls-private-key",
		// "env.security-proxy.tls-sni":                       "apps.security-proxy-setup.tls-sni",
	}

	migrated := false
	for k, v := range namespaceMap {
		setting, err := snapctl.Get(k).Run()
		if err != nil {
			return err
		}
		if setting != "" {
			if err := snapctl.Unset(k).Run(); err != nil {
				return err
			}
			switch {
			case strings.HasSuffix(v, proxyUserAttributes):
				prefix := strings.TrimSuffix(v, proxyUserAttributes)

				// split into multiple options
				parts := strings.Split(setting, ",")
				user, id, algorithm := parts[0], parts[1], parts[2]

				if err := snapctl.Set(prefix+"user", user).Run(); err != nil {
					return err
				}
				if err := snapctl.Set(prefix+"id", id).Run(); err != nil {
					return err
				}
				if err := snapctl.Set(prefix+"algorithm", algorithm).Run(); err != nil {
					return err
				}

				log.Debugf("Migrated %s to %s", k, proxyUserAttributes)
			default:
				if err := snapctl.Set(v, setting).Run(); err != nil {
					return err
				}
				log.Debugf("Migrated %s to %s", k, v)
			}

			migrated = true
		}
	}

	if migrated {
		for _, s := range clear {
			if err := snapctl.Unset(s).Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

func processAppCustomOptions(service, key string, value configOptions) error {
	switch service {
	case "secrets-config":
		return processSecretsConfigOptions(key, value)
	default:
		return fmt.Errorf("Unknown custom option %s for service %s", key, service)
	}
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

	// TODO
	// reject unknown servies

	// iterate through the known services in this snap
	for _, service := range services {
		log.Debugf("Processing service: %s", service)

		// get the configuration specified for each service
		// and create the environment override file
		appConfig := options.Apps[service]
		log.Debugf("Processing appConfig: %v", appConfig)
		if appConfig != nil {
			for k, v := range appConfig {
				if k == "config" { // config overrides
					config := v
					log.Debugf("Processing config: %v", config)
					if config != nil {
						configuration, err := getConfigMap(config)

						log.Debugf("Processing configuration: %v", configuration)
						if err != nil {
							return err
						}
						overrides := getEnvVarFile(service)

						log.Debugf("Processing overrides: %v", overrides)
						for env, value := range configuration {

							log.Debugf("Processing overrides setEnvVariable: %v=%v", env, value)
							overrides.setEnvVariable(env, value)
						}
						overrides.writeEnvFile(true)
					}
				} else { // non-config options
					if err := processAppCustomOptions(service, k, v); err != nil {
						return err
					}
				}
			}
			// config := appConfig["config"]
			// log.Debugf("Processing config: %v", config)
			// if config != nil {
			// 	configuration, err := getConfigMap(config)

			// 	log.Debugf("Processing configuration: %v", configuration)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	overrides := getEnvVarFile(service)

			// 	log.Debugf("Processing overrides: %v", overrides)
			// 	for env, value := range configuration {

			// 		log.Debugf("Processing overrides setEnvVariable: %v %v", env, value)
			// 		overrides.setEnvVariable(env, value)
			// 	}
			// 	overrides.writeEnvFile(true)
			// }
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

	configEnabled, err := snapctl.Get("config.enabled").Run()
	if err != nil {
		return err
	}

	isSet := func(v string) bool {
		return !(v == "" || v == "{}")
	}

	if configEnabled != "true" {
		appsOptions, err := snapctl.Get("apps").Run()
		if err != nil {
			return err
		}
		globalOptions, err := snapctl.Get("config").Run()
		if err != nil {
			return err
		}
		if isSet(appsOptions) || isSet(globalOptions) {
			return fmt.Errorf(`'config.' and 'app.' options are allowed only when config.enabled is true.
Note: Setting config.enabled=true will convert the following legacy 'env.' options:
	- env.security-secret-store.add-secretstore-tokens
	- env.security-secret-store.add-known-secrets
	- env.security-bootstrapper.add-registry-acl-roles
	- env.security-proxy.user
	- env.security-proxy.public-key
	- env.security-proxy.tls-certificate
	- env.security-proxy.tls-private-key
	- env.security-proxy.tls-sni
All other legacy 'env.' options will be unset!`)
		} else {
			// do nothing
			return nil
		}
	}

	err = migrateLegacyInternalOptions()
	if err != nil {
		return err
	}

	if err := processGlobalConfigOptions(services); err != nil {
		return err
	}

	// The app-specific options have higher precedence.
	// They should be processed last to end up at the bottom of the .env file
	// 	and override global environment variables.
	if err := processAppConfigOptions(services); err != nil {
		return err
	}

	return nil

}
