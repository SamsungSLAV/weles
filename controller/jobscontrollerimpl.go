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

// File controller/jobscontrollerimpl.go contains JobsController interface
// implementation.

package controller

import (
	"sync"
	"time"

	"git.tizen.org/tools/weles"
)

// JobsControllerImpl structure stores Weles' Jobs data. It controls
// collision-free JobID creation. It stores state of Jobs' execution and saves
// data to DB. It implements JobsController interface.
type JobsControllerImpl struct {
	JobsController
	// mutex protects JobsControllerImpl structure.
	mutex *sync.RWMutex
	// lastID is the last used ID for the Job.
	lastID weles.JobID
	// jobs stores information about Weles' Jobs.
	jobs map[weles.JobID]*Job
}

// setupLastID initializes last used ID. Value is read from DB meta data.
func (js *JobsControllerImpl) setupLastID() {
	// TODO initialize with meta data read from DB.
	// Current implementation starts with seconds from Epoch to avoid problems with
	// artifacts database.

	js.lastID = weles.JobID(time.Now().Unix())
}

// NewJobsController creates and initializes a new instance of Jobs structure.
// It is the only valid way of creating it.
func NewJobsController() JobsController {
	js := &JobsControllerImpl{
		mutex: new(sync.RWMutex),
		jobs:  make(map[weles.JobID]*Job),
	}

	js.setupLastID()

	// TODO load Jobs data from DB.

	return js
}

// nextID generates and returns ID assigned to a new Job.
// It also updates lastID and saves the information in DB meta data.
func (js *JobsControllerImpl) nextID() weles.JobID {
	js.lastID++

	// TODO save new lastID in DB.

	return js.lastID
}

// NewJob creates and initializes a new Job.
func (js *JobsControllerImpl) NewJob(yaml []byte) (weles.JobID, error) {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	j := js.nextID()

	now := time.Now()
	js.jobs[j] = &Job{
		JobInfo: weles.JobInfo{
			JobID:   j,
			Created: now,
			Updated: now,
			Status:  weles.JOB_NEW,
		},
		yaml: yaml,
	}

	// TODO save struct in DB

	return j, nil
}

// GetYaml returns yaml Job description.
func (js *JobsControllerImpl) GetYaml(j weles.JobID) ([]byte, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	job, ok := js.jobs[j]
	if !ok {
		return nil, weles.ErrJobNotFound
	}

	return job.yaml, nil
}

// SetConfig stores config in Jobs structure.
func (js *JobsControllerImpl) SetConfig(j weles.JobID, conf weles.Config) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.ErrJobNotFound
	}

	job.config = conf
	job.Updated = time.Now()
	return nil
}

// isStatusChangeValid verifies if Job's status change is valid.
// It is a helper function for SetStatusAndInfo.
func isStatusChangeValid(oldStatus, newStatus weles.JobStatus) bool {
	if oldStatus == newStatus {
		return true
	}
	switch oldStatus {
	case weles.JOB_NEW:
		switch newStatus {
		case weles.JOB_PARSING, weles.JOB_CANCELED, weles.JOB_FAILED:
			return true
		}
	case weles.JOB_PARSING:
		switch newStatus {
		case weles.JOB_DOWNLOADING, weles.JOB_CANCELED, weles.JOB_FAILED:
			return true
		}
	case weles.JOB_DOWNLOADING:
		switch newStatus {
		case weles.JOB_WAITING, weles.JOB_CANCELED, weles.JOB_FAILED:
			return true
		}
	case weles.JOB_WAITING:
		switch newStatus {
		case weles.JOB_RUNNING, weles.JOB_CANCELED, weles.JOB_FAILED:
			return true
		}
	case weles.JOB_RUNNING:
		switch newStatus {
		case weles.JOB_COMPLETED, weles.JOB_CANCELED, weles.JOB_FAILED:
			return true
		}
	}
	return false
}

// SetStatusAndInfo changes status of the Job and updates info. Only valid
// changes are allowed.
// There are 3 terminal statuses: JOB_FAILED, JOB_CANCELED, JOB_COMPLETED;
// and 5 non-terminal statuses: JOB_NEW, JOB_PARSING, JOB_DOWNLOADING,
// JOB_WAITING, JOB_RUNNING.
// Only below changes of statuses are allowed:
// * JOB_NEW --> {JOB_PARSING, JOB_CANCELED, JOB_FAILED}
// * JOB_PARSING --> {JOB_DOWNLOADING, JOB_CANCELED, JOB_FAILED}
// * JOB_DOWNLOADING --> {JOB_WAITING, JOB_CANCELED, JOB_FAILED}
// * JOB_WAITING --> {JOB_RUNNING, JOB_CANCELED, JOB_FAILED}
// * JOB_RUNNING --> {JOB_COMPLETED, JOB_CANCELED, JOB_FAILED}
func (js *JobsControllerImpl) SetStatusAndInfo(j weles.JobID, newStatus weles.JobStatus, msg string) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.ErrJobNotFound
	}

	if !isStatusChangeValid(job.Status, newStatus) {
		return weles.ErrJobStatusChangeNotAllowed
	}

	job.Status = newStatus
	job.Info = msg
	job.Updated = time.Now()
	return nil
}

// GetConfig returns Job's config.
func (js *JobsControllerImpl) GetConfig(j weles.JobID) (weles.Config, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.Config{}, weles.ErrJobNotFound
	}

	return job.config, nil
}

// SetDryad saves access info for acquired Dryad.
func (js *JobsControllerImpl) SetDryad(j weles.JobID, d weles.Dryad) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.ErrJobNotFound
	}

	job.dryad = d
	return nil
}

// GetDryad returns Dryad acquired for the Job.
func (js *JobsControllerImpl) GetDryad(j weles.JobID) (weles.Dryad, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.Dryad{}, weles.ErrJobNotFound
	}

	return job.dryad, nil
}

// List returns information about Jobs. If argument is a nil/empty slice
// information about all Jobs is returned. Otherwise result is filtered
// and contains information about requested Jobs only.
func (js *JobsControllerImpl) List(filter []weles.JobID) ([]weles.JobInfo, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()
	ret := make([]weles.JobInfo, 0, len(js.jobs))
	if len(filter) == 0 {
		// Get all Jobs.
		for _, job := range js.jobs {
			ret = append(ret, job.JobInfo)
		}
	} else {
		// Get filtered Jobs.
		for _, j := range filter {
			job, ok := js.jobs[j]
			if ok {
				ret = append(ret, job.JobInfo)
			}
		}
	}
	return ret, nil
}
