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
package options

/*

 const userSemaphoreFile = ".secrets-config-user"
 const tlsSemaphoreFile = ".secrets-config-tls"
 const kongAdminTokenFile = "kong-admin-jwt"

 // Write a security-proxy file
 func securityProxyWriteFile(filename, contents string) (path string, err error) {
	 path = fmt.Sprintf("%s/secrets/security-proxy-setup/%s", SnapData, filename)
	 err = ioutil.WriteFile(path, []byte(contents), 0600)
	 if err == nil {
		 Debug(fmt.Sprintf("Wrote file '%s'", path))
	 } else {
		 err = fmt.Errorf("failed to write file %s - %v", path, err)
	 }
	 return
 }

 // Read a security-proxy file
 func securityProxyReadFile(filename string) (contents string, err error) {
	 path := fmt.Sprintf("%s/secrets/security-proxy-setup/%s", SnapData, filename)
	 bytes, err := ioutil.ReadFile(path)
	 if err == nil {
		 contents = string(bytes)
		 Debug(fmt.Sprintf("Read file '%s'", path))
	 } else {
		 err = fmt.Errorf("failed to read file %s - %v", path, err)
	 }
	 return
 }

 // Delete a security-proxy semaphore file
 func securityProxyRemoveSemaphore(filename string) (err error) {
	 path := fmt.Sprintf("%s/secrets/security-proxy-setup/%s", SnapData, filename)
	 err = os.Remove(path)
	 if err == nil {
		 Debug("Removed file '" + path + "'")
	 } else { var
	 return
 }

 // Execute the secrets-config tool with the given arguments
 func securityProxyExecSecretsConfig(args []string) error {
	 service := "security-proxy-setup"
	 cmdSecretsConfig := exec.Command("secrets-config", args...)
	 cmdSecretsConfig.Dir = fmt.Sprintf("%s/%s", SnapDataConf, service)
	 out, err := cmdSecretsConfig.Output()
	 Debug("Executed \"secrets-config " + fmt.Sprint(args) + "\" result=" + string(out))
	 return err
 }

 // Remove the semaphore, so that we can set the certificate again
 func securityProxyDeleteCurrentTLSCertIfSet() error {

	 return securityProxyRemoveSemaphore(tlsSemaphoreFile)
 }

 // Delete the current user - if one has been set up
 func securityProxyDeleteCurrentUserIfSet() error {
	 service := "security-proxy-setup"
	 // if no user has been set up, then ignore the request
	 username, err := securityProxyReadFile(userSemaphoreFile)
	 if err != nil {
		 Debug("proxy: No user has been set up")
		 return nil
	 }

	 args := []string{"proxy", "deluser", "--user", username}
	 cmdSecretsConfig := exec.Command("secrets-config", args...)
	 cmdSecretsConfig.Dir = fmt.Sprintf("%s/%s", SnapDataConf, service)
	 out, err := cmdSecretsConfig.Output()
	 if err != nil {
		 return err
	 }

	 securityProxyRemoveSemaphore(userSemaphoreFile)
	 Debug("Executed \"secrets-config " + fmt.Sprint(args) + "\" result=" + string(out))
	 Info("proxy: Removed current user")
	 return nil
 }

 // Set up the proxy with the specified user.
 func securityProxyAddUser(jwtUsername, jwtUserID, jwtAlgorithm, jwtPublicKey string) error {
	 currentUser, err := securityProxyReadFile(userSemaphoreFile)
	 if err == nil && currentUser != "" {
		 if currentUser == jwtUsername {
			 //	If a user has already been set - and it's the same user - then silently ignore the request
			 Debug("proxy: Ignoring request to set up same user again")
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
		 return fmt.Errorf("failed to create proxy user - %v", err)
	 }
	 _, err = securityProxyWriteFile(userSemaphoreFile, jwtUsername)
	 if err != nil {
		 return err
	 }
	 Info("proxy: Added new user")
	 return nil
 }

 // Set the TLS certificate. If a certificate has already been set then silently ignore the request
 func securityProxySetTLSCertificate(tlsCertificate, tlsPrivateKey, tlsSNI string) error {
	 _, err := securityProxyReadFile(tlsSemaphoreFile)
	 if err == nil {
		 Debug("The TLS certificate has already been set. To set it again, first set tls-certificate and tls-private-key to an empty string")
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
	 if tlsSNI != "" {
		 args = append(args, "--snis", tlsSNI)
	 }
	 err = securityProxyExecSecretsConfig(args)

	 if err != nil {
		 return fmt.Errorf("failed to set TLS certificate - %v", err)
	 }
	 _, err = securityProxyWriteFile(tlsSemaphoreFile, "TLS certificate set")
	 if err != nil {
		 return err
	 }
	 Info("New TLS Certificate and private key set")
	 return nil
 }
*/

