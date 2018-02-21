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

// File jobmanager.go provides JobManager interface.

package weles

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
