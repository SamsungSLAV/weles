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

	"git.tizen.org/tools/weles"
	"git.tizen.org/tools/weles/server/operations/artifacts"
)

// ArtifactLister is a handler which passess requests for listing artifacts to artifactmanager.
func (a *APIDefaults) ArtifactLister(params artifacts.ArtifactListerParams) middleware.Responder {
	if (params.After != nil) && (params.Before != nil) {
		return artifacts.NewArtifactListerBadRequest().WithPayload(&weles.ErrResponse{
			Message: weles.ErrBeforeAfterNotAllowed.Error()})
	}

	var artifactInfoReceived []weles.ArtifactInfo
	var listInfo weles.ListInfo
	var err error
	paginator := weles.ArtifactPagination{}
	if a.PageLimit != 0 {
		paginator = setArtifactPaginator(params, a.PageLimit)
	}

	if params.ArtifactFilterAndSort != nil {
		artifactInfoReceived, listInfo, err = a.Managers.AM.ListArtifact(
			*params.ArtifactFilterAndSort.Filter, *params.ArtifactFilterAndSort.Sorter, paginator)
	} else {
		artifactInfoReceived, listInfo, err = a.Managers.AM.ListArtifact(
			weles.ArtifactFilter{}, weles.ArtifactSorter{}, paginator)
	}

	// TODO: remove this when artifactmanager will return this.
	if len(artifactInfoReceived) == 0 {
		return artifacts.NewArtifactListerNotFound().WithPayload(
			&weles.ErrResponse{Message: weles.ErrArtifactNotFound.Error()})
	}

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

func responderArtifact206(
	listInfo weles.ListInfo,
	paginator weles.ArtifactPagination,
	artifactInfoReturned []*weles.ArtifactInfo,
	defaultPageLimit int32) (responder *artifacts.ArtifactListerPartialContent) {
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

func responderArtifact200(
	listInfo weles.ListInfo,
	paginator weles.ArtifactPagination,
	artifactInfoReturned []*weles.ArtifactInfo,
	defaultPageLimit int32) (responder *artifacts.ArtifactListerOK) {

	var artifactListerURL artifacts.ArtifactListerURL
	responder = artifacts.NewArtifactListerOK()
	responder.SetTotalRecords(listInfo.TotalRecords)
	if paginator.ID != 0 { //not the first page
		// keep in mind that ArtifactPath in paginator is taken from query parameter,
		// not ArtifactManager
		if paginator.Forward == true {
			tmp := artifactInfoReturned[0].ID
			artifactListerURL.Before = &tmp
			if defaultPageLimit != paginator.Limit {
				tmp := paginator.Limit
				artifactListerURL.Limit = &tmp
			}
			responder.SetPrevious(artifactListerURL.String())
		}
		if paginator.Forward == false {
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

func setArtifactPaginator(
	params artifacts.ArtifactListerParams,
	defaultPageLimit int32) (paginator weles.ArtifactPagination) {
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

// artifactInfoReceivedToReturn does the same thing as jobInfoReceivedToReturn.
// TODO:make ArtifactInfos and JobInfos types implement interface with a function that will return
// slice of pointers. Will probably need to use reflect which I'm not familiar with thus not done now.
func artifactInfoReceivedToReturn(artifactInfoReceived []weles.ArtifactInfo) []*weles.ArtifactInfo {
	artifactInfoReturned := make([]*weles.ArtifactInfo, len(artifactInfoReceived))
	for i := range artifactInfoReceived {
		artifactInfoReturned[i] = &artifactInfoReceived[i]
	}
	return artifactInfoReturned

}
