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

// Package artifacts is responsible for Weles system's job artifact management.
package artifacts

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/artifacts/database"
	"github.com/SamsungSLAV/weles/artifacts/downloader"
)

// ArtifactDownloader downloads requested file if there is need to.
type ArtifactDownloader interface {
	// Download starts downloading requested artifact.
	Download(URI weles.ArtifactURI, path weles.ArtifactPath, ch chan weles.ArtifactStatusChange,
	) error

	// CheckInCache checks if file already exists in ArtifactDB.
	CheckInCache(URI weles.ArtifactURI) (weles.ArtifactInfo, error)

	// Close waits for all jobs to finish, and gracefully closes ArtifactDownloader.
	Close()
}

// Storage should be used by Weles' subsystems that need access to ArtifactDB
// or information about artifacts stored there.
// Storage implements ArtifactManager interface.
type Storage struct {
	weles.ArtifactManager
	db         database.ArtifactDB
	dir        string
	downloader ArtifactDownloader
	notifier   chan weles.ArtifactStatusChange
}

func newArtifactManager(db, dir string, notifierCap, workersCount, queueCap int,
) (weles.ArtifactManager, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	notifier := make(chan weles.ArtifactStatusChange, notifierCap)

	am := Storage{
		dir:        dir,
		downloader: downloader.NewDownloader(notifier, workersCount, queueCap),
		notifier:   notifier,
	}
	err = am.db.Open(db)
	if err != nil {
		return nil, err
	}

	go am.listenToChanges()

	return &am, nil
}

// NewArtifactManager returns initialized Storage implementing ArtifactManager interface.
// If db or dir is empy, default value will be used.
func NewArtifactManager(db, dir string, notifierCap, workersCount, queueCap int,
) (weles.ArtifactManager, error) {
	return newArtifactManager(filepath.Join(dir, db), dir, notifierCap, workersCount, queueCap)
}

// ListArtifact is part of implementation of ArtifactManager interface.
func (s *Storage) ListArtifact(filter weles.ArtifactFilter, sorter weles.ArtifactSorter,
	paginator weles.ArtifactPagination) ([]weles.ArtifactInfo, weles.ListInfo, error) {

	return s.db.Filter(filter, sorter, paginator)
}

// PushArtifact is part of implementation of ArtifactManager interface.
func (s *Storage) PushArtifact(artifact weles.ArtifactDescription,
	ch chan weles.ArtifactStatusChange) (weles.ArtifactPath, error) {

	path, err := s.CreateArtifact(artifact)
	if err != nil {
		return "", err
	}

	err = s.downloader.Download(artifact.URI, path, ch)
	if err != nil {
		err2 := s.db.SetStatus(weles.ArtifactStatusChange{
			Path:      path,
			NewStatus: weles.ArtifactStatusFAILED,
		})
		if err2 != nil {
			return "", errors.New(
				"failed to download artifact: " + err.Error() +
					" and failed to set artifacts status to failed: " + err2.Error())
		}
		return "", errors.New("failed to download artifact: " + err.Error())
	}
	return path, nil
}

// CreateArtifact is part of implementation of ArtifactManager interface.
func (s *Storage) CreateArtifact(artifact weles.ArtifactDescription) (weles.ArtifactPath, error) {
	path, err := s.getNewPath(artifact)
	if err != nil {
		return "", err
	}

	err = s.db.InsertArtifactInfo(&weles.ArtifactInfo{
		ArtifactDescription: artifact,
		Path:                path,
		Status:              "",
		Timestamp:           strfmt.DateTime(time.Now().UTC()),
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// GetArtifactInfo is part of implementation of ArtifactManager interface.
func (s *Storage) GetArtifactInfo(path weles.ArtifactPath) (weles.ArtifactInfo, error) {
	return s.db.SelectPath(path)
}

// Close closes Storage's ArtifactDB.
func (s *Storage) Close() error {
	s.downloader.Close()
	close(s.notifier)
	return s.db.Close()
}

// getNewPath prepares new path for artifact.
func (s *Storage) getNewPath(ad weles.ArtifactDescription) (weles.ArtifactPath, error) {
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

	defer func() {
		if err = f.Close(); err != nil {
			log.Println("failed to close file")
			//TODO: aalexanderr log
		}
	}()
	return weles.ArtifactPath(f.Name()), err
}

// listenToChanges updates artifact's status in db every time Storage is notified
// about status change.
func (s *Storage) listenToChanges() {
	for change := range s.notifier {
		// Error handled in SetStatus function.
		_ = s.db.SetStatus(change) //nolint: gas, gosec
	}
}
