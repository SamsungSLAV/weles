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
	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/server/operations/jobs"
	middleware "github.com/go-openapi/runtime/middleware"

	"io/ioutil"
)

// JobCreator is a handler which passes yaml file with job description to jobmanager.
func (m *Managers) JobCreator(params jobs.JobCreatorParams) middleware.Responder {
	byteContainer, err := ioutil.ReadAll(params.Yamlfile)
	if err != nil {
		return jobs.NewJobCreatorUnprocessableEntity().WithPayload(
			&weles.ErrResponse{Message: err.Error(), Type: ""})
	}

	jobID, err := m.JM.CreateJob(byteContainer)
	if err != nil {
		return jobs.NewJobCreatorInternalServerError().WithPayload(
			&weles.ErrResponse{Message: err.Error(), Type: ""})
	}

	return jobs.NewJobCreatorCreated().WithPayload(jobID)
}
