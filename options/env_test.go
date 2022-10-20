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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigKeyToEnvVar(t *testing.T) {

	t.Run("hierarchy disabled", func(t *testing.T) {
		var cp configProcessor

		v, err := cp.configKeyToEnvVar("x-y", "_", "_", false)
		require.NoError(t, err)
		require.Equal(t, "X_Y", v)

		v, err = cp.configKeyToEnvVar("x-y", "_", "__", false)
		require.NoError(t, err)
		require.Equal(t, "X__Y", v)

		_, err = cp.configKeyToEnvVar("x.y", "_", "_", false)
		require.Error(t, err)
	})

	t.Run("hierarchy enabled", func(t *testing.T) {
		var cp configProcessor

		v, err := cp.configKeyToEnvVar("x.y", "_", "_", true)
		require.NoError(t, err)
		require.Equal(t, "X_Y", v)

		v, err = cp.configKeyToEnvVar("x.y-z", "_", "__", true)
		require.NoError(t, err)
		require.Equal(t, "X_Y__Z", v)

		v, err = cp.configKeyToEnvVar("x.y-z", "___", "__", true)
		require.NoError(t, err)
		require.Equal(t, "X___Y__Z", v)
	})
}
