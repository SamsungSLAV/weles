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

// File controller/downloaderimpl.go contains Downloader implementation.

package controller

import (
	"fmt"
	"sync"

	"github.com/SamsungSLAV/slav/logger"
	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/controller/notifier"
)

const (
	formatJobStatus = "Internal Weles error while changing Job status : %s"
	formatJobConfig = "Internal Weles error while getting Job config : %s"
	formatURI       = "Internal Weles error while registering URI:<%s> in ArtifactManager : %s"
	formatPath      = "Internal Weles error while creating a new path in ArtifactManager : %s"
	formatConfig    = "Internal Weles error while setting config : %s"
	formatDownload  = "Failed to download some artifacts for the Job"
	formatReady     = "%d / %d artifacts ready"
)

// jobArtifactsInfo contains information about progress of downloading
// artifacts required by a single Job.
type jobArtifactsInfo struct {
	paths       int
	ready       int
	failed      int
	configSaved bool
}

// DownloaderImpl implements delegating downloading of artifacts required
// by Jobs to ArtifactsManager, monitors progress and notifies
// Controller, when all files are ready.
type DownloaderImpl struct {
	// Notifier provides channel for communication with Controller.
	notifier.Notifier
	// jobs references module implementing Jobs management.
	jobs JobsController
	// artifacts references Weles module implementing ArtifactManager for
	// managing ArtifactsDB.
	artifacts weles.ArtifactManager
	// collector gathers artifact status changes from ArtifactManager
	collector chan weles.ArtifactStatusChange

	// path2Job identifies Job related to the artifact path.
	path2Job map[string]weles.JobID
	// info contains information about progress of downloading artifacts
	// for Jobs.
	info map[weles.JobID]*jobArtifactsInfo
	//mutex protects access to path2Job and info maps.
	mutex *sync.Mutex
}

// NewDownloader creates a new DownloaderImpl structure setting up references
// to used Weles modules.
func NewDownloader(j JobsController, a weles.ArtifactManager) Downloader {
	ret := &DownloaderImpl{
		Notifier:  notifier.NewNotifier(),
		jobs:      j,
		artifacts: a,
		collector: make(chan weles.ArtifactStatusChange),
		path2Job:  make(map[string]weles.JobID),
		info:      make(map[weles.JobID]*jobArtifactsInfo),
		mutex:     new(sync.Mutex),
	}
	go ret.loop()
	return ret
}

// pathStatusChange reacts on notification from ArtifactManager and updates
// path and job structures.
func (h *DownloaderImpl) pathStatusChange(path string, status weles.ArtifactStatus,
) (changed bool, j weles.JobID, info string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	j, ok := h.path2Job[path]
	if !ok {
		logger.WithProperty("path", path).Error("Failed to match path with JobID.")
		return
	}
	i, ok := h.info[j]
	if !ok {
		logger.WithProperty("JobID", j).Error("Failed to match ArtifactInfo with JobID.")
		delete(h.path2Job, path)
		return
	}
	switch status {
	case weles.ArtifactStatusREADY:
		i.ready++
		info = fmt.Sprintf(formatReady, i.ready, i.paths)
	case weles.ArtifactStatusFAILED:
		i.failed++
		info = "Failed to download artifact"
	default:
		return
	}
	delete(h.path2Job, path)
	changed = true
	return
}

// removePath removes mapping from the path to related Job.
func (h *DownloaderImpl) removePath(path string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.path2Job, path)
}

// loop handles all notifications from ArtifactManager, updates
// jobArtifactsInfo and send notification to Controller, when the answer
// is ready.
//
// It is run in a separate goroutine - one for all jobs.
func (h *DownloaderImpl) loop() {
	for {
		change, open := <-h.collector
		if !open {
			return
		}
		update, j, info := h.pathStatusChange(string(change.Path), change.NewStatus)
		if !update {
			continue
		}

		err := h.jobs.SetStatusAndInfo(j, weles.JobStatusDOWNLOADING, info)
		if err != nil {
			logger.WithProperty("JobID", j).Error("Failed to set job status to DOWNLOADING.")
			h.removePath(string(change.Path))
			h.fail(j, fmt.Sprintf(formatJobStatus, err.Error()))
		}
		h.sendIfReady(j)
	}
}

// fail responses failure to Controller.
func (h *DownloaderImpl) fail(j weles.JobID, msg string) {
	if h.removeJobInfo(j) == nil {
		h.SendFail(j, msg)
	}
}

// succeed responses success to Controller.
func (h *DownloaderImpl) succeed(j weles.JobID) {
	if h.removeJobInfo(j) == nil {
		h.SendOK(j)
	}
}

// initializeJobInfo creates a jobArtifactInfo structure.
func (h *DownloaderImpl) initializeJobInfo(j weles.JobID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	_, ok := h.info[j]
	if !ok {
		h.info[j] = new(jobArtifactsInfo)
	}
}

// removeJobInfo removes a jobArtifactInfo structure.
func (h *DownloaderImpl) removeJobInfo(j weles.JobID) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, ok := h.info[j]
	if !ok {
		logger.WithProperty("JobID", j).Error("Failed to match JobInfo with JobID.")
		return weles.ErrJobNotFound
	}
	delete(h.info, j)
	return nil
}

