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
	"bytes"
	"fmt"

	"github.com/canonical/edgex-snap-hooks/v2/log"
)

func isEnvironmentVariable(k string) bool {
	//	apps := (len(k) > 0 && k[0] >= 'A' && k[0] <= 'Z')
	//	return apps
	return true
}

func processAppSettings(jsonMap string) error {
	/*
		The schema agreed on is:
		apps.<app>.<ENV_KEY> -> setting env variable for an app
		apps.<app>.<option> -> setting another option for CLI executation or CLI arg override
		apps.<app>.auto-start (boolean) -> turn auto start on/off by seting to true/false
		<ENV_KEY> -> setting env variable for all apps (e.g. DEBUG=true, SERVICE_SERVERBINDADDRESS=0.0.0.0)
	*/

	log.Info("Processing app settings")
	serviceMap, err := unmarshal(jsonMap)
	if err != nil {
		return err
	}
	if serviceMap == nil {
		return nil
	}
	for k, v := range serviceMap {
		log.Infof("Setting configuration for app %s", k)
		if isValidService(k) {
			settings, err := getServiceSettingMap(v)
			if err != nil {
				return fmt.Errorf("Invalid configuration map: %v", err)
			}

			b := bytes.Buffer{}

			for env, value := range settings {

				if isEnvironmentVariable(env) {
					setEnvVariable(&b, env, value)
				} else {
					return fmt.Errorf("Invalid setting, %s = %s", env, value)
				}
			}
			log.Infof("Got settings for %s: %v", k, settings)
			writeEnvFile(&b, k)

		} else {
			return fmt.Errorf("%s is not a valid app name in this snap", k)
		}
	}

	return nil
}

/*


	// If autostart is not explicitly set, default to "no"
	// as only example service configuration and profiles
	// are provided by default.
	autostart, err := cli.Config(hooks.AutostartConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'autostart' failed: %v", err))
		os.Exit(1)
	}
	if autostart == "" {
		hooks.Debug("edgex-device-mqtt autostart is NOT set, initializing to 'no'")
		autostart = "no"
	}
	autostart = strings.ToLower(autostart)

	hooks.Debug(fmt.Sprintf("edgex-device-mqtt autostart is %s", autostart))

	// service is stopped/disabled by default in the install hook
	switch autostart {
	case "true":
		fallthrough
	case "yes":
		err = cli.Start("device-mqtt", true)
		if err != nil {
			hooks.Error(fmt.Sprintf("Can't start service - %v", err))
			os.Exit(1)
		}
	case "false":
		// no action necessary
	case "no":
		// no action necessary
	default:
		hooks.Error(fmt.Sprintf("Invalid value for 'autostart' : %s", autostart))
		os.Exit(1)
	}

*/

/*


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

v
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
				 securityProxyDeleteCurrentUserIfSet()
			 } else if jwtUsername != "" && jwtPublicKey != "" {
				 // else add a new user
				 err = securityProxyAddUser(jwtUsername, jwtUserID, jwtAlgorithm, jwtPublicKey)
				 if err != nil {
					 return err
				 }
			 }

			 if tlsCertificate == "" && tlsPrivateKey == "" {
				 // if the values have been set to "" then clear the semaphore so that a new cert can be set
				 securityProxyDeleteCurrentTLSCertIfSet()
			 } else if tlsCertificate != "" && tlsPrivateKey != "" {
				 // Set the TLS certificate and private key
				 err = securityProxySetTLSCertificate(tlsCertificate, tlsPrivateKey, tlsSNI)
				 if err != nil {
					 return err
				 }
			 }
		 }
	 }

 }
*/
