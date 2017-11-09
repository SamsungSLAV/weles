/*
 *  Copyright (c) 2017 Samsung Electronics Co., Ltd All Rights Reserved
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

// Package artifacts is responsible for Weles system's job artifact management.
package artifacts

import (
	. "git.tizen.org/tools/weles"
)

// ArtifactDownloader downloads requested file if there is need to.
type ArtifactDownloader interface {
	// Download starts downloading requested artifact.
	Download(URI ArtifactURI, path ArtifactPath, ch chan ArtifactStatusChange) error

	// CheckInCache checks if file already exists in ArtifactDB.
	CheckInCache(URI ArtifactURI) (ArtifactInfo, error)

	// Close waits for all jobs to finish, and gracefully closes ArtifactDownloader.
	Close()
}

// Storage should be used by Weles' subsystems that need access to ArtifactDB
// or information about artifacts stored there.
// Storage implements ArtifactManager interface.
type Storage struct {
	ArtifactManager
}

// ListArtifact is part of implementation of ArtifactManager interface.
func (s *Storage) ListArtifact(filter ArtifactFilter) ([]ArtifactInfo, error) {
	return nil, ErrNotImplemented
}

// PushArtifact is part of implementation of ArtifactManager interface.
func (s *Storage) PushArtifact(artifact ArtifactDescription, ch chan ArtifactStatusChange) (ArtifactPath, error) {
	return "", ErrNotImplemented
}

// CreateArtifact is part of implementation of ArtifactManager interface.
func (s *Storage) CreateArtifact(artifact ArtifactDescription) (ArtifactPath, error) {
	return "", ErrNotImplemented
}

// GetArtifactInfo is part of implementation of ArtifactManager interface.
func (s *Storage) GetArtifactInfo(path ArtifactPath) (ArtifactInfo, error) {
	return ArtifactInfo{}, ErrNotImplemented
}
