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
		fEmpty  = weles.JobFilter{}
		fFilled = weles.JobFilter{
			JobID: []weles.JobID{10, 100, 131},
			Info:  []string{"something", "and something else"},
			Name:  []string{"name123"},
			// time.Date nsec arg must be 0 as it is 0ed out when transported via api
			CreatedAfter: strfmt.DateTime(time.Date(2017, time.May, 3, 11, 34, 55, 0, time.UTC)),
		}
		sEmpty    = weles.JobSorter{}
		sDescNoBy = weles.JobSorter{
			Order: enums.SortOrderDescending,
		}
		sAscNoBy = weles.JobSorter{
			Order: enums.SortOrderAscending,
		}
		sNoOrderID = weles.JobSorter{
			By: enums.JobSortByID,
		}
		sNoOrderCreatedDate = weles.JobSorter{
			By: enums.JobSortByCreatedDate,
		}
		sDescID = weles.JobSorter{
			Order: enums.SortOrderDescending,
			By:    enums.JobSortByID,
		}
		sAscID = weles.JobSorter{
			Order: enums.SortOrderAscending,
			By:    enums.JobSortByID,
		}
		// default value
		sDefault = sAscID
		// when pagination is on and no query params are set. When used, limit should also be set.
		pJFwDefaultLimit = weles.JobPaginator{Forward: true, Limit: defaultPageLimit}
		pEmpty           = weles.Paginator{}

		pJDefault = pJFwDefaultLimit

		jobInfo420        = createJobInfoSlice(420)
		jobInfoFirstPage  = jobInfo420[:defaultPageLimit]
		listInfoFirstPage = weles.ListInfo{
			TotalRecords:     uint64(len(jobInfo420)),
			RemainingRecords: uint64(len(jobInfo420) - defaultPageLimit),
		}
	)

	BeforeEach(func() {
		mockCtrl, mockJobManager, _, apiDefaults, testserver = testServerSetup()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		testserver.Close()
	})

	// helper functions
	newHTTPRequest := func(reqBody io.Reader, contentH, acceptH string) (req *http.Request) {
		req, err := http.NewRequest(http.MethodPost, testserver.URL+basePath+listJobsPath, reqBody)
		Expect(err).ToNot(HaveOccurred())
		req.Header.Set("Content-Type", contentH)
		req.Header.Set("Accept", acceptH)
		req.Close = true
		return req
	}

	newBody := func(f weles.JobFilter, s weles.JobSorter, p weles.Paginator) *bytes.Reader {
		data := jobs.JobListerBody{
			Filter:    &f,
			Sorter:    &s,
			Paginator: &p,
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
	Describe("client sends correct request", func() {
		It("server should accept empty post request", func() {
			mockJobManager.EXPECT().ListJobs(fEmpty, sDefault, pJDefault).
				Return(jobInfoFirstPage, listInfoFirstPage, nil)

			resp, err := testserver.Client().Do(newHTTPRequest(nil, JSON, JSON))
			Expect(err).ToNot(HaveOccurred())

			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()
			checkJobInfoMarshalling(respBody, jobInfoFirstPage)
		})

		DescribeTable("server should pass filter to JobManager without modifications",
			func(filter weles.JobFilter) {
				mockJobManager.EXPECT().ListJobs(filter, sDefault, pJDefault).
					Return(jobInfoFirstPage, listInfoFirstPage, nil)

				_, err := testserver.Client().
					Do(newHTTPRequest(newBody(filter, sEmpty, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())
			},
			Entry("when receiving empty filter", fEmpty),
			Entry("when receiving filled filter", fFilled),
		)

		DescribeTable("server should pass sorter to JobManager, but set default values "+
			"on empty fields",
			func(sSent, sExpected weles.JobSorter) {
				mockJobManager.EXPECT().ListJobs(fEmpty, sExpected, pJFwDefaultLimit).
					Return(jobInfoFirstPage, listInfoFirstPage, nil)

				_, err := testserver.Client().Do(
					newHTTPRequest(newBody(fEmpty, sSent, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())
			},
			Entry("should set default order and by", sEmpty, sDefault),
			Entry("should pass ascending order and by ID", sAscID, sAscID),
			Entry("should pass descending order and by ID", sDescID, sDescID),
			Entry("should pass descending order and set default by",
				sDescNoBy, weles.JobSorter{
					Order: sDescNoBy.Order,
					By:    sDefault.By,
				}),
			Entry("should pass ascending order and set default by",
				sAscNoBy, weles.JobSorter{
					Order: sAscNoBy.Order,
					By:    sDefault.By,
				}),
			Entry("should pass by ID and set default order",
				sNoOrderID, weles.JobSorter{
					Order: sDefault.Order,
					By:    sNoOrderID.By,
				}),
			Entry("should pass by CreatedDate and set default order",
				sNoOrderCreatedDate, weles.JobSorter{
					Order: sDefault.Order,
					By:    sNoOrderCreatedDate.By,
				}),
		)

		DescribeTable("server should pass paginator object to JobManager, but set default values "+
			"on empty fields",
			func(globalLimit int32, pSent weles.Paginator, pExpected weles.JobPaginator) {
				apiDefaults.PageLimit = globalLimit
				listInfo := weles.ListInfo{
					TotalRecords:     uint64(len(jobInfo420)),
					RemainingRecords: 0,
				}
				mockJobManager.EXPECT().ListJobs(
					fEmpty, sDefault, pExpected).Return(
					jobInfo420, listInfo, nil)

				_, err := testserver.Client().
					Do(newHTTPRequest(newBody(fEmpty, sDefault, pSent), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())
			},
			Entry("when empty Paginator is sent",
				int32(defaultPageLimit), weles.Paginator{}, pJDefault),
			Entry("should set pagination direction to Forward when no direction is supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: defaultPageLimit},
				weles.JobPaginator{Forward: true, Limit: defaultPageLimit}),
			Entry("should pass Forward direction when supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: defaultPageLimit, Direction: enums.DirectionForward},
				weles.JobPaginator{Forward: true, Limit: defaultPageLimit}),
			Entry("should pass Backward direction when supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: defaultPageLimit, Direction: enums.DirectionBackward},
				weles.JobPaginator{Forward: false, Limit: defaultPageLimit}),
			Entry("should pass Limit when supplied",
				int32(defaultPageLimit),
				weles.Paginator{Limit: 69},
				weles.JobPaginator{Limit: 69, Forward: true}),
			Entry("should set Limit to globalLimit when not supplied",
				int32(defaultPageLimit),
				weles.Paginator{},
				weles.JobPaginator{Limit: defaultPageLimit, Forward: true}),
			Entry("should set Limit to globalLimit when not supplied",
				int32(69),
				weles.Paginator{},
				weles.JobPaginator{Limit: 69, Forward: true}),
			Entry("should pass ID when supplied",
				int32(defaultPageLimit),
				weles.Paginator{ID: 50},
				weles.JobPaginator{JobID: 50, Limit: defaultPageLimit, Forward: true}),
		)

		DescribeTable("server should respond with 200/206 depending on "+
			"ListInfo.RemainingRecords returned by JobManager",
			func(listInfo weles.ListInfo, statusCode int) {
				mockJobManager.EXPECT().ListJobs(fEmpty, sDefault, pJDefault).
					Return(jobInfo420, listInfo, nil)

				resp, err := testserver.Client().
					Do(newHTTPRequest(newBody(fEmpty, sDefault, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(statusCode))
			},
			Entry("first and last page",
				weles.ListInfo{RemainingRecords: 0}, 200),
			Entry("first page out of n (n>3)",
				weles.ListInfo{RemainingRecords: 20}, 206),
		)

		DescribeTable("Should set Weles-List-{Total,Remaining,Batch-Size} "+
			"based on ListInfo and JobInfo",
			func(jobInfo []weles.JobInfo, listInfo weles.ListInfo) {
				mockJobManager.EXPECT().ListJobs(fEmpty, sDefault, pJDefault).
					Return(jobInfo, listInfo, nil)

				resp, err := testserver.Client().Do(
					newHTTPRequest(newBody(fEmpty, sDefault, pEmpty), JSON, JSON))
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
			Entry("case 3",
				jobInfo420[:50], weles.ListInfo{TotalRecords: 100, RemainingRecords: 0}),
		)
	})

	Describe("JobManager returns error", func() {
		DescribeTable("Server should return appropriate status code and error message",
			func(statusCode int, jmErr error) {
				mockJobManager.EXPECT().ListJobs(fEmpty, sDefault, pJDefault).
					Return(jobInfoFirstPage, listInfoFirstPage, jmErr)

				resp, err := testserver.Client().Do(
					newHTTPRequest(newBody(fEmpty, sDefault, pEmpty), JSON, JSON))
				Expect(err).ToNot(HaveOccurred())

				respBody, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				Expect(err).ToNot(HaveOccurred())

				checkErrorMarshalling(respBody, jmErr)
				Expect(resp.StatusCode).To(Equal(statusCode))
				// should not set headers on error
				Expect(resp.Header.Get(NextPageHdr)).To(Equal(""))
				Expect(resp.Header.Get(PreviousPageHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListTotalHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListRemainingHdr)).To(Equal(""))
				Expect(resp.Header.Get(ListBatchSizeHdr)).To(Equal(""))

			},
			Entry("404 status, Job not found error",
				404, weles.ErrJobNotFound),
			Entry("500 status, Unexpected error",
				500, errors.New("This is unexpected error")),
		)
	})
})
