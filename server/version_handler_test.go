// Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
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
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SamsungSLAV/weles"
)

var _ = Describe("VersionHandler", func() {
	var (
		testserver *httptest.Server
	)

	BeforeEach(func() {
		_, _, _, _, testserver = testServerSetup()
	})

	AfterEach(func() {
		testserver.Close()
	})

	Describe("obtaining information on API and server version", func() {
		getClientResp := func() (resp *http.Response) {
			client := testserver.Client()

			req, err := http.NewRequest(http.MethodGet, testserver.URL+"/api/v1/version", nil)
			Expect(err).ToNot(HaveOccurred())

			resp, err = client.Do(req)
			Expect(err).ToNot(HaveOccurred())

			return resp
		}

		Context("request to v0.1.0 API on v0.1.0 server", func() {
			It("should respond with proper body and 200 Status Code", func() {
				resp := getClientResp()
				defer resp.Body.Close()

				respBody, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				versionEncoded, err := json.Marshal(weles.Version{
					API:    "v1",
					State:  "devel",
					Server: "0.1.0",
				})
				Expect(err).ToNot(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(200))
				Expect(string(respBody)).To(MatchJSON(string(versionEncoded)))
			})
		})
	})
})
