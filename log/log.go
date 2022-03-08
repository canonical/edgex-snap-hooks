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

	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

var (
	slog            *syslog.Writer
	debug           bool
	snapInstanceKey string // used as default syslog tag and tag prefix
	tag             string // syslog tag and staderr prefix
)

func init() {
	debug = (os.Getenv("DEBUG") == "true")
	// snap config option overrides environment variable
	if !debug {
		snapctl.Unset("debug").Run()
		value, err := snapctl.Get("debug").Run()
		if err != nil {
			Stderr(err)
			os.Exit(1)
		}
		debug = (value == "true")
	}

	snapInstanceKey = os.Getenv("SNAP_INSTANCE_NAME")
	if snapInstanceKey == "" {
		Stderr("SNAP_INSTANCE_NAME environment variable not set.")
		os.Exit(1)
	}
	tag = snapInstanceKey

	if err := setupSyslogWriter(tag); err != nil {
		Stderr(err)
		os.Exit(1)
	}

	Debugf("debug=%t", debug)
}

func setupSyslogWriter(tag string) error {
	writer, err := syslog.New(syslog.LOG_INFO, tag)
	if err != nil {
		return err
	}
	// switch to new global writer only if no error
	slog = writer
	return nil
}

// SetComponentName adds a component name to syslog tag as "my-snap.<component>"
// The default tag is just "my-snap", read from the snap environment.
// This function is NOT thread-safe. It should not be called concurrently with
// the other logging functions of this package.
func SetComponentName(component string) {
	// update global value
	tag = snapInstanceKey + "." + component
	Debugf("Changing syslog tag to: %s", tag)

	if err := setupSyslogWriter(tag); err != nil {
		Errorf("Error changing syslog tag: %s", err)
	}
}

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
	Stderr(a...)
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

// Stderr writes the given input to standard error.
// It formats similar to fmt.Sprintf, adds the global tag as prefix, and appends
// a newline
func Stderr(a ...interface{}) {
	fmt.Fprintf(os.Stderr, tag+": %s\n", a...)
}

// Stderrf writes the given input to standard error.
// It formats similar to fmt.Sprintf and calls log.Stderr which appends a newline
func Stderrf(format string, a ...interface{}) {
	Stderr(fmt.Sprintf(format, a...))
}
