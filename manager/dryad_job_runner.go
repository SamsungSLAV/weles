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

// File manager/dryad_job_runner.go provides implementation of DryadJobRunner.

package manager

import (
	"context"

	"git.tizen.org/tools/weles/manager/dryad"
)

// dryadJobRunner implements DryadJobRunner interface.
type dryadJobRunner struct {
	DryadJobRunner
	ctx     context.Context
	rusalka dryad.SessionProvider
	device  dryad.DeviceCommunicationProvider
}

// newDryadJobRunner prepares a new instance of dryadJobRunner
// and returns DryadJobRunner interface to it.
func newDryadJobRunner(ctx context.Context, rusalka dryad.SessionProvider, device dryad.DeviceCommunicationProvider) DryadJobRunner {
	return &dryadJobRunner{
		ctx:     ctx,
		rusalka: rusalka,
		device:  device,
	}
}

// Deploy is part of DryadJobRunner interface.
func (d *dryadJobRunner) Deploy() error {
	// TODO(amistewicz): implement.
	return nil
}

// Boot is part of DryadJobRunner interface.
func (d *dryadJobRunner) Boot() error {
	// TODO(amistewicz): implement.
	return nil
}

// Test is part of DryadJobRunner interface.
func (d *dryadJobRunner) Test() error {
	// TODO(amistewicz): implement.
	return nil
}
