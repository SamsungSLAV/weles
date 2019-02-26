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
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/enums"
	"github.com/SamsungSLAV/weles/fixtures"
	"github.com/SamsungSLAV/weles/mock"
	"github.com/SamsungSLAV/weles/server"
	"github.com/SamsungSLAV/weles/server/operations/artifacts"
)

var _ = Describe("Listing artifacts with server initialized", func() {
	var (
		mockCtrl            *gomock.Controller
		apiDefaults         *server.APIDefaults
		mockArtifactManager *mock.MockArtifactManager
		testserver          *httptest.Server
	)

	// data to test against
	var (
		fEmpty = weles.ArtifactFilter{}

		fFilled = weles.ArtifactFilter{
			Alias: []weles.ArtifactAlias{"sdaaa", "aalliass"},
			JobID: []weles.JobID{1, 43, 3},
			Status: []enums.ArtifactStatus{
				enums.ArtifactStatusDOWNLOADING,
				enums.ArtifactStatusREADY,
			},
			Type: []enums.ArtifactType{
				enums.ArtifactTypeRESULT,
				enums.ArtifactTypeYAML,
			},
		}

		sEmpty = weles.ArtifactSorter{}

		sDescNoBy = weles.ArtifactSorter{
			Order: enums.SortOrderDescending,
		}

		sAscNoBy = weles.ArtifactSorter{
			Order: enums.SortOrderAscending,
		}

		sNoOrderID = weles.ArtifactSorter{
			By: enums.ArtifactSortByID,
		}

		sDescID = weles.ArtifactSorter{
			Order: enums.SortOrderDescending,
			By:    enums.ArtifactSortByID,
		}

		sAscID = weles.ArtifactSorter{
			Order: enums.SortOrderAscending,
			By:    enums.ArtifactSortByID,
		}

		// default value
		sorterDefault = sAscID

		// when pagination is on and no query params are set. When used, limit should also be set.
		pAFwDefaultLimit = weles.ArtifactPaginator{Forward: true, Limit: defaultPageLimit}
		pEmpty           = weles.Paginator{}

		pADefault = pAFwDefaultLimit

		artifactInfo420       = fixtures.CreateArtifactInfoSlice(420)
		artifactInfoFirstPage = artifactInfo420[:defaultPageLimit]
		listInfoFirstPage     = weles.ListInfo{
			TotalRecords:     uint64(len(artifactInfo420)),
			RemainingRecords: uint64(len(artifactInfo420) - defaultPageLimit),
		}
	)

	BeforeEach(func() {
		mockCtrl, _, mockArtifactManager, apiDefaults, testserver = testServerSetup()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		testserver.Close()
	})

	// helper functions
	newHTTPRequest := func(reqBody io.Reader, contentH, acceptH string) (req *http.Request) {
		req, err := http.NewRequest(
			http.MethodPost, testserver.URL+basePath+listArtifactsPath, reqBody)
		Expect(err).ToNot(HaveOccurred())
		req.Header.Set("Content-Type", contentH)
		req.Header.Set("Accept", acceptH)
		req.Close = true
		return req
	}

	newBody := func(f weles.ArtifactFilter, s weles.ArtifactSorter, p weles.Paginator,
	) *bytes.Reader {
		data := artifacts.ArtifactListerBody{
			Filter:    &f,
			Sorter:    &s,
			Paginator: &p,
		}
		marshalled, err := json.Marshal(data)
		Expect(err).ToNot(HaveOccurred())
		return bytes.NewReader(marshalled)
	}

	checkArtifactInfoMarshalling := func(respBody []byte, artifactInfo []weles.ArtifactInfo) {
		marshalled, err := json.Marshal(artifactInfo)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(respBody)).To(MatchJSON(string(marshalled)))
	}

	checkErrorMarshalling := func(respBody []byte, e error) {
		errMarshalled, err := json.Marshal(weles.ErrResponse{
			Message: e.Error(),
			Type:    "",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(string(respBody)).To(MatchJSON(string(errMarshalled)))
	}

	Describe("client sends correct request", func() {
		It("server should accept empty post request", func() {
			mockArtifactManager.EXPECT().ListArtifact(fEmpty, sorterDefault, pADefault).
				Return(artifactInfoFirstPage, listInfoFirstPage, nil)

			_, err := testserver.Client().Do(newHTTPRequest(nil, JSON, JSON))
			Expect(err).ToNot(HaveOccurred())
		})

		DescribeTable("server should pass filter to ArtifactManager",
			func(filter weles.ArtifactFilter) {
				mockArtifactManager.EXPECT().ListArtifact(filter, sorterDefault, pADefault).
					Return(artifactInfoFirstPage, listInfoFirstPage, nil)

				resp, err := testserver.Client().
					Do(newHTTPRequest(newBody(filter, sEmpty, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

				respBody, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				checkArtifactInfoMarshalling(respBody, artifactInfoFirstPage)
			},
			Entry("when receiving empty filter", fEmpty),
			Entry("when receiving filled filter", fFilled),
		)

		DescribeTable("server should pass sorter to ArtifactManager, but set default values "+
			"on empty fields",
			func(sent, expected weles.ArtifactSorter) {
				mockArtifactManager.EXPECT().ListArtifact(fEmpty, expected, pADefault).
					Return(artifactInfoFirstPage, listInfoFirstPage, nil)

				_, err := testserver.Client().
					Do(newHTTPRequest(newBody(fEmpty, sent, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())
			},
			Entry("should set default order and by",
				sEmpty, sorterDefault),
			Entry("should pass ascending order and by ID",
				sAscID, sAscID),
			Entry("should pass descending order and by ID",
				sDescID, sDescID),
			Entry("should pass descending order and set default by",
				sDescNoBy, weles.ArtifactSorter{
					Order: sDescNoBy.Order,
					By:    sorterDefault.By,
				}),
			Entry("should pass ascending order and set default by",
				sAscNoBy, weles.ArtifactSorter{
					Order: sAscNoBy.Order,
					By:    sorterDefault.By,
				}),
			Entry("should pass by ID and set default order",
				sNoOrderID, weles.ArtifactSorter{
					Order: sorterDefault.Order,
					By:    sNoOrderID.By,
				}),
		)

		DescribeTable("server should set paginator object to ArtifactManager, "+
			"but set default values on empty fields",
			func(globalLimit int32, pSent weles.Paginator, pExpected weles.ArtifactPaginator) {
				apiDefaults.PageLimit = globalLimit

				mockArtifactManager.EXPECT().ListArtifact(fEmpty, sorterDefault, pExpected).
					Return(artifactInfoFirstPage, listInfoFirstPage, nil)

				_, err := testserver.Client().Do(newHTTPRequest(
					newBody(fEmpty, sorterDefault, pSent), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

			},
			Entry("when empty Paginator is sent",
				int32(defaultPageLimit), weles.Paginator{}, pADefault),
			Entry("should set pagination direction to Forward when no direction is supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: defaultPageLimit},
				weles.ArtifactPaginator{Forward: true, Limit: defaultPageLimit}),
			Entry("should pass Forward direction when supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: defaultPageLimit, Direction: enums.DirectionForward},
				weles.ArtifactPaginator{Forward: true, Limit: defaultPageLimit}),
			Entry("should pass Backward direction when supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: defaultPageLimit, Direction: enums.DirectionBackward},
				weles.ArtifactPaginator{Forward: false, Limit: defaultPageLimit}),
			Entry("should pass Limit when supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: 69},
				weles.ArtifactPaginator{Limit: 69, Forward: true}),
			Entry("should set Limit to globalLimit when not supplied",
				int32(defaultPageLimit),
				weles.Paginator{},
				weles.ArtifactPaginator{Limit: defaultPageLimit, Forward: true}),
			Entry("should set Limit to globalLimit when not supplied",
				int32(69),
				weles.Paginator{},
				weles.ArtifactPaginator{Limit: 69, Forward: true}),
			Entry("should pass ID when supplied",
				int32(defaultPageLimit),
				weles.Paginator{ID: 50},
				weles.ArtifactPaginator{ID: 50, Limit: defaultPageLimit, Forward: true}),
		)

		DescribeTable("server should respond with 200/206 depending on "+
			"ListInfo.RemainingRecords returned by ArtifactManager",
			func(listInfo weles.ListInfo, statusCode int) {
				mockArtifactManager.EXPECT().
					ListArtifact(fEmpty, sorterDefault, pADefault).
					Return(artifactInfo420, listInfo, nil)
				resp, err := testserver.Client().Do(
					newHTTPRequest(newBody(fEmpty, sorterDefault, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(statusCode))
			},
			Entry("No more artifacts",
				weles.ListInfo{RemainingRecords: 0}, 200),
			Entry("More artifacts to show",
				weles.ListInfo{RemainingRecords: 320}, 206),
		)

		DescribeTable("Should set Weles-List-{Total,Remaining,Batch-Size} "+
			"based on listinfo and artifactlist",
			func(artifactInfo []weles.ArtifactInfo, listInfo weles.ListInfo) {
				apiDefaults.PageLimit = 100

				mockArtifactManager.EXPECT().
					ListArtifact(fEmpty, sorterDefault, pADefault).
					Return(artifactInfo, listInfo, nil)

				resp, err := testserver.Client().Do(
					newHTTPRequest(newBody(fEmpty, sorterDefault, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

				Expect(resp.Header.Get(ListTotalHdr)).
					To(Equal(strconv.FormatUint(listInfo.TotalRecords, 10)))
				Expect(resp.Header.Get(ListRemainingHdr)).
					To(Equal(strconv.FormatUint(listInfo.RemainingRecords, 10)))
				Expect(resp.Header.Get(ListBatchSizeHdr)).
					To(Equal(strconv.Itoa(len(artifactInfo))))
			},
			Entry("case 1",
				artifactInfo420, weles.ListInfo{TotalRecords: 420, RemainingRecords: 50}),
			Entry("case 2",
				artifactInfo420[:100], weles.ListInfo{TotalRecords: 100, RemainingRecords: 10}),
			Entry("case 3",
				artifactInfo420[:50], weles.ListInfo{TotalRecords: 100, RemainingRecords: 0}),
		)
	})

	Describe("ArtifactManager returns error", func() {
		DescribeTable("Server should return appropriate status code and error message",
			func(statusCode int, amErr error) {
				mockArtifactManager.EXPECT().ListArtifact(fEmpty, sorterDefault, pADefault).
					Return(artifactInfoFirstPage, listInfoFirstPage, amErr)

				resp, err := testserver.Client().
					Do(newHTTPRequest(newBody(fEmpty, sorterDefault, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

				respBody, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				checkErrorMarshalling(respBody, amErr)
				Expect(resp.StatusCode).To(Equal(statusCode))
				// should not set headers on error
				Expect(resp.Header.Get(NextPageHdr)).To(Equal(""))
				Expect(resp.Header.Get(PreviousPageHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListTotalHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListRemainingHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListBatchSizeHdr)).To(Equal(""))

			},
			Entry("404 status, Artifact not found error",
				404, weles.ErrArtifactNotFound),
			Entry("500 status, Unexpected error",
				500, errors.New("This is unexpected error")),
		)
	})
})
