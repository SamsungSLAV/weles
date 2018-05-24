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
	"net/http/httptest"
	"testing"

	"git.tizen.org/tools/weles/mock"
	"git.tizen.org/tools/weles/server"
	"git.tizen.org/tools/weles/server/operations"
	"github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	JSON = "application/json"
	OMIT = "omit"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

func testServerSetup() (mockCtrl *gomock.Controller, mockJobManager *mock.MockJobManager, mockArtifactManager *mock.MockArtifactManager, mockManagers *server.Managers, testserver *httptest.Server) {
	mockCtrl = gomock.NewController(GinkgoT())
	mockJobManager = mock.NewMockJobManager(mockCtrl)
	mockArtifactManager = mock.NewMockArtifactManager(mockCtrl)
	swaggerSpec, _ := loads.Analyzed(server.SwaggerJSON, "")
	api := operations.NewWelesAPI(swaggerSpec)
	srv := server.NewServer(api)
	mockManagers = server.NewManagers(mockJobManager, mockArtifactManager)
	srv.WelesConfigureAPI(mockManagers)
	testserver = httptest.NewServer(srv.GetHandler())
	return
}
