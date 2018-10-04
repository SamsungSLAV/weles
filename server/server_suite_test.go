// Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package server_test

import (
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tideland/golib/audit"

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/mock"
	"git.tizen.org/tools/weles/server"
	"git.tizen.org/tools/weles/server/operations"
)

const (
	JSON = "application/json"
	OMIT = "omit"

	dateLayout         = "Mon Jan 2 15:04:05 -0700 MST 2006"
	someDate           = "Tue Jan 2 15:04:05 +0100 CET 1900"
	durationIncrement1 = "25h"
	durationIncrement2 = "+100h"

	basePath          = "/api/v1"
	listArtifactsPath = "/artifacts/list"
	listJobsPath      = "/jobs/list"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

func testServerSetup() (mockCtrl *gomock.Controller, mockJobManager *mock.MockJobManager,
	mockArtifactManager *mock.MockArtifactManager, apiDefaults *server.APIDefaults,
	testserver *httptest.Server) {

	mockCtrl = gomock.NewController(GinkgoT())
	mockJobManager = mock.NewMockJobManager(mockCtrl)
	mockArtifactManager = mock.NewMockArtifactManager(mockCtrl)
	swaggerSpec, _ := loads.Analyzed(server.SwaggerJSON, "")
	api := operations.NewWelesAPI(swaggerSpec)
	srv := server.NewServer(api)
	apiDefaults = &server.APIDefaults{
		Managers: server.NewManagers(mockJobManager, mockArtifactManager),
	}
	srv.WelesConfigureAPI(apiDefaults)
	testserver = httptest.NewServer(srv.GetHandler())
	return
}

// createJobInfoSlice is a function to create random data for tests of JobLister
func createJobInfoSlice(sliceLenght int) (ret []weles.JobInfo) {
	// checking for errors omitted due to fixed input.
	dateTimeIter, _ := time.Parse(dateLayout, someDate)
	durationIncrement, _ := time.ParseDuration(durationIncrement1)
	durationIncrement2, _ := time.ParseDuration(durationIncrement2)
	jobInfo := make([]weles.JobInfo, sliceLenght)
	gen := audit.NewGenerator(rand.New(rand.NewSource(time.Now().UTC().UnixNano())))
	for i := range jobInfo {
		tmp := weles.JobInfo{}
		createdTime := gen.Time(time.Local, dateTimeIter, durationIncrement)
		tmp.Created = strfmt.DateTime(createdTime)
		tmp.Updated = strfmt.DateTime(gen.Time(time.Local, createdTime, durationIncrement2))
		tmp.Info = gen.Sentence()
		tmp.Name = gen.Word()
		tmp.Status = weles.JobStatus(gen.OneStringOf(string(weles.JobStatusNEW),
			string(weles.JobStatusPARSING), string(weles.JobStatusDOWNLOADING),
			string(weles.JobStatusWAITING), string(weles.JobStatusRUNNING),
			string(weles.JobStatusCOMPLETED), string(weles.JobStatusFAILED),
			string(weles.JobStatusCANCELED)))
		tmp.JobID = weles.JobID(i + 1)
		dateTimeIter = dateTimeIter.Add(durationIncrement)
		jobInfo[i] = tmp
	}
	return jobInfo
}
