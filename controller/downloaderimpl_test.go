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

package controller

import (
	"errors"
	"fmt"
	"sync"

	"github.com/SamsungSLAV/weles"
	cmock "github.com/SamsungSLAV/weles/controller/mock"
	"github.com/SamsungSLAV/weles/controller/notifier"
	mock "github.com/SamsungSLAV/weles/mock"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DownloaderImpl", func() {
	var r <-chan notifier.Notification
	var jc *cmock.MockJobsController
	var am *mock.MockArtifactManager
	var h *DownloaderImpl
	var ctrl *gomock.Controller
	j := weles.JobID(0xCAFE)
	paths := []string{}
	for i := 0; i < 9; i++ {
		paths = append(paths, fmt.Sprintf("path_%d", i))
	}
	infos := []string{""}
	for i := 1; i <= 7; i++ {
		infos = append(infos, fmt.Sprintf("%d / 7 artifacts ready", i))
	}
	config := weles.Config{Action: weles.Action{
		Deploy: weles.Deploy{Images: []weles.ImageDefinition{
			{URI: "image_0", ChecksumURI: "md5_0"},
			{URI: "image_1"},
			{ChecksumURI: "md5_2"},
		}},
		Test: weles.Test{TestCases: []weles.TestCase{
			{TestActions: []weles.TestAction{
				weles.Push{URI: "uri_0", Alias: "alias_0"},
				weles.Push{URI: "uri_1", Alias: "alias_1"},
				weles.Pull{Alias: "alias_2"},
			}},
			{TestActions: []weles.TestAction{
				weles.Push{URI: "uri_3", Alias: "alias_3"},
			}},
			{TestActions: []weles.TestAction{
				weles.Pull{Alias: "alias_4"},
			}},
		}},
	}}
	updatedConfig := weles.Config{Action: weles.Action{
		Deploy: weles.Deploy{Images: []weles.ImageDefinition{
			{URI: "image_0", ChecksumURI: "md5_0", Path: paths[0], ChecksumPath: paths[1]},
			{URI: "image_1", Path: paths[2]},
			{ChecksumURI: "md5_2", ChecksumPath: paths[3]},
		}},
		Test: weles.Test{TestCases: []weles.TestCase{
			{TestActions: []weles.TestAction{
				weles.Push{URI: "uri_0", Alias: "alias_0", Path: paths[4]},
				weles.Push{URI: "uri_1", Alias: "alias_1", Path: paths[5]},
				weles.Pull{Alias: "alias_2", Path: paths[7]},
			}},
			{TestActions: []weles.TestAction{
				weles.Push{URI: "uri_3", Alias: "alias_3", Path: paths[6]},
			}},
			{TestActions: []weles.TestAction{
				weles.Pull{Alias: "alias_4", Path: paths[8]},
			}},
		}},
	}}
	err := errors.New("test error")

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		jc = cmock.NewMockJobsController(ctrl)
		am = mock.NewMockArtifactManager(ctrl)

		h = NewDownloader(jc, am).(*DownloaderImpl)
		r = h.Listen()
	})
	AfterEach(func() {
		ctrl.Finish()
	})
	Describe("NewDownloader", func() {
		It("should create a new object", func() {
			Expect(h).NotTo(BeNil())
			Expect(h.jobs).To(Equal(jc))
			Expect(h.artifacts).To(Equal(am))
			Expect(h.collector).NotTo(BeNil())
			Expect(h.path2Job).NotTo(BeNil())
			Expect(h.info).NotTo(BeNil())
			Expect(h.mutex).NotTo(BeNil())
		})
	})
	Describe("Loop", func() {
		It("should stop loop function after closing collector channel", func() {
			close(h.collector)
		})
	})
	Describe("DispatchDownloads", func() {
		sendChange := func(from, to int, status weles.ArtifactStatus) {
			for i := from; i < to; i++ {
				h.collector <- weles.ArtifactStatusChange{
					Path:      weles.ArtifactPath(paths[i]),
					NewStatus: status,
				}
			}
		}
		eventuallyNoti := func(offset int, ok bool, msg string) {
			expectedNotification := notifier.Notification{
				JobID: j,
				OK:    ok,
				Msg:   msg,
			}
			EventuallyWithOffset(offset, r).Should(Receive(Equal(expectedNotification)))
		}
		eventuallyPathEmpty := func(offset int) {
			EventuallyWithOffset(offset, func() int {
				h.mutex.Lock()
				defer h.mutex.Unlock()
				return len(h.path2Job)
			}).Should(BeZero())
		}
		eventuallyInfoEmpty := func(offset int) {
			EventuallyWithOffset(offset, func() int {
				h.mutex.Lock()
				defer h.mutex.Unlock()
				return len(h.info)
			}).Should(BeZero())
		}
		eventuallyEmpty := func(offset int) {
			eventuallyPathEmpty(offset + 1)
			eventuallyInfoEmpty(offset + 1)
		}
		expectInfo := func(offset int, config bool, paths int) {
			h.mutex.Lock()
			defer h.mutex.Unlock()
			ExpectWithOffset(offset, len(h.info)).To(Equal(1))
			v, ok := h.info[j]
			ExpectWithOffset(offset, ok).To(BeTrue())
			ExpectWithOffset(offset, v.configSaved).To(Equal(config))
			ExpectWithOffset(offset, v.failed).To(Equal(0))
			ExpectWithOffset(offset, v.paths).To(Equal(paths))
			ExpectWithOffset(offset, v.ready).To(Equal(0))
		}
		expectPath := func(offset, from, to int) {
			h.mutex.Lock()
			defer h.mutex.Unlock()
			ExpectWithOffset(offset, len(h.path2Job)).To(Equal(to - from))
			for i := from; i < to; i++ {
				v, ok := h.path2Job[paths[i]]
				ExpectWithOffset(offset, ok).To(BeTrue(), "i = %d", i)
				ExpectWithOffset(offset, v).To(Equal(j), "i = %d", i)
			}
		}
		expectFail := func(offset int, pathsNo int, msg string) {
			eventuallyNoti(offset+1, false, msg)
			expectPath(offset+1, 0, pathsNo)
			eventuallyInfoEmpty(offset + 1)
			sendChange(0, pathsNo, weles.ArtifactStatusREADY)
			eventuallyPathEmpty(offset + 1)
		}
		defaultSetStatusAndInfo := func(successfulEntries int, fail bool) *gomock.Call {
			var i int
			var prev, call *gomock.Call

			for i = 0; i < successfulEntries; i++ {
				call = jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusDOWNLOADING, infos[i])
				if prev != nil {
					call.After(prev)
				}
				prev = call
			}
			if fail {
				call = jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusDOWNLOADING,
					infos[i]).Return(err)
				if prev != nil {
					call.After(prev)
				}
			}
			return call
		}
		defaultGetConfig := func() {
			jc.EXPECT().GetConfig(j).Return(config, nil)
		}
		defaultPush := func(successfulEntries int, fail bool) *gomock.Call {
			types := []weles.ArtifactType{
				weles.ArtifactTypeIMAGE,
				weles.ArtifactTypeIMAGE,
				weles.ArtifactTypeIMAGE,
				weles.ArtifactTypeIMAGE,
				weles.ArtifactTypeTEST,
				weles.ArtifactTypeTEST,
				weles.ArtifactTypeTEST,
			}
			aliases := []weles.ArtifactAlias{"Image_0", "ImageMD5_0", "Image_1", "ImageMD5_2",
				"alias_0", "alias_1", "alias_3"}
			uris := []weles.ArtifactURI{"image_0", "md5_0", "image_1", "md5_2", "uri_0", "uri_1",
				"uri_3"}
			var i int
			var prev, call *gomock.Call

			for i = 0; i < successfulEntries; i++ {
				call = am.EXPECT().Download(
					weles.ArtifactDescription{
						JobID: j,
						Type:  types[i],
						Alias: aliases[i],
						URI:   uris[i],
					},
					h.collector).Return(weles.ArtifactPath(paths[i]), nil)
				if prev != nil {
					call.After(prev)
				}
				prev = call
			}
			if fail {
				call = am.EXPECT().Download(
					weles.ArtifactDescription{
						JobID: j,
						Type:  types[i],
						Alias: aliases[i],
						URI:   uris[i],
					},
					h.collector).Return(weles.ArtifactPath(""), err)
				if prev != nil {
					call.After(prev)
				}
			}
			return call
		}
		defaultCreate := func(successfulEntries int, fail bool) *gomock.Call {
			types := []weles.ArtifactType{weles.ArtifactTypeTEST, weles.ArtifactTypeTEST}
			aliases := []weles.ArtifactAlias{"alias_2", "alias_4"}
			returnPaths := []weles.ArtifactPath{weles.ArtifactPath(paths[7]),
				weles.ArtifactPath(paths[8])}
			var i int
			var prev, call *gomock.Call

			for i = 0; i < successfulEntries; i++ {
				call = am.EXPECT().CreateArtifact(
					weles.ArtifactDescription{JobID: j, Type: types[i], Alias: aliases[i]}).
					Return(returnPaths[i], nil)
				if prev != nil {
					call.After(prev)
				}
				prev = call
			}
			if fail {
				call = am.EXPECT().CreateArtifact(
					weles.ArtifactDescription{JobID: j, Type: types[i], Alias: aliases[i]}).
					Return(weles.ArtifactPath(""), err)
				if prev != nil {
					call.After(prev)
				}
			}
			return call
		}
		defaultSetConfig := func() {
			jc.EXPECT().SetConfig(j, updatedConfig)
		}
		It("should delegate downloading of all artifacts successfully", func() {
			defaultSetStatusAndInfo(8, false)
			defaultGetConfig()
			defaultPush(7, false)
			defaultCreate(2, false)
			defaultSetConfig()

			h.DispatchDownloads(j)

			expectPath(1, 0, 7)
			expectInfo(1, true, 7)

			sendChange(0, 7, weles.ArtifactStatusREADY)

			eventuallyNoti(1, true, "")
			eventuallyEmpty(1)
		})
		It("should fail if cannot set config", func() {
			defaultSetStatusAndInfo(1, false)
			defaultGetConfig()
			defaultPush(7, false)
			defaultCreate(2, false)
			jc.EXPECT().SetConfig(j, updatedConfig).Return(err)

			h.DispatchDownloads(j)

			expectFail(1, 7,
				"Internal Weles error while setting config : test error")
		})
		It("should fail if pull fails", func() {
			defaultSetStatusAndInfo(1, false)
			defaultGetConfig()
			defaultPush(6, false)
			defaultCreate(0, true)

			h.DispatchDownloads(j)

			expectFail(1, 6,
				"Internal Weles error while creating a new path in ArtifactManager : "+
					"test error")
		})
		It("should fail if push for TESTFILE fails", func() {
			defaultSetStatusAndInfo(1, false)
			jc.EXPECT().GetConfig(j).Return(config, nil)
			defaultPush(4, true)

			h.DispatchDownloads(j)

			expectFail(1, 4,
				"Internal Weles error while registering URI:<uri_0> in ArtifactManager : "+
					"test error")
		})
		It("should fail if push for MD5 fails", func() {
			defaultSetStatusAndInfo(1, false)
			jc.EXPECT().GetConfig(j).Return(config, nil)
			defaultPush(1, true)

			h.DispatchDownloads(j)

			expectFail(
				1, 1, "Internal Weles error while registering URI:<md5_0> in ArtifactManager : "+
					"test error")
		})
		It("should fail if push for image fails", func() {
			defaultSetStatusAndInfo(1, false)
			jc.EXPECT().GetConfig(j).Return(config, nil)
			defaultPush(2, true)

			h.DispatchDownloads(j)

			expectFail(1, 2, "Internal Weles error while registering URI:<image_1> in "+
				"ArtifactManager : test error")
		})
		It("should fail if getting config fails", func() {
			defaultSetStatusAndInfo(1, false)
			jc.EXPECT().GetConfig(j).Return(weles.Config{}, err)

			h.DispatchDownloads(j)

			expectFail(1, 0, "Internal Weles error while getting Job config : "+
				"test error")
		})
		It("should fail if setting status fails", func() {
			defaultSetStatusAndInfo(0, true)

			h.DispatchDownloads(j)

			expectFail(1, 0, "Internal Weles error while changing Job status : test error")
		})
		It("should succeed when there is nothing to download", func() {
			emptyConfig := weles.Config{Action: weles.Action{
				Deploy: weles.Deploy{Images: []weles.ImageDefinition{}},
				Test: weles.Test{TestCases: []weles.TestCase{
					{TestActions: []weles.TestAction{
						weles.Pull{Alias: "alias_2"},
					}},
					{TestActions: []weles.TestAction{
						weles.Pull{Alias: "alias_4"},
					}},
				}},
			}}
			emptyUpdatedConfig := weles.Config{Action: weles.Action{
				Deploy: weles.Deploy{Images: []weles.ImageDefinition{}},
				Test: weles.Test{TestCases: []weles.TestCase{
					{TestActions: []weles.TestAction{
						weles.Pull{Alias: "alias_2", Path: paths[7]},
					}},
					{TestActions: []weles.TestAction{
						weles.Pull{Alias: "alias_4", Path: paths[8]},
					}},
				}},
			}}

			defaultSetStatusAndInfo(1, false)
			jc.EXPECT().GetConfig(j).Return(emptyConfig, nil)

			defaultCreate(2, false)
			jc.EXPECT().SetConfig(j, emptyUpdatedConfig)

			h.DispatchDownloads(j)

			eventuallyEmpty(1)
			eventuallyNoti(1, true, "")
		})
		It("should handle downloading failure", func() {
			c := defaultSetStatusAndInfo(4, false)
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusDOWNLOADING,
				"Failed to download artifact").After(c)
			defaultGetConfig()
			defaultPush(7, false)
			defaultCreate(2, false)
			defaultSetConfig()

			h.DispatchDownloads(j)

			expectPath(1, 0, 7)
			expectInfo(1, true, 7)

			sendChange(0, 3, weles.ArtifactStatusREADY)
			sendChange(3, 4, weles.ArtifactStatusFAILED)

			eventuallyNoti(1, false, formatDownload)
			expectPath(1, 4, 7)
			eventuallyInfoEmpty(1)

			sendChange(4, 7, weles.ArtifactStatusDOWNLOADING)
			eventuallyPathEmpty(1)
		})
		It("should block reply until configuration is saved and all artifacts are downloaded",
			func() {
				defaultSetStatusAndInfo(8, false)
				defaultGetConfig()
				defaultPush(7, false)
				defaultCreate(2, false)

				holdDownload := sync.WaitGroup{}
				holdDownload.Add(1)
				setConfigReached := sync.WaitGroup{}
				setConfigReached.Add(1)

				jc.EXPECT().SetConfig(j, updatedConfig).Do(func(weles.JobID, weles.Config) {
					setConfigReached.Done()
					holdDownload.Wait()
				})

				go h.DispatchDownloads(j)

				setConfigReached.Wait()

				expectPath(1, 0, 7)
				expectInfo(1, false, 7)

				sendChange(0, 7, weles.ArtifactStatusREADY)
				holdDownload.Done()

				eventuallyNoti(1, true, "")
				eventuallyEmpty(1)
			})
		It("should handle failure in updating info", func() {
			defaultSetStatusAndInfo(5, true)
			defaultGetConfig()
			defaultPush(7, false)
			defaultCreate(2, false)
			defaultSetConfig()

			h.DispatchDownloads(j)

			expectPath(1, 0, 7)
			expectInfo(1, true, 7)

			sendChange(0, 7, weles.ArtifactStatusREADY)

			eventuallyNoti(1, false, "Internal Weles error while changing Job status : test error")
			eventuallyEmpty(1)
		})
		It("should leave no data left if failure response is sent while still processing config",
			func() {
				defaultSetStatusAndInfo(5, true)
				defaultGetConfig()
				defaultPush(7, false)
				defaultCreate(2, false)

				holdDownload := sync.WaitGroup{}
				holdDownload.Add(1)
				setConfigReached := sync.WaitGroup{}
				setConfigReached.Add(1)

				jc.EXPECT().SetConfig(j, updatedConfig).Do(func(weles.JobID, weles.Config) {
					setConfigReached.Done()
					holdDownload.Wait()
				})

				go h.DispatchDownloads(j)
				setConfigReached.Wait()

				expectPath(1, 0, 7)
				expectInfo(1, false, 7)

				sendChange(0, 7, weles.ArtifactStatusREADY)

				eventuallyNoti(1, false,
					"Internal Weles error while changing Job status : test error")

				holdDownload.Done()

				eventuallyEmpty(1)
			})
		It("should leave no data left if failure response is sent while pushing", func() {
			c := defaultSetStatusAndInfo(1, false)
			jc.EXPECT().SetStatusAndInfo(
				j, weles.JobStatusDOWNLOADING, "1 / 1 artifacts ready").Return(err).After(c)
			defaultGetConfig()
			holdDownload := sync.WaitGroup{}
			holdDownload.Add(1)
			pushReached := sync.WaitGroup{}
			pushReached.Add(1)

			defaultPush(2, false).Do(func(weles.ArtifactDescription,
				chan weles.ArtifactStatusChange) {
				pushReached.Done()
				holdDownload.Wait()
			})

			go h.DispatchDownloads(j)
			pushReached.Wait()

			expectPath(1, 0, 1)
			expectInfo(1, false, 1)

			sendChange(0, 1, weles.ArtifactStatusREADY)

			eventuallyNoti(1, false, "Internal Weles error while changing Job status : test error")

			holdDownload.Done()
			sendChange(1, 2, weles.ArtifactStatusREADY)

			eventuallyEmpty(1)
		})
		It("should ignore changes to non-terminal states", func() {
			defaultSetStatusAndInfo(8, false)
			defaultGetConfig()
			defaultPush(7, false)
			defaultCreate(2, false)
			defaultSetConfig()

			h.DispatchDownloads(j)

			expectPath(1, 0, 7)
			expectInfo(1, true, 7)

			sendChange(0, 7, weles.ArtifactStatusDOWNLOADING)
			sendChange(0, 7, weles.ArtifactStatusPENDING)

			expectPath(1, 0, 7)
			expectInfo(1, true, 7)

			sendChange(0, 7, weles.ArtifactStatusREADY)

			eventuallyNoti(1, true, "")
			eventuallyEmpty(1)
		})
	})
})
