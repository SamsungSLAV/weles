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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	. "git.tizen.org/tools/weles"
	. "git.tizen.org/tools/weles/artifacts/database"
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
	db  ArtifactDB
	dir string
}

const (
	// defaultDb is default ArtifactDB name.
	defaultDb = "weles.db"
	// defaultDir is default directory for ArtifactManager storage.
	defaultDir = "/tmp/weles/"
)

func newArtifactManager(db, dir string) (ArtifactManager, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	am := Storage{dir: dir}
	err = am.db.Open(db)

	if err != nil {
		return nil, err
	}
	return &am, nil
}

// NewArtifactManager returns initialized Storage implementing ArtifactManager interface.
// If db or dir is empy, default value will be used.
func NewArtifactManager(db, dir string) (ArtifactManager, error) {
	if db == "" {
		db = defaultDb
	}
	if dir == "" {
		dir = defaultDir
	}
	return newArtifactManager(filepath.Join(dir, db), dir)
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
	path, err := s.getNewPath(artifact)
	if err != nil {
		return "", err
	}

	err = s.db.InsertArtifactInfo(&ArtifactInfo{artifact, path, "", time.Now().UTC()})
	if err != nil {
		return "", err
	}
	return path, nil
}

// GetArtifactInfo is part of implementation of ArtifactManager interface.
func (s *Storage) GetArtifactInfo(path ArtifactPath) (ArtifactInfo, error) {
	return s.db.SelectPath(path)
}

// Close closes Storage's ArtifactDB.
func (s *Storage) Close() error {
	return s.db.Close()
}

// getNewPath prepares new path for artifact.
func (s *Storage) getNewPath(ad ArtifactDescription) (ArtifactPath, error) {
	var (
		jobDir  = filepath.Join(s.dir, strconv.FormatUint(uint64(ad.JobID), 10))
		typeDir = filepath.Join(jobDir, string(ad.Type))
		err     error
	)

	// Organize by filetypes
	err = os.MkdirAll(typeDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Add human readable prefix
	f, err := ioutil.TempFile(typeDir, string(ad.Alias))
	if err != nil {
		return "", err
	}
	defer f.Close()
	return ArtifactPath(f.Name()), err
}
