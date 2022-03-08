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
	"fmt"
	"log/syslog"
	"os"
)

var (
	slog  *syslog.Writer
	debug bool = false
)

// Debug writes the given input to syslog (sev=LOG_DEBUG) if snap `debug`
// configuration option is set to `true`.
// It formats similar to fmt.Sprint
func Debug(a ...interface{}) {
	if debug {
		slog.Debug(fmt.Sprint(a...))
	}
}

// Debugf writes the given input to syslog (sev=LOG_DEBUG) if snap `debug`
// configuration option is set to `true`.
// It formats similar to fmt.Sprintf
func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

// Error writes the given input to syslog (sev=LOG_ERROR).
// It formats similar to fmt.Sprint
func Error(a ...interface{}) {
	msg := fmt.Sprint(a...)
	slog.Err(msg)
	// print to stderr as well for snap command error output
	fmt.Fprint(os.Stderr, msg)
}

// Errorf writes the given input to syslog (sev=LOG_ERROR).
// It formats similar to fmt.Sprintf
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

// Info writes the given input to syslog (sev=LOG_INFO).
// It formats similar to fmt.Sprint
func Info(a ...interface{}) {
	slog.Info(fmt.Sprint(a...))
}

// Infof writes the given input to syslog (sev=LOG_INFO).
// It formats similar to fmt.Sprintf
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

// Warn writes the given input to syslog (sev=LOG_WARNING).
// It formats similar to fmt.Sprint
func Warn(a ...interface{}) {
	slog.Err(fmt.Sprint(a...))
}

// Warnf writes the given input to syslog (sev=LOG_WARNING).
// It formats similar to fmt.Sprintf
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

func init() {
	// TODO: use snapctl get debug to read this
	// Depends on https://github.com/canonical/edgex-snap-hooks/pull/26
	debug = true

	snap := os.Getenv("SNAP")
	if snap == "" {
		panic("SNAP environment variable is not set")
	}

	var err error
	slog, err = syslog.New(syslog.LOG_INFO, snap+":hook")
	if err != nil {
		panic(err)
	}
}
