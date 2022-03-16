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
	"errors"
	"os"

	"github.com/canonical/edgex-snap-hooks/v2/log"
)

var (
	// Snap contains the value of the SNAP environment variable.
	Snap string
	// SnapConf contains the expanded path '$SNAP/config'.
	SnapConf string
	// SnapCommon contains the value of the SNAP_COMMON environment variable.
	SnapCommon string
	// SnapData contains the value of the SNAP_DATA environment variable.
	SnapData string
	// SnapDataConf contains the expanded path '$SNAP_DATA/config'.
	SnapDataConf string
	// SnapInst contains the value of the SNAP_INSTANCE_NAME environment variable.
	SnapInst string
	// SnapName contains the value of the SNAP_NAME environment variable.
	SnapName string
	// SnapRev contains the value of the SNAP_REVISION environment variable.
	SnapRev string
)

const (
	// AutostartConfig is a configuration key used indicate that a
	// service (application or device) should be autostarted on install
	AutostartConfig = "autostart"
	// EnvConfig is the prefix used for configure hook keys used for
	// EdgeX configuration overrides.
	EnvConfig = "env"
	// ProfileConfig is a configuration key that specifies a named
	// configuration profile
	ProfileConfig = "profile"

	snapEnv         = "SNAP"
	snapCommonEnv   = "SNAP_COMMON"
	snapDataEnv     = "SNAP_DATA"
	snapInstNameEnv = "SNAP_INSTANCE_NAME"
	snapNameEnv     = "SNAP_NAME"
	snapRevEnv      = "SNAP_REVISION"
)

// getEnvVars populates global variables for each of the SNAP*
// variables defined in the snap's environment
func getEnvVars() error {
	Snap = os.Getenv(snapEnv)
	if Snap == "" {
		return errors.New("SNAP is not set")
	}

	SnapCommon = os.Getenv(snapCommonEnv)
	if SnapCommon == "" {
		return errors.New("SNAP_COMMON is not set")
	}

	SnapData = os.Getenv(snapDataEnv)
	if SnapData == "" {
		return errors.New("SNAP_DATA is not set")
	}

	SnapInst = os.Getenv(snapInstNameEnv)
	if SnapInst == "" {
		return errors.New("SNAP_INSTANCE_NAME is not set")
	}

	SnapName = os.Getenv(snapNameEnv)
	if SnapName == "" {
		return errors.New("SNAP_NAME is not set")
	}

	SnapRev = os.Getenv(snapRevEnv)
	if SnapRev == "" {
		return errors.New("SNAP_REVISION_NAME is not set")
	}

	SnapConf = Snap + "/config"
	SnapDataConf = SnapData + "/config"

	return nil
}

func init() {
	if err := getEnvVars(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
