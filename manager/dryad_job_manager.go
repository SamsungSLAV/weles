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

// File manager/dryad_job_manager.go provides implementation of DryadJobManager.

package manager

import (
	"sync"

	. "git.tizen.org/tools/weles"
)

// DryadJobs implements DryadJobManager interface.
type DryadJobs struct {
	DryadJobManager
	jobs      map[JobID]*dryadJob
	jobsMutex *sync.RWMutex
}

// NewDryadJobManager returns DryadJobManager interface of a new instance of DryadJobs.
func NewDryadJobManager() DryadJobManager {
	return &DryadJobs{
		jobs:      make(map[JobID]*dryadJob),
		jobsMutex: new(sync.RWMutex),
	}
}

// Create is part of DryadJobManager interface.
func (d *DryadJobs) Create(job JobID, rusalka Dryad, changes chan<- DryadJobStatusChange) error {
	_, ok := d.jobs[job]
	if ok {
		return ErrDuplicated
	}
	d.jobsMutex.Lock()
	defer d.jobsMutex.Unlock()
	// FIXME(amistewicz): dryadJobs should not be stored indefinitely.
	d.jobs[job] = newDryadJob(job, rusalka, changes)
	return nil
}

// Cancel is part of DryadJobManager interface.
func (d *DryadJobs) Cancel(job JobID) error {
	d.jobsMutex.RLock()
	defer d.jobsMutex.RUnlock()
	dJob, ok := d.jobs[job]
	if !ok {
		return ErrNotExist
	}
	dJob.cancel()
	return nil
}

// List is part of DryadJobManager interface.
func (d *DryadJobs) List(filter *DryadJobFilter) ([]DryadJobInfo, error) {
	if filter == nil {
		d.jobsMutex.RLock()
		defer d.jobsMutex.RUnlock()
		ret := make([]DryadJobInfo, 0, len(d.jobs))
		for _, job := range d.jobs {
			info := job.GetJobInfo()
			ret = append(ret, info)
		}
		return ret, nil
	}
	// TODO(amistewicz): implement.
	panic("not implemented")
}
