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

package manager

import (
	"context"
	"fmt"
	"sync"

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/manager/dryad"
)

type dryadJob struct {
	info       weles.DryadJobInfo
	mutex      *sync.Mutex
	runner     DryadJobRunner
	notify     chan<- weles.DryadJobStatusChange
	cancel     context.CancelFunc
	failReason string
}

// newDryadJobWithCancel creates an instance of dryadJob without a goroutine.
// It is intended to be used by tests and newDryadJob only.
func newDryadJobWithCancel(job weles.JobID, changes chan<- weles.DryadJobStatusChange,
	runner DryadJobRunner, cancel context.CancelFunc) *dryadJob {

	dJob := &dryadJob{
		mutex:  new(sync.Mutex),
		runner: runner,
		info: weles.DryadJobInfo{
			Job: job,
		},
		notify: changes,
		cancel: cancel,
	}
	dJob.changeStatus(weles.DryadJobStatusNEW)
	return dJob
}

// newDryadJob creates an instance of dryadJob and starts a goroutine
// executing phases of given job implemented by provider of DryadJobRunner interface.
func newDryadJob(job weles.JobID, rusalka weles.Dryad, conf weles.Config,
	changes chan<- weles.DryadJobStatusChange) *dryadJob {

	// FIXME: It should use the proper path to the artifactory.
	session := dryad.NewSessionProvider(rusalka, "")
	device := dryad.NewDeviceCommunicationProvider(session)

	ctx, cancel := context.WithCancel(context.Background())
	runner := newDryadJobRunner(ctx, session, device, conf)

	dJob := newDryadJobWithCancel(job, changes, runner, cancel)

	go dJob.run(ctx)
	return dJob
}

// GetJobInfo returns DryadJobInfo of dryadJob and prevents race condition.
func (d *dryadJob) GetJobInfo() weles.DryadJobInfo {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.info
}

// changeState updates Status and sends DryadJobStatusChange to the notify channel.
func (d *dryadJob) changeStatus(state weles.DryadJobStatus) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.info.Status = state
	select {
	case d.notify <- weles.DryadJobStatusChange{Job: d.info.Job, Status: state}:
	default:
	}
}

func (d *dryadJob) executePhase(name weles.DryadJobStatus, f func() error) {
	d.changeStatus(name)
	err := f()
	if err != nil {
		panic(fmt.Errorf("%s phase failed: %s", name, err))
	}
}

// run executes stages of dryadJob in order.
func (d *dryadJob) run(_ context.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				d.failReason = err.Error()
			} else {
				d.failReason = fmt.Sprintf("run panicked: %v", r)
			}
			d.changeStatus(weles.DryadJobStatusFAIL)
			return
		}
		d.changeStatus(weles.DryadJobStatusOK)
	}()
	d.executePhase(weles.DryadJobStatusDEPLOY, d.runner.Deploy)
	d.executePhase(weles.DryadJobStatusBOOT, d.runner.Boot)
	d.executePhase(weles.DryadJobStatusTEST, d.runner.Test)
}
