// -*- Mode: Go; indent-tabs-mode: t -*-

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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/canonical/edgex-snap-hooks/v2/log"
)

const (
	userSemaphoreFile  = ".secrets-config-user"
	tlsSemaphoreFile   = ".secrets-config-tls"
	kongAdminTokenFile = "kong-admin-jwt"
	// proxy admin user
	defaultProxyUsername = "admin" // identical to object key
	defaultProxyUserId   = "1"
	defaultProxyUserAlgo = "ES256"
)

type proxyOptions struct {
	Admin struct {
		PublicKey string `json:"public-key"`
	} `json:"admin"` // identical to username

	TLS struct {
		Cert string `json:"cert"`
		Key  string `json:"key"`
		SNIs string `json:"snis"`
	} `json:"tls"`
}

// Write a security-proxy file
func securityProxyWriteFile(filename, contents string) (path string, err error) {
	path = fmt.Sprintf("%s/secrets/security-proxy-setup/%s", env.SnapData, filename)
	err = ioutil.WriteFile(path, []byte(contents), 0600)
	if err == nil {
		log.Debugf("Wrote file '%s'", path)
	} else {
		err = fmt.Errorf("failed to write file %s - %v", path, err)
	}
	return
}

// Read a security-proxy file
func securityProxyReadFile(filename string) (contents string, err error) {
	path := fmt.Sprintf("%s/secrets/security-proxy-setup/%s", env.SnapData, filename)
	bytes, err := ioutil.ReadFile(path)
	if err == nil {
		contents = string(bytes)
		log.Debugf("Read file '%s'", path)
	} else {
		err = fmt.Errorf("failed to read file %s - %v", path, err)
	}
	return
}

// Delete a security-proxy semaphore file
func securityProxyRemoveSemaphore(filename string) (err error) {
	path := fmt.Sprintf("%s/secrets/security-proxy-setup/%s", env.SnapData, filename)
	err = os.Remove(path)
	if err == nil {
		log.Debug("Removed file '" + path + "'")
	} else {
		log.Debugf("Could not remove file '%s' : %v", path, err)
	}
	return
}

// Execute the secrets-config tool with the given arguments
func securityProxyExecSecretsConfig(args []string) error {
	service := "security-proxy-setup"
	cmdSecretsConfig := exec.Command("secrets-config", args...)
	cmdSecretsConfig.Dir = fmt.Sprintf("%s/%s", env.SnapDataConf, service)
	out, err := cmdSecretsConfig.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing 'secrets-config %s': error=%s: output=%s",
			strings.Join(args, " "), err, out)
	}
	log.Debugf("Executed 'secrets-config %s': output=%s", strings.Join(args, " "), out)
	return nil
}

// Remove the semaphore, so that we can set the certificate again
func securityProxyDeleteCurrentTLSCertIfSet() error {
	return securityProxyRemoveSemaphore(tlsSemaphoreFile)
}

// Delete the current user - if one has been set up
func securityProxyDeleteCurrentUserIfSet() error {
	// if no user has been set up, then ignore the request
	username, err := securityProxyReadFile(userSemaphoreFile)
	if err != nil {
		log.Debug("proxy: No user has been set up")
		return nil
	}

	kongAdminToken, err := securityProxyReadFile(kongAdminTokenFile)
	if err != nil {
		return err
	}

	args := []string{"proxy", "deluser", "--user", username, "--jwt", kongAdminToken}
	err = securityProxyExecSecretsConfig(args)
	if err != nil {
		return fmt.Errorf("error executing secrets-config command: %s", err)
	}

	securityProxyRemoveSemaphore(userSemaphoreFile)
	// log.Debug("Executed \"secrets-config " + fmt.Sprint(args) + "\" result=" + string(out))
	log.Info("proxy: Removed current user")
	return nil
}

// Set up the proxy with the specified user.
func securityProxyAddUser(jwtUsername, jwtUserID, jwtAlgorithm, jwtPublicKey string) error {
	currentUser, err := securityProxyReadFile(userSemaphoreFile)
	if err == nil && currentUser != "" {
		if currentUser == jwtUsername {
			//	If a user has already been set - and it's the same user - then silently ignore the request
			log.Debug("proxy: Ignoring request to set up same user again")
			return nil
		} else {
			// if this is a different user, then return an error
			return fmt.Errorf("the proxy user has already been set. To add a new user, first delete the current user by setting 'user' and 'public-key' to an empy string")
		}
	}

	publicKeyFilePath, err := securityProxyWriteFile("jwt-user-public-key.pem", jwtPublicKey)
	if err != nil {
		return err
	}

	kongAdminToken, err := securityProxyReadFile(kongAdminTokenFile)
	if err != nil {
		return err
	}

	args := []string{"proxy", "adduser", "--token-type", "jwt", "--user", jwtUsername, "--id", jwtUserID, "--algorithm", jwtAlgorithm, "--public_key", publicKeyFilePath, "--jwt", kongAdminToken}
	err = securityProxyExecSecretsConfig(args)
	if err != nil {
		return fmt.Errorf("error executing secrets-config command: %s", err)
	}
	_, err = securityProxyWriteFile(userSemaphoreFile, jwtUsername)
	if err != nil {
		return err
	}
	log.Info("proxy: Added new user")
	return nil
}