// push delegates downloading single uri to ArtifactDB.
func (h *DownloaderImpl) push(j weles.JobID, t weles.ArtifactType, alias, uri string,
) (string, error) {
	ad := weles.ArtifactDescription{
		JobID: j,
		Type:  t,
		Alias: weles.ArtifactAlias(alias),
		URI:   weles.ArtifactURI(uri),
	}
	p, err := h.artifacts.PushArtifact(ad, h.collector)
	if err != nil {
		logger.WithError(err).WithProperties(logger.Properties{"JobID": j, "URI": uri}).
			Error("Failed to push artifact to db")
		return "", err
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	i, ok := h.info[j]
	if !ok {
		logger.WithProperty("JobID", j).Errorf("Failed to match jobsArtifactInfo with JobID.")
		return "", weles.ErrJobNotFound
	}
	i.paths++
	h.path2Job[string(p)] = j

	return string(p), nil
}

// pullCreate creates a new path for pull artifact.
func (h *DownloaderImpl) pullCreate(j weles.JobID, alias string) (string, error) {
	p, err := h.artifacts.CreateArtifact(weles.ArtifactDescription{
		JobID: j,
		Type:  weles.ArtifactTypeTEST,
		Alias: weles.ArtifactAlias(alias),
	})
	return string(p), err
}

// configSaved updates info structure.
func (h *DownloaderImpl) configSaved(j weles.JobID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	i, ok := h.info[j]
	if !ok {
		logger.WithProperty("JobID", j).Errorf("Failed to match jobsArtifactInfo with JobID.")
		return
	}

	i.configSaved = true
}

// verify if an answer to the Controller should be send.
func (h *DownloaderImpl) verify(j weles.JobID) (success, send bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	i, ok := h.info[j]
	if !ok { // Job is not monitored (maybe it has been responded already).
		return false, false
	}
	if i.failed > 0 { // Some artifacts fail to be downloaded.
		return false, true
	}
	if !i.configSaved { // Config is not yet fully analyzed and saved.
		return false, false
	}
	if i.ready == i.paths { // All artifacts are ready.
		return true, true
	}
	return false, false
}

// sendIfReady sends an answer to the Controller if it is ready.
func (h *DownloaderImpl) sendIfReady(j weles.JobID) {
	success, send := h.verify(j)

	if !send {
		return
	}

	if success {
		h.succeed(j)
	} else {
		h.fail(j, formatDownload)
	}
}

// DispatchDownloads parses Job's config and delegates to ArtifactManager downloading
// of all images and files to be pushed during Job execution. It also creates
// ArtifactDB paths for files that will be pulled from Dryad.
func (h *DownloaderImpl) DispatchDownloads(j weles.JobID) {
	h.initializeJobInfo(j)

	err := h.jobs.SetStatusAndInfo(j, weles.JobStatusDOWNLOADING, "")
	if err != nil {
		logger.WithError(err).WithProperty("JobID", j).
			Error("Failed to set jobs' status to DOWNLOADING")
		h.fail(j, fmt.Sprintf(formatJobStatus, err.Error()))
		return
	}

	config, err := h.jobs.GetConfig(j)
	if err != nil {
		logger.WithError(err).WithProperty("JobID", j).Error("Failed to get jobs' config.")
		h.fail(j, fmt.Sprintf(formatJobConfig, err.Error()))
		return
	}

	for i, image := range config.Action.Deploy.Images {
		if image.URI != "" {
			var path string
			path, err = h.push(j, weles.ArtifactTypeIMAGE, fmt.Sprintf("Image_%d", i), image.URI)
			if err != nil {
				logger.WithError(err).
					WithProperties(logger.Properties{"URI": image.URI, "JobID": j}).
					Error("Failed to create path for IMAGE artifact.")
				h.fail(j, fmt.Sprintf(formatURI, image.URI, err.Error()))
				return
			}
			config.Action.Deploy.Images[i].Path = path
		}
		if image.ChecksumURI != "" {
			var path string
			path, err = h.push(j, weles.ArtifactTypeIMAGE, fmt.Sprintf("ImageMD5_%d", i),
				image.ChecksumURI)
			if err != nil {
				logger.WithError(err).
					WithProperties(logger.Properties{"URI": image.ChecksumURI, "JobID": j}).
					Errorf("Failed to create path for IMAGEMD5 artifact.")
				h.fail(j, fmt.Sprintf(formatURI, image.ChecksumURI, err.Error()))
				return
			}
			config.Action.Deploy.Images[i].ChecksumPath = path
		}
	}
	var path string
	for i, tc := range config.Action.Test.TestCases {
		for k, ta := range tc.TestActions {
			switch ta.(type) {
			case weles.Push:
				action := ta.(weles.Push)
				path, err = h.push(j, weles.ArtifactTypeTEST, action.Alias, action.URI)
				if err != nil {
					logger.WithError(err).
						WithProperties(logger.Properties{"URI": action.URI, "JobID": j}).
						Error("Failed to create path for push/TEST artifact.")
					h.fail(j, fmt.Sprintf(formatURI, action.URI, err.Error()))
					return
				}
				action.Path = path
				config.Action.Test.TestCases[i].TestActions[k] = action
			case weles.Pull:
				action := ta.(weles.Pull)
				path, err = h.pullCreate(j, action.Alias)
				if err != nil {
					logger.Error("Failed to create new path for pull/TEST artifact")
					h.fail(j, fmt.Sprintf(formatPath, err.Error()))
					return
				}
				action.Path = path
				config.Action.Test.TestCases[i].TestActions[k] = action
			}
		}
	}

	err = h.jobs.SetConfig(j, config)
	if err != nil {
		logger.WithError(err).WithProperty("JobID", j).Errorf("Failed to set jobs' config")
		h.fail(j, fmt.Sprintf(formatConfig, err.Error()))
		return
	}

	h.configSaved(j)
	h.sendIfReady(j)
}
