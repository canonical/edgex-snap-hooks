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

package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvVars(t *testing.T) {
	// Arrange
	os.Setenv(snapEnv, "/snap/testsnap/x1")
	os.Setenv(snapCommonEnv, "/snap/testsnap/common")
	os.Setenv(snapDataEnv, "/var/snap/testsnap/x1")
	os.Setenv(snapInstNameEnv, "testsnap")
	os.Setenv(snapRevEnv, "2112")

	// Test
	err := getEnvVars()

	// Assert values
	assert.Nil(t, err)
	assert.Equal(t, Snap, "/snap/testsnap/x1")
	assert.Equal(t, SnapCommon, "/snap/testsnap/common")
	assert.Equal(t, snapNameEnv, "SNAP_NAME")
	assert.Equal(t, SnapData, "/var/snap/testsnap/x1")
	assert.Equal(t, SnapInst, "testsnap")
	assert.Equal(t, SnapRev, "2112")
	// assert.Equal(t, SnapConf, "/snap/testsnap/x1/config")
	// assert.Equal(t, SnapDataConf, "/var/snap/testsnap/x1/config")
}
