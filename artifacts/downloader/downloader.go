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

// Package downloader is responsible for Weles system's job artifact downloading.
package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"git.tizen.org/tools/weles"
)

// Downloader implements ArtifactDownloader interface.
type Downloader struct {
	notification chan weles.ArtifactStatusChange // can be used to monitor ArtifactStatusChanges.
	queue        chan downloadJob
	wg           sync.WaitGroup
}

// downloadJob provides necessary info for download to be done.
type downloadJob struct {
	path weles.ArtifactPath
	uri  weles.ArtifactURI
	ch   chan weles.ArtifactStatusChange
}

// newDownloader returns initilized Downloader.
func newDownloader(notification chan weles.ArtifactStatusChange, workers int, queueSize int) *Downloader {

	d := &Downloader{
		notification: notification,
		queue:        make(chan downloadJob, queueSize),
	}

	// Start all workers.
	d.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go d.work()
	}
	return d
}

// NewDownloader returns Downloader initialized  with default queue length
func NewDownloader(notification chan weles.ArtifactStatusChange, workerCount, queueCap int) *Downloader {
	return newDownloader(notification, workerCount, queueCap)
}

// Close is part of implementation of ArtifactDownloader interface.
// It waits for running download jobs to stop and closes used channels.
func (d *Downloader) Close() {
	close(d.queue)
	d.wg.Wait()
}

// getData downloads file from provided location and saves it in a prepared path.
func (d *Downloader) getData(URI weles.ArtifactURI, path weles.ArtifactPath) error {
	resp, err := http.Get(string(URI))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error %v %v", URI, resp.Status)
	}

	file, err := os.Create(string(path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// download downloads artifact from provided URI and saves it to specified path.
// It sends notification about status changes to two channels - Downloader's notification
// channel, and other one, that can be specified passed as an argument.
func (d *Downloader) download(URI weles.ArtifactURI, path weles.ArtifactPath, ch chan weles.ArtifactStatusChange) {
	if path == "" {
		return
	}

	change := weles.ArtifactStatusChange{
		Path:      path,
		NewStatus: weles.ArtifactStatusDOWNLOADING,
	}
	channels := []chan weles.ArtifactStatusChange{ch, d.notification}
	notify(change, channels)

	err := d.getData(URI, path)
	if err != nil {
		os.Remove(string(path))
		change.NewStatus = weles.ArtifactStatusFAILED
	} else {
		change.NewStatus = weles.ArtifactStatusREADY
	}
	notify(change, channels)
}

// Download is part of implementation of ArtifactDownloader interface.
// It puts new downloadJob on the queue.
func (d *Downloader) Download(URI weles.ArtifactURI, path weles.ArtifactPath, ch chan weles.ArtifactStatusChange) error {
	channels := []chan weles.ArtifactStatusChange{ch, d.notification}
	notify(weles.ArtifactStatusChange{Path: path, NewStatus: weles.ArtifactStatusPENDING}, channels)

	job := downloadJob{
		path: path,
		uri:  URI,
		ch:   ch,
	}

	select {
	case d.queue <- job:
	default:
		return ErrQueueFull
	}
	return nil
}

func (d *Downloader) work() {
	defer d.wg.Done()
	for job := range d.queue {
		d.download(job.uri, job.path, job.ch)
	}
}

// CheckInCache is part of implementation of ArtifactDownloader interface.
// TODO implement.
func (d *Downloader) CheckInCache(URI weles.ArtifactURI) (weles.ArtifactInfo, error) {
	return weles.ArtifactInfo{}, weles.ErrNotImplemented
}

// notify sends ArtifactStatusChange to all specified channels.
func notify(change weles.ArtifactStatusChange, channels []chan weles.ArtifactStatusChange) {
	for _, ch := range channels {
		ch <- change
	}
}
