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

// File artifactmanager.go defines ArtifactManager interface and structures related to it.

package weles

import "time"

// ArtifactType denotes type and function of an artifact.
type ArtifactType string

const (
	// AM_IMAGEFILE - image file.
	AM_IMAGEFILE ArtifactType = "IMAGE"
	// AM_RESULTFILE - all outputs, files built during tests, etc.
	AM_RESULTFILE ArtifactType = "RESULT"
	// AM_TESTFILE - additional files uploaded by user for conducting test.
	AM_TESTFILE ArtifactType = "TESTFILE"
	// AM_YAMLFILE - yaml file describing Weles Job.
	AM_YAMLFILE ArtifactType = "YAMLFILE"
)

// ArtifactPath describes path to artifact in ArtifactDB filesystem.
type ArtifactPath string

// ArtifactStatus describes artifact status and availability.
type ArtifactStatus string

const (
	// AM_DOWNLOADING - artifact is currently being downloaded.
	AM_DOWNLOADING ArtifactStatus = "DOWNLOADING"
	// AM_READY - artifact has been downloaded and is ready to use.
	AM_READY ArtifactStatus = "READY"
	// AM_FAILED - file is not available for use (e.g. download failed).
	AM_FAILED ArtifactStatus = "FAILED"
	// AM_PENDING - artifact download has not started yet.
	AM_PENDING ArtifactStatus = "PENDING"
)

// ArtifactURI is used to identify artifact's source.
type ArtifactURI string

// ArtifactAlias is used to identify artifact's alias.
type ArtifactAlias string

// ArtifactDescription contains information needed to create new artifact in ArtifactDB.
type ArtifactDescription struct {
	JobID JobID
	Type  ArtifactType
	Alias ArtifactAlias
	URI   ArtifactURI
}

// ArtifactInfo describes single artifact stored in ArtifactDB.
type ArtifactInfo struct {
	ArtifactDescription
	Path      ArtifactPath
	Status    ArtifactStatus
	Timestamp time.Time
}

// ArtifactFilter is used to filter results from ArtifactDB.
type ArtifactFilter struct {
	JobID  []JobID
	Type   []ArtifactType
	Status []ArtifactStatus
	Alias  []ArtifactAlias
}

// ArtifactStatusChange contains information about new status of an artifact.
// It is used to monitor status changes.
type ArtifactStatusChange struct {
	Path      ArtifactPath
	NewStatus ArtifactStatus
}

// ArtifactManager provides access to content in ArtifactDB required for Job execution.
// It provides data from ArtifactDB for lookup and retrieval.
// It is responsible for downloading job artifacts to ArtifactDB.
type ArtifactManager interface {
	// List filters ArtifactDB and returns list of all matching artifacts.
	ListArtifact(filter ArtifactFilter) ([]ArtifactInfo, error)

	// Push inserts artifact to ArtifactDB and returns its path.
	PushArtifact(artifact ArtifactDescription, ch chan ArtifactStatusChange) (ArtifactPath, error)

	// Create constructs ArtifactPath in ArtifactDB, but no file is created.
	CreateArtifact(artifact ArtifactDescription) (ArtifactPath, error)

	// GetFileInfo retrieves information about an artifact from ArtifactDB.
	GetArtifactInfo(path ArtifactPath) (ArtifactInfo, error)
}
