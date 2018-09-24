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

// File manager/dryad_job_manager.go provides implementation of DryadJobManager.

package manager

import (
	"sync"

	"github.com/SamsungSLAV/slav/logger"
	"github.com/SamsungSLAV/weles"
)

// DryadJobs implements DryadJobManager interface.
type DryadJobs struct {
	weles.DryadJobManager
	jobs           map[weles.JobID]*dryadJob
	jobsMutex      *sync.RWMutex
	artifactDBPath string
}

// NewDryadJobManager returns DryadJobManager interface of a new instance of DryadJobs.
func NewDryadJobManager(artifactDBPath string) weles.DryadJobManager {
	return &DryadJobs{
		jobs:           make(map[weles.JobID]*dryadJob),
		jobsMutex:      new(sync.RWMutex),
		artifactDBPath: artifactDBPath,
	}
}

// Create is part of DryadJobManager interface.
func (d *DryadJobs) Create(job weles.JobID, rusalka weles.Dryad, conf weles.Config,
	changes chan<- weles.DryadJobStatusChange) error {

	_, ok := d.jobs[job]
	if ok {
		logger.WithProperty("JobID", job).Errorf("Tried to create job that already exists.")
		return ErrDuplicated
	}
	d.jobsMutex.Lock()
	defer d.jobsMutex.Unlock()
	// FIXME(amistewicz): dryadJobs should not be stored indefinitely.
	d.jobs[job] = newDryadJob(job, rusalka, conf, changes, d.artifactDBPath)
	return nil
}

// Cancel is part of DryadJobManager interface.
func (d *DryadJobs) Cancel(job weles.JobID) error {
	d.jobsMutex.RLock()
	defer d.jobsMutex.RUnlock()
	dJob, ok := d.jobs[job]
	if !ok {
		logger.WithProperty("JobID", job).Errorf("Tried to cancel nonexistent job.")
		return ErrNotExist
	}
	dJob.cancel()
	return nil
}

// createStatusMatcher creates a matcher for DryadJobStatus.
// It is a helper function of List.
func createStatusMatcher(statuses []weles.DryadJobStatus) func(weles.DryadJobStatus) bool {
	if len(statuses) == 0 {
		return func(s weles.DryadJobStatus) bool {
			return true
		}
	}
	m := make(map[weles.DryadJobStatus]interface{})
	for _, status := range statuses {
		m[status] = nil
	}
	return func(s weles.DryadJobStatus) bool {
		_, ok := m[s]
		return ok
	}
}

// List is part of DryadJobManager interface.
func (d *DryadJobs) List(filter *weles.DryadJobFilter) ([]weles.DryadJobInfo, error) {
	d.jobsMutex.RLock()
	defer d.jobsMutex.RUnlock()

	// Trivial case - return all.
	if filter == nil {
		ret := make([]weles.DryadJobInfo, 0, len(d.jobs))
		for _, job := range d.jobs {
			info := job.GetJobInfo()
			ret = append(ret, info)
		}
		return ret, nil
	}

	ret := make([]weles.DryadJobInfo, 0)
	statusMatcher := createStatusMatcher(filter.Statuses)

	// References undefined - check only Status.
	if filter.References == nil {
		for _, job := range d.jobs {
			info := job.GetJobInfo()
			if statusMatcher(info.Status) {
				ret = append(ret, info)
			}
		}
		return ret, nil
	}

	// References defined - iterate only over requested keys.
	for _, id := range filter.References {
		job, ok := d.jobs[id]
		if !ok {
			continue
		}
		info := job.GetJobInfo()
		if statusMatcher(info.Status) {
			ret = append(ret, info)
		}
	}
	return ret, nil
}
