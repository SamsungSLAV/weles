/*
 *  Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

package dryad

import (
	"fmt"

	"github.com/SamsungSLAV/slav/logger"
)

// prefixPath is a parent directory DUT scripts.
const prefixPath = "/usr/local/bin/"

// deviceCommunicationProvider implements DeviceCommunicationProvider interface.
type deviceCommunicationProvider struct {
	DeviceCommunicationProvider
	credentials     Credentials
	sessionProvider SessionProvider
}

// NewDeviceCommunicationProvider returns new instance of DeviceCommunicationProvider.
func NewDeviceCommunicationProvider(session SessionProvider) DeviceCommunicationProvider {
	return &deviceCommunicationProvider{
		sessionProvider: session,
	}
}

// Boot function is a part of DeviceCommunicationProvider interface.
func (d *deviceCommunicationProvider) Boot() (err error) {
	_, _, err = d.sessionProvider.Exec(prefixPath + "dut_boot.sh")
	return
}

// Login is a part of DeviceCommunicationProvider interface.
func (d *deviceCommunicationProvider) Login(credentials Credentials) error {
	d.credentials = credentials
	_, _, err := d.sessionProvider.Exec(prefixPath+"dut_login.sh", d.credentials.Username,
		d.credentials.Password)
	if err != nil {
		logger.Error("Failed to login", err)
	}

	return err
}

// CopyFilesTo function is a part of DeviceCommunicationProvider interface.
// It uses tmpfs of MuxPi so caller must take into consideration size of all files
// that are to be copied.
func (d *deviceCommunicationProvider) CopyFilesTo(src []string, dest string) error {
	for _, path := range src {
		_, _, err := d.sessionProvider.Exec(prefixPath+"dut_copyto.sh", path, dest)
		if err != nil {
			logger.Errorf("Failed to copy %s to %s: %s", path, dest, err.Error())
			return fmt.Errorf("failed to copy %s to %s: %v", path, dest, err)
		}
	}
	return nil
}

// CopyFilesFrom function is a part of DeviceCommunicationProvider interface.
// It uses tmpfs of MuxPi so caller must take into consideration size of all files
// that are to be copied.
func (d *deviceCommunicationProvider) CopyFilesFrom(src []string, dest string) error {
	for _, path := range src {
		_, _, err := d.sessionProvider.Exec(prefixPath+"dut_copyfrom.sh", path, dest)
		if err != nil {
			logger.Errorf("Failed to copy %s to %s: %s", path, dest, err.Error())
			return fmt.Errorf("failed to copy %s to %s: %v", path, dest, err)
		}
	}
	return nil
}

// Exec function is a part of DeviceCommunicationProvider interface.
func (d *deviceCommunicationProvider) Exec(cmd ...string) (stdout, stderr []byte, err error) {
	return d.sessionProvider.Exec(append([]string{prefixPath + "dut_exec.sh"}, cmd...)...)
}

// Close function is a part of DeviceCommunicationProvider interface.
func (d *deviceCommunicationProvider) Close() error {
	return nil // Nothing to do for the time.
}
