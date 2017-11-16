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

// File dryadjobmanager.go defines DryadJobManager interface and structures related to it.

package weles

// DryadJobStatus is a representation of current state of DryadJob.
type DryadJobStatus string

const (
	// DJ_NEW - initial status of DryadJob after call to Create.
	DJ_NEW DryadJobStatus = "CREATED"
	// DJ_DEPLOY - DryadJob is executing deploy section of job definition.
	DJ_DEPLOY DryadJobStatus = "DEPLOYING"
	// DJ_BOOT - DryadJob is executing boot section of job definition.
	DJ_BOOT DryadJobStatus = "BOOTING"
	// DJ_TEST - DryadJob is executing test section of job definition.
	DJ_TEST DryadJobStatus = "EXECUTING TESTS"
	// DJ_FAIL - an irrecoverable error has been encountered
	// and execution had to be stopped early.
	DJ_FAIL DryadJobStatus = "ERROR OCCURRED"
	// DJ_OK - DryadJob has finished execution successfully.
	DJ_OK DryadJobStatus = "DONE"
)

// DryadJobInfo contains information about DryadJob.
type DryadJobInfo struct {
	Job    JobID
	Status DryadJobStatus
}

// DryadJobStatusChange is information passed on the channel to the caller of Create.
type DryadJobStatusChange DryadJobInfo

// Dryad contains information about device allocated for Job
// and credentials required to use it.
type Dryad struct{}

// DryadJobFilter is used by List to access only jobs of interest.
//
// Job is matching DryadJobFilter if References contain value of
// its Job field and Statuses - Status.
type DryadJobFilter struct {
	References []JobID
	Statuses   []DryadJobStatus
}

// DryadJobManager organizes running Jobs on allocated Dryad.
type DryadJobManager interface {
	// Create starts execution of Job definition on allocated Dryad.
	//
	// Slow read from a channel may miss some events.
	Create(JobID, Dryad, chan<- DryadJobStatusChange) error

	// Cancel stops DryadJob associated with Job.
	//
	// It has no effect if Cancel has been called before
	// or job has already terminated.
	Cancel(JobID) error

	// List returns information about DryadJobs matching DryadJobFilter
	// or all if it is not specified.
	List(*DryadJobFilter) ([]DryadJobInfo, error)
}
