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

// File controller/boruterimpl.go implements Boruter interface
// for communication with Boruta. Communication is used for acquiring,
// monitoring and releasing Dryads for Weles' Jobs.

package controller

import (
	"fmt"
	"sync"
	"time"

	"git.tizen.org/tools/boruta"
	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/controller/notifier"
)

// TODO ProlongAccess to Dryad in Boruta, before time expires.

// jobBorutaInfo contains information about status of acquiring Dryad from
// Boruta for running a single Job.
type jobBorutaInfo struct {
	// rid is the Boruta's request ID for the Job.
	rid boruta.ReqID
	// status is the current state of the request.
	status boruta.ReqState
	// timeout defines Dryad acquirement duration.
	timeout time.Time
}

// BoruterImpl is a Handler that is responsible for managing communication
// with Boruta, acquiring Dryads, prolonging access and releasing them.
type BoruterImpl struct {
	// Notifier provides channel for communication with Controller.
	notifier.Notifier
	// jobs references module implementing Jobs management.
	jobs JobsController
	// boruta is Boruta's client.
	boruta boruta.Requests

	// info contains information about status of acquiring Dryad from Boruta.
	info map[weles.JobID]*jobBorutaInfo
	// rid2Job maps Boruta's RequestID to Weles' JobID.
	rid2Job map[boruta.ReqID]weles.JobID
	// mutex protects access to info and rid2Job maps.
	mutex *sync.Mutex
	// borutaCheckPeriod defines how often Boruta is asked for requests' status.
	borutaCheckPeriod time.Duration
	// finish is channel for stopping internal goroutine.
	finish chan int
	// looper waits for internal goroutine running loop to finish.
	looper sync.WaitGroup
}

// NewBoruter creates a new BoruterImpl structure setting up references
// to used Weles and Boruta modules.
func NewBoruter(j JobsController, b boruta.Requests, period time.Duration) Boruter {
	ret := &BoruterImpl{
		Notifier:          notifier.NewNotifier(),
		jobs:              j,
		boruta:            b,
		info:              make(map[weles.JobID]*jobBorutaInfo),
		rid2Job:           make(map[boruta.ReqID]weles.JobID),
		mutex:             new(sync.Mutex),
		borutaCheckPeriod: period,
		finish:            make(chan int),
	}
	ret.looper.Add(1)
	go ret.loop()
	return ret
}

// Finish internal goroutine.
func (h *BoruterImpl) Finish() {
	h.finish <- 1
	h.looper.Wait()
}

// add registers new Boruta's request ID for the Job to be monitored.
func (h *BoruterImpl) add(j weles.JobID, r boruta.ReqID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.info[j] = &jobBorutaInfo{
		rid: r,
	}
	h.rid2Job[r] = j
}

// remove Boruta's request ID for the Job from monitored requests.
func (h *BoruterImpl) remove(j weles.JobID, r boruta.ReqID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.rid2Job, r)
	delete(h.info, j)
}

// pop gets and removes Boruta's request ID for the Job and a Job from monitored set.
// It returns request ID related to the removed Job
func (h *BoruterImpl) pop(j weles.JobID) (r boruta.ReqID, err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	rinfo, ok := h.info[j]
	if !ok {
		return r, weles.ErrJobNotFound
	}
	r = rinfo.rid
	delete(h.rid2Job, r)
	delete(h.info, j)
	return
}

// setProlongTime stores time until Dryad is acquired from Boruta.
func (h *BoruterImpl) setProlongTime(j weles.JobID, rinfo boruta.ReqInfo) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.info[j].timeout = rinfo.Job.Timeout
}

// updateStatus analyzes single Boruta's request info and verifies if it is
// related to any of Weles' Jobs. If so, method returns new status of request
// and ID of related Job. Otherwise zero-value status is returned.
func (h *BoruterImpl) updateStatus(rinfo boruta.ReqInfo) (newState boruta.ReqState, j weles.JobID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var ok bool
	j, ok = h.rid2Job[rinfo.ID]
	if !ok {
		return
	}
	info := h.info[j]
	if info.status == rinfo.State {
		return
	}
	info.status = rinfo.State
	newState = rinfo.State

	return
}

