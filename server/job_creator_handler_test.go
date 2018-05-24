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
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/mock"
	"git.tizen.org/tools/weles/server"
	"git.tizen.org/tools/weles/server/operations/jobs"
)

var _ = Describe("JobCreatorHandler", func() {

	var (
		mockCtrl       *gomock.Controller
		mockJobManager *mock.MockJobManager
		mockManagers   *server.Managers
		testserver     *httptest.Server
	)

	BeforeEach(func() {
		mockCtrl, mockJobManager, _, mockManagers, testserver = testServerSetup()
	})

	AfterEach(func() {
		testserver.Close()
		mockCtrl.Finish()
	})

	Describe("Creating a job", func() {
		requestBody := func(fileName string, fieldName string, acceptH string) (req *http.Request) {
			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			//create new form-data header with provided key-value pair
			fileWriter, err := bodyWriter.CreateFormFile(fieldName, fileName)
			Expect(err).ToNot(HaveOccurred())

			file, err := os.OpenFile(fileName, os.O_RDONLY, 0755)
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()
			_, err = io.Copy(fileWriter, file)
			Expect(err).ToNot(HaveOccurred())
			bodyWriter.Close()

			req, err = http.NewRequest(http.MethodPost, testserver.URL+"/api/v1/jobs/", bodyBuf)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
			if acceptH != OMIT {
				req.Header.Set("Accept", acceptH)
			}
			return req
		}
		mockInput := func(fileName string) (orgBody []byte) {
			file, err := os.OpenFile(fileName, os.O_RDONLY, 0755)
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()
			orgBody, err = ioutil.ReadAll(file)
			Expect(err).ToNot(HaveOccurred())
			return orgBody
		}
		Context("server receives correct POST request with accept header XML/JSON", func() {
			DescribeTable("should respond with 201 and JobID in body (XML/JSON)",
				func(accept string, expect string) {

					req := requestBody("test_sample.yml", "yamlfile", accept)
					orgBody := mockInput("test_sample.yml")
					client := testserver.Client()
					mockJobManager.EXPECT().CreateJob(orgBody).Return(weles.JobID(1234), nil)

					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					defer resp.Body.Close()
					resp_body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					Expect(resp_body).To(MatchJSON(expect))
					Expect(resp.StatusCode).To(Equal(201))

				},
				Entry("json", JSON, "1234\n"),
				Entry("default json", OMIT, "1234\n"),
			)
		})
		Context("server receives correct POST request but CreateJob returns error", func() {
			DescribeTable("should respond with 500 and error message in body (XML/JSON)",
				func(accept string, expect string) {

					req := requestBody("test_sample.yml", "yamlfile", accept)
					orgBody := mockInput("test_sample.yml")
					client := testserver.Client()
					mockJobManager.EXPECT().CreateJob(orgBody).Return(weles.JobID(0), errors.New("Unparsable"))

					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					defer resp.Body.Close()
					resp_body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(resp_body)).To(MatchJSON(expect))
					Expect(resp.StatusCode).To(Equal(500))
				},
				Entry("json", JSON, "{\"message\":\"Unparsable\"}\n"),
				Entry("default json", OMIT, "{\"message\":\"Unparsable\"}\n"),
			)
		})

		Context("handler receives nil instead of file", func() {
			It("should return unprocessable entity object", func() {
				req, err := http.NewRequest(http.MethodPost, testserver.URL+"/api/v1/jobs/", errReader(0))
				Expect(err).ToNot(HaveOccurred())
				params := jobs.JobCreatorParams{Yamlfile: errReader(0), HTTPRequest: req}

				ret := mockManagers.JobCreator(params)
				Expect(ret.(*jobs.JobCreatorUnprocessableEntity).Payload).To(Equal(&weles.ErrResponse{Message: "reader error"}))

			})
		})
	})
})

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("reader error")
}
func (errReader) Close() (err error) {
	return errors.New("close error")
}