/*

 const userSemaphoreFile = ".secrets-config-user"
 const tlsSemaphoreFile = ".secrets-config-tls"
 const kongAdminTokenFile = "kong-admin-jwt"

 // Write a security-proxy file
 func securityProxyWriteFile(filename, contents string) (path string, err error) {
	 path = fmt.Sprintf("%s/secrets/security-proxy-setup/%s", SnapData, filename)
	 err = ioutil.WriteFile(path, []byte(contents), 0600)
	 if err == nil {
		 Debug(fmt.Sprintf("Wrote file '%s'", path))
	 } else {
		 err = fmt.Errorf("failed to write file %s - %v", path, err)
	 }
	 return
 }

 // Read a security-proxy file
 func securityProxyReadFile(filename string) (contents string, err error) {
	 path := fmt.Sprintf("%s/secrets/security-proxy-setup/%s", SnapData, filename)
	 bytes, err := ioutil.ReadFile(path)
	 if err == nil {
		 contents = string(bytes)
		 Debug(fmt.Sprintf("Read file '%s'", path))
	 } else {
		 err = fmt.Errorf("failed to read file %s - %v", path, err)
	 }
	 return
 }

 // Delete a security-proxy semaphore file
 func securityProxyRemoveSemaphore(filename string) (err error) {
	 path := fmt.Sprintf("%s/secrets/security-proxy-setup/%s", SnapData, filename)
	 err = os.Remove(path)
	 if err == nil {
		 Debug("Removed file '" + path + "'")
	 } else {
		 Debug(fmt.Sprintf("Could not remove file '%s' : %v", path, err))
	 }
	 return
 }

 // Execute the secrets-config tool with the given arguments
 func securityProxyExecSecretsConfig(args []string) error {
	 service := "security-proxy-setup"
	 cmdSecretsConfig := exec.Command("secrets-config", args...)
	 cmdSecretsConfig.Dir = fmt.Sprintf("%s/%s", SnapDataConf, service)
	 out, err := cmdSecretsConfig.Output()
	 Debug("Executed \"secrets-config " + fmt.Sprint(args) + "\" result=" + string(out))
	 return err
 }

 // Remove the semaphore, so that we can set the certificate again
 func securityProxyDeleteCurrentTLSCertIfSet() error {

	 return securityProxyRemoveSemaphore(tlsSemaphoreFile)
 }

 // Delete the current user - if one has been set up
 func securityProxyDeleteCurrentUserIfSet() error {
	 service := "security-proxy-setup"
	 // if no user has been set up, then ignore the request
	 username, err := securityProxyReadFile(userSemaphoreFile)
	 if err != nil {
		 Debug("proxy: No user has been set up")
		 return nil
	 }

	 args := []string{"proxy", "deluser", "--user", username}
	 cmdSecretsConfig := exec.Command("secrets-config", args...)
	 cmdSecretsConfig.Dir = fmt.Sprintf("%s/%s", SnapDataConf, service)
	 out, err := cmdSecretsConfig.Output()
	 if err != nil {
		 return err
	 }

	 securityProxyRemoveSemaphore(userSemaphoreFile)
	 Debug("Executed \"secrets-config " + fmt.Sprint(args) + "\" result=" + string(out))
	 Info("proxy: Removed current user")
	 return nil
 }

 // Set up the proxy with the specified user.
 func securityProxyAddUser(jwtUsername, jwtUserID, jwtAlgorithm, jwtPublicKey string) error {
	 currentUser, err := securityProxyReadFile(userSemaphoreFile)
	 if err == nil && currentUser != "" {
		 if currentUser == jwtUsername {
			 //	If a user has already been set - and it's the same user - then silently ignore the request
			 Debug("proxy: Ignoring request to set up same user again")
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
		 return fmt.Errorf("failed to create proxy user - %v", err)
	 }
	 _, err = securityProxyWriteFile(userSemaphoreFile, jwtUsername)
	 if err != nil {
		 return err
	 }
	 Info("proxy: Added new user")
	 return nil
 }

 // Set the TLS certificate. If a certificate has already been set then silently ignore the request
 func securityProxySetTLSCertificate(tlsCertificate, tlsPrivateKey, tlsSNI string) error {
	 _, err := securityProxyReadFile(tlsSemaphoreFile)
	 if err == nil {
		 Debug("The TLS certificate has already been set. To set it again, first set tls-certificate and tls-private-key to an empty string")
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
	 if tlsSNI != "" {
		 args = append(args, "--snis", tlsSNI)
	 }
	 err = securityProxyExecSecretsConfig(args)

	 if err != nil {
		 return fmt.Errorf("failed to set TLS certificate - %v", err)
	 }
	 _, err = securityProxyWriteFile(tlsSemaphoreFile, "TLS certificate set")
	 if err != nil {
		 return err
	 }
	 Info("New TLS Certificate and private key set")
	 return nil
 }
*/
