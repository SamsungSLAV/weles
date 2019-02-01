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
	"github.com/SamsungSLAV/weles/enums"
	"github.com/SamsungSLAV/weles/server/operations/artifacts"
)

// ArtifactLister is a handler which passess requests for listing artifacts to ArtifactManager.
func (a *APIDefaults) ArtifactLister(params artifacts.ArtifactListerParams) middleware.Responder {
	paginator := weles.ArtifactPaginator{}
	if a.PageLimit != 0 {
		paginator = *params.ArtifactListBody.Paginator
	}
	filter := setArtifactFilter(params.ArtifactListBody.Filter)
	sorter := setArtifactSorter(params.ArtifactListBody.Sorter)

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
		return artifacts.NewArtifactListerOK().
			WithWelesListTotal(listInfo.TotalRecords).
			WithWelesListBatchSize(int32(len(artifactInfoReturned))).
			WithPayload(artifactInfoReturned)
	} else {
		return artifacts.NewArtifactListerPartialContent().
			WithWelesListTotal(listInfo.TotalRecords).
			WithWelesListRemaining(listInfo.RemainingRecords).
			WithWelesListBatchSize(int32(len(artifactInfoReturned))).
			WithPayload(artifactInfoReturned)
	}
	// should never happen but better be safe.
	return artifacts.NewArtifactListerInternalServerError().
		WithPayload(&weles.ErrResponse{Message: "Unknown internal error ocurred."})
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
			Order: enums.SortOrderAscending,
			By:    enums.ArtifactSortByID,
		}
	}
	if si.Order == "" {
		so.Order = enums.SortOrderAscending
	} else {
		so.Order = si.Order
	}
	if si.By == "" {
		so.By = enums.ArtifactSortByID
	} else {
		so.By = si.By
	}
	return
}

// artifactInfoReceivedToReturn creates slice of pointers from slice of values of ArtifactInfo
// struct. Very similiar function can be found in job_lister_handler.go. Separate functions are
// present as generic one would need to use reflect which affects performance.
func artifactInfoReceivedToReturn(artifactInfoReceived []weles.ArtifactInfo,
) []*weles.ArtifactInfoExt {
	artifactInfoReturned := make([]*weles.ArtifactInfoExt, len(artifactInfoReceived))
	for i := range artifactInfoReceived {
		artifactInfoReturned[i] = &artifactInfoReceived[i].ArtifactInfoExt
	}
	return artifactInfoReturned

}
