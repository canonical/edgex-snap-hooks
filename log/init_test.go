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

package log

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {

	t.Run("debug", func(t *testing.T) {
		// should be false by default
		require.False(t, debug)

		// set it to true and check
		output, err := exec.Command("snapctl", "set", "debug=true").CombinedOutput()
		assert.NoError(t, err, "Error setting config value via snapctl: %s", output)
		initialize()
		require.True(t, debug)

		// unset and re-check
		output, err = exec.Command("snapctl", "unset", "debug").CombinedOutput()
		assert.NoError(t, err, "Error setting config value via snapctl: %s", output)
		initialize()
		require.False(t, debug)
	})

	t.Run("global instance key", func(t *testing.T) {
		require.NotEmpty(t, snapInstanceKey)
	})

	t.Run("global syslog writer", func(t *testing.T) {
		require.NotNil(t, slog)
	})
}
