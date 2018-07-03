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

package downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"git.tizen.org/tools/weles"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Downloader", func() {

	const pigs = `The pig, if I am not mistaken,
Supplies us sausage, ham, and bacon.
Let others say his heart is big --
I call it stupid of the pig.

-Ogden Nash`

	var platinumKoala *Downloader

	var (
		tmpDir     string
		validDir   string
		invalidDir string
		validURL   weles.ArtifactURI = "validURL"
		invalidURL weles.ArtifactURI = "invalidURL"
		ts         *httptest.Server
		ch         chan weles.ArtifactStatusChange
	)

	var (
		notifyCap    int = 100 // notitication channel capacity.
		notification chan weles.ArtifactStatusChange
		workersCount = 8
	)

	checkChannels := func(ch1, ch2 chan weles.ArtifactStatusChange, change weles.ArtifactStatusChange) {
		Eventually(ch1).Should(Receive(Equal(change)))
		Eventually(ch2).Should(Receive(Equal(change)))
	}

	BeforeEach(func() {

		var err error
		// prepare Downloader.
		notification = make(chan weles.ArtifactStatusChange, notifyCap)
		platinumKoala = NewDownloader(notification, workersCount)

		// prepare temporary directories.
		tmpDir, err = ioutil.TempDir("", "weles-")
		Expect(err).ToNot(HaveOccurred())
		validDir = filepath.Join(tmpDir, "valid")
		err = os.MkdirAll(validDir, os.ModePerm)
		Expect(err).ToNot(HaveOccurred())
		// directory is not created therefore path will be invalid.
		invalidDir = filepath.Join(tmpDir, "invalid")

		ch = make(chan weles.ArtifactStatusChange, 5)
	})

	AfterEach(func() {
		platinumKoala.Close()
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())

	})

	prepareServer := func(url weles.ArtifactURI) *httptest.Server {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if url == validURL {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, pigs)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		return ts
	}

	DescribeTable("getData(): Notify channels and save data to file",
		func(url weles.ArtifactURI, valid bool, finalResult weles.ArtifactStatus) {
			ts = prepareServer(url)
			defer ts.Close()

			dir := validDir
			if !valid {
				dir = invalidDir
			}
			filename := weles.ArtifactPath(filepath.Join(dir, "test"))

			err := platinumKoala.getData(weles.ArtifactURI(ts.URL), weles.ArtifactPath(filename))

			if valid && url != invalidURL {
				Expect(err).ToNot(HaveOccurred())
				content, err := ioutil.ReadFile(string(filename))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal(pigs))
			} else {
				Expect(string(filename)).NotTo(BeAnExistingFile())
				_, err := ioutil.ReadFile(string(filename))
				Expect(err).To(HaveOccurred())
			}

		},
		Entry("download valid file to valid path", validURL, true, weles.AM_READY),
		Entry("fail when url is invalid", invalidURL, true, weles.AM_FAILED),
		Entry("fail when path is invalid", validURL, false, weles.AM_FAILED),
		Entry("fail when url and path are invalid", invalidURL, false, weles.AM_FAILED),
	)

	DescribeTable("download(): Notify channels and save data to file",
		func(url weles.ArtifactURI, valid bool, finalResult weles.ArtifactStatus) {
			ts = prepareServer(url)
			defer ts.Close()

			dir := validDir
			if !valid {
				dir = invalidDir
			}
			filename := weles.ArtifactPath(filepath.Join(dir, "test"))

			status := weles.ArtifactStatusChange{filename, weles.AM_DOWNLOADING}

			platinumKoala.download(weles.ArtifactURI(ts.URL), weles.ArtifactPath(filename), ch)

			status.NewStatus = weles.AM_DOWNLOADING
			checkChannels(ch, platinumKoala.notification, status)

			status.NewStatus = finalResult
			checkChannels(ch, platinumKoala.notification, status)

			if valid && url != invalidURL {
				content, err := ioutil.ReadFile(string(filename))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal(pigs))
			} else {
				Expect(string(filename)).NotTo(BeAnExistingFile())
			}

		},
		Entry("download valid file to valid path", validURL, true, weles.AM_READY),
		Entry("fail when url is invalid", invalidURL, true, weles.AM_FAILED),
		Entry("fail when path is invalid", validURL, false, weles.AM_FAILED),
		Entry("fail when url and path are invalid", invalidURL, false, weles.AM_FAILED),
	)

	DescribeTable("Download(): Notify ch channel about any changes",
		func(url weles.ArtifactURI, valid bool, finalResult weles.ArtifactStatus) {
			ts = prepareServer(url)
			defer ts.Close()

			dir := validDir
			if !valid {
				dir = invalidDir
			}
			path := weles.ArtifactPath(filepath.Join(dir, "animal"))

			err := platinumKoala.Download(weles.ArtifactURI(ts.URL), path, ch)
			Expect(err).ToNot(HaveOccurred())

			status := weles.ArtifactStatusChange{path, weles.AM_PENDING}
			Eventually(ch).Should(Receive(Equal(status)))

			status.NewStatus = weles.AM_DOWNLOADING
			Eventually(ch).Should(Receive(Equal(status)))

			status.NewStatus = finalResult
			Eventually(ch).Should(Receive(Equal(status)))
		},
		Entry("download valid file to valid path", validURL, true, weles.AM_READY),
		Entry("fail when url is invalid", invalidURL, true, weles.AM_FAILED),
		Entry("fail when path is invalid", validURL, false, weles.AM_FAILED),
		Entry("fail when url and path are invalid", invalidURL, false, weles.AM_FAILED),
	)

	DescribeTable("Download(): Download files to specified path.",
		func(url weles.ArtifactURI, filename string, poem string) {
			ts = prepareServer(url)
			defer ts.Close()

			path := weles.ArtifactPath(filepath.Join(validDir, filename))

			err := platinumKoala.Download(weles.ArtifactURI(ts.URL), path, ch)
			Expect(err).ToNot(HaveOccurred())

			Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{path, weles.AM_PENDING})))
			Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{path, weles.AM_DOWNLOADING})))

			if poem != "" {
				Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{path, weles.AM_READY})))
				content, err := ioutil.ReadFile(string(path))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(BeIdenticalTo(poem))
			} else {
				Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{path, weles.AM_FAILED})))
				content, err := ioutil.ReadFile(string(path))
				Expect(err).To(HaveOccurred())
				Expect(content).To(BeNil())

			}
		},
		Entry("download valid file to valid path", validURL, "pigs", pigs),
		Entry("fail when url is invalid", invalidURL, "cows", nil),
	)

	Describe("DownloadJob queue capacity", func() {
		It("should return error if queue if full.", func() {
			ts = prepareServer(validURL)

			notification := make(chan weles.ArtifactStatusChange, notifyCap)
			ironGopher := newDownloader(notification, 0, 0)
			defer ironGopher.Close()

			path := weles.ArtifactPath(filepath.Join(validDir, "file"))

			err := ironGopher.Download(weles.ArtifactURI(ts.URL), path, ch)
			Expect(err).To(Equal(ErrQueueFull))
		})
	})
})
