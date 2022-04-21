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

package options_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	hooks "github.com/canonical/edgex-snap-hooks/v2"
	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/options"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testService   = "test-service"
	testService2  = "test-service2"
	configEnabled = "config-enabled"
)

func TestProcessAppConfig(t *testing.T) {

	configDir := fmt.Sprintf("%s/%s/res/", env.SnapDataConf, testService)
	envFile := path.Join(configDir, testService+".env")
	os.MkdirAll(configDir, os.ModePerm)

	configDir2 := fmt.Sprintf("%s/%s/res/", env.SnapDataConf, testService2)
	envFile2 := path.Join(configDir2, testService2+".env")
	os.MkdirAll(configDir2, os.ModePerm)

	require.NoError(t, snapctl.Set("debug", "true").Run())
	log.Init()

	t.Cleanup(func() {
		assert.NoError(t, snapctl.Unset("apps", "config", "env").Run())
		assert.NoError(t, snapctl.Unset("debug").Run())
		assert.NoError(t, os.RemoveAll(configDir))
		assert.NoError(t, os.RemoveAll(configDir2))
	})

	t.Run("reject empty service list", func(t *testing.T) {
		require.Error(t, options.ProcessAppConfig())
	})

	t.Run("global options", func(t *testing.T) {
		const key, value = "config.x.y", "value"

		t.Cleanup(func() {
			assert.NoError(t, snapctl.Unset(key).Run())

			assert.NoError(t, os.RemoveAll(envFile))
			assert.NoError(t, os.RemoveAll(envFile2))
		})

		t.Run("reject without enabling", func(t *testing.T) {
			require.NoError(t, snapctl.Set(key, value).Run())

			require.Error(t, options.ProcessAppConfig(testService, testService2))
		})

		t.Run("set", func(t *testing.T) {
			require.NoError(t, snapctl.Set(configEnabled, "true").Run())
			t.Cleanup(func() {
				require.NoError(t, snapctl.Unset(configEnabled).Run())
				require.NoError(t, snapctl.Unset("config").Run())

				require.NoError(t, options.ProcessAppConfig(testService, testService2))

				// it should be removed from both env files
				require.Error(t, isInFile(envFile, "export X_Y=value"),
					"File content:\n%s", readFile(envFile))
				require.Error(t, isInFile(envFile2, "export X_Y=value"),
					"File content:\n%s", readFile(envFile2))
			})

			require.NoError(t, snapctl.Set(key, value).Run())

			require.NoError(t, options.ProcessAppConfig(testService, testService2))

			// both env files should have it
			require.NoError(t, isInFile(envFile, "export X_Y=value"),
				"File content:\n%s", readFile(envFile))
			require.NoError(t, isInFile(envFile2, "export X_Y=value"),
				"File content:\n%s", readFile(envFile2))
		})

		t.Run("unset", func(t *testing.T) {

		})
	})

	t.Run("single app options", func(t *testing.T) {
		const key, value = "apps." + testService + ".config.x.y", "value"

		t.Cleanup(func() {
			assert.NoError(t, snapctl.Unset(key).Run())
			assert.NoError(t, os.RemoveAll(envFile))
		})

		t.Run("set", func(t *testing.T) {
			require.NoError(t, snapctl.Set(configEnabled, "true").Run())
			t.Cleanup(func() {
				require.NoError(t, snapctl.Unset("apps").Run())
				require.NoError(t, snapctl.Unset(configEnabled).Run())

				require.NoError(t, options.ProcessAppConfig(testService, testService2))

				// it should be removed from the env file
				require.Error(t, isInFile(envFile, "export X_Y=value"),
					"File content:\n%s", readFile(envFile))
			})

			require.NoError(t, snapctl.Set(key, value).Run())

			require.NoError(t, options.ProcessAppConfig(testService, testService2))

			// first env file should have it
			require.NoError(t, isInFile(envFile, "export X_Y=value"),
				"File content:\n%s", readFile(envFile))

			// second env file should NOT have it
			require.Error(t, isInFile(envFile2, "export X_Y=value"),
				"File content:\n%s", readFile(envFile2))
		})

		// t.Run("unset", func(t *testing.T) {
		// 	require.NoError(t, snapctl.Unset(key, value).Run())

		// 	require.NoError(t, options.ProcessAppConfig(testService, testService2))

		// 	// it should be removed from the env file
		// 	require.Error(t, isInFile(envFile, "export X_Y=value"),
		// 		"File content:\n%s", readFile(envFile))
		// })
	})

	t.Run("Set mixed legacy options", func(t *testing.T) {
		const (
			legacyKey, legacyValue = "env.security-bootstrapper.add-registry-acl-roles", "legacy1,legacy2"
			key, value             = "apps.security-bootstrapper.config.add-registry-acl-roles", "legacy1,legacy2"
		)

		t.Cleanup(func() {
			//			assert.NoError(t, snapctl.Unset("env").Run())
			assert.NoError(t, snapctl.Unset(legacyKey).Run())
			assert.NoError(t, snapctl.Unset(key).Run())
		})
		t.Run("set", func(t *testing.T) {
			require.NoError(t, snapctl.Set(configEnabled, "true").Run())
			t.Cleanup(func() {
				assert.NoError(t, snapctl.Unset(configEnabled).Run())
			})

			require.NoError(t, snapctl.Set(legacyKey, legacyValue).Run())
			require.NoError(t, options.ProcessAppConfig("security-bootstrapper"))
			k, err := snapctl.Get(key).Run()
			require.Equal(t, k, value)
			require.NoError(t, err)

		})
	})

	t.Run("reject mixed legacy options", func(t *testing.T) {
		const (
			legacyKey, legacyValue = "env.core-data.service.host", "legacy"
			key, value             = "apps.core-data.config.x.y", "value"
		)

		configCoreDataDir := fmt.Sprintf("%s/core-data/res/", env.SnapDataConf)
		// envFileCoreData := path.Join(configCoreDataDir, "core-data.env")
		os.MkdirAll(configCoreDataDir, os.ModePerm)

		t.Cleanup(func() {
			assert.NoError(t, snapctl.Unset(legacyKey).Run())
			assert.NoError(t, snapctl.Unset(key).Run())

			assert.NoError(t, os.RemoveAll(envFile))
			assert.NoError(t, os.RemoveAll(configDir2))
		})

		t.Run("set", func(t *testing.T) {
			require.NoError(t, snapctl.Set(legacyKey, legacyValue).Run())
			require.NoError(t, snapctl.Set(key, value).Run())

			require.NoError(t, applyLegacyOptions("core-data"))
			require.Error(t, options.ProcessAppConfig(testService, "core-data"))
		})

	})

}

// utility testing functions

func isInFile(file string, line string) error {
	// read the whole file at once
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if strings.Contains(string(b), line) {
		return nil
	} else {
		return fmt.Errorf("Line %s not found in %s", line, file)
	}
}

func readFile(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func applyLegacyOptions(service string) error {
	envJSON, err := hooks.NewSnapCtl().Config(hooks.EnvConfig + "." + service)
	if err != nil {
		return fmt.Errorf("failed to read config options for %s: %v", service, err)
	}

	if envJSON != "" {
		if err := hooks.HandleEdgeXConfig(service, envJSON, nil); err != nil {
			return err
		}
	}
	return nil
}