// Set the TLS certificate. If a certificate has already been set then silently ignore the request
func securityProxySetTLSCertificate(tlsCertificate, tlsPrivateKey, tlsSNIs string) error {
	_, err := securityProxyReadFile(tlsSemaphoreFile)
	if err == nil {
		log.Debug("The TLS certificate has already been set. To set it again, first set tls-certificate and tls-private-key to an empty string")
		return nil
	}
	tlsCertFilename, err := securityProxyWriteFile("tls-certificate.pem", tlsCertificate)
	if err != nil {
		return err
	}
	tlsPrivateKeyFilename, err := securityProxyWriteFile("tls-private-key.pem", tlsPrivateKey)
	if err != nil {
		return err
	}

	kongAdminToken, err := securityProxyReadFile(kongAdminTokenFile)
	if err != nil {
		return err
	}

	args := []string{"proxy", "tls", "--incert", tlsCertFilename, "--inkey", tlsPrivateKeyFilename, "--admin_api_jwt", kongAdminToken}
	if tlsSNIs != "" {
		args = append(args, "--snis", tlsSNIs)
	}

	err = securityProxyExecSecretsConfig(args)
	if err != nil {
		return fmt.Errorf("failed to set TLS certificate: %s", err)
	}

	_, err = securityProxyWriteFile(tlsSemaphoreFile, "TLS certificate set")
	if err != nil {
		return err
	}
	log.Info("New TLS Certificate and private key set")
	return nil
}

func validateSecretsConfigProxyOptions(proxyObj proxyOptions) error {

	// TODO
	// flatten and reject unknown keys

	// const (
	// 	keyAdminPubKey = "admin.public-key"
	// 	keyTLSCert     = "tls.cert"
	// 	keyTLSKey      = "tls.key"
	// 	keyTLSSNIs     = "tls.snis"
	// )

	// supportedKeys := map[string]bool{
	// 	keyAdmin:       true,
	// 	keyAdminPubKey: true,
	// 	keyTLS:         true,
	// 	keyTLSCert:     true,
	// 	keyTLSKey:      true,
	// 	keyTLSSNIs:     true,
	// }

	// for k := range proxyOptions {
	// 	supported := supportedKeys[k]
	// 	if !supported {
	// 		return fmt.Errorf("Unsupported secrets-config proxy key: %s", k)
	// 	}
	// }

	// validate with object
	o := &proxyObj

	if (o.TLS.Cert != "" && o.TLS.Key == "") ||
		(o.TLS.Cert == "" && o.TLS.Key != "") {
		return fmt.Errorf("proxy cert and key must be set together")
	}

	if (o.TLS.Cert == "" || o.TLS.Key == "") && o.TLS.SNIs != "" {
		return fmt.Errorf("proxy snis must be set together with proxy cert and key")
	}

	return nil
}

func processSecretsConfigProxyOptions(proxy *proxyOptions) error {
	if proxy == nil {
		// no proxy options
		return nil
	}

	log.Debugf("Processing secrets-config proxy: %v", proxy)

	err := validateSecretsConfigProxyOptions(*proxy)
	if err != nil {
		return fmt.Errorf("error validating secrets-config proxy options: %s", err)
	}

	// process the admin user
	if proxy.Admin.PublicKey != "" {
		if err := securityProxyAddUser(
			defaultProxyUsername,
			defaultProxyUserId,
			defaultProxyUserAlgo,
			proxy.Admin.PublicKey); err != nil {
			return fmt.Errorf("error adding user: %s", err)
		}
	} else {
		if err := securityProxyDeleteCurrentUserIfSet(); err != nil {
			return fmt.Errorf("error deleting user: %s", err)
		}
	}

	// process TLS certificate
	if proxy.TLS.Cert != "" && proxy.TLS.Key != "" {
		err := securityProxySetTLSCertificate(
			proxy.TLS.Cert,
			proxy.TLS.Key,
			proxy.TLS.SNIs)
		if err != nil {
			return err
		}
	} else {
		securityProxyDeleteCurrentTLSCertIfSet()
	}

	return nil
}

func processSecretsConfigOptions(ao appOptions) error {
	return processSecretsConfigProxyOptions(ao.Proxy)
}
