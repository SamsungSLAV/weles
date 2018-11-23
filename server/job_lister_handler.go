// Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package server

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/server/operations/jobs"
)

// JobLister is a handler which passess requests for listing jobs to jobmanager.
func (a *APIDefaults) JobLister(params jobs.JobListerParams) middleware.Responder {
	paginator := weles.JobPagination{}
	if a.PageLimit != 0 {
		if (params.After != nil) && (params.Before != nil) {
			return jobs.NewJobListerBadRequest().WithPayload(
				&weles.ErrResponse{Message: weles.ErrBeforeAfterNotAllowed.Error()})
		}
		paginator = setJobPaginator(params, a.PageLimit)
	}
	filter := setJobFilter(params.JobFilterAndSort.Filter)
	sorter := setJobSorter(params.JobFilterAndSort.Sorter)

	jobInfoReceived, listInfo, err := a.Managers.JM.ListJobs(filter, sorter, paginator)
	if err != nil {
		// due to weles.ErrInvalidArgument implementing error interface rather than being error
		// (which is intentional as we want to pass underlying error) switch err.(type) checks only
		// weles.ErrInvalidArgument. Rest of error handling should be in default clause of the type
		// switch.
		switch err.(type) {
		default:
			if err == weles.ErrJobNotFound {
				return jobs.NewJobListerNotFound().WithPayload(
					&weles.ErrResponse{Message: weles.ErrJobNotFound.Error()})
			}
			return jobs.NewJobListerInternalServerError().WithPayload(
				&weles.ErrResponse{Message: err.Error()})
		case weles.ErrInvalidArgument:
			return jobs.NewJobListerBadRequest().WithPayload(
				&weles.ErrResponse{Message: err.Error()})
		}
	}
	jobInfoReturned := jobInfoReceivedToReturned(jobInfoReceived)

	if (listInfo.RemainingRecords == 0) || (a.PageLimit == 0) {
		return responder200(listInfo, paginator, jobInfoReturned, a.PageLimit)
	}
	return responder206(listInfo, paginator, jobInfoReturned, a.PageLimit)
}

func responder206(listInfo weles.ListInfo, paginator weles.JobPagination,
	jobInfoReturned []*weles.JobInfo, defaultPageLimit int32,
) (responder *jobs.JobListerPartialContent) {
	var jobListerURL jobs.JobListerURL

	responder = jobs.NewJobListerPartialContent()
	responder.SetTotalRecords(listInfo.TotalRecords)
	responder.SetRemainingRecords(listInfo.RemainingRecords)

	tmp := uint64(jobInfoReturned[len(jobInfoReturned)-1].JobID)
	jobListerURL.After = &tmp

	if defaultPageLimit != paginator.Limit {
		tmp := paginator.Limit
		jobListerURL.Limit = &tmp
	}
	responder.SetNext(jobListerURL.String())
	if paginator.JobID != 0 { // not the first page
		var jobListerURL jobs.JobListerURL
		tmp = uint64(jobInfoReturned[0].JobID)
		jobListerURL.Before = &tmp
		if defaultPageLimit != paginator.Limit {
			tmp := paginator.Limit
			jobListerURL.Limit = &tmp
		}
		responder.SetPrevious(jobListerURL.String())
	}
	responder.SetPayload(jobInfoReturned)
	return
}

