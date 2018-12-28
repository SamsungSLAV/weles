/*
 *  Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

// File artifactmanager.go defines ArtifactManager interface.

package weles

// ArtifactManager provides access to content in ArtifactDB required for Job execution.
// It provides data from ArtifactDB for lookup and retrieval.
// It is responsible for downloading job artifacts to ArtifactDB.
type ArtifactManager interface {
	// List filters ArtifactDB and returns list of all matching artifacts.
	ListArtifact(filter ArtifactFilter, sorter ArtifactSorter, paginator ArtifactPaginator,
	) ([]ArtifactInfo, ListInfo, error)

	// Push inserts artifact to ArtifactDB and returns its path.
	PushArtifact(artifact ArtifactDescription, ch chan ArtifactStatusChange) (ArtifactPath, error)

	// Create constructs ArtifactPath in ArtifactDB, but no file is created.
	CreateArtifact(artifact ArtifactDescription) (ArtifactPath, error)

	// GetFileInfo retrieves information about an artifact from ArtifactDB.
	GetArtifactInfo(path ArtifactPath) (ArtifactInfo, error)

	// Close gracefully closes ArtifactManager.
	Close() error
}
