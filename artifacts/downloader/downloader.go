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

// Package downloader is responsible for Weles system's job artifact downloading.
package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"

	. "git.tizen.org/tools/weles"
)

// Downloader implements ArtifactDownloader interface.
type Downloader struct {
	notification chan ArtifactStatusChange // can be used to monitor ArtifactStatusChanges.
	queue        chan downloadJob
}

// downloadJob provides necessary info for download to be done.
type downloadJob struct {
}

// queueCap is the default length of download queue.
const queueCap = 100

// newDownloader returns initilized Downloader.
func newDownloader(notification chan ArtifactStatusChange, workerCount int, queueSize int) *Downloader {

	return &Downloader{
		notification: notification,
		queue:        make(chan downloadJob, queueSize),
	}
}

// NewDownloader returns Downloader initialized  with default queue length
func NewDownloader(notification chan ArtifactStatusChange, workerCount int) *Downloader {
	return newDownloader(notification, workerCount, queueCap)
}

// Close is part of implementation of ArtifactDownloader interface.
// It closes used channels.
func (d *Downloader) Close() {
	close(d.queue)
}

// getData downloads file from provided location and saves it in a prepared path.
func (d *Downloader) getData(URI ArtifactURI, path ArtifactPath) error {

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
func (d *Downloader) download(URI ArtifactURI, path ArtifactPath, ch chan ArtifactStatusChange) {
	if path == "" {
		return
	}

	change := ArtifactStatusChange{
		Path:      path,
		NewStatus: AM_DOWNLOADING,
	}
	channels := []chan ArtifactStatusChange{ch, d.notification}
	notify(change, channels)

	err := d.getData(URI, path)
	if err != nil {
		os.Remove(string(path))
		change.NewStatus = AM_FAILED
	} else {
		change.NewStatus = AM_READY
	}

	notify(change, channels)
}

// Download is part of implementation of ArtifactDownloader interface.
// TODO implement.
func (d *Downloader) Download(URI ArtifactURI, path ArtifactPath, ch chan ArtifactStatusChange) error {
	return ErrNotImplemented

}

// CheckInCache is part of implementation of ArtifactDownloader interface.
// TODO implement.
func (d *Downloader) CheckInCache(URI ArtifactURI) (ArtifactInfo, error) {
	return ArtifactInfo{}, ErrNotImplemented
}

// notify sends ArtifactStatusChange to all specified channels.
func notify(change ArtifactStatusChange, channels []chan ArtifactStatusChange) {
	for _, ch := range channels {
		ch <- change
	}
}
