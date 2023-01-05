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
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/options"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testService  = "test-service"
	testService2 = "test-service2"
	appOptions   = "app-options"
)

func TestProcessConfig(t *testing.T) {
	// uncomment to cleanup previous mess
	// assert.NoError(t, snapctl.Unset("app-options", "config-enabled", "apps", "config").Run())

	configDir := fmt.Sprintf("%s/config/%s/res/", env.SnapData, testService)
	envFile := path.Join(configDir, testService+".env")
	os.MkdirAll(configDir, os.ModePerm)

	configDir2 := fmt.Sprintf("%s/config/%s/res/", env.SnapData, testService2)
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
		require.Error(t, options.ProcessConfig())
	})

	t.Run("global options", func(t *testing.T) {
		const key, value = "config.x-y", "value"

		t.Cleanup(func() {
			assert.NoError(t, snapctl.Unset("config").Run())

			assert.NoError(t, os.RemoveAll(envFile))
			assert.NoError(t, os.RemoveAll(envFile2))
		})

		t.Run("set+unset", func(t *testing.T) {
			require.NoError(t, snapctl.Set(appOptions, "true").Run())
			t.Cleanup(func() {
				require.NoError(t, snapctl.Unset("config").Run())
				require.NoError(t, options.ProcessConfig(testService, testService2))
				// disable config after processing once, otherwise the env files won't get cleaned up
				require.NoError(t, snapctl.Unset(appOptions).Run())

				require.False(t, fileExists(t, envFile), "Env file should not exist.")
				require.False(t, fileExists(t, envFile2), "Env file should not exist.")
			})

			require.NoError(t, snapctl.Set(key, value).Run())

			require.NoError(t, options.ProcessConfig(testService, testService2))

			// both env files should have it
			require.NoError(t, fileContains(t, envFile, `X_Y="value"`),
				"File content:\n%s", readFile(t, envFile))
			require.NoError(t, fileContains(t, envFile2, `X_Y="value"`),
				"File content:\n%s", readFile(t, envFile2))
		})

		t.Run("unset", func(t *testing.T) {

		})
	})

	t.Run("single app options", func(t *testing.T) {
		const key, value = "apps." + testService + ".config.x-y", "value"

		t.Cleanup(func() {
			assert.NoError(t, snapctl.Unset("apps").Run())
			assert.NoError(t, os.RemoveAll(envFile))
		})

		t.Run("set+unset", func(t *testing.T) {
			require.NoError(t, snapctl.Set(appOptions, "true").Run())
			t.Cleanup(func() {
				require.NoError(t, snapctl.Unset("apps").Run())
				require.NoError(t, options.ProcessConfig(testService, testService2))
				// disable config after processing once, otherwise the env files won't get cleaned up
				require.NoError(t, snapctl.Unset(appOptions).Run())

				require.False(t, fileExists(t, envFile), "Env file should not exist.")
			})

			require.NoError(t, snapctl.Set(key, value).Run())

			require.NoError(t, options.ProcessConfig(testService, testService2))

			// first env file should have it
			require.NoError(t, fileContains(t, envFile, `X_Y="value"`),
				"File content:\n%s", readFile(t, envFile))

			// second env file should NOT have it
			require.Error(t, fileContains(t, envFile2, `X_Y="value"`),
				"File content:\n%s", readFile(t, envFile2))
		})
	})

	t.Run("reject unknown app", func(t *testing.T) {
		const key, value = "apps.unknown.config.x-y", "value"

		require.NoError(t, snapctl.Set(appOptions, "true").Run())
		require.NoError(t, snapctl.Set(key, value).Run())
		t.Cleanup(func() {
			assert.NoError(t, snapctl.Unset(appOptions).Run())
			assert.NoError(t, snapctl.Unset(key).Run())
			require.NoError(t, snapctl.Unset("apps").Run())
		})

		err := options.ProcessConfig(testService, "core-data")
		assert.Error(t, err)
		require.Contains(t, err.Error(), "unsupported")

	})

	t.Run("reject bad keys", func(t *testing.T) {
		const app = "test-service"
		const key, value = "apps." + app + ".config.x-y", "value"

		require.NoError(t, snapctl.Set(appOptions, "true").Run())

		t.Cleanup(func() {
			require.NoError(t, snapctl.Unset("apps").Run())
			require.NoError(t, options.ProcessConfig(app))
			// disable config after processing once, otherwise the env files won't get cleaned up
			require.NoError(t, snapctl.Unset(appOptions).Run())
		})

		require.NoError(t, snapctl.Set(key, value).Run())
		require.NoError(t, options.ProcessConfig(app))

		// env file should have the X_Y
		require.NoError(t, fileContains(t, envFile, `X_Y="value"`),
			"File content:\n%s", readFile(t, envFile))

		// set something bad
		require.NoError(t, snapctl.Set("apps."+app+".config.dots.disallowed", value).Run())
		require.Error(t, options.ProcessConfig(app))

		// env file should still have the X_Y
		require.Error(t, fileContains(t, envFile, `DOTS_DISALLOWED="value"`),
			"File content:\n%s", readFile(t, envFile))
		require.NoError(t, fileContains(t, envFile, `X_Y="value"`),
			"File content:\n%s", readFile(t, envFile))
	})

	t.Run("hierarchy enabled", func(t *testing.T) {
		const app = "test-service"
		const key, value = "apps." + app + ".config.p-a-r-e-n-t.child", "value"

		require.NoError(t, snapctl.Set(appOptions, "true").Run())

		t.Cleanup(func() {
			require.NoError(t, snapctl.Unset("apps").Run())
			require.NoError(t, options.ProcessConfig(app))
			// disable config after processing once, otherwise the env files won't get cleaned up
			require.NoError(t, snapctl.Unset(appOptions).Run())
		})

		require.NoError(t, snapctl.Set(key, value).Run())

		// perform processing with custom configuration
		options.EnableConfigHierarchy()
		options.SetHierarchySeparator("__")
		require.NoError(t, options.ProcessConfig(app))

		// env file should have P_A_R_E_N_T__CHILD
		require.NoError(t, fileContains(t, envFile, `P_A_R_E_N_T__CHILD="value"`),
			"File content:\n%s", readFile(t, envFile))
	})
}

// utility testing functions

func fileExists(t *testing.T, file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		// file does not exist
		return false
	}
	t.Fatalf("Error checking if file exists: %s", err)
	return false
}

func fileContains(t *testing.T, file string, line string) error {
	if strings.Contains(readFile(t, file), line) {
		return nil
	} else {
		return fmt.Errorf("Line %s not found in %s", line, file)
	}
}

func readFile(t *testing.T, file string) string {
	b, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return "File not found: " + file
	} else if err != nil {
		t.Fatalf("Error reading file: %s", err)
	}
	return string(b)
}
