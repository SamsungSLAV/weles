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

package server

import (
	"fmt"
	"os"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/server/operations/artifacts"
	middleware "github.com/go-openapi/runtime/middleware"
)

// ArtifactDownloader is a handler which passess JobID to JobManager to cancel a job.
func (m *Managers) ArtifactDownloader(params artifacts.ArtifactDownloaderParams) middleware.Responder {
	ai, li, err := m.AM.ListArtifact(
		weles.ArtifactFilter{ID: []int64{params.ArtifactID}},
		weles.ArtifactSorter{SortBy: weles.ArtifactSortByID, SortOrder: weles.SortOrderAscending},
		weles.ArtifactPagination{ID: 0, Forward: true})

	if err != nil {
		switch err {
		case weles.ErrArtifactNotFound:
			return artifacts.NewArtifactDownloaderNotFound().WithPayload(&weles.ErrResponse{
				Message: err.Error()})
		default:
			return artifacts.NewArtifactDownloaderInternalServerError().WithPayload(&weles.ErrResponse{
				Message: err.Error()})
		}
	}
	if li.TotalRecords != 1 {
		return artifacts.NewArtifactDownloaderInternalServerError().WithPayload(&weles.ErrResponse{
			Message: fmt.Sprintf("expected to receive one artifact, received %d", li.TotalRecords)})
	}
	artifact, err := os.Open(string(ai[0].Path)) // close will be handled by swagger
	if err != nil {
		return artifacts.NewArtifactDownloaderInternalServerError().WithPayload(&weles.ErrResponse{
			Message: err.Error()})
	}
	return artifacts.NewArtifactDownloaderOK().WithPayload(artifact)
}
