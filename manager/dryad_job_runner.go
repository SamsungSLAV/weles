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

// File manager/dryad_job_runner.go provides implementation of DryadJobRunner.

package manager

import (
	"context"

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/manager/dryad"
)

// dryadJobRunner implements DryadJobRunner interface.
type dryadJobRunner struct {
	DryadJobRunner
	ctx     context.Context
	rusalka dryad.SessionProvider
	device  dryad.DeviceCommunicationProvider
	conf    weles.Config
}

// newDryadJobRunner prepares a new instance of dryadJobRunner
// and returns DryadJobRunner interface to it.
func newDryadJobRunner(ctx context.Context, rusalka dryad.SessionProvider,
	device dryad.DeviceCommunicationProvider, conf weles.Config) DryadJobRunner {
	return &dryadJobRunner{
		ctx:     ctx,
		rusalka: rusalka,
		device:  device,
		conf:    conf,
	}
}

// Deploy is part of DryadJobRunner interface.
func (d *dryadJobRunner) Deploy() (err error) {
	err = d.rusalka.TS()
	if err != nil {
		return
	}

	// Generate partition mapping for FOTA and store it on Dryad.
	urls := make([]string, 0, len(d.conf.Action.Deploy.Images))
	for _, image := range d.conf.Action.Deploy.Images {
		if p := image.Path; p != "" {
			urls = append(urls, p)
		}
	}
	partLayout := make([]fotaMap, 0, len(d.conf.Action.Deploy.PartitionLayout))
	for _, layout := range d.conf.Action.Deploy.PartitionLayout {
		if name, part := layout.ImageName, layout.ID; name != "" && part != 0 {
			partLayout = append(partLayout, fotaMap{name, part})
		}

	}
	mapping := newMapping(partLayout)
	_, _, err = d.rusalka.Exec("echo", "'"+string(mapping)+"'", ">", fotaFilePath)
	if err != nil {
		return
	}

	// Run FOTA.
	_, _, err = d.rusalka.Exec(newFotaCmd(fotaSDCardPath, fotaFilePath, urls).GetCmd()...)
	return err
}

// Boot is part of DryadJobRunner interface.
func (d *dryadJobRunner) Boot() (err error) {
	// Start DUT.
	err = d.device.Boot()
	if err != nil {
		return
	}

	// Login to the device only if credentials were specified.
	if username, password := d.conf.Action.Boot.Login, d.conf.Action.Boot.Password; username !=
		"" && password != "" {
		return d.device.Login(dryad.Credentials{Username: username, Password: password})
	}
	return nil
}

// Test is part of DryadJobRunner interface.
func (d *dryadJobRunner) Test() error {
	for _, testcase := range d.conf.Action.Test.TestCases {
		for _, testaction := range testcase.TestActions {
			switch action := testaction.(type) {
			case weles.Push:
				d.device.CopyFilesTo([]string{action.Path}, action.Dest)
			case weles.Run:
				// Exec joins arguments in a single string.
				// Split and then Join are avoided.
				d.device.Exec(action.Name)
			case weles.Pull:
				d.device.CopyFilesFrom([]string{action.Src}, action.Path)
			default:
				panic("unknown test action type")
			}
		}
	}
	return nil
}
