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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/mock"
)

var _ = Describe("JobCancelerHandler", func() {

	var (
		mockCtrl       *gomock.Controller
		mockJobManager *mock.MockJobManager
		testserver     *httptest.Server
	)

	BeforeEach(func() {
		mockCtrl, mockJobManager, _, _, testserver = testServerSetup()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		testserver.Close()
	})

	Describe("cancelling a job", func() {
		getClientResp := func(accept string) (resp *http.Response) {
			client := testserver.Client()
			req, err := http.NewRequest(http.MethodPost, testserver.URL+"/api/v1/jobs/1234/cancel",
				nil)
			Expect(err).ToNot(HaveOccurred())
			if accept != OMIT {
				req.Header.Set("Accept", accept)
			}
			resp, err = client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			return resp
		}
		Context("correct request", func() {
			It("should respond with 204 Status Code", func() {
				mockJobManager.EXPECT().CancelJob(weles.JobID(1234))
				resp := getClientResp("")
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(204))
			})
		})
		Context("server should respond", func() {
			DescribeTable("with appropriate error",
				func(accept string, erro error, statuscode int) {

					mockJobManager.EXPECT().CancelJob(weles.JobID(1234)).Return(erro)
					resp := getClientResp(accept)
					defer resp.Body.Close()

					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					errorEncoded, err := json.Marshal(weles.ErrResponse{
						Message: erro.Error(),
						Type:    ""})
					Expect(err).ToNot(HaveOccurred())
					Expect(string(respBody)).To(MatchJSON(string(errorEncoded)))

					Expect(resp.StatusCode).To(Equal(statuscode))
				},
				Entry("job does not exist - 404",
					JSON, weles.ErrJobNotFound, 404),
				Entry("job does not exist - 404",
					OMIT, weles.ErrJobNotFound, 404),
				Entry("job already has final status - 403",
					JSON, weles.ErrJobStatusChangeNotAllowed, 403),
				Entry("job already has final status - 403",
					OMIT, weles.ErrJobStatusChangeNotAllowed, 403),
				Entry("unexpected error - 500",
					JSON, errors.New("Some other error"), 500),
				Entry("unexpected error - 500",
					OMIT, errors.New("Some other error"), 500),
			)
		})

	})
})
