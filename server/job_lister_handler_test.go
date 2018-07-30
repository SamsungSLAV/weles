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
	"time"

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/mock"
	"git.tizen.org/tools/weles/server"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobListerHandler", func() {

	var (
		mockCtrl       *gomock.Controller
		apiDefaults    *server.APIDefaults
		mockJobManager *mock.MockJobManager
		testserver     *httptest.Server
	)

	BeforeEach(func() {
		mockCtrl, mockJobManager, _, apiDefaults, testserver = testServerSetup()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		testserver.Close()
	})

	Describe("Listing jobs", func() {
		createRequest := func(
			reqBody io.Reader,
			path string,
			query string,
			contentH string,
			acceptH string) (req *http.Request) {
			if path == "" {
				path = "/api/v1/jobs/list"
			}
			req, err := http.NewRequest(http.MethodPost, testserver.URL+path+query, reqBody)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", contentH)
			req.Header.Set("Accept", acceptH)
			return req
		}

		filterSorterReqBody := func(
			filter weles.JobFilter,
			sorter weles.JobSorter,
			contentH string) (
			rb *bytes.Reader) {

			jobFilterSort := weles.JobFilterAndSort{
				Filter: &filter,
				Sorter: &sorter}

			jobFilterSortMarshalled, err := json.Marshal(jobFilterSort)
			Expect(err).ToNot(HaveOccurred())

			return bytes.NewReader(jobFilterSortMarshalled)

		}

		checkReceivedJobInfo := func(respBody []byte, jobInfo []weles.JobInfo, acceptH string) {
			jobInfoMarshalled, err := json.Marshal(jobInfo)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(respBody)).To(MatchJSON(string(jobInfoMarshalled)))

		}

		checkReceivedErr := func(respBody []byte, jmerr error, acceptH string) {
			errMarshalled, err := json.Marshal(weles.ErrResponse{
				Message: jmerr.Error(),
				Type:    ""})
			Expect(err).ToNot(HaveOccurred())
			Expect(string(respBody)).To(MatchJSON(string(errMarshalled)))
		}

		//few structs to test against
		someTime := strfmt.DateTime(time.Date(2017, time.May, 3, 11, 34, 55, 0, time.UTC))
		// time.Date nsec argument must be 0 as it is 0ed out when transported via api.

		filledFilter1 := weles.JobFilter{
			CreatedAfter: someTime,
			Status: []weles.JobStatus{
				weles.JobStatusNEW,
				weles.JobStatusPARSING,
				weles.JobStatusDOWNLOADING}}

		filledFilter2 := weles.JobFilter{
			JobID: []weles.JobID{10, 100, 131},
			Info: []string{
				"something",
				"something else",
				"some really different thing"}}

		filledFilter3 := weles.JobFilter{
			UpdatedBefore: someTime,
			Name: []string{
				"daass",
				"sdasa",
				"asdasf32qw;;dq"}}

		emptyFilter := weles.JobFilter{}

		filledSorter1 := weles.JobSorter{
			SortBy:    weles.JobSortByCreatedDate,
			SortOrder: weles.SortOrderAscending}

		filledSorter2 := weles.JobSorter{
			SortBy:    weles.JobSortByJobStatus,
			SortOrder: weles.SortOrderDescending}

		emptySorter := weles.JobSorter{}

		emptyPaginator := weles.JobPagination{}
		emptyPaginator2 := weles.JobPagination{Forward: true}
		after100 := "?after=100"
		pAfter100 := weles.JobPagination{
			JobID:   weles.JobID(100),
			Forward: true}
		after100Limit50 := "?after=100&limit=50"
		pAfter100Limit50 := weles.JobPagination{
			JobID:   weles.JobID(100),
			Forward: true,
			Limit:   int32(50)}

		before100 := "?before=100"
		pBefore100 := weles.JobPagination{
			JobID:   weles.JobID(100),
			Forward: false}

		before100Limit50 := "?before=100&limit=50"
		pBefore100Limit50 := weles.JobPagination{
			JobID:   weles.JobID(100),
			Forward: false,
			Limit:   int32(50)}

		limit50 := "?limit=50"
		pLimit50 := weles.JobPagination{
			Forward: true,
			Limit:   int32(50)}

		type queryPaginator struct {
			query     string
			paginator weles.JobPagination
		}

		queryPaginatorOK := []queryPaginator{
			{
				query:     "",
				paginator: emptyPaginator2},
			{
				query:     before100,
				paginator: pBefore100},
			{
				query:     before100Limit50,
				paginator: pBefore100Limit50},
			{
				query:     after100,
				paginator: pAfter100},
			{
				query:     after100Limit50,
				paginator: pAfter100Limit50},
			{
				query:     limit50,
				paginator: pLimit50}}

		queriesOK := []string{"", before100, before100Limit50, after100, after100Limit50, limit50}

		jobInfoSlice420 := createJobInfoSlice(420)

		Context("j: Server receives correct request and has pagination turned off", func() {
			for _, currQuery := range queriesOK { // expected behaviour - handler should ignore
				// different queries when pagination is turned off globally.
				DescribeTable("should respond with all avaliable jobs, ignoring query params",
					func(filter weles.JobFilter, sorter weles.JobSorter, query string) {
						apiDefaults.PageLimit = 0

						listInfo := weles.ListInfo{
							TotalRecords:     uint64(len(jobInfoSlice420)),
							RemainingRecords: 0}

						mockJobManager.EXPECT().ListJobs(
							filter, sorter, emptyPaginator).Return(
							jobInfoSlice420, listInfo, nil)

						reqBody := filterSorterReqBody(filter, sorter, JSON)

						client := testserver.Client()
						req := createRequest(reqBody, "", query, JSON, JSON)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())

						defer resp.Body.Close()
						respBody, err := ioutil.ReadAll(resp.Body)
						Expect(err).ToNot(HaveOccurred())

						checkReceivedJobInfo(respBody, jobInfoSlice420, JSON)

						Expect(resp.StatusCode).To(Equal(200))
						Expect(resp.Header.Get("Next")).To(Equal(""))
						Expect(resp.Header.Get("Previous")).To(Equal(""))
						Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))
						Expect(resp.Header.Get("TotalRecords")).To(Equal(strconv.Itoa(len(jobInfoSlice420))))

					},

					Entry("given empty filter and sorter", emptyFilter, emptySorter, currQuery),

					Entry("given filled filter and sorter", filledFilter2, filledSorter2, currQuery),

					Entry("given filled filter and empty sorter", filledFilter1, emptySorter, currQuery),

					Entry("given empty filter and filled sorter", emptyFilter, filledSorter1, currQuery),
				)
			}
		})

		Context("server receives correct request and has pagination turned on", func() {
			jobInfoAll := createJobInfoSlice(100)
			globalLimit := []int32{70, 100, 111}
			for _, currgl := range globalLimit {
				for _, curr := range queryPaginatorOK {
					// expected behaviour - handler passes query parameters to JM.
					// this should impact data returned by JM. It is not reflected in the
					// below mock of JM as it is out of scope of server unit tests.
					DescribeTable("should respond with all jobs",
						func(jobInfo []weles.JobInfo,
							filter weles.JobFilter,
							sorter weles.JobSorter,
							paginator weles.JobPagination,
							query string, gl int32) {

							apiDefaults.PageLimit = gl

							if !strings.Contains(query, "limit") {
								paginator.Limit = apiDefaults.PageLimit
							}

							listInfo := weles.ListInfo{
								TotalRecords:     uint64(len(jobInfo)),
								RemainingRecords: 0}

							mockJobManager.EXPECT().ListJobs(
								filter, sorter, paginator).Return(
								jobInfo, listInfo, nil)

							reqBody := filterSorterReqBody(filter, sorter, JSON)
							req := createRequest(reqBody, "", query, JSON, JSON)

							client := testserver.Client()
							resp, err := client.Do(req)
							Expect(err).ToNot(HaveOccurred())

							defer resp.Body.Close()
							respBody, err := ioutil.ReadAll(resp.Body)
							Expect(err).ToNot(HaveOccurred())

							checkReceivedJobInfo(respBody, jobInfo, JSON)

							Expect(resp.StatusCode).To(Equal(200))
							// Next and Previous headers are ignored here as they are tested in other context.
							Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))
							Expect(resp.Header.Get("TotalRecords")).To(
								Equal(strconv.Itoa(len(jobInfo))))
						},
						Entry("given empty request, when JM has less jobs than page size",
							jobInfoAll[:40], emptyFilter, emptySorter, curr.paginator, curr.query, currgl),

						Entry("given filled filter, when JM returns less jobs than Default Page size",
							jobInfoAll[10:67], filledFilter3, emptySorter, curr.paginator, curr.query, currgl),

						Entry("given filled filter, when JM returns same amount of filtered jobs as Default Page Size",
							jobInfoAll, filledFilter1, filledSorter2, curr.paginator, curr.query, currgl),
					)
				}
			}
		})

		Context("Pagination on", func() {
			jobInfoAll := createJobInfoSlice(400)
			DescribeTable("paginating forward",
				func(jobInfo []weles.JobInfo,
					startingPageNo int,
					filter weles.JobFilter,
					sorter weles.JobSorter) {

					apiDefaults.PageLimit = 100

					//prepare data for first call
					jobInfoStartingPage := jobInfo[(startingPageNo-1)*
						int(apiDefaults.PageLimit) : startingPageNo*
						int(apiDefaults.PageLimit)] //first page of data

					startingPageQuery := ""

					paginator := weles.JobPagination{Limit: apiDefaults.PageLimit}
					paginator.Forward = true
					if startingPageNo != 1 {
						paginator.JobID = jobInfoStartingPage[0].JobID
						startingPageQuery = "?after=" +
							jobInfo[(startingPageNo-1)*int(apiDefaults.PageLimit)].JobID.String()
					}

					listInfo := weles.ListInfo{TotalRecords: uint64(len(jobInfo)), RemainingRecords: uint64(len(jobInfo) - startingPageNo*len(jobInfoStartingPage))}

					first := mockJobManager.EXPECT().ListJobs(filter, sorter, paginator).Return(jobInfoStartingPage, listInfo, nil)

					reqBody := filterSorterReqBody(filter, sorter, JSON)
					req := createRequest(reqBody, "", startingPageQuery, JSON, JSON)
					req.Close = true
					client := testserver.Client()
					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedJobInfo(respBody, jobInfoStartingPage, JSON)
					Expect(resp.StatusCode).To(Equal(206))

					Expect(resp.Header.Get("Next")).To(Equal("/api/v1/jobs/list" + "?after=" + jobInfoStartingPage[apiDefaults.PageLimit-1].JobID.String()))
					prevCheck := ""
					if startingPageNo != 1 {
						prevCheck = "/api/v1/jobs/list" + "?before=" + jobInfoStartingPage[0].JobID.String()
					}

					Expect(resp.Header.Get("Previous")).To(Equal(prevCheck))
					Expect(resp.Header.Get("RemainingRecords")).To(Equal(strconv.Itoa(len(jobInfo) - (startingPageNo * len(jobInfoStartingPage)))))
					Expect(resp.Header.Get("TotalRecords")).To(Equal(strconv.Itoa(len(jobInfo))))

					nextPage := resp.Header.Get("Next")
					resp.Body.Close()
					testserver.CloseClientConnections()

					//prepare data for second call based on previous
					var jobInfo2 []weles.JobInfo
					var secondReturnCode int
					if (len(jobInfo) - startingPageNo*int(apiDefaults.PageLimit)) <= int(apiDefaults.PageLimit) {

						jobInfo2 = jobInfo[startingPageNo*int(apiDefaults.PageLimit):] //next page is not full
						secondReturnCode = 200
					} else {
						jobInfo2 = jobInfo[startingPageNo*int(apiDefaults.PageLimit) : (startingPageNo+1)*int(apiDefaults.PageLimit)] //last page is full
						secondReturnCode = 206
					}

					paginator2 := weles.JobPagination{Limit: apiDefaults.PageLimit, Forward: true, JobID: jobInfoStartingPage[int(apiDefaults.PageLimit)-1].JobID}
					listInfo2 := weles.ListInfo{TotalRecords: listInfo.TotalRecords}

					if tmp := len(jobInfo) - (startingPageNo+1)*int(apiDefaults.PageLimit); tmp < 0 {
						listInfo2.RemainingRecords = 0
					} else {
						listInfo2.RemainingRecords = uint64(tmp)
					}
					//filter and sorter should stay the same.
					mockJobManager.EXPECT().ListJobs(filter, sorter, paginator2).Return(jobInfo2, listInfo2, nil).After(first)

					client2 := testserver.Client()
					req2 := createRequest(reqBody, nextPage, "", JSON, JSON)
					req2.Close = true
					resp2, err := client2.Do(req2)
					Expect(err).ToNot(HaveOccurred())

					defer resp2.Body.Close()
					respBody2, err := ioutil.ReadAll(resp2.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedJobInfo(respBody2, jobInfo2, JSON)

					Expect(resp2.StatusCode).To(Equal(secondReturnCode))

					if secondReturnCode == 200 {
						Expect(resp2.Header.Get("Next")).To(Equal(""))
						prevCheck = "/api/v1/jobs/list" + "?before=" + jobInfo2[0].JobID.String()
						Expect(resp2.Header.Get("Previous")).To(Equal(prevCheck))
					} else {
						prevCheck = "/api/v1/jobs/list" + "?before=" + jobInfo2[0].JobID.String()
						nextCheck := "/api/v1/jobs/list" + "?after=" + jobInfo2[int(apiDefaults.PageLimit)-1].JobID.String()
						Expect(resp2.Header.Get("Next")).To(Equal(nextCheck))
						Expect(resp2.Header.Get("Previous")).To(Equal(prevCheck))
					}
					if tmp := strconv.Itoa(len(jobInfo) - startingPageNo*len(jobInfoStartingPage) - len(jobInfo2)); tmp != "0" {
						Expect(resp2.Header.Get("RemainingRecords")).To(Equal(tmp))
					} else {

						Expect(resp2.Header.Get("RemainingRecords")).To(Equal(""))
					}
					Expect(resp2.Header.Get("TotalRecords")).To(Equal(strconv.Itoa(len(jobInfo))))

				},
				Entry("1->2/2", // from 1 to 2 out of 2 (pages)
					jobInfoAll[:170], 1, emptyFilter, emptySorter),
				Entry("1->2/3",
					jobInfoAll[:270], 1, emptyFilter, emptySorter),
				Entry("2->3/3",
					jobInfoAll[:300], 2, emptyFilter, emptySorter),
				Entry("2->3/4",
					jobInfoAll[:350], 2, emptyFilter, emptySorter),
			)
			DescribeTable("paginating backward",
				func(jobInfo []weles.JobInfo,
					startingPageNo int, pages int,
					filter weles.JobFilter,
					sorter weles.JobSorter) {

					apiDefaults.PageLimit = 100
					paginator := weles.JobPagination{Limit: apiDefaults.PageLimit}
					//prepare data for first call
					var jobInfoStartingPage []weles.JobInfo
					var startingPageQuery string
					listInfo := weles.ListInfo{}
					if startingPageNo == pages {
						jobInfoStartingPage = jobInfo[(startingPageNo-1)*
							int(apiDefaults.PageLimit) : len(jobInfo)-1]
						paginator.Forward = true
						paginator.JobID = jobInfoStartingPage[len(jobInfoStartingPage)-1].JobID
						startingPageQuery = "?after=" + paginator.JobID.String()
						listInfo = weles.ListInfo{
							TotalRecords:     uint64(len(jobInfo)),
							RemainingRecords: 0}

					} else {
						jobInfoStartingPage = jobInfo[(startingPageNo)*int(apiDefaults.PageLimit) : (startingPageNo+1)*
							int(apiDefaults.PageLimit)] //first page of data
						paginator.Forward = false
						paginator.JobID = jobInfo[(startingPageNo*int(apiDefaults.PageLimit))-1].JobID
						startingPageQuery = "?before=" + paginator.JobID.String()
						listInfo = weles.ListInfo{
							TotalRecords:     uint64(len(jobInfo)),
							RemainingRecords: uint64(len(jobInfo) - (int(apiDefaults.PageLimit) * startingPageNo))}

					}

					first := mockJobManager.EXPECT().ListJobs(filter, sorter, paginator).Return(jobInfoStartingPage, listInfo, nil)

					reqBody := filterSorterReqBody(filter, sorter, JSON)
					req := createRequest(reqBody, "", startingPageQuery, JSON, JSON)
					req.Close = true
					client := testserver.Client()
					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedJobInfo(respBody, jobInfoStartingPage, JSON)

					if startingPageNo == pages {
						Expect(resp.StatusCode).To(Equal(200))
						Expect(resp.Header.Get("Previous")).To(Equal("/api/v1/jobs/list?before=" + jobInfoStartingPage[0].JobID.String()))
						Expect(resp.Header.Get("Next")).To(Equal(""))
						Expect(resp.Header.Get("TotalRecords")).To(Equal(strconv.Itoa(len(jobInfo))))
						Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))
					} else {
						Expect(resp.StatusCode).To(Equal(206))
						Expect(resp.Header.Get("Previous")).To(Equal("/api/v1/jobs/list?before=" + jobInfoStartingPage[0].JobID.String()))
						Expect(resp.Header.Get("Next")).To(Equal("/api/v1/jobs/list?after=" + jobInfoStartingPage[len(jobInfoStartingPage)-1].JobID.String()))
						Expect(resp.Header.Get("TotalRecords")).To(Equal(strconv.Itoa(len(jobInfo))))
						Expect(resp.Header.Get("RemainingRecords")).To(Equal(strconv.FormatUint(listInfo.RemainingRecords, 10)))
					}

					prevPage := resp.Header.Get("Previous")

					resp.Body.Close()
					testserver.CloseClientConnections()

					//prepare data for second call based on previous

					var jobInfo2 []weles.JobInfo
					paginator2 := weles.JobPagination{Limit: apiDefaults.PageLimit, Forward: false, JobID: jobInfoStartingPage[0].JobID}
					listInfo2 := weles.ListInfo{TotalRecords: listInfo.TotalRecords}
					if startingPageNo == pages {
						jobInfo2 = jobInfo[(startingPageNo-2)*int(apiDefaults.PageLimit) : (startingPageNo-1)*int(apiDefaults.PageLimit)]
						if startingPageNo-1 == 1 {
							listInfo2.RemainingRecords = 0
						} else {
							listInfo2.RemainingRecords = uint64((pages - (startingPageNo - 1)) * int(apiDefaults.PageLimit))
						}
					} else {
						jobInfo2 = jobInfo[(startingPageNo-1)*int(apiDefaults.PageLimit) : startingPageNo*int(apiDefaults.PageLimit)]
						listInfo2.RemainingRecords = uint64(apiDefaults.PageLimit)
					}

					mockJobManager.EXPECT().ListJobs(filter, sorter, paginator2).Return(jobInfo2, listInfo2, nil).After(first)
					client2 := testserver.Client()
					req2 := createRequest(reqBody, prevPage, "", JSON, JSON)
					req2.Close = true
					resp2, err := client2.Do(req2)
					Expect(err).ToNot(HaveOccurred())

					defer resp2.Body.Close()
					respBody2, err := ioutil.ReadAll(resp2.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedJobInfo(respBody2, jobInfo2, JSON)
					if startingPageNo == pages {

						if startingPageNo-1 == 1 {
							Expect(resp2.StatusCode).To(Equal(200))
							Expect(resp2.Header.Get("RemainingRecords")).To(Equal(""))
							Expect(resp2.Header.Get("Previous")).To(Equal(""))

						} else {
							Expect(resp2.StatusCode).To(Equal(206))
							Expect(resp2.Header.Get("RemainingRecords")).To(Equal(strconv.FormatInt(int64(apiDefaults.PageLimit), 10)))
							Expect(resp2.Header.Get("Previous")).To(Equal("/api/v1/jobs/list?before=" + jobInfo2[0].JobID.String()))
						}
					} else {
						Expect(resp2.StatusCode).To(Equal(206))
						Expect(resp2.Header.Get("RemainingRecords")).To(Equal(strconv.FormatInt(int64(apiDefaults.PageLimit), 10)))
					}
					Expect(resp2.Header.Get("Next")).To(Equal("/api/v1/jobs/list?after=" + jobInfo2[len(jobInfo2)-1].JobID.String()))
					Expect(resp2.Header.Get("TotalRecords")).To(Equal(strconv.Itoa(len(jobInfo))))
				},
				// ginkgo does not like It clauses with the same name. To avoid conflicts with
				// tests of listing artifacts, each it clause not referring to jobs exlicitly
				// is prefixed with j:
				Entry("j: 2->1/2", //meaning: from 1 to 2 out of 2 (pages)
					jobInfoAll[:170], 2, 2, emptyFilter, emptySorter),
				Entry("j: 2->1/3",
					jobInfoAll[:270], 2, 3, emptyFilter, emptySorter),
				Entry("j: 3->2/4",
					jobInfoAll[:350], 3, 4, emptyFilter, emptySorter),
				Entry("j: 3->2/3",
					jobInfoAll[:300], 3, 3, emptyFilter, emptySorter),
			)
		})

		Context("There is an error", func() {
			DescribeTable("Server should respond with error from JobManager",
				func(pageLimit int, aviJobs int, filter weles.JobFilter, sorter weles.JobSorter, statusCode int, jmerr error) {
					apiDefaults.PageLimit = int32(pageLimit)
					jobInfo := createJobInfoSlice(aviJobs)
					paginator := weles.JobPagination{Limit: apiDefaults.PageLimit}
					if pageLimit == 0 {
						paginator.Forward = false
					} else {
						paginator.Forward = true
					}
					listInfo := weles.ListInfo{TotalRecords: uint64(aviJobs), RemainingRecords: 0}

					mockJobManager.EXPECT().ListJobs(filter, sorter, paginator).Return(jobInfo, listInfo, jmerr)
					reqBody := filterSorterReqBody(filter, sorter, JSON)
					client := testserver.Client()
					req := createRequest(reqBody, "", "", JSON, JSON)
					resp, err := client.Do(req)
					Expect(err).ToNot(HaveOccurred())

					defer resp.Body.Close()
					respBody, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())

					checkReceivedErr(respBody, jmerr, JSON)

					Expect(resp.StatusCode).To(Equal(statusCode))
					Expect(resp.Header.Get("Next")).To(Equal(""))
					Expect(resp.Header.Get("Previous")).To(Equal(""))
					Expect(resp.Header.Get("TotalRecords")).To(Equal(""))
					Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))

				},
				Entry("404 status, Job not found error  when server has 0 jobs avaliable,pagination off",
					0, 0, emptyFilter, emptySorter, 404, weles.ErrJobNotFound),
				Entry("404 status, Job not found error  when server has 100 jobs but none fulfilling filter, pagination off",
					0, 100, filledFilter1, emptySorter, 404, weles.ErrJobNotFound),
				Entry("500 status, JobManager unexpected error when server has 100 jobs, pagination off",
					0, 100, emptyFilter, emptySorter, 500, errors.New("This is some errors string")),
				Entry("404 status, Job not found error  when server has 0 jobs avaliable,pagination on",
					100, 0, emptyFilter, emptySorter, 404, weles.ErrJobNotFound),
				Entry("404 status, Job not found error  when server has 100 jobs but none fulfilling filter, pagination on",
					100, 100, filledFilter1, emptySorter, 404, weles.ErrJobNotFound),
				Entry("500 status, JobManager unexpected error when server has 100 jobs, pagination on",
					100, 100, emptyFilter, emptySorter, 500, errors.New("This is some errors string")),
			)
		})

		DescribeTable("error returned by server due to both before and after query params set",
			func(defaultPageLimit int32, query string, acceptH string, contentH string, filter weles.JobFilter, sorter weles.JobSorter) {
				apiDefaults.PageLimit = defaultPageLimit

				reqBody := filterSorterReqBody(filter, sorter, contentH)
				req := createRequest(reqBody, "", query, contentH, acceptH)

				client := testserver.Client()
				resp, err := client.Do(req)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				respBody, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				checkReceivedErr(respBody, weles.ErrBeforeAfterNotAllowed, acceptH)

				Expect(resp.StatusCode).To(Equal(400))
				Expect(resp.Header.Get("Next")).To(Equal(""))
				Expect(resp.Header.Get("Previous")).To(Equal(""))
				Expect(resp.Header.Get("TotalRecords")).To(Equal(""))
				Expect(resp.Header.Get("RemainingRecords")).To(Equal(""))

			},
			Entry("json, pagination off",
				int32(0), "?before=10&after=20", JSON, JSON, emptyFilter, emptySorter),
			Entry("json, pagination on",
				int32(100), "?before=10&after=20", JSON, JSON, emptyFilter, emptySorter),
		)

	})
})
