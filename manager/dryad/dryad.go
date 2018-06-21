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

// Package dryad provides Dryad Manager.
package dryad

// SessionProvider is used to execute steps
// from job definition that require communication with Dryad.
//
// It should automatically reconnect if connection has been lost unless Close is called.
type SessionProvider interface {
	// Exec runs a cmd on Dryad.
	// Execution time is limited by default timeout of the session.
	Exec(cmd ...string) (stdout, stderr []byte, err error)

	// DUT switches connections of SDcard and power supply to Device Under Test (DUT).
	//
	// Additional actions, notable PowerTick, may be required for successful device boot.
	DUT() error

	// TS switches connections of SDcard to Test Server (TS).
	// Power is cut from the device and it has no longer access to SDcard.
	TS() error

	// PowerTick switches voltage input on and off in order to cause a device reboot.
	// Moreover it may temporarily change state of dypers.
	//
	// Tick length and dyper actions are defined in device configuration.
	PowerTick() error

	// Close terminates session to Dryad.
	Close() error

	// SendFile sends file to Dryad.
	SendFile(src, dst string) error

	// ReceiveFile receives file from Dryad.
	ReceiveFile(src, dst string) error
}

// Credentials are used to login to device.
type Credentials struct {
	Username string
	Password string
}

// DeviceCommunicationProvider is used to execute steps
// from job definition that require communication with DUT.
type DeviceCommunicationProvider interface {
	// Login changes user which is used by remaining methods of this interface.
	//
	// In case of serial it is simple login to terminal.
	//
	// In case of SSH it is used to initialize the connection.
	//
	// In case of SDB it corresponds to `sdb root on` if username is "root"
	// or `sdb root off` otherwise.
	//
	// If shell has appropriate permissions, `su - username` may be used.
	Login(Credentials) error

	// CopyFilesTo transfers data from src to dest present on device.
	// All non-existing directories in dest will be created.
	//
	// It corresponds to command `sdb push src dest`.
	CopyFilesTo(src []string, dest string) error

	// CopyFilesFrom transfers data from src present on device to dest.
	// All non-existing directories in dest will be created.
	//
	// It corresponds to command `sdb pull src dest`.
	CopyFilesFrom(src []string, dest string) error

	// Exec runs a command on device until it exits or timeout occurs.
	// error occurs also when a cmd has non-zero return value.
	// command may be terminated if the stdout and stderr is too large, err will be set.
	//
	// Large outputs should be redirected to files.
	Exec(cmd ...string) (stdout, stderr []byte, err error)

	// Close terminates session to Device.
	Close() error
}
