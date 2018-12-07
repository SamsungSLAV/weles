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

package downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/SamsungSLAV/slav/logger"
	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/testutil"
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
		queueCap     = 100
	)

	var (
		ws        *testutil.WriterString
		log       *logger.Logger = logger.NewLogger()
		stderrLog *logger.Logger = logger.NewLogger()
	)

	stderrLog.AddBackend("default", logger.Backend{
		Filter:     logger.NewFilterPassAll(),
		Serializer: logger.NewSerializerText(),
		Writer:     logger.NewWriterStderr(),
	})

	checkChannels := func(ch1, ch2 chan weles.ArtifactStatusChange,
		change weles.ArtifactStatusChange) {

		Eventually(ch1).Should(Receive(Equal(change)))
		Eventually(ch2).Should(Receive(Equal(change)))
	}

	BeforeEach(func() {

		var err error
		// prepare Downloader.
		notification = make(chan weles.ArtifactStatusChange, notifyCap)
		platinumKoala = NewDownloader(notification, workersCount, queueCap)

		// prepare temporary directories.
		tmpDir, err = ioutil.TempDir("", "weles-")
		Expect(err).ToNot(HaveOccurred())
		validDir = filepath.Join(tmpDir, "valid")
		err = os.MkdirAll(validDir, os.ModePerm)
		Expect(err).ToNot(HaveOccurred())
		// directory is not created therefore path will be invalid.
		invalidDir = filepath.Join(tmpDir, "invalid")

		ch = make(chan weles.ArtifactStatusChange, 5)

		ws = testutil.NewWriterString()
		log.AddBackend("string", logger.Backend{
			Filter:     logger.NewFilterPassAll(),
			Serializer: logger.NewSerializerText(),
			Writer:     ws,
		})
		logger.SetDefault(log)
	})

	AfterEach(func() {
		platinumKoala.Close()
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())

		logger.SetDefault(stderrLog)
	})

	prepareServer := func(url weles.ArtifactURI) *httptest.Server {
		testServer := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if url == validURL {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, pigs)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
		return testServer
	}

	Describe("getData(): Notify channels and save data to file", func() {
		It("should download valid file to valid path", func() {
			ts = prepareServer(validURL)
			defer ts.Close()

			dir := validDir
			filename := weles.ArtifactPath(filepath.Join(dir, "test"))

			err := platinumKoala.getData(weles.ArtifactURI(ts.URL), weles.ArtifactPath(filename))
			Expect(err).ToNot(HaveOccurred())

			content, err := ioutil.ReadFile(string(filename))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal(pigs))

			Consistently(func() string {
				return ws.GetString()
			}).Should(BeEmpty())
		})
		It("should fail when path is invalid", func() {
			ts = prepareServer(validURL)
			defer ts.Close()

			dir := invalidDir
			filename := weles.ArtifactPath(filepath.Join(dir, "test"))

			err := platinumKoala.getData(weles.ArtifactURI(ts.URL), weles.ArtifactPath(filename))

			Expect(string(filename)).NotTo(BeAnExistingFile())
			_, err = ioutil.ReadFile(string(filename))
			Expect(err).To(HaveOccurred())

			Eventually(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to create file."))
		})
		DescribeTable("response to invalid url",
			func(valid bool) {
				ts = prepareServer(invalidURL)
				defer ts.Close()

				dir := validDir
				if !valid {
					dir = invalidDir
				}
				filename := weles.ArtifactPath(filepath.Join(dir, "test"))

				err := platinumKoala.getData(weles.ArtifactURI(ts.URL),
					weles.ArtifactPath(filename))

				Expect(string(filename)).NotTo(BeAnExistingFile())
				_, err = ioutil.ReadFile(string(filename))
				Expect(err).To(HaveOccurred())

				Eventually(func() string {
					return ws.GetString()
				}).Should(ContainSubstring(
					"Received wrong response from server after downloading artifact."))
			},
			Entry("fail when url is invalid", true),
			Entry("fail when url and path are invalid", false),
		)
	})

	Describe("download(): Notify channels and save data to file", func() {
		It("should download valid file to valid path", func() {
			ts = prepareServer(validURL)
			defer ts.Close()

			dir := validDir
			filename := weles.ArtifactPath(filepath.Join(dir, "test"))

			status := weles.ArtifactStatusChange{
				Path:      filename,
				NewStatus: weles.ArtifactStatusDOWNLOADING,
			}

			platinumKoala.download(weles.ArtifactURI(ts.URL), weles.ArtifactPath(filename), ch)

			status.NewStatus = weles.ArtifactStatusDOWNLOADING
			checkChannels(ch, platinumKoala.notification, status)

			status.NewStatus = weles.ArtifactStatusREADY
			checkChannels(ch, platinumKoala.notification, status)

			content, err := ioutil.ReadFile(string(filename))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal(pigs))

			Consistently(func() string {
				return ws.GetString()
			}).Should(BeEmpty())
		})
		It("should fail when path is invalid", func() {
			ts = prepareServer(validURL)
			defer ts.Close()

			dir := invalidDir
			filename := weles.ArtifactPath(filepath.Join(dir, "test"))

			status := weles.ArtifactStatusChange{
				Path:      filename,
				NewStatus: weles.ArtifactStatusDOWNLOADING,
			}

			platinumKoala.download(weles.ArtifactURI(ts.URL), weles.ArtifactPath(filename), ch)

			status.NewStatus = weles.ArtifactStatusDOWNLOADING
			checkChannels(ch, platinumKoala.notification, status)

			status.NewStatus = weles.ArtifactStatusFAILED
			checkChannels(ch, platinumKoala.notification, status)

			Expect(string(filename)).NotTo(BeAnExistingFile())

			Eventually(func() string {
				return ws.GetString()
			}).Should(SatisfyAll(
				ContainSubstring("Failed to create file."),
				ContainSubstring("Failed to remove artifact."),
			))
		})
		DescribeTable("response to invalid url",
			func(valid bool) {
				ts = prepareServer(invalidURL)
				defer ts.Close()

				dir := validDir
				if !valid {
					dir = invalidDir
				}
				filename := weles.ArtifactPath(filepath.Join(dir, "test"))

				status := weles.ArtifactStatusChange{
					Path:      filename,
					NewStatus: weles.ArtifactStatusDOWNLOADING,
				}

				platinumKoala.download(weles.ArtifactURI(ts.URL), weles.ArtifactPath(filename), ch)

				status.NewStatus = weles.ArtifactStatusDOWNLOADING
				checkChannels(ch, platinumKoala.notification, status)

				status.NewStatus = weles.ArtifactStatusFAILED
				checkChannels(ch, platinumKoala.notification, status)

				Expect(string(filename)).NotTo(BeAnExistingFile())

				Eventually(func() string {
					return ws.GetString()
				}).Should(SatisfyAll(
					ContainSubstring(
						"Received wrong response from server after downloading artifact."),
					ContainSubstring("Failed to remove artifact."),
				))
			},
			Entry("fail when url is invalid", true),
			Entry("fail when url and path are invalid", false),
		)
	})

	Describe("Download(): Notify ch channel about any changes", func() {
		It("should download valid file to valid path", func() {
			ts = prepareServer(validURL)
			defer ts.Close()

			dir := validDir
			path := weles.ArtifactPath(filepath.Join(dir, "animal"))

			err := platinumKoala.Download(weles.ArtifactURI(ts.URL), path, ch)
			Expect(err).ToNot(HaveOccurred())

			status := weles.ArtifactStatusChange{
				Path:      path,
				NewStatus: weles.ArtifactStatusPENDING,
			}
			Eventually(ch).Should(Receive(Equal(status)))

			status.NewStatus = weles.ArtifactStatusDOWNLOADING
			Eventually(ch).Should(Receive(Equal(status)))

			status.NewStatus = weles.ArtifactStatusREADY
			Eventually(ch).Should(Receive(Equal(status)))

			Consistently(func() string {
				return ws.GetString()
			}).Should(BeEmpty())
		})
		It("should fail when path is invalid", func() {
			ts = prepareServer(validURL)
			defer ts.Close()

			dir := invalidDir
			path := weles.ArtifactPath(filepath.Join(dir, "animal"))

			err := platinumKoala.Download(weles.ArtifactURI(ts.URL), path, ch)
			Expect(err).ToNot(HaveOccurred())

			status := weles.ArtifactStatusChange{
				Path:      path,
				NewStatus: weles.ArtifactStatusPENDING,
			}
			Eventually(ch).Should(Receive(Equal(status)))

			status.NewStatus = weles.ArtifactStatusDOWNLOADING
			Eventually(ch).Should(Receive(Equal(status)))

			status.NewStatus = weles.ArtifactStatusFAILED
			Eventually(ch).Should(Receive(Equal(status)))

			Eventually(func() string {
				return ws.GetString()
			}).Should(SatisfyAll(
				ContainSubstring("Failed to create file."),
				ContainSubstring("Failed to remove artifact."),
			))
		})
		DescribeTable("response to invalid url",
			func(valid bool) {
				ts = prepareServer(invalidURL)
				defer ts.Close()

				dir := validDir
				if !valid {
					dir = invalidDir
				}
				path := weles.ArtifactPath(filepath.Join(dir, "animal"))

				err := platinumKoala.Download(weles.ArtifactURI(ts.URL), path, ch)
				Expect(err).ToNot(HaveOccurred())

				status := weles.ArtifactStatusChange{
					Path:      path,
					NewStatus: weles.ArtifactStatusPENDING,
				}
				Eventually(ch).Should(Receive(Equal(status)))

				status.NewStatus = weles.ArtifactStatusDOWNLOADING
				Eventually(ch).Should(Receive(Equal(status)))

				status.NewStatus = weles.ArtifactStatusFAILED
				Eventually(ch).Should(Receive(Equal(status)))

				Eventually(func() string {
					return ws.GetString()
				}).Should(SatisfyAll(
					ContainSubstring(
						"Received wrong response from server after downloading artifact."),
					ContainSubstring("Failed to remove artifact."),
				))
			},
			Entry("fail when url is invalid", true),
			Entry("fail when url and path are invalid", false),
		)
	})

	DescribeTable("Download(): Download files to specified path.",
		func(url weles.ArtifactURI, filename string, poem string) {
			ts = prepareServer(url)
			defer ts.Close()

			path := weles.ArtifactPath(filepath.Join(validDir, filename))

			err := platinumKoala.Download(weles.ArtifactURI(ts.URL), path, ch)
			Expect(err).ToNot(HaveOccurred())

			Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{
				Path:      path,
				NewStatus: weles.ArtifactStatusPENDING,
			})))
			Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{
				Path:      path,
				NewStatus: weles.ArtifactStatusDOWNLOADING,
			})))

			if poem != "" {
				Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{
					Path:      path,
					NewStatus: weles.ArtifactStatusREADY,
				})))
				content, err := ioutil.ReadFile(string(path))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(BeIdenticalTo(poem))

				Consistently(func() string {
					return ws.GetString()
				}).Should(BeEmpty())
			} else {
				Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{
					Path:      path,
					NewStatus: weles.ArtifactStatusFAILED,
				})))
				content, err := ioutil.ReadFile(string(path))
				Expect(err).To(HaveOccurred())
				Expect(content).To(BeNil())

				Eventually(func() string {
					return ws.GetString()
				}).Should(SatisfyAll(
					ContainSubstring(
						"Received wrong response from server after downloading artifact."),
					ContainSubstring("Failed to remove artifact."),
				))
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

			Consistently(func() string {
				return ws.GetString()
			}).Should(BeEmpty())
		})
	})
})