func responder200(listInfo weles.ListInfo, paginator weles.JobPagination,
	jobInfoReturned []*weles.JobInfo, defaultPageLimit int32,
) (responder *jobs.JobListerOK) {
	var jobListerURL jobs.JobListerURL

	responder = jobs.NewJobListerOK()
	responder.SetTotalRecords(listInfo.TotalRecords)

	if paginator.JobID != 0 { //not the first page
		// keep in mind that JobID in paginator is taken from query parameter, not jobmanager
		if paginator.Forward {
			tmp := uint64(jobInfoReturned[0].JobID)
			jobListerURL.Before = &tmp
			if defaultPageLimit != paginator.Limit {
				tmp := paginator.Limit
				jobListerURL.Limit = &tmp
			}
			responder.SetPrevious(jobListerURL.String())
		} else {
			tmp := uint64(jobInfoReturned[len(jobInfoReturned)-1].JobID)
			jobListerURL.After = &tmp
			if defaultPageLimit != paginator.Limit {
				tmp := paginator.Limit
				jobListerURL.Limit = &tmp
			}
			responder.SetNext(jobListerURL.String())
		}
	}
	responder.SetPayload(jobInfoReturned)
	return
}

// normalizeDate is a helper function - adjusts 0 value to 0001-01-01T00:00:00.000Z instead of
// Unix 0- 1970-01-01T00:00:00.000Z. This is required by controller.
func normalizeDate(i strfmt.DateTime) strfmt.DateTime {
	if time.Time(i).Unix() != 0 {
		return i
	}
	return strfmt.DateTime{}
}

// setJobFilter adjusts filter's 0 values to be consistent and acceptable by controller.
// Controller treats slices with 0 len as empty, slices with lenght of 1 and empty value should not
// be passed to controller.
func setJobFilter(i *weles.JobFilter) (o weles.JobFilter) {
	if i != nil {
		o.CreatedBefore = normalizeDate(i.CreatedBefore)
		o.CreatedAfter = normalizeDate(i.CreatedAfter)
		o.UpdatedBefore = normalizeDate(i.UpdatedBefore)
		o.UpdatedAfter = normalizeDate(i.UpdatedAfter)

		if len(i.JobID) > 0 {
			if !(len(i.JobID) == 1 && i.JobID[0] == 0) {
				o.JobID = i.JobID
			}
		}
		if len(i.Info) > 0 {
			if !(len(i.Info) == 1 && i.Info[0] == "") {
				o.Info = i.Info
			}
		}
		if len(i.Name) > 0 {
			if !(len(i.Name) == 1 && i.Name[0] == "") {
				o.Name = i.Name
			}
		}
		if len(i.Status) > 0 {
			if !(len(i.Status) == 1 && i.Status[0] == "") {
				o.Status = i.Status
			}
		}
	}
	return
}

func setJobPaginator(params jobs.JobListerParams, defaultPageLimit int32,
) (paginator weles.JobPagination) {
	paginator.Forward = true
	if params.After != nil {
		paginator.JobID = weles.JobID(*params.After)
	} else if params.Before != nil {
		paginator.JobID = weles.JobID(*params.Before)
		paginator.Forward = false
	}

	if params.Limit == nil {
		paginator.Limit = defaultPageLimit
	} else {
		paginator.Limit = *params.Limit
	}
	return
}

// setJobSorter sets default sorter values.
func setJobSorter(si *weles.JobSorter) (so weles.JobSorter) {
	if si == nil {
		return weles.JobSorter{
			SortOrder: weles.SortOrderAscending,
			SortBy:    weles.JobSortByID,
		}
	}
	if si.SortOrder == "" {
		so.SortOrder = weles.SortOrderAscending
	} else {
		so.SortOrder = si.SortOrder
	}
	if si.SortBy == "" {
		so.SortBy = weles.JobSortByID
	} else {
		so.SortBy = si.SortBy
	}
	return
}

//jobInfoReceivedToReturned is a function which changes the slice of values to slice of pointers.
//It is required due to fact that swagger generates responses as slices of pointers rather than
//slices of values that the interface provides.
func jobInfoReceivedToReturned(jobInfoReceived []weles.JobInfo) []*weles.JobInfo {
	jobInfoReturned := make([]*weles.JobInfo, len(jobInfoReceived))
	for i := range jobInfoReceived {
		jobInfoReturned[i] = &jobInfoReceived[i]
	}
	return jobInfoReturned
}
