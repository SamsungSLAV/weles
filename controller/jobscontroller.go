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

// File controller/jobscontroller.go defines JobsController interface for
// managing Jobs inside Controller. It is defined to provide additional layer
// for strict managing Job structures only. This allows mocking up
// the interface in tests.

package controller

import (
	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/enums"
)

// Job contains all information about Job embedding public part - JobInfo.
type Job struct {
	weles.JobInfo
	config weles.Config
	yaml   []byte
	dryad  weles.Dryad
}

// JobsController defines methods for Jobs structures operations inside
// Controller.
type JobsController interface {
	// NewJob creates a new Job and returns newly assigned JobID.
	NewJob(yaml []byte) (weles.JobID, error)
	// GetYaml returns yaml Job description.
	GetYaml(weles.JobID) ([]byte, error)
	// SetConfig sets config in Job.
	SetConfig(weles.JobID, weles.Config) error
	// SetStatusAndInfo changes status and info of the Job.
	SetStatusAndInfo(weles.JobID, enums.JobStatus, string) error
	// GetConfig gets Job's config.
	GetConfig(weles.JobID) (weles.Config, error)
	// SetDryad saves access info for acquired Dryad.
	SetDryad(weles.JobID, weles.Dryad) error
	// GetDryad returns Dryad acquired for the Job.
	GetDryad(weles.JobID) (weles.Dryad, error)
	// List returns information on Jobs. It takes 3 arguments:
	// - JobFilter containing filters
	// - JobSorter containing sorting key and sorting direction
	// - JobPagination containing element after/before which a page should be returned. It also
	// contains information about direction of pagination and the size of the returned page which
	// must always be set.
	List(filter weles.JobFilter, sorter weles.JobSorter, paginator weles.JobPaginator) ([]weles.JobInfo, weles.ListInfo, error) // nolint:lll
}
