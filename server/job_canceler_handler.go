// Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package server

import (
	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/server/operations/jobs"
	middleware "github.com/go-openapi/runtime/middleware"
)

// JobCanceller is a handler which passess JobID to JobManager to cancel a job.
func (m *Managers) JobCanceller(params jobs.JobCancelerParams) middleware.Responder {
	err := m.JM.CancelJob(weles.JobID(params.JobID))
	switch err {
	case nil:
		return jobs.NewJobCancelerNoContent()

	case weles.ErrJobNotFound:
		return jobs.NewJobCancelerNotFound().WithPayload(&weles.ErrResponse{Message: err.Error(), Type: ""})

	case weles.ErrJobStatusChangeNotAllowed:
		return jobs.NewJobCancelerForbidden().WithPayload(&weles.ErrResponse{Message: err.Error(), Type: ""})

	default:
		return jobs.NewJobCancelerInternalServerError().WithPayload(&weles.ErrResponse{Message: err.Error(), Type: ""})
	}

}
