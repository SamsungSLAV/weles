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
	"github.com/go-openapi/runtime/middleware"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/server/operations/artifacts"
)

// ArtifactLister is a handler which passess requests for listing artifacts to ArtifactManager.
func (a *APIDefaults) ArtifactLister(params artifacts.ArtifactListerParams) middleware.Responder {
	paginator := weles.ArtifactPagination{}
	if a.PageLimit != 0 {
		if (params.After != nil) && (params.Before != nil) {
			return artifacts.NewArtifactListerBadRequest().WithPayload(&weles.ErrResponse{
				Message: weles.ErrBeforeAfterNotAllowed.Error()})
		}
		paginator = setArtifactPaginator(params, a.PageLimit)
	}
	filter := setArtifactFilter(params.ArtifactFilterAndSort.Filter)
	sorter := setArtifactSorter(params.ArtifactFilterAndSort.Sorter)

	artifactInfoReceived, listInfo, err := a.Managers.AM.ListArtifact(filter, sorter, paginator)

	switch err {
	default:
		return artifacts.NewArtifactListerInternalServerError().WithPayload(
			&weles.ErrResponse{Message: err.Error()})
	case weles.ErrArtifactNotFound:
		return artifacts.NewArtifactListerNotFound().WithPayload(
			&weles.ErrResponse{Message: weles.ErrArtifactNotFound.Error()})
	case nil:
	}

	artifactInfoReturned := artifactInfoReceivedToReturn(artifactInfoReceived)

	if (listInfo.RemainingRecords == 0) || (a.PageLimit == 0) { //last page...
		return responderArtifact200(listInfo, paginator, artifactInfoReturned, a.PageLimit)
	} //not last page...
	return responderArtifact206(listInfo, paginator, artifactInfoReturned, a.PageLimit)
}

// responderArtifact206 builds 206 HTTP response with appropriate headers and body.
func responderArtifact206(listInfo weles.ListInfo, paginator weles.ArtifactPagination,
	artifactInfoReturned []*weles.ArtifactInfo, defaultPageLimit int32,
) (responder *artifacts.ArtifactListerPartialContent) {
	var artifactListerURL artifacts.ArtifactListerURL

	responder = artifacts.NewArtifactListerPartialContent()
	responder.SetTotalRecords(listInfo.TotalRecords)
	responder.SetRemainingRecords(listInfo.RemainingRecords)

	tmp := artifactInfoReturned[len(artifactInfoReturned)-1].ID
	artifactListerURL.After = &tmp

	if defaultPageLimit != paginator.Limit {
		tmp := paginator.Limit
		artifactListerURL.Limit = &tmp
	}
	responder.SetNext(artifactListerURL.String())

	if paginator.ID != 0 { //... and not the first
		//paginator.ID is from query parameter not artifactmanager
		var artifactListerURL artifacts.ArtifactListerURL
		tmp = artifactInfoReturned[0].ID
		artifactListerURL.Before = &tmp
		if defaultPageLimit != paginator.Limit {
			tmp := paginator.Limit
			artifactListerURL.Limit = &tmp
		}
		responder.SetPrevious(artifactListerURL.String())
	}
	responder.SetPayload(artifactInfoReturned)
	return
}

// responderArtifact200 builds 200 HTTP response with appropriate headers and body.
func responderArtifact200(listInfo weles.ListInfo, paginator weles.ArtifactPagination,
	artifactInfoReturned []*weles.ArtifactInfo, defaultPageLimit int32,
) (responder *artifacts.ArtifactListerOK) {
	var artifactListerURL artifacts.ArtifactListerURL
	responder = artifacts.NewArtifactListerOK()
	responder.SetTotalRecords(listInfo.TotalRecords)
	if paginator.ID != 0 { //not the first page
		// keep in mind that ArtifactPath in paginator is taken from query parameter,
		// not ArtifactManager
		if paginator.Forward {
			tmp := artifactInfoReturned[0].ID
			artifactListerURL.Before = &tmp
			if defaultPageLimit != paginator.Limit {
				tmp := paginator.Limit
				artifactListerURL.Limit = &tmp
			}
			responder.SetPrevious(artifactListerURL.String())
		} else {
			tmp := artifactInfoReturned[len(artifactInfoReturned)-1].ID
			artifactListerURL.After = &tmp
			if defaultPageLimit != paginator.Limit {
				tmp2 := paginator.Limit
				artifactListerURL.Limit = &tmp2
			}
			responder.SetNext(artifactListerURL.String())
		}
	}
	responder.SetPayload(artifactInfoReturned)
	return
}

// setArtifactFilter adjusts filter's 0 values to be consistent and acceptable by the artifacts db
// That is []string with only 1 empty element should be removed.
func setArtifactFilter(fi *weles.ArtifactFilter) (fo weles.ArtifactFilter) {
	if fi != nil {
		if len(fi.JobID) > 0 {
			fo.JobID = fi.JobID
		}
		if len(fi.Alias) > 0 {
			if !(len(fi.Alias) == 1 && fi.Alias[0] == "") {
				fo.Alias = fi.Alias
			}
		}
		if len(fi.Status) > 0 {
			if !(len(fi.Status) == 1 && fi.Status[0] == "") {
				fo.Status = fi.Status
			}
		}
		if len(fi.Type) > 0 {
			if !(len(fi.Type) == 1 && fi.Type[0] == "") {
				fo.Type = fi.Type
			}
		}
	}
	return
}

// setArtifactSorter sets default sorter values.
func setArtifactSorter(si *weles.ArtifactSorter) (so weles.ArtifactSorter) {
	if si == nil {
		return weles.ArtifactSorter{
			SortOrder: weles.SortOrderAscending,
			SortBy:    weles.ArtifactSortByID,
		}
	}
	if si.SortOrder == "" {
		so.SortOrder = weles.SortOrderAscending
	} else {
		so.SortOrder = si.SortOrder
	}
	if si.SortBy == "" {
		so.SortBy = weles.ArtifactSortByID
	} else {
		so.SortBy = si.SortBy
	}
	return
}

// setArtifactPaginator creates and fills paginator object with default values.
func setArtifactPaginator(params artifacts.ArtifactListerParams, defaultPageLimit int32,
) (paginator weles.ArtifactPagination) {
	paginator.Forward = true
	if params.After != nil {
		paginator.ID = *params.After
	} else if params.Before != nil {
		paginator.ID = *params.Before
		paginator.Forward = false
	}
	if params.Limit == nil {
		paginator.Limit = defaultPageLimit
	} else {
		paginator.Limit = *params.Limit
	}
	return paginator
}

// artifactInfoReceivedToReturn creates slice of pointers from slice of values of ArtifactInfo
// struct. Very similiar function can be found in job_lister_handler.go. Separate functions are
// present as generic one would need to use reflect which affects performance.
func artifactInfoReceivedToReturn(artifactInfoReceived []weles.ArtifactInfo) []*weles.ArtifactInfo {
	artifactInfoReturned := make([]*weles.ArtifactInfo, len(artifactInfoReceived))
	for i := range artifactInfoReceived {
		artifactInfoReturned[i] = &artifactInfoReceived[i]
	}
	return artifactInfoReturned

}
