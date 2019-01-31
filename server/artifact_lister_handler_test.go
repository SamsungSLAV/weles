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
		emptyFilter = weles.ArtifactFilter{}

		filledFilter = weles.ArtifactFilter{
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

		sorterEmpty = weles.ArtifactSorter{}

		sorterDescNoBy = weles.ArtifactSorter{
			Order: enums.SortOrderDescending,
		}

		sorterAscNoBy = weles.ArtifactSorter{
			Order: enums.SortOrderAscending,
		}

		sorterNoOrderID = weles.ArtifactSorter{
			By: enums.ArtifactSortByID,
		}

		sorterDescID = weles.ArtifactSorter{
			Order: enums.SortOrderDescending,
			By:    enums.ArtifactSortByID,
		}

		sorterAscID = weles.ArtifactSorter{
			Order: enums.SortOrderAscending,
			By:    enums.ArtifactSortByID,
		}

		// default value
		sorterDefault = sorterAscID

		// when pagination is on and no query params are set. When used, limit should also be set.
		emptyPaginatorOn = weles.ArtifactPaginator{Forward: true}
		// when pagination is off
		emptyPaginatorOff = weles.ArtifactPaginator{}

		artifactInfo420 = fixtures.CreateArtifactInfoSlice(420)
	)

	BeforeEach(func() {
		mockCtrl, _, mockArtifactManager, apiDefaults, testserver = testServerSetup()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		testserver.Close()
	})

	// helper functions
	newHTTPRequest := func(reqBody io.Reader, query, contentH, acceptH string) (req *http.Request) {
		req, err := http.NewRequest(
			http.MethodPost, testserver.URL+basePath+listArtifactsPath+query, reqBody)
		Expect(err).ToNot(HaveOccurred())
		req.Header.Set("Content-Type", contentH)
		req.Header.Set("Accept", acceptH)
		req.Close = true
		return req
	}

	newBody := func(f weles.ArtifactFilter, s weles.ArtifactSorter) *bytes.Reader {
		data := artifacts.ArtifactListerBody{
			Filter: &f,
			Sorter: &s,
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

	Describe("Pagination is turned off", func() {
		Describe("client sends correct request", func() {
			It("server should accept empty post request", func() {
				apiDefaults.PageLimit = 0
				listInfo := weles.ListInfo{
					TotalRecords:     uint64(len(artifactInfo420)),
					RemainingRecords: 0,
				}
				mockArtifactManager.EXPECT().ListArtifact(
					emptyFilter, sorterDefault, emptyPaginatorOff).Return(
					artifactInfo420, listInfo, nil)

				_, err := testserver.Client().Do(newHTTPRequest(nil, "", JSON, JSON))
				Expect(err).ToNot(HaveOccurred())
			})

			DescribeTable("server should ignore query params",
				func(query string) {
					apiDefaults.PageLimit = 0
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(artifactInfo420)),
						RemainingRecords: 0,
					}
					mockArtifactManager.EXPECT().ListArtifact(emptyFilter,
						sorterDefault, emptyPaginatorOff).Return(
						artifactInfo420, listInfo, nil)

					_, err := testserver.Client().Do(newHTTPRequest(nil, query, JSON, JSON))
					Expect(err).ToNot(HaveOccurred())
				},

				Entry("no query params set", ""),
				Entry("after query set", "?after=50"),
				Entry("after and limit query set", "?after=50&limit=10"),
				Entry("after and before query set", "?after=50&before=20"),
				Entry("after and before and limit query set", "?after=50&before=30&limit=13"),
				Entry("before query set", "?before=100"),
				Entry("before and limit query set", "?before=100&limit=12"),
			)

			DescribeTable("server should pass filter to ArtifactManager",
				func(filter weles.ArtifactFilter) {
					apiDefaults.PageLimit = 0
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(artifactInfo420)),
						RemainingRecords: 0,
					}
					mockArtifactManager.EXPECT().ListArtifact(
						filter, sorterDefault, emptyPaginatorOff).Return(
						artifactInfo420, listInfo, nil)

					_, err := testserver.Client().Do(
						newHTTPRequest(newBody(filter, sorterEmpty), "", JSON, JSON))
					Expect(err).ToNot(HaveOccurred())
				},
				Entry("when receiving empty filter", emptyFilter),
				Entry("when receiving filled filter", filledFilter),
			)

			DescribeTable("server should pass sorter to ArtifactManager, but set default values "+
				"on empty fields",
				func(sent, expected weles.ArtifactSorter) {
					apiDefaults.PageLimit = 0
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(artifactInfo420)),
						RemainingRecords: 0,
					}
					mockArtifactManager.EXPECT().ListArtifact(
						emptyFilter, expected, emptyPaginatorOff).Return(
						artifactInfo420, listInfo, nil)

					_, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sent), "", JSON, JSON))
					Expect(err).ToNot(HaveOccurred())
				},
				Entry("should set default order and by",
					sorterEmpty, sorterDefault),
				Entry("should pass ascending order and by ID",
					sorterAscID, sorterAscID),
				Entry("should pass descending order and by ID",
					sorterDescID, sorterDescID),
				Entry("should pass descending order and set default by",
					sorterDescNoBy, weles.ArtifactSorter{
						Order: sorterDescNoBy.Order,
						By:    sorterDefault.By,
					}),
				Entry("should pass ascending order and set default by",
					sorterAscNoBy, weles.ArtifactSorter{
						Order: sorterAscNoBy.Order,
						By:    sorterDefault.By,
					}),
				Entry("should pass by ID and set default order",
					sorterNoOrderID, weles.ArtifactSorter{
						Order: sorterDefault.Order,
						By:    sorterNoOrderID.By,
					}),
			)

			DescribeTable("should respond with all artifacts and correct headers",
				func(recordCount int) {
					apiDefaults.PageLimit = 0
					artifactInfo := fixtures.CreateArtifactInfoSlice(recordCount)
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(artifactInfo)),
						RemainingRecords: 0,
					}
					mockArtifactManager.EXPECT().ListArtifact(emptyFilter,
						sorterDefault, emptyPaginatorOff).Return(
						artifactInfo, listInfo, nil)

					resp, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
					defer resp.Body.Close()
					Expect(err).ToNot(HaveOccurred())

					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					checkArtifactInfoMarshalling(respBody, artifactInfo)

					By("Response must have 200 statuscode")
					Expect(resp.StatusCode).To(Equal(200))
					By("Next, Previous, RemainingRecords Headers should not be set")
					Expect(resp.Header.Get(NextPageHdr)).To(Equal(""))
					Expect(resp.Header.Get(PreviousPageHdr)).To(Equal(""))
					Expect(resp.Header.Get(ListRemainingHdr)).To(Equal(""))
					By("TotalRecords should be set to length of list")
					Expect(resp.Header.Get(ListTotalHdr)).To(Equal(strconv.Itoa(
						len(artifactInfo))))
					Expect(resp.Header.Get(ListBatchSizeHdr)).To(Equal(strconv.Itoa(
						len(artifactInfo))))
				},
				Entry("20 records avaliable", 20),
				Entry("420 records avaliable", 420),
			)
		})
	})

	Describe("ArtifactManager returns error", func() {
		DescribeTable("Server should return appropriate status code and error message",
			func(pageLimit int32, statusCode int, amerr error) {

				apiDefaults.PageLimit = pageLimit

				listInfo := weles.ListInfo{
					TotalRecords:     uint64(len(artifactInfo420)),
					RemainingRecords: 0,
				}
				var paginator weles.ArtifactPaginator
				if pageLimit == 0 {
					paginator = emptyPaginatorOff
				} else {
					paginator = emptyPaginatorOn
					paginator.Limit = pageLimit
				}
				mockArtifactManager.EXPECT().ListArtifact(
					emptyFilter, sorterDefault, paginator).Return(
					artifactInfo420, listInfo, amerr)

				resp, err := testserver.Client().Do(
					newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

				defer resp.Body.Close()
				respBody, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				checkErrorMarshalling(respBody, amerr)
				Expect(resp.StatusCode).To(Equal(statusCode))
				// should not set headers on error
				Expect(resp.Header.Get(NextPageHdr)).To(Equal(""))
				Expect(resp.Header.Get(PreviousPageHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListTotalHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListRemainingHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListBatchSizeHdr)).To(Equal(""))

			},
			Entry("pagination off, 404 status, Artifact not found error",
				int32(0), 404, weles.ErrArtifactNotFound),
			Entry("pagination on, 404 status, Artifact not found error",
				int32(100), 404, weles.ErrArtifactNotFound),
			Entry("pagination off, 500 status, Unexpected error",
				int32(0), 500, errors.New("This is unexpected error")),
			Entry("pagination on, 500 status, Unexpected error",
				int32(100), 500, errors.New("This is unexpected error")),
		)
	})

	Describe("Pagination turned on", func() {
		Describe("Correct request", func() {
			DescribeTable("server should set paginator object depending on query params",
				func(query string, expectedPaginator weles.ArtifactPaginator) {
					apiDefaults.PageLimit = 500
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(artifactInfo420)),
						RemainingRecords: 0,
					}
					mockArtifactManager.EXPECT().ListArtifact(
						emptyFilter, sorterDefault, expectedPaginator).Return(
						artifactInfo420, listInfo, nil)

					_, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), query, JSON, JSON))
					Expect(err).ToNot(HaveOccurred())

				},
				Entry("when no query params set", "",
					weles.ArtifactPaginator{Forward: true, Limit: 500}),
				Entry("when after param is set", "?after=30",
					weles.ArtifactPaginator{Forward: true, Limit: 500, ID: 30}),
				Entry("when after and limit params are set", "?after=30&limit=20",
					weles.ArtifactPaginator{Forward: true, Limit: 20, ID: 30}),
				Entry("when before param is set", "?before=30",
					weles.ArtifactPaginator{Forward: false, Limit: 500, ID: 30}),
				Entry("when before and limit params are set", "?before=30&limit=15",
					weles.ArtifactPaginator{Forward: false, Limit: 15, ID: 30}),
				Entry("when limit param is set", "?limit=30",
					weles.ArtifactPaginator{Forward: true, Limit: 30}),
			)

			DescribeTable("server should respond with 200/206 depending on "+
				"ListInfo.RemainingRecords returned by ArtifactManager",
				func(listInfo weles.ListInfo, statusCode int) {
					apiDefaults.PageLimit = 100
					paginator := emptyPaginatorOn
					paginator.Limit = apiDefaults.PageLimit

					mockArtifactManager.EXPECT().
						ListArtifact(emptyFilter, sorterDefault, paginator).
						Return(artifactInfo420, listInfo, nil)
					resp, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
					Expect(err).ToNot(HaveOccurred())

					Expect(resp.StatusCode).To(Equal(statusCode))
				},
				Entry("No more artifacts",
					weles.ListInfo{RemainingRecords: 0}, 200),
				Entry("More artifacts to show",
					weles.ListInfo{RemainingRecords: 320}, 206),
			)

			DescribeTable("on last page, server should NOT set headers: Weles-List-Remaining and",
				func(paginator weles.ArtifactPaginator, query string) {
					apiDefaults.PageLimit = 100
					paginator.Limit = apiDefaults.PageLimit

					mockArtifactManager.EXPECT().
						ListArtifact(emptyFilter, sorterDefault, paginator).
						Return(artifactInfo420,
							weles.ListInfo{TotalRecords: 420, RemainingRecords: 0},
							nil)

					resp, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), query, JSON, JSON))
					Expect(err).ToNot(HaveOccurred())

					Expect(resp.StatusCode).To(Equal(200))
					if paginator.Forward {
						Expect(resp.Header.Get(NextPageHdr)).To(Equal(""))
					} else {
						Expect(resp.Header.Get(PreviousPageHdr)).To(Equal(""))
					}
					Expect(resp.Header.Get(ListRemainingHdr)).To(Equal(""))

				},
				Entry("Weles-Next-Page when paginating forward (page n/n, n!=0)",
					weles.ArtifactPaginator{ID: 400, Forward: true}, "?after=400"),
				Entry("Weles-Previous-Page when paginating forward (page 0/n, n!=0)",
					weles.ArtifactPaginator{ID: 100, Forward: false}, "?before=100"),
			)

			// Checking Weles-{Next,Previous}-Page is omitted due to implementation
			// being generated swagger code.

			DescribeTable("Should set Weles-List-{Total,Remaining,Batch-Size} "+
				"based on listinfo and artifactlist",
				func(artifactInfo []weles.ArtifactInfo, listInfo weles.ListInfo) {
					apiDefaults.PageLimit = 100
					paginator := emptyPaginatorOn
					paginator.Limit = apiDefaults.PageLimit

					mockArtifactManager.EXPECT().
						ListArtifact(emptyFilter, sorterDefault, paginator).
						Return(artifactInfo, listInfo, nil)

					resp, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
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
			)
		})

		Describe("Error ", func() {
			DescribeTable("returned by server due to both before and after query params set",
				func(query string) {
					apiDefaults.PageLimit = 100

					resp, err := testserver.Client().Do(newHTTPRequest(nil, query, JSON, JSON))
					Expect(err).ToNot(HaveOccurred())

					respBody, err := ioutil.ReadAll(resp.Body)
					defer resp.Body.Close()
					Expect(err).ToNot(HaveOccurred())
					checkErrorMarshalling(respBody, weles.ErrBeforeAfterNotAllowed)

					Expect(resp.StatusCode).To(Equal(400))
					// headers should not be set on error
					Expect(resp.Header.Get(NextPageHdr)).To(Equal(""))
					Expect(resp.Header.Get(PreviousPageHdr)).To(Equal(""))
					Expect(resp.Header.Get(ListTotalHdr)).To(Equal(""))
					Expect(resp.Header.Get(ListRemainingHdr)).To(Equal(""))

				},
				Entry("empty body", "?before=10&after=20"),
				Entry("empty body, additional limit query set", "?before=10&after=20&limit=10"),
			)
		})
	})
})
