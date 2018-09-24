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

// File controller/dryaderimpl.go implements Dryader interface. It delegates
// and controls Job execution by DryadJobManager.

package controller

import (
	"fmt"
	"sync"

	"github.com/SamsungSLAV/slav/logger"
	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/controller/notifier"
)

// DryaderImpl implements Dryader. It delegates and controls Job execution
// by DryadJobManager.
type DryaderImpl struct {
	// Notifier provides channel for communication with Controller.
	notifier.Notifier
	// jobs references module implementing Jobs management.
	jobs JobsController
	// djm manages DryadJobs.
	djm weles.DryadJobManager
	// info contains Jobs delegated to DryadJobManager and not completed yet
	// - active Jobs collection.
	info map[weles.JobID]bool
	// mutex protects access to info map.
	mutex *sync.Mutex
	// listener listens on notifications from DryadJobManager.
	listener chan weles.DryadJobStatusChange
	// finish is channel for stopping internal goroutine.
	finish chan int
	// looper waits for internal goroutine running loop to finish.
	looper sync.WaitGroup
}

// NewDryader creates a new DryaderImpl structure setting up references
// to used Weles modules.
func NewDryader(j JobsController, d weles.DryadJobManager) Dryader {
	ret := &DryaderImpl{
		Notifier: notifier.NewNotifier(),
		jobs:     j,
		djm:      d,
		info:     make(map[weles.JobID]bool),
		mutex:    new(sync.Mutex),
		listener: make(chan weles.DryadJobStatusChange),
		finish:   make(chan int),
	}
	ret.looper.Add(1)
	go ret.loop()
	return ret
}

// Finish internal goroutine.
func (h *DryaderImpl) Finish() {
	h.finish <- 1
	h.looper.Wait()
}

// add adds a new Job delegated to DryadJobManager to active Jobs collection.
func (h *DryaderImpl) add(j weles.JobID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.info[j] = true
}

// remove Job from active Jobs collection.
func (h *DryaderImpl) remove(j weles.JobID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.info, j)
}

// setStatus sets Jobs status to RUNNING and updates info.
func (h *DryaderImpl) setStatus(j weles.JobID, msg string) {
	err := h.jobs.SetStatusAndInfo(j, weles.JobStatusRUNNING, msg)
	if err != nil {
		logger.WithError(err).WithProperty("JobID", j).
			Error("Failed to change job state to RUNNING")
		h.remove(j)
		h.SendFail(j, fmt.Sprintf("Internal Weles error while changing Job status : %s",
			err.Error()))
	}
}

// loop monitors DryadJob's status.
func (h *DryaderImpl) loop() {
	defer h.looper.Done()
	for {
		select {
		case <-h.finish:
			return
		case recv := <-h.listener:
			change := weles.DryadJobInfo(recv)
			h.mutex.Lock()
			_, ok := h.info[change.Job]
			h.mutex.Unlock()
			if !ok {
				continue
			}

			switch change.Status {
			case weles.DryadJobStatusNEW:
				h.setStatus(change.Job, "Started")
			case weles.DryadJobStatusDEPLOY:
				h.setStatus(change.Job, "Deploying")
			case weles.DryadJobStatusBOOT:
				h.setStatus(change.Job, "Booting")
			case weles.DryadJobStatusTEST:
				h.setStatus(change.Job, "Testing")
			case weles.DryadJobStatusFAIL:
				h.remove(change.Job)
				h.SendFail(change.Job, "Failed to execute test on Dryad.")
			case weles.DryadJobStatusOK:
				h.remove(change.Job)
				h.SendOK(change.Job)
			}
		}
	}
}

// StartJob registers new Job to be executed in DryadJobManager.
func (h *DryaderImpl) StartJob(j weles.JobID) {
	d, err := h.jobs.GetDryad(j)
	if err != nil {
		logger.WithError(err).WithProperty("JobID", j).Error("Failed to get Dryad for job.")
		h.SendFail(j, fmt.Sprintf("Internal Weles error while getting Dryad for Job : %s",
			err.Error()))
		return
	}

	config, err := h.jobs.GetConfig(j)
	if err != nil {
		logger.WithError(err).WithProperty("JobID", j).Error("Failed to get config for job.")
		h.SendFail(j, fmt.Sprintf("Internal Weles error while getting Job config : %s",
			err.Error()))
		return
	}

	h.add(j)

	err = h.djm.Create(j, d, config, h.listener)
	if err != nil {
		logger.WithError(err).WithProperty("JobID", j).Error("Failed to run job on dryad.")
		h.remove(j)
		h.SendFail(j, fmt.Sprintf("Cannot delegate Job to Dryad : %s", err.Error()))
		return
	}
}

// CancelJob breaks Job execution in DryadJobManager.
func (h *DryaderImpl) CancelJob(j weles.JobID) {
	h.mutex.Lock()
	_, ok := h.info[j]
	h.mutex.Unlock()
	if !ok {
		return
	}

	h.remove(j)
	_ = h.djm.Cancel(j)
}
