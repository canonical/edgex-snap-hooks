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
	"strconv"
)

func unmarshal(envJSON string) (map[string]interface{}, error) {
	if envJSON == "" {
		return nil, nil
	}

	var m map[string]interface{}
	err := json.Unmarshal([]byte(envJSON), &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall EdgeX config - %v : %v", envJSON, err)
	}

	return m, nil

}

func getServiceSettingMap(config interface{}) (map[string]string, error) {
	result := make(map[string]string)

	if err := flattenConfigJSON("", "", config, result); err != nil {
		return nil, err
	}

	return result, nil
}

// p is the current prefix of the config key being processed (e.g. "service", "security.auth")
// k is the key name of the current JSON object being processed
// vJSON is the current object
// flatConf is a map containing the configuration keys/values processed thus far
func flattenConfigJSON(p string, k string, vJSON interface{}, flatConf map[string]string) error {
	var mk string

	// top level keys don't include "env", so no separator needed
	if p == "" {
		mk = k
	} else {
		mk = fmt.Sprintf("%s.%s", p, k)
	}

	switch t := vJSON.(type) {
	case string:
		flatConf[mk] = t
	case bool:
		flatConf[mk] = strconv.FormatBool(t)
	case float64:
		flatConf[mk] = strconv.FormatFloat(t, 'f', -1, 64)
	case map[string]interface{}:

		for k, v := range t {
			err := flattenConfigJSON(mk, k, v, flatConf)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("internal error: invalid JSON configuration from snapd - prefix: %s key: %s obj: %v", p, k, t)
	}
	return nil
}
