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

package hooks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/options"
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

// CtlCli is the test obj for overridding functions
type CtlCli struct{}

// SnapCtl interface provides abstration for unit testing
type SnapCtl interface {
	Config(key string) (string, error)
	SetConfig(key string, val string) error
	UnsetConfig(key string) error
	Stop(svc string, disable bool) error
}

// CopyFile copies a file within the snap
func CopyFile(srcPath, destPath string) error {

	inFile, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// TODO: check file perm
	err = ioutil.WriteFile(destPath, inFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CopyDir copies a whole directory recursively
// snippet from https://blog.depa.do/post/copy-files-and-directories-in-go
func CopyDir(srcPath string, dstPath string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	srcinfo, err = os.Stat(srcPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dstPath, srcinfo.Mode())
	if err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(srcPath); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(srcPath, fd.Name())
		dstfp := path.Join(srcPath, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFileReplace copies a file within the snap and replaces strings using
// the string/replace values in the rStrings parameter.
func CopyFileReplace(srcPath, destPath string, rStrings map[string]string) error {

	inFile, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	rStr := string(inFile)
	for k, v := range rStrings {
		rStr = strings.Replace(rStr, k, v, 1)
	}

	// TODO: check file perm
	outBytes := []byte(rStr)
	err = ioutil.WriteFile(destPath, outBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Deprecated: use log.Debug or log.Debugf
func Debug(msg string) {
	log.Debug(msg)
}

// Deprecated: use log.Error or log.Errorf
func Error(msg string) {
	log.Error(msg)
}

// Deprecated: use log.Info or log.Infof
func Info(msg string) {
	log.Info(msg)
}

// Deprecated: use log.Warn or log.Warnf
func Warn(msg string) {
	log.Warn(msg)
}

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

// Deprecated: init function is called on package import.
func Init(setDebug bool, snapName string) error {
	return nil
}

func init() {
	if err := getEnvVars(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// NewSnapCtl returns a normal runtime client
func NewSnapCtl() *CtlCli {
	return &CtlCli{}
}

// Config uses snapctl to get a value from a key, or returns error.
func (cc *CtlCli) Config(key string) (string, error) {
	output, err := exec.Command("snapctl", "get", key).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("snapctl get failed for %s: %s: %s", key, err, output)
	}
	return strings.TrimSpace(string(output)), nil
}

// SetConfig uses snapctl to set a config value from a key, or returns error.
func (cc *CtlCli) SetConfig(key string, val string) error {
	output, err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, val)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl set failed for %s: %s: %s", key, err, output)
	}
	return nil
}

// UnsetConfig uses snapctl to unset a config value from a key
func (cc *CtlCli) UnsetConfig(key string) error {
	output, err := exec.Command("snapctl", "unset", key).CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl unset failed for %s: %s: %s", key, err, output)
	}
	return nil
}

// Start uses snapctrl to start a service and optionally enable it
func (cc *CtlCli) Start(svc string, enable bool) error {
	var cmd *exec.Cmd

	name := SnapName + "." + svc
	if enable {
		cmd = exec.Command("snapctl", "start", "--enable", name)
	} else {
		cmd = exec.Command("snapctl", "start", name)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl start %s failed: %s: %s", name, err, output)
	}

	return nil
}

// StartMultiple uses snapctl to start one or more services and optionally enable all
func (cc *CtlCli) StartMultiple(enable bool, services ...string) error {
	if len(services) == 0 {
		return fmt.Errorf("no services set to start")
	}

	args := []string{"start"}

	if enable {
		args = append(args, "--enable")
	}

	for _, s := range services {
		args = append(args, SnapName+"."+s)
	}

	output, err := exec.Command("snapctl", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl start failed: %s: %s", err, output)
	}

	return nil
}

// Stop uses snapctrl to stop a service and optionally disable it
func (cc *CtlCli) Stop(svc string, disable bool) error {
	var cmd *exec.Cmd

	name := SnapName + "." + svc
	if disable {
		cmd = exec.Command("snapctl", "stop", "--disable", name)
	} else {
		cmd = exec.Command("snapctl", "stop", name)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("snapctl stop %s failed: %s: %s", name, err, output)
	}

	return nil
}

// service status object
type service struct {
	name    string
	enabled bool
	active  bool
	notes   string
}

// services uses snapctl to get the list of services
func (cc *CtlCli) services() ([]service, error) {

	cmd := exec.Command("snapctl", "services")

	std, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("snapctl services failed: %s: %s", err, std)
	}

	scanner := bufio.NewScanner(bytes.NewReader(std))

	// throw away the header:
	// Service   Startup   Current   Notes
	scanner.Scan()

	var services []service
	for scanner.Scan() {
		line := scanner.Text()

		// Split by whitespaces up to four parts.
		// The last part is for notes which may contain spaces in itself.
		cells := regexp.MustCompile("[[:space:]]+").Split(line, 4)
		if len(cells) != 4 {
			return nil, fmt.Errorf("snapctl services: error parsing output: unexpected number of columns")
		}
		service := service{
			name:  cells[0],
			notes: cells[3],
		}
		if cells[1] == "enabled" {
			service.enabled = true
		}
		if cells[2] == "active" {
			service.active = true
		}
		services = append(services, service)
	}

	Info(fmt.Sprintf("snapctl services: %#v", services))

	return services, nil
}

// EnabledServices uses snapctl to get the list of enabled services
func (cc *CtlCli) EnabledServices() ([]string, error) {
	services, err := cc.services()
	if err != nil {
		return nil, err
	}

	var enabledServices []string
	for _, service := range services {
		if service.enabled {
			enabledServices = append(enabledServices, service.name)
		}
	}

	return enabledServices, nil
}

// p is the current prefix of the config key being processed (e.g. "service", "security.auth")
// k is the key name of the current JSON object being processed
// vJSON is the current object
// flatConf is a map containing the configuration keys/values processed thus far
func flattenConfigJSON(p string, k string, vJSON interface{}, flatConf map[string]string) {
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
			flattenConfigJSON(mk, k, v, flatConf)
		}
	default:
		panic(fmt.Sprintf("internal error: invalid JSON configuration from snapd - prefix: %s key: %s obj: %v", p, k, t))
	}
}

// This func checks the given key for a service-specific prefix
// delimited by a '/'. The prefix can either be a service name or a
// CSV service list. If found, the prefix is compared against
// the specified service parameter. If the prefix doesn't match
// one of the given service names, then false is returned. The string
// retval is the incoming k parameter stripped of any prefix.
func checkForServiceSpecificKey(k, service string) (bool, string) {
	var noPrefixEnv = k

	subStrs := strings.Split(k, "/")
	if len(subStrs) == 2 {
		noPrefixEnv = subStrs[1]
		servicesPrefix := strings.Split(subStrs[0], ",")
		for _, servicePrefix := range servicesPrefix {
			if servicePrefix == service {
				return true, noPrefixEnv
			}
		}
		return false, noPrefixEnv
	}

	return true, noPrefixEnv
}

func getConfigEnvVar(k string, extraConf map[string]string) (string, bool) {
	var env string

	env, ok := ConfToEnv[k]
	if ok {
		return env, true
	}

	if extraConf != nil {
		env, ok = extraConf[k]
		if ok {
			return env, true
		}
	}

	return env, false
}

// HandleEdgeXConfig processes snap configuration which can be used to override
// edgexfoundry configuration via environment variables sourced by the snap
// service wrapper script. The parameter service is used to create a new service
// specific file (named <service>.env) in the $SNAP_DATA/config/res directory of
// the service. The parameter envJSON is a JSON document holding the service's
// configuration as returned by snapd. The parameter extraConfig is a map of
// additional configuration keys supported by the snap. For example the following
// configuration option:
//
// [Driver]
// MyDriverOption = "foo"
//
// ...would require an entry in extraConf like this:
//
// extraConf["driver.mydriveroption"]"DRIVER_MYDRIVEROPTION"
//
// .. and would be set like this for a device or application service:
//
// ```
// $ sudo snap set mysnap env.driver.mydriveroption="foo"
// ```
//
func HandleEdgeXConfig(service, envJSON string, extraConf map[string]string) error {

	if envJSON == "" {
		return nil
	}

	var m map[string]interface{}
	var flatConf = make(map[string]string)

	err := json.Unmarshal([]byte(envJSON), &m)
	if err != nil {
		return fmt.Errorf("failed to unmarshall EdgeX config - %v", err)
	}

	for k, v := range m {
		flattenConfigJSON("", k, v, flatConf)
	}

	b := bytes.Buffer{}

	var jwtUsername, jwtUserID, jwtAlgorithm, jwtPublicKey string
	var tlsCertificate, tlsPrivateKey, tlsSNI string

	for k, v := range flatConf {

		// TODO: extract the security-proxy logic into its own function

		// a couple of special cases for security-proxy, to create an user/token and set the TLS cert.
		// This uses the standard naming schema but doesn't actually use environment variables
		if service == "security-proxy" {
			value := strings.TrimSpace(v)
			// These config options are read and validated in the loop,
			// but handled collectively afterwards
			switch k {
			case "user":
				if value != "" {
					s := strings.Split(value, ",")
					if len(s) != 3 {
						return fmt.Errorf("security-proxy.user expects a value containing 'username,userID,algorithm'. Example: 'me,1234,ES256' but got " + fmt.Sprint(s))
					}
					jwtUsername = strings.TrimSpace(s[0])
					jwtUserID = strings.TrimSpace(s[1])
					jwtAlgorithm = strings.ToUpper(strings.TrimSpace(s[2]))
					if jwtAlgorithm != "ES256" && jwtAlgorithm != "RS256" {
						return fmt.Errorf("invalid algorithm value - should be ES256 or RS256")
					}
				}
				continue
			case "public-key":
				jwtPublicKey = value
				continue
			case "tls-certificate":
				tlsCertificate = value
				continue
			case "tls-private-key":
				tlsPrivateKey = value
				continue
			case "tls-sni":
				tlsSNI = value
				continue
			}
		}

		env, ok := getConfigEnvVar(k, extraConf)
		if !ok {
			return errors.New("invalid EdgeX config option - " + k)
		}

		// checkForServiceSpecificKey() checks the env var for a service
		// prefix, and if it finds one, ensures that it matches the
		// current service being handled. If the match fails, then the
		// key is ignored
		ok, envNoPrefix := checkForServiceSpecificKey(env, service)
		if !ok {
			// TODO: should this be an error or warn OK?
			Warn(fmt.Sprintf("Invalid key %s specified for %s", k, service))
			continue
		}

		_, err := fmt.Fprintf(&b, "export %s=%s\n", envNoPrefix, v)
		if err != nil {
			return err
		}
	}

	// install-mode is set in the install hook of edgexfoundry
	installMode, err := NewSnapCtl().Config("install-mode")
	if err != nil {
		return fmt.Errorf("failed to read 'install-mode': %s", err)
	}

	// post-startup config handling
	if installMode != "defer-startup" {
		if service == "security-proxy" {
			if jwtUsername == "" && jwtPublicKey == "" {
				// if the values have been set to "" then delete the current user
				options.SecurityProxyDeleteCurrentUserIfSet()
			} else if jwtUsername != "" && jwtPublicKey != "" {
				// else add a new user
				err = options.SecurityProxyAddUser(jwtUsername, jwtUserID, jwtAlgorithm, jwtPublicKey)
				if err != nil {
					return err
				}
			}

			if tlsCertificate == "" && tlsPrivateKey == "" {
				// if the values have been set to "" then clear the semaphore so that a new cert can be set
				options.SecurityProxyDeleteCurrentTLSCertIfSet()
			} else if tlsCertificate != "" && tlsPrivateKey != "" {
				// Set the TLS certificate and private key
				err = options.SecurityProxySetTLSCertificate(tlsCertificate, tlsPrivateKey, tlsSNI)
				if err != nil {
					return err
				}
			}
		}
	}

	// Handle security-* service naming. The service names in this
	// hook historically do not align with the actual binary commands.
	// As such, when handling configuration settings for them, we need
	// to translate the hook name to the actual binary name.
	if service == "security-proxy" {
		service = "security-proxy-setup"
	} else if service == "security-secret-store" {
		service = "security-secretstore-setup"
	}

	// The app-service-configurable snap is the one outlier snap that doesn't
	// include the service name in it's configuration path.
	var path string
	if SnapName == "edgex-app-service-configurable" {
		path = fmt.Sprintf("%s/res/%s.env", SnapDataConf, service)
	} else {
		path = fmt.Sprintf("%s/%s/res/%s.env", SnapDataConf, service, service)
	}

	tmp := path + ".tmp"
	err = ioutil.WriteFile(tmp, b.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s.env file - %v", service, err)
	}

	err = os.Rename(tmp, path)
	if err != nil {
		return fmt.Errorf("failed to rename %s.env.tmp file - %v", service, err)
	}

	return nil
}
