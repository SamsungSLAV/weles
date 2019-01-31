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
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/enums"
	"github.com/SamsungSLAV/weles/mock"
	"github.com/SamsungSLAV/weles/server"
	"github.com/SamsungSLAV/weles/server/operations/jobs"
)

var _ = Describe("Listing jobs with server initialized", func() {
	var (
		mockCtrl       *gomock.Controller
		apiDefaults    *server.APIDefaults
		mockJobManager *mock.MockJobManager
		testserver     *httptest.Server
	)

	// data to test against
	var (
		emptyFilter = weles.JobFilter{}

		filletFilter = weles.JobFilter{
			JobID: []weles.JobID{10, 100, 131},
			Info:  []string{"something", "and something else"},
			Name:  []string{"name123"},
			// time.Date nsec arg must be 0 as it is 0ed out when transported via api
			CreatedAfter: strfmt.DateTime(time.Date(2017, time.May, 3, 11, 34, 55, 0, time.UTC)),
		}

		sorterEmpty = weles.JobSorter{}

		sorterDescNoBy = weles.JobSorter{
			Order: enums.SortOrderDescending,
		}

		sorterAscNoBy = weles.JobSorter{
			Order: enums.SortOrderAscending,
		}

		sorterNoOrderID = weles.JobSorter{
			By: enums.JobSortByID,
		}

		sorterNoOrderCreatedDate = weles.JobSorter{
			By: enums.JobSortByCreatedDate,
		}

		sorterDescID = weles.JobSorter{
			Order: enums.SortOrderDescending,
			By:    enums.JobSortByID,
		}

		sorterAscID = weles.JobSorter{
			Order: enums.SortOrderAscending,
			By:    enums.JobSortByID,
		}

		// default value
		sorterDefault = sorterAscID

		// when pagination is on and no query params are set. When used, limit should also be set.
		emptyPaginatorOn = weles.JobPaginator{Forward: true}
		// when pagination is off
		emptyPaginatorOff = weles.JobPaginator{}

		jobInfo420 = createJobInfoSlice(420)
	)

	BeforeEach(func() {
		mockCtrl, mockJobManager, _, apiDefaults, testserver = testServerSetup()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		testserver.Close()
	})

	// helper functions
	newHTTPRequest := func(reqBody io.Reader, query, contentH, acceptH string) (req *http.Request) {
		req, err := http.NewRequest(http.MethodPost, testserver.URL+basePath+listJobsPath+query,
			reqBody)
		Expect(err).ToNot(HaveOccurred())
		req.Header.Set("Content-Type", contentH)
		req.Header.Set("Accept", acceptH)
		req.Close = true
		return req
	}

	newBody := func(f weles.JobFilter, s weles.JobSorter) *bytes.Reader {
		data := jobs.JobListerBody{
			Filter: &f,
			Sorter: &s,
		}
		marshalled, err := json.Marshal(data)
		Expect(err).ToNot(HaveOccurred())
		return bytes.NewReader(marshalled)
	}

	checkJobInfoMarshalling := func(respBody []byte, jobInfo []weles.JobInfo) {
		marshalled, err := json.Marshal(jobInfo)
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
					TotalRecords:     uint64(len(jobInfo420)),
					RemainingRecords: 0,
				}
				mockJobManager.EXPECT().ListJobs(
					emptyFilter, sorterDefault, emptyPaginatorOff).Return(
					jobInfo420, listInfo, nil)

				_, err := testserver.Client().Do(newHTTPRequest(nil, "", JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

			})

			DescribeTable("server should ignore query params",
				func(query string) {
					apiDefaults.PageLimit = 0
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(jobInfo420)),
						RemainingRecords: 0,
					}
					mockJobManager.EXPECT().ListJobs(emptyFilter,
						sorterDefault, emptyPaginatorOff).Return(
						jobInfo420, listInfo, nil)

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

			DescribeTable("server should pass filter to JobManager",
				func(filter weles.JobFilter) {
					apiDefaults.PageLimit = 0
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(jobInfo420)),
						RemainingRecords: 0,
					}
					mockJobManager.EXPECT().ListJobs(
						filter, sorterDefault, emptyPaginatorOff).Return(
						jobInfo420, listInfo, nil)

					_, err := testserver.Client().Do(
						newHTTPRequest(newBody(filter, sorterEmpty), "", JSON, JSON))
					Expect(err).ToNot(HaveOccurred())
				},
				Entry("when receiving empty filter", emptyFilter),
				Entry("when receiving filled filter", filletFilter),
			)

			DescribeTable("server should pass sorter to JobManager, but set default values "+
				"on empty fields",
				func(sent, expected weles.JobSorter) {
					apiDefaults.PageLimit = 0
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(jobInfo420)),
						RemainingRecords: 0,
					}
					mockJobManager.EXPECT().ListJobs(
						emptyFilter, expected, emptyPaginatorOff).Return(
						jobInfo420, listInfo, nil)

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
					sorterDescNoBy, weles.JobSorter{
						Order: sorterDescNoBy.Order,
						By:    sorterDefault.By,
					}),
				Entry("should pass ascending order and set default by",
					sorterAscNoBy, weles.JobSorter{
						Order: sorterAscNoBy.Order,
						By:    sorterDefault.By,
					}),
				Entry("should pass by ID and set default order",
					sorterNoOrderID, weles.JobSorter{
						Order: sorterDefault.Order,
						By:    sorterNoOrderID.By,
					}),
				Entry("should pass by CreatedDate and set default order",
					sorterNoOrderCreatedDate, weles.JobSorter{
						Order: sorterDefault.Order,
						By:    sorterNoOrderCreatedDate.By,
					}),
			)
			DescribeTable("should respond with all jobs and correct headers",
				func(recordCount int) {
					apiDefaults.PageLimit = 0
					jobInfo := createJobInfoSlice(recordCount)
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(jobInfo)),
						RemainingRecords: 0,
					}
					mockJobManager.EXPECT().ListJobs(emptyFilter,
						sorterDefault, emptyPaginatorOff).Return(
						jobInfo, listInfo, nil)

					resp, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
					Expect(err).ToNot(HaveOccurred())

					respBody, err := ioutil.ReadAll(resp.Body)
					defer resp.Body.Close()
					Expect(err).ToNot(HaveOccurred())

					checkJobInfoMarshalling(respBody, jobInfo)

					Expect(resp.StatusCode).To(Equal(200))
					Expect(resp.Header.Get(NextPageHdr)).To(Equal(""))
					Expect(resp.Header.Get(PreviousPageHdr)).To(Equal(""))
					Expect(resp.Header.Get(ListRemainingHdr)).To(Equal("0"))
					Expect(resp.Header.Get(ListTotalHdr)).To(Equal(strconv.Itoa(
						len(jobInfo))))
					Expect(resp.Header.Get(ListTotalHdr)).To(Equal(strconv.Itoa(
						len(jobInfo))))

				},
				Entry("20 records avaliable", 20),
				Entry("420 records avaliable", 420),
			)

		})

	})

	Describe("JobManager returns error", func() {
		DescribeTable("Server should return appropriate status code and error message",
			func(pageLimit int32, statusCode int, amerr error) {
				apiDefaults.PageLimit = pageLimit
				listInfo := weles.ListInfo{
					TotalRecords:     uint64(len(jobInfo420)),
					RemainingRecords: 0,
				}
				var paginator weles.JobPaginator
				if pageLimit == 0 {
					paginator = emptyPaginatorOff
				} else {
					paginator = emptyPaginatorOn
					paginator.Limit = pageLimit
				}
				mockJobManager.EXPECT().ListJobs(
					emptyFilter, sorterDefault, paginator).Return(
					jobInfo420, listInfo, amerr)

				resp, err := testserver.Client().Do(
					newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

				respBody, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
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
			Entry("pagination off, 404 status, Job not found error",
				int32(0), 404, weles.ErrJobNotFound),
			Entry("pagination on, 404 status, Job not found error",
				int32(100), 404, weles.ErrJobNotFound),
			Entry("pagination off, 500 status, Unexpected error",
				int32(0), 500, errors.New("This is unexpected error")),
			Entry("pagination on, 500 status, Unexpected error",
				int32(100), 500, errors.New("This is unexpected error")),
		)
	})
	Describe("Pagination turned on", func() {
		Describe("Correct request", func() {
			DescribeTable("server should set paginator object depending on query params",
				func(query string, expectedPaginator weles.JobPaginator) {
					apiDefaults.PageLimit = 500
					listInfo := weles.ListInfo{
						TotalRecords:     uint64(len(jobInfo420)),
						RemainingRecords: 0,
					}
					mockJobManager.EXPECT().ListJobs(
						emptyFilter, sorterDefault, expectedPaginator).Return(
						jobInfo420, listInfo, nil)

					_, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), query, JSON, JSON))
					Expect(err).ToNot(HaveOccurred())

				},
				Entry("when no query params set", "",
					weles.JobPaginator{Forward: true, Limit: 500}),
				Entry("when after param is set", "?after=30",
					weles.JobPaginator{Forward: true, Limit: 500, JobID: 30}),
				Entry("when after and limit params are set", "?after=30&limit=20",
					weles.JobPaginator{Forward: true, Limit: 20, JobID: 30}),
				Entry("when before param is set", "?before=30",
					weles.JobPaginator{Forward: false, Limit: 500, JobID: 30}),
				Entry("when before and limit params are set", "?before=30&limit=15",
					weles.JobPaginator{Forward: false, Limit: 15, JobID: 30}),
				Entry("when limit param is set", "?limit=30",
					weles.JobPaginator{Forward: true, Limit: 30}),
			)

			DescribeTable("server should respond with 200/206 depending on "+
				"ListInfo.RemainingRecords returned by JobManager",
				func(listInfo weles.ListInfo, statusCode int) {
					apiDefaults.PageLimit = 100
					paginator := emptyPaginatorOn
					paginator.Limit = apiDefaults.PageLimit

					mockJobManager.EXPECT().ListJobs(
						emptyFilter, sorterDefault, paginator).Return(
						jobInfo420, listInfo, nil)

					resp, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
					Expect(err).ToNot(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(statusCode))
				},
				Entry("first and last page",
					weles.ListInfo{RemainingRecords: 0}, 200),
				Entry("first page out of n (n>3)",
					weles.ListInfo{RemainingRecords: 20}, 206),
			)

			DescribeTable("on last page, server should NOT set headers:",
				func(paginator weles.JobPaginator, query string) {
					apiDefaults.PageLimit = 100
					paginator.Limit = apiDefaults.PageLimit

					mockJobManager.EXPECT().
						ListJobs(emptyFilter, sorterDefault, paginator).
						Return(jobInfo420,
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

				},
				Entry("Weles-Next-Page when paginating forward (page n/n, n!=0)",
					weles.JobPaginator{JobID: 400, Forward: true}, "?after=400"),
				Entry("Weles-Previous-Page when paginating forward (page 0/n, n!=0)",
					weles.JobPaginator{JobID: 100, Forward: false}, "?before=100"),
			)

			// Correctness of setting Weles-{Next,Previous}-Page is omitted due to implementation
			// being generated swagger code.

			DescribeTable("Should set Weles-List-{Total,Remaining,Batch-Size} "+
				"based on ListInfo and JobInfo",
				func(jobInfo []weles.JobInfo, listInfo weles.ListInfo) {
					apiDefaults.PageLimit = 100
					paginator := emptyPaginatorOn
					paginator.Limit = apiDefaults.PageLimit

					mockJobManager.EXPECT().
						ListJobs(emptyFilter, sorterDefault, paginator).
						Return(jobInfo, listInfo, nil)

					resp, err := testserver.Client().Do(
						newHTTPRequest(newBody(emptyFilter, sorterDefault), "", JSON, JSON))
					Expect(err).ToNot(HaveOccurred())
					defer resp.Body.Close()

					Expect(resp.Header.Get(ListTotalHdr)).
						To(Equal(strconv.FormatUint(listInfo.TotalRecords, 10)))
					Expect(resp.Header.Get(ListRemainingHdr)).
						To(Equal(strconv.FormatUint(listInfo.RemainingRecords, 10)))
					Expect(resp.Header.Get(ListBatchSizeHdr)).
						To(Equal(strconv.Itoa(len(jobInfo))))
				},
				Entry("case 1",
					jobInfo420, weles.ListInfo{TotalRecords: 420, RemainingRecords: 50}),
				Entry("case 2",
					jobInfo420[:100], weles.ListInfo{TotalRecords: 100, RemainingRecords: 10}),
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
