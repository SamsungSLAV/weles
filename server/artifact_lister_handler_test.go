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
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/fixtures"
	"git.tizen.org/tools/weles/mock"
	"git.tizen.org/tools/weles/server"
	"git.tizen.org/tools/weles/server/operations/artifacts"
)

var _ = Describe("ArtifactListerHandler", func() {
	var (
		mockCtrl            *gomock.Controller
		apiDefaults         *server.APIDefaults
		mockArtifactManager *mock.MockArtifactManager
		testserver          *httptest.Server
	)

	BeforeEach(func() {
		mockCtrl, _, mockArtifactManager, apiDefaults, testserver = testServerSetup()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		testserver.Close()
	})
	Describe("Listing artifacts", func() {
		createRequest := func(
			reqBody io.Reader,
			path string,
			query string,
			contentH string,
			acceptH string) (req *http.Request) {
			if path == "" {
				path = "/api/v1/artifacts/list"
			}
			req, err := http.NewRequest(http.MethodPost, testserver.URL+path+query, reqBody)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", contentH)
			req.Header.Set("Accept", acceptH)
			return req
		}

		filterSorterReqBody := func(filter weles.ArtifactFilter, sorter weles.ArtifactSorter,
			contentH string) (rb *bytes.Reader) {

			artifactFilterSort := artifacts.ArtifactListerBody{
				Filter: &filter,
				Sorter: &sorter}

			artifactFilterSortMarshalled, err := json.Marshal(artifactFilterSort)
			Expect(err).ToNot(HaveOccurred())

			return bytes.NewReader(artifactFilterSortMarshalled)

		}

		checkReceivedArtifactInfo := func(respBody []byte, artifactInfo []weles.ArtifactInfo,
			acceptH string) {

			artifactInfoMarshalled, err := json.Marshal(artifactInfo)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(respBody)).To(MatchJSON(string(artifactInfoMarshalled)))

		}

		checkReceivedArtifactErr := func(respBody []byte, e error, acceptH string) {
			errMarshalled, err := json.Marshal(weles.ErrResponse{
				Message: e.Error(),
				Type:    ""})
			Expect(err).ToNot(HaveOccurred())
			Expect(string(respBody)).To(MatchJSON(string(errMarshalled)))
		}

		//few structs to test against
		emptyFilterA := weles.ArtifactFilter{}

		filledFilterA1 := weles.ArtifactFilter{
			Alias: []weles.ArtifactAlias{"sdaaa", "aalliass"},
			JobID: []weles.JobID{1, 43, 3},
			Status: []weles.ArtifactStatus{
				weles.ArtifactStatusDOWNLOADING,
				weles.ArtifactStatusREADY},
			Type: []weles.ArtifactType{
				weles.ArtifactTypeRESULT,
				weles.ArtifactTypeYAML}}

		filledFilterA2 := weles.ArtifactFilter{
			Alias: []weles.ArtifactAlias{"aalliass"},
			JobID: []weles.JobID{1, 43, 3, 9, 2, 10404},
			Status: []weles.ArtifactStatus{
				weles.ArtifactStatusFAILED},
			Type: []weles.ArtifactType{
				weles.ArtifactTypeRESULT,
				weles.ArtifactTypeYAML}}

		//empty sorter marshalls to default sorter.
		defaultSorterA := weles.ArtifactSorter{
			SortBy:    weles.ArtifactSortByID,
			SortOrder: weles.SortOrderAscending}

		filledSorterA2 := weles.ArtifactSorter{
			SortBy: weles.ArtifactSortByID}

		emptySorterA := weles.ArtifactSorter{}

		emptyPaginatorA := weles.ArtifactPagination{}
		emptyPaginatorFw := weles.ArtifactPagination{Forward: true}
		artifactInfo420 := fixtures.CreateArtifactInfoSlice(420)

		type queryPaginator struct {
			query     string
			paginator weles.ArtifactPagination
		}

		queryPaginatorOK := []queryPaginator{
			{
				query:     "",
				paginator: emptyPaginatorFw},
			{
				query: "?before=10",
				paginator: weles.ArtifactPagination{
					ID:      int64(10),
					Forward: false}},
			{
				query: "?before=30&limit=50",
				paginator: weles.ArtifactPagination{
					ID:      int64(30),
					Forward: false,
					Limit:   50}},
			{
				query: "?after=40",
				paginator: weles.ArtifactPagination{
					ID:      int64(40),
					Forward: true}},
			{
				query: "?after=70",
				paginator: weles.ArtifactPagination{
					ID:      int64(70),
					Forward: true,
					Limit:   200}},
			{
				query: "?limit=50",
				paginator: weles.ArtifactPagination{
					ID:      int64(0),
					Forward: true,
					Limit:   50}}}

		Context("a: Server receives correct request and has pagination turned off", func() {
			// ginkgo does not like It clauses with the same name. To avoid conflicts with
			// tests of listing jobs, each it clause not referring to artifacts explicitly
			// is prefixed with a:
			for _, curr := range queryPaginatorOK { //expected behaviour- handler should ignore
				// different queries when pagination is turned off globally.
				DescribeTable("should respond with all avaliable artifacts, ignoring query params",
					func(
						filter weles.ArtifactFilter,
						sorter weles.ArtifactSorter,
						query string) {
						apiDefaults.PageLimit = 0

						listInfo := weles.ListInfo{
							TotalRecords:     uint64(len(artifactInfo420)),
							RemainingRecords: 0}

						if sorter.SortOrder == "" {
							sorter.SortOrder = weles.SortOrderAscending
						}
						if sorter.SortBy == "" {
							sorter.SortBy = weles.ArtifactSortByID
						}
						mockArtifactManager.EXPECT().ListArtifact(
							filter, sorter, emptyPaginatorA).Return(
							artifactInfo420, listInfo, nil)

						reqBody := filterSorterReqBody(filter, sorter, JSON)

						client := testserver.Client()
						req := createRequest(reqBody, "", query, JSON, JSON)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())

						defer resp.Body.Close()
						respBody, err := ioutil.ReadAll(resp.Body)
						Expect(err).ToNot(HaveOccurred())

						checkReceivedArtifactInfo(respBody, artifactInfo420, JSON)

						Expect(resp.StatusCode).To(Equal(200))
						Expect(resp.Header.Get("Next")).To(Equal(""))
						Expect(resp.Header.Get("Previous")).To(Equal(""))
						Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))
						Expect(resp.Header.Get("TotalRecords")).To(Equal(strconv.Itoa(
							len(artifactInfo420))))

					},

					Entry("a: given empty filter and sorter",
						emptyFilterA, emptySorterA, curr.query),

					Entry("a: given filled filter and sorter",
						filledFilterA2, filledSorterA2, curr.query),

					Entry("a: given filled filter and empty sorter",
						filledFilterA1, emptySorterA, curr.query),

					Entry("a: given empty filter and filled sorter",
						emptyFilterA, defaultSorterA, curr.query),
				)
			}
		})

		Context("a: Server receives correct request and has pagination turned on", func() {
			artifactInfoAll := fixtures.CreateArtifactInfoSlice(100)
			globalLimit := []int32{70, 100, 111}
			for _, currgl := range globalLimit {
				for _, curr := range queryPaginatorOK {
					// expected behaviour- handler passes query parameters to AM.
					// this should impact data returned by JM. It is not reflected in the
					// below mock of JM as it is out of scope of server unit tests.
					DescribeTable("a: should respond with all artifacts",
						func(artifactInfo []weles.ArtifactInfo,
							filter weles.ArtifactFilter,
							sorter weles.ArtifactSorter,
							paginator weles.ArtifactPagination,
							query string, gl int32) {

							apiDefaults.PageLimit = gl

							if !strings.Contains(query, "limit") {
								paginator.Limit = apiDefaults.PageLimit
							}

							listInfo := weles.ListInfo{
								TotalRecords:     uint64(len(artifactInfo)),
								RemainingRecords: 0}

							mockArtifactManager.EXPECT().ListArtifact(
								filter, defaultSorterA, paginator).Return(
								artifactInfo, listInfo, nil)

							reqBody := filterSorterReqBody(filter, sorter, JSON)
							req := createRequest(reqBody, "", query, JSON, JSON)

							client := testserver.Client()
							resp, err := client.Do(req)
							Expect(err).ToNot(HaveOccurred())

							defer resp.Body.Close()
							respBody, err := ioutil.ReadAll(resp.Body)
							Expect(err).ToNot(HaveOccurred())

							checkReceivedArtifactInfo(respBody, artifactInfo, JSON)

							Expect(resp.StatusCode).To(Equal(200))
							// Next and Previous headers are ignored here as they are tested
							// in other context.
							Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))
							Expect(resp.Header.Get("TotalRecords")).To(
								Equal(strconv.Itoa(len(artifactInfo))))
						},
						Entry("given empty request when AM has less jobs than page size",
							artifactInfoAll[:40], emptyFilterA, emptySorterA, curr.paginator,
							curr.query, currgl),

						Entry("given filled filter, when AM returns less jobs than "+
							"Default Page size",
							artifactInfoAll[10:67], filledFilterA2, emptySorterA, curr.paginator,
							curr.query, currgl),

						Entry("given filled filter, when AM returns same amount of filtered jobs"+
							" as Default Page Size",
							artifactInfoAll, filledFilterA1, emptySorterA, curr.paginator,
							curr.query, currgl),
					)
				}
			}
		})

		Context("a: Pagination on", func() {
			artifactInfoAll := fixtures.CreateArtifactInfoSlice(400)
			DescribeTable("a: paginating forward",
				func(artifactInfo []weles.ArtifactInfo,
					startingPageNo int,
					filter weles.ArtifactFilter,
					sorter weles.ArtifactSorter) {

					apiDefaults.PageLimit = 100

					//prepare data for first call
					artifactInfoStartingPage := artifactInfo[(startingPageNo-1)*
						int(apiDefaults.PageLimit) : startingPageNo*
						int(apiDefaults.PageLimit)] //first page of data

					startingPageQuery := ""

					paginator := weles.ArtifactPagination{
						Limit:   apiDefaults.PageLimit,
						Forward: true,
					}
					if startingPageNo != 1 {
						paginator.Forward = true
						paginator.ID = artifactInfoStartingPage[0].ID
						startingPageQuery = "?after=" +
							strconv.FormatInt(artifactInfo[(startingPageNo-1)*int(
								apiDefaults.PageLimit)].ID, 10)
					}

					listInfo := weles.ListInfo{
						TotalRecords: uint64(len(artifactInfo)),
						RemainingRecords: uint64(len(artifactInfo) - startingPageNo*len(
							artifactInfoStartingPage)),
					}

					first := mockArtifactManager.EXPECT().ListArtifact(
						filter, defaultSorterA, paginator).Return(artifactInfoStartingPage, listInfo, nil)

					reqBody := filterSorterReqBody(filter, sorter, JSON)
					req := createRequest(reqBody, "", startingPageQuery, JSON, JSON)
					req.Close = true
					client := testserver.Client()
					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedArtifactInfo(respBody, artifactInfoStartingPage, JSON)
					Expect(resp.StatusCode).To(Equal(206))

					Expect(resp.Header.Get("Next")).To(Equal("/api/v1/artifacts/list" + "?after=" +
						strconv.FormatInt(artifactInfoStartingPage[apiDefaults.PageLimit-1].ID,
							10)))
					prevCheck := ""
					if startingPageNo != 1 {
						prevCheck = "/api/v1/artifacts/list" + "?before=" +
							strconv.FormatInt(artifactInfoStartingPage[0].ID, 10)
					}
					Expect(resp.Header.Get("Previous")).To(Equal(prevCheck))
					Expect(resp.Header.Get("RemainingRecords")).To(
						Equal(strconv.Itoa(len(artifactInfo) - (startingPageNo * len(
							artifactInfoStartingPage)))))
					Expect(resp.Header.Get("TotalRecords")).To(
						Equal(strconv.Itoa(len(artifactInfo))))

					nextPage := resp.Header.Get("Next")
					resp.Body.Close()
					testserver.CloseClientConnections()

					//prepare data for second call based on previous
					var artifactInfo2 []weles.ArtifactInfo
					var secondReturnCode int
					if (len(artifactInfo) - startingPageNo*
						int(apiDefaults.PageLimit)) <= int(apiDefaults.PageLimit) {

						artifactInfo2 = artifactInfo[startingPageNo*int(apiDefaults.PageLimit):]
						//next page is not full
						secondReturnCode = 200
					} else {
						artifactInfo2 = artifactInfo[startingPageNo*
							int(apiDefaults.PageLimit) : (startingPageNo+1)*
							int(apiDefaults.PageLimit)] //last page is full
						secondReturnCode = 206
					}

					paginator2 := weles.ArtifactPagination{
						Limit:   apiDefaults.PageLimit,
						Forward: true,
						ID: artifactInfoStartingPage[int(apiDefaults.PageLimit)-
							1].ID}
					listInfo2 := weles.ListInfo{TotalRecords: listInfo.TotalRecords}

					if tmp := len(artifactInfo) - (startingPageNo+1)*int(
						apiDefaults.PageLimit); tmp < 0 {
						listInfo2.RemainingRecords = 0
					} else {
						listInfo2.RemainingRecords = uint64(tmp)
					}
					//filter and sorter should stay the same.
					mockArtifactManager.EXPECT().ListArtifact(filter, defaultSorterA, paginator2).Return(
						artifactInfo2, listInfo2, nil).After(first)

					client2 := testserver.Client()
					reqBody = filterSorterReqBody(filter, sorter, JSON)
					req2 := createRequest(reqBody, nextPage, "", JSON, JSON)
					req2.Close = true
					resp2, err := client2.Do(req2)
					Expect(err).ToNot(HaveOccurred())

					defer resp2.Body.Close()
					respBody2, err := ioutil.ReadAll(resp2.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedArtifactInfo(respBody2, artifactInfo2, JSON)

					Expect(resp2.StatusCode).To(Equal(secondReturnCode))

					if secondReturnCode == 200 {
						Expect(resp2.Header.Get("Next")).To(Equal(""))
						prevCheck = "/api/v1/artifacts/list" + "?before=" +
							strconv.FormatInt(artifactInfo2[0].ID, 10)
						Expect(resp2.Header.Get("Previous")).To(Equal(prevCheck))
					} else {
						prevCheck = "/api/v1/artifacts/list" + "?before=" +
							strconv.FormatInt(artifactInfo2[0].ID, 10)
						nextCheck := "/api/v1/artifacts/list" + "?after=" +
							strconv.FormatInt(artifactInfo2[int(apiDefaults.PageLimit)-1].ID, 10)
						Expect(resp2.Header.Get("Next")).To(Equal(nextCheck))
						Expect(resp2.Header.Get("Previous")).To(Equal(prevCheck))
					}
					if tmp := strconv.Itoa(len(artifactInfo) -
						startingPageNo*len(artifactInfoStartingPage) -
						len(artifactInfo2)); tmp != "0" {
						Expect(resp2.Header.Get("RemainingRecords")).To(Equal(tmp))
					} else {

						Expect(resp2.Header.Get("RemainingRecords")).To(Equal(""))
					}
					Expect(resp2.Header.Get("TotalRecords")).To(
						Equal(strconv.Itoa(len(artifactInfo))))

				},
				Entry("a: 1->2/2", // from 1 to 2 out of 2 (pages)
					artifactInfoAll[:170], 1, emptyFilterA, emptySorterA),
				Entry("a: 1->2/3",
					artifactInfoAll[:270], 1, emptyFilterA, emptySorterA),
				Entry("a: 2->3/3",
					artifactInfoAll[:300], 2, emptyFilterA, emptySorterA),
				Entry("a: 2->3/4",
					artifactInfoAll[:350], 2, emptyFilterA, emptySorterA),
			)

			DescribeTable("a: paginating backward",
				func(artifactInfo []weles.ArtifactInfo,
					startingPageNo int, pages int,
					filter weles.ArtifactFilter,
					sorter weles.ArtifactSorter) {

					apiDefaults.PageLimit = 100
					paginator := weles.ArtifactPagination{Limit: apiDefaults.PageLimit}
					//prepare data for first call
					var artifactInfoStartingPage []weles.ArtifactInfo
					var startingPageQuery string
					listInfo := weles.ListInfo{}
					if startingPageNo == pages {
						artifactInfoStartingPage = artifactInfo[(startingPageNo-1)*
							int(apiDefaults.PageLimit) : len(artifactInfo)-1]
						paginator.Forward = true
						paginator.ID = artifactInfoStartingPage[len(artifactInfoStartingPage)-1].ID
						startingPageQuery = "?after=" + strconv.FormatInt(paginator.ID, 10)
						listInfo = weles.ListInfo{
							TotalRecords:     uint64(len(artifactInfo)),
							RemainingRecords: 0}

					} else {
						artifactInfoStartingPage = artifactInfo[(startingPageNo)*
							int(apiDefaults.PageLimit) : (startingPageNo+1)*
							int(apiDefaults.PageLimit)] //first page of data
						paginator.Forward = false
						paginator.ID = artifactInfo[(startingPageNo*
							int(apiDefaults.PageLimit))-1].ID
						startingPageQuery = "?before=" + strconv.FormatInt(paginator.ID, 10)
						listInfo = weles.ListInfo{
							TotalRecords: uint64(len(artifactInfo)),
							RemainingRecords: uint64(len(artifactInfo) - (int(
								apiDefaults.PageLimit) * startingPageNo))}

					}

					first := mockArtifactManager.EXPECT().ListArtifact(filter, defaultSorterA,
						paginator).Return(artifactInfoStartingPage, listInfo, nil)

					reqBody := filterSorterReqBody(filter, sorter, JSON)
					req := createRequest(reqBody, "", startingPageQuery, JSON, JSON)
					req.Close = true
					client := testserver.Client()
					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedArtifactInfo(respBody, artifactInfoStartingPage, JSON)

					if startingPageNo == pages {
						Expect(resp.StatusCode).To(Equal(200))
						Expect(resp.Header.Get("Previous")).To(
							Equal("/api/v1/artifacts/list?before=" +
								strconv.FormatInt(artifactInfoStartingPage[0].ID, 10)))
						Expect(resp.Header.Get("Next")).To(Equal(""))
						Expect(resp.Header.Get("TotalRecords")).To(
							Equal(strconv.Itoa(len(artifactInfo))))
						Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))

					} else {
						Expect(resp.StatusCode).To(Equal(206))
						Expect(resp.Header.Get("Previous")).To(
							Equal("/api/v1/artifacts/list?before=" +
								strconv.FormatInt(artifactInfoStartingPage[0].ID, 10)))
						Expect(resp.Header.Get("Next")).To(
							Equal("/api/v1/artifacts/list?after=" +
								strconv.FormatInt(artifactInfoStartingPage[len(
									artifactInfoStartingPage)-1].ID, 10)))
						Expect(resp.Header.Get("TotalRecords")).To(
							Equal(strconv.Itoa(len(artifactInfo))))
						Expect(resp.Header.Get("RemainingRecords")).To(
							Equal(strconv.FormatUint(listInfo.RemainingRecords, 10)))
					}

					prevPage := resp.Header.Get("Previous")

					resp.Body.Close()
					testserver.CloseClientConnections()

					//prepare data for second call based on previous

					var artifactInfo2 []weles.ArtifactInfo
					paginator2 := weles.ArtifactPagination{
						Limit:   apiDefaults.PageLimit,
						Forward: false,
						ID:      artifactInfoStartingPage[0].ID,
					}

					listInfo2 := weles.ListInfo{TotalRecords: listInfo.TotalRecords}
					if startingPageNo == pages {
						artifactInfo2 = artifactInfo[(startingPageNo-2)*
							int(apiDefaults.PageLimit) : (startingPageNo-1)*
							int(apiDefaults.PageLimit)]

						if startingPageNo-1 == 1 {
							listInfo2.RemainingRecords = 0
						} else {
							listInfo2.RemainingRecords = uint64((pages - (startingPageNo - 1)) *
								int(apiDefaults.PageLimit))
						}
					} else {
						artifactInfo2 = artifactInfo[(startingPageNo-1)*
							int(apiDefaults.PageLimit) : startingPageNo*
							int(apiDefaults.PageLimit)]

						listInfo2.RemainingRecords = uint64(apiDefaults.PageLimit)
					}

					mockArtifactManager.EXPECT().ListArtifact(filter, defaultSorterA, paginator2).Return(
						artifactInfo2, listInfo2, nil).After(first)
					client2 := testserver.Client()
					reqBody = filterSorterReqBody(filter, sorter, JSON)
					req2 := createRequest(reqBody, prevPage, "", JSON, JSON)
					req2.Close = true
					resp2, err := client2.Do(req2)
					Expect(err).ToNot(HaveOccurred())

					defer resp2.Body.Close()
					respBody2, err := ioutil.ReadAll(resp2.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedArtifactInfo(respBody2, artifactInfo2, JSON)
					if startingPageNo == pages {

						if startingPageNo-1 == 1 {
							Expect(resp2.StatusCode).To(Equal(200))
							Expect(resp2.Header.Get("RemainingRecords")).To(Equal(""))
							Expect(resp2.Header.Get("Previous")).To(Equal(""))

						} else {
							Expect(resp2.StatusCode).To(Equal(206))
							Expect(resp2.Header.Get("RemainingRecords")).To(Equal("100"))
							Expect(resp2.Header.Get("Previous")).To(
								Equal("/api/v1/artifacts/list?before=" +
									strconv.FormatInt(artifactInfo2[0].ID, 10)))
						}
					} else {
						Expect(resp2.StatusCode).To(Equal(206))
						Expect(resp2.Header.Get("RemainingRecords")).To(Equal("100"))
					}

					Expect(resp2.Header.Get("Next")).To(
						Equal("/api/v1/artifacts/list?after=" +
							strconv.FormatInt(artifactInfo2[len(artifactInfo2)-1].ID, 10)))
					Expect(resp2.Header.Get("TotalRecords")).To(
						Equal(strconv.Itoa(len(artifactInfo))))
				},
				Entry("a: 2->1/2",
					artifactInfoAll[:170], 2, 2, emptyFilterA, emptySorterA),
				Entry("a: 2->1/3",
					artifactInfoAll[:270], 2, 3, emptyFilterA, emptySorterA),
				Entry("a: 3->2/4",
					artifactInfoAll[:350], 3, 4, emptyFilterA, emptySorterA),
				Entry("a: 3->2/3",
					artifactInfoAll[:300], 3, 3, emptyFilterA, emptySorterA),
			)
		})

		Context("a: There is an error", func() {
			DescribeTable("Server should respond with error from ArtifactManager",
				func(pageLimit, aviArtifacts int, filter weles.ArtifactFilter,
					sorter weles.ArtifactSorter, statusCode int, amerr error) {

					apiDefaults.PageLimit = int32(pageLimit)
					artifactInfo := fixtures.CreateArtifactInfoSlice(aviArtifacts)
					paginator := weles.ArtifactPagination{Limit: apiDefaults.PageLimit}

					if pageLimit == 0 {
						paginator.Forward = false
					} else {
						paginator.Forward = true
					}
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(aviArtifacts),
						RemainingRecords: 0,
					}

					mockArtifactManager.EXPECT().ListArtifact(filter, defaultSorterA, paginator).Return(
						artifactInfo, listInfo, amerr)
					reqBody := filterSorterReqBody(filter, sorter, JSON)
					client := testserver.Client()
					req := createRequest(reqBody, "", "", JSON, JSON)
					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					defer resp.Body.Close()
					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedArtifactErr(respBody, amerr, JSON)

					Expect(resp.StatusCode).To(Equal(statusCode))
					Expect(resp.Header.Get("Next")).To(Equal(""))
					Expect(resp.Header.Get("Previous")).To(Equal(""))
					Expect(resp.Header.Get("TotalRecords")).To(Equal(""))
					Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))

				},
				Entry("404 status, Artifact not found error, "+
					"when server has 0 artifacts avaliable,pagination off",
					0, 0, emptyFilterA, emptySorterA,
					404, weles.ErrArtifactNotFound),

				Entry("404 status, Artifact not found error "+
					"when server has 0 jobs avaliable, pagination on",
					100, 0, emptyFilterA, emptySorterA,
					404, weles.ErrArtifactNotFound),
				Entry("404 status, Artifact not found error, "+
					"when server has 100 artifacts but none fulfilling filter, pagination off",
					0, 100, filledFilterA1, emptySorterA,
					404, weles.ErrArtifactNotFound),
				Entry("404status, Artifact not found error "+
					"when server has 100 artifacts but none fulfilling filter, pagination on",
					100, 100, filledFilterA1, emptySorterA,
					404, weles.ErrArtifactNotFound),
				Entry("500 status, ArtifactManager unexpected error "+
					"when server has 100 artifacts, pagination off",
					0, 100, emptyFilterA, emptySorterA,
					500, errors.New("This is some errors string")),
				Entry("500 status, ArtifactManager unexpected error "+
					"when server has 100 artifacts, pagination on",
					100, 100, emptyFilterA, emptySorterA,
					500, errors.New("This is some errors string")),
			)
		})

		DescribeTable("a: Error returned by server due to both before and after query params set",
			func(defaultPageLimit int32, query string, filter weles.ArtifactFilter,
				sorter weles.ArtifactSorter) {

				apiDefaults.PageLimit = defaultPageLimit

				reqBody := filterSorterReqBody(filter, sorter, JSON)
				req := createRequest(reqBody, "", query, JSON, JSON)

				client := testserver.Client()
				resp, err := client.Do(req)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				respBody, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				checkReceivedArtifactErr(respBody, weles.ErrBeforeAfterNotAllowed, JSON)

				Expect(resp.StatusCode).To(Equal(400))
				Expect(resp.Header.Get("Next")).To(Equal(""))
				Expect(resp.Header.Get("Previous")).To(Equal(""))
				Expect(resp.Header.Get("TotalRecords")).To(Equal(""))
				Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))

			},
			Entry("a: json, pagination on",
				int32(100), "?before=10&after=20", emptyFilterA, emptySorterA),
		)

	})
})