// acquire gets Dryad from Boruta and sets information about it in
// JobsController. It stores acquired Dryad's expiration time
// and notifies Controller about getting Dryad for the Job.
func (h *BoruterImpl) acquire(j weles.JobID, rinfo boruta.ReqInfo) {
	ai, err := h.boruta.AcquireWorker(rinfo.ID)
	if err != nil {
		h.remove(j, rinfo.ID)
		h.SendFail(j, fmt.Sprintf("Cannot acquire worker from Boruta : %s", err.Error()))
		return
	}
	// TODO acquire username from Boruta.
	err = h.jobs.SetDryad(j, weles.Dryad{Addr: ai.Addr, Key: ai.Key, Username: "boruta-user"})
	if err != nil {
		h.remove(j, rinfo.ID)
		h.SendFail(j, fmt.Sprintf("Internal Weles error while setting Dryad : %s", err.Error()))
		return
	}
	h.setProlongTime(j, rinfo)
	h.SendOK(j)
}

// loop monitors Boruta's requests.
func (h *BoruterImpl) loop() {
	defer h.looper.Done()
	for {
		select {
		case <-h.finish:
			return
		case <-time.After(h.borutaCheckPeriod):
		}

		// TODO use filter with slice of ReqIDs when implemented in Boruta.
		requests, err := h.boruta.ListRequests(nil)
		if err != nil {
			// TODO log error
			continue
		}

		for _, rinfo := range requests {
			status, j := h.updateStatus(rinfo)

			switch status {
			case boruta.INPROGRESS:
				h.acquire(j, rinfo)
			case boruta.CANCEL:
				h.remove(j, rinfo.ID)
			case boruta.DONE:
				h.remove(j, rinfo.ID)
			case boruta.TIMEOUT:
				h.remove(j, rinfo.ID)
				h.SendFail(j, "Timeout in Boruta.")
			case boruta.INVALID:
				h.remove(j, rinfo.ID)
				h.SendFail(j, "No suitable device in Boruta to run test.")
			case boruta.FAILED:
				h.remove(j, rinfo.ID)
				h.SendFail(j, "Boruta failed during request execution.")
			}
		}
	}
}

// getCaps prepares Capabilities for registering new request in Boruta.
func (h *BoruterImpl) getCaps(config weles.Config) boruta.Capabilities {
	if config.DeviceType == "" {
		return boruta.Capabilities{}
	}

	return boruta.Capabilities{
		"DeviceType": config.DeviceType,
	}
}

// getPriority prepares Priority for registering new request in Boruta.
func (h *BoruterImpl) getPriority(config weles.Config) boruta.Priority {
	switch config.Priority {
	case weles.LOW:
		return 11
	case weles.MEDIUM:
		return 7
	case weles.HIGH:
		return 3
	default:
		return 7
	}
}

// getOwner prepares Owner for registering new request in Boruta.
func (h *BoruterImpl) getOwner() boruta.UserInfo {
	return boruta.UserInfo{}
}

// getValidAfter prepares ValidAfter time for registering new request in Boruta.
func (h *BoruterImpl) getValidAfter(config weles.Config) time.Time {
	return time.Now()
}

// getDeadline prepares Deadline time for registering new request in Boruta.
func (h *BoruterImpl) getDeadline(config weles.Config) time.Time {
	const defaultDelay = 24 * time.Hour
	if config.Timeouts.JobTimeout == weles.ValidPeriod(0) {
		return time.Now().Add(defaultDelay)
	}

	return time.Now().Add(time.Duration(config.Timeouts.JobTimeout))
}

// Request registers new request in Boruta and adds it to monitored requests.
func (h *BoruterImpl) Request(j weles.JobID) {
	err := h.jobs.SetStatusAndInfo(j, weles.JobStatusWAITING, "")
	if err != nil {
		h.SendFail(j, fmt.Sprintf("Internal Weles error while changing Job status : %s", err.Error()))
		return
	}

	config, err := h.jobs.GetConfig(j)
	if err != nil {
		h.SendFail(j, fmt.Sprintf("Internal Weles error while getting Job config : %s", err.Error()))
		return
	}

	caps := h.getCaps(config)
	priority := h.getPriority(config)
	owner := h.getOwner()
	validAfter := h.getValidAfter(config)
	deadline := h.getDeadline(config)

	r, err := h.boruta.NewRequest(caps, priority, owner, validAfter, deadline)
	if err != nil {
		h.SendFail(j, fmt.Sprintf("Failed to create request in Boruta : %s", err.Error()))
		return
	}

	h.add(j, r)
}

// Release returns Dryad to Boruta's pool and closes Boruta's request.
func (h *BoruterImpl) Release(j weles.JobID) {
	r, err := h.pop(j)
	if err != nil {
		return
	}
	h.boruta.CloseRequest(r)
}
