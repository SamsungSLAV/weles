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

// File jobmanager.go provides JobManager interface with Job related
// structures.

package weles

import "time"

// JobStatus specifies state of the Job.
type JobStatus string

const (
	// JOB_NEW - The new Job has been created.
	JOB_NEW JobStatus = "NEW"
	// JOB_PARSING - Provided yaml file is being parsed and interpreted.
	JOB_PARSING JobStatus = "PARSING"
	// JOB_DOWNLOADING - Images and/or files required for the test are being
	// downloaded.
	JOB_DOWNLOADING JobStatus = "DOWNLOADING"
	// JOB_WAITING - Job is waiting for Boruta worker.
	JOB_WAITING JobStatus = "WAITING"
	// JOB_RUNNING - Job is being executed.
	JOB_RUNNING JobStatus = "RUNNING"
	// JOB_COMPLETED - Job is completed.
	// This is a terminal state.
	JOB_COMPLETED JobStatus = "COMPLETED"
	// JOB_FAILED - Job execution has failed.
	// This is a terminal state.
	JOB_FAILED JobStatus = "FAILED"
	// JOB_CANCELED - Job has been canceled with API call.
	// This is a terminal state.
	JOB_CANCELED JobStatus = "CANCELED"
)

// JobInfo contains Job information available for public API.
type JobInfo struct {
	// JobID is a unique Job identifier.
	JobID JobID
	// Name is the Job name acquired from yaml file during Job creation.
	Name string
	// Created is the Job creation time in UTC.
	Created time.Time
	// Updated is the time of latest Jobs' status modification.
	Updated time.Time
	// Status specifies current state of the Job.
	Status JobStatus
	// Info provides additional information about current state,
	// e.g. cause of failure.
	Info string
}

// JobManager interface defines API for actions that can be called on Weles' Jobs
// by external modules. These methods are intended to be used by HTTP server.
type JobManager interface {
	// CreateJob creates a new Job in Weles using recipe passed in YAML format.
	// It returns ID of created Job or error.
	CreateJob(yaml []byte) (JobID, error)
	// CancelJob stops execution of Job identified by JobID.
	CancelJob(JobID) error
	// ListJobs returns information on Jobs. If argument is a nil/empty slice
	// information about all Jobs is returned. Otherwise result is filtered
	// and contains information about requested Jobs only.
	ListJobs([]JobID) ([]JobInfo, error)
}
