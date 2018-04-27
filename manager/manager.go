/*
 *  Copyright (c) 2017 Samsung Electronics Co., Ltd All Rights Reserved
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

// Package manager provides Dryad Job Manager.
package manager

// DryadJobRunner executes DryadJob on allocated Dryad.
// SessionProvider is used for actions on Dryad, DeviceCommunicationProvider - device.
type DryadJobRunner interface {
	// Deploy prepares device for Boot.
	//
	// It usually formats SDcard and copies image to it.
	Deploy() error

	// Boot starts up a device and prepares environment for Test.
	//
	// It usually attempts to log in to console.
	Boot() error

	// Test runs tests on a device.
	//
	// Test usually consists of following actions:
	// Push - deploy additional content,
	// Execute - run requested commands,
	// Collect - gather results.
	Test() error
}
