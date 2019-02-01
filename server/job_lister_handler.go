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
	"github.com/SamsungSLAV/weles/enums"
	"github.com/SamsungSLAV/weles/server/operations/jobs"
)

// JobLister is a handler which passess requests for listing jobs to jobmanager.
func (a *APIDefaults) JobLister(params jobs.JobListerParams) middleware.Responder {
	paginator := weles.JobPaginator{}
	if a.PageLimit != 0 {
		paginator = *params.JobListBody.Paginator
	}
	filter := setJobFilter(params.JobListBody.Filter)
	sorter := setJobSorter(params.JobListBody.Sorter)

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
		return jobs.NewJobListerOK().
			WithWelesListTotal(listInfo.TotalRecords).
			WithWelesListBatchSize(int32(len(jobInfoReturned))).
			WithPayload(jobInfoReturned)
	} else {
		return jobs.NewJobListerPartialContent().
			WithWelesListTotal(listInfo.TotalRecords).
			WithWelesListRemaining(listInfo.RemainingRecords).
			WithWelesListBatchSize(int32(len(jobInfoReturned))).
			WithPayload(jobInfoReturned)
	}
	// should never happen but better be safe.
	return jobs.NewJobListerInternalServerError().
		WithPayload(&weles.ErrResponse{Message: "Unkown internal error occurred."})

}

// normalizeDate is a helper function - adjusts 0 value to "0001-01-01T00:00:00.000Z" instead of
// Unix 0 "1970-01-01T00:00:00.000Z". This is required by controller.
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

// setJobSorter sets default sorter values.
func setJobSorter(si *weles.JobSorter) (so weles.JobSorter) {
	if si == nil {
		return weles.JobSorter{
			Order: enums.SortOrderAscending,
			By:    enums.JobSortByID,
		}
	}
	if si.Order == "" {
		so.Order = enums.SortOrderAscending
	} else {
		so.Order = si.Order
	}
	if si.By == "" {
		so.By = enums.JobSortByID
	} else {
		so.By = si.By
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
