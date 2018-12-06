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
	"net"
	"sync"
	"time"

	"github.com/SamsungSLAV/boruta"
	"github.com/SamsungSLAV/slav/logger"
	"github.com/SamsungSLAV/weles"
	cmock "github.com/SamsungSLAV/weles/controller/mock"
	"github.com/SamsungSLAV/weles/controller/notifier"
	"github.com/SamsungSLAV/weles/testutil"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BoruterImpl", func() {
	var r <-chan notifier.Notification
	var jc *cmock.MockJobsController
	var req *cmock.MockRequests
	var h Boruter
	var ctrl *gomock.Controller
	var config weles.Config
	var caps boruta.Capabilities
	var priority boruta.Priority
	var ws *testutil.WriterString
	j := weles.JobID(0xCAFE)
	rid := boruta.ReqID(0xD0DA)
	period := 50 * time.Millisecond
	jobTimeout := time.Hour
	owner := boruta.UserInfo{}
	err := errors.New("test error")

	log := logger.NewLogger()
	stderrLog := logger.NewLogger()
	stderrLog.AddBackend("default", logger.Backend{
		Filter:     logger.NewFilterPassAll(),
		Serializer: logger.NewSerializerText(),
		Writer:     logger.NewWriterStderr(),
	})

	expectRegistered := func(offset int) {
		h.(*BoruterImpl).mutex.Lock()
		defer h.(*BoruterImpl).mutex.Unlock()

		ExpectWithOffset(offset, len(h.(*BoruterImpl).info)).To(Equal(1))
		info, ok := h.(*BoruterImpl).info[j]
		ExpectWithOffset(offset, ok).To(BeTrue())
		ExpectWithOffset(offset, info.rid).To(Equal(rid))

		ExpectWithOffset(offset, len(h.(*BoruterImpl).rid2Job)).To(Equal(1))
		job, ok := h.(*BoruterImpl).rid2Job[rid]
		ExpectWithOffset(offset, ok).To(BeTrue())
		ExpectWithOffset(offset, job).To(Equal(j))
	}
	eventuallyEmpty := func(offset int) {
		EventuallyWithOffset(offset, func() int {
			h.(*BoruterImpl).mutex.Lock()
			defer h.(*BoruterImpl).mutex.Unlock()
			return len(h.(*BoruterImpl).info)
		}).Should(BeZero())
		EventuallyWithOffset(offset, func() int {
			h.(*BoruterImpl).mutex.Lock()
			defer h.(*BoruterImpl).mutex.Unlock()
			return len(h.(*BoruterImpl).rid2Job)
		}).Should(BeZero())
	}
	eventuallyNoti := func(offset int, ok bool, msg string) {
		expectedNotification := notifier.Notification{
			JobID: j,
			OK:    ok,
			Msg:   msg,
		}
		EventuallyWithOffset(offset, r).Should(Receive(Equal(expectedNotification)))
	}

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		jc = cmock.NewMockJobsController(ctrl)
		req = cmock.NewMockRequests(ctrl)

		h = NewBoruter(jc, req, period)
		r = h.Listen()

		config = weles.Config{
			DeviceType: "TestDeviceType",
			Priority:   "medium",
			Timeouts: weles.Timeouts{
				JobTimeout: weles.ValidPeriod(jobTimeout),
			},
		}
		caps = boruta.Capabilities{"device_type": "TestDeviceType"}
		priority = boruta.Priority(7)

		ws = testutil.NewWriterString()
		log.AddBackend("string", logger.Backend{
			Filter:     logger.NewFilterPassAll(),
			Serializer: logger.NewSerializerText(),
			Writer:     ws,
		})
		logger.SetDefault(log)
	})
	AfterEach(func() {
		h.(*BoruterImpl).Finish()
		ctrl.Finish()
		logger.SetDefault(stderrLog)
	})
	Describe("NewBoruter", func() {
		It("should create a new object", func() {
			Expect(h).NotTo(BeNil())
			Expect(h.(*BoruterImpl).jobs).To(Equal(jc))
			Expect(h.(*BoruterImpl).boruta).To(Equal(req))
			Expect(h.(*BoruterImpl).info).NotTo(BeNil())
			Expect(h.(*BoruterImpl).rid2Job).NotTo(BeNil())
			Expect(h.(*BoruterImpl).mutex).NotTo(BeNil())
			Expect(h.(*BoruterImpl).borutaCheckPeriod).To(Equal(period))
		})
	})
	Describe("loop", func() {
		It("should ignore ListRequests errors", func() {
			counter := 5
			mutex := &sync.Mutex{}
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(caps, priority, owner, gomock.Any(), gomock.Any()).Return(
				rid, nil)
			req.EXPECT().ListRequests(nil).AnyTimes().Return([]boruta.ReqInfo{}, err).Do(
				func(boruta.ListFilter) {
					mutex.Lock()
					defer mutex.Unlock()
					counter--
				})

			h.Request(j)
			Eventually(func() int {
				mutex.Lock()
				defer mutex.Unlock()
				return counter
			}).Should(BeNumerically("<", 0))

			expectRegistered(1)

			Consistently(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to list Boruta requests."))
		})
	})
	Describe("Request", func() {
		It("should register job successfully", func() {
			var va, dl time.Time
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(caps, priority, owner, gomock.Any(), gomock.Any()).Return(
				rid, nil).Do(
				func(c boruta.Capabilities, p boruta.Priority, ui boruta.UserInfo,
					validAfter, deadline time.Time) {

					va = validAfter
					dl = deadline
				})
			req.EXPECT().ListRequests(nil).AnyTimes()

			before := time.Now()
			h.Request(j)
			after := time.Now()

			Expect(va).To(BeTemporally(">=", before))
			Expect(va).To(BeTemporally("<=", after))
			Expect(dl).To(BeTemporally(">=", before.Add(jobTimeout)))
			Expect(dl).To(BeTemporally("<=", after.Add(jobTimeout)))

			expectRegistered(1)

			Consistently(func() string {
				return ws.GetString()
			}).Should(BeEmpty())
		})
		It("should register job successfully even when JobTimout is not defined in Config", func() {
			var va, dl time.Time
			config.Timeouts.JobTimeout = weles.ValidPeriod(0)
			defaultDelay := 24 * time.Hour

			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(caps, priority, owner, gomock.Any(), gomock.Any()).Return(
				rid, nil).Do(
				func(c boruta.Capabilities, p boruta.Priority, ui boruta.UserInfo,
					validAfter, deadline time.Time) {

					va = validAfter
					dl = deadline
				})
			req.EXPECT().ListRequests(nil).AnyTimes()

			before := time.Now()
			h.Request(j)
			after := time.Now()

			Expect(va).To(BeTemporally(">=", before))
			Expect(va).To(BeTemporally("<=", after))
			Expect(dl).To(BeTemporally(">=", before.Add(defaultDelay)))
			Expect(dl).To(BeTemporally("<=", after.Add(defaultDelay)))

			expectRegistered(1)

			Consistently(func() string {
				return ws.GetString()
			}).Should(BeEmpty())
		})
		It("should fail if NewRequest fails", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(caps, priority, owner, gomock.Any(), gomock.Any()).Return(
				boruta.ReqID(0), err)
			req.EXPECT().ListRequests(nil).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Failed to create request in Boruta : test error")
			eventuallyEmpty(1)

			Eventually(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to create request in Boruta."))
		})
		It("should fail if GetConfig fails", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(weles.Config{}, err)
			req.EXPECT().ListRequests(nil).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Internal Weles error while getting Job config : test error")
			eventuallyEmpty(1)

			Eventually(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to get Job config."))
		})
		It("should fail if SetStatusAndInfo fails", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "").Return(err)
			req.EXPECT().ListRequests(nil).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Internal Weles error while changing Job status : test error")
			eventuallyEmpty(1)

			Eventually(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to change JobStatus."))
		})
		It("should call NewRequest with empty caps if no device type provided", func() {
			config.DeviceType = ""
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(boruta.Capabilities{}, priority, owner, gomock.Any(),
				gomock.Any()).Return(boruta.ReqID(0), err)
			req.EXPECT().ListRequests(nil).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Failed to create request in Boruta : test error")
			eventuallyEmpty(1)

			Eventually(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to create request in Boruta."))
		})
		It("should call NewRequest with proper priority", func() {
			m := map[weles.Priority]boruta.Priority{
				weles.LOW:                 boruta.Priority(11),
				weles.MEDIUM:              boruta.Priority(7),
				weles.HIGH:                boruta.Priority(3),
				weles.Priority("unknown"): boruta.Priority(7),
			}
			for k, v := range m {
				config.Priority = k
				jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
				jc.EXPECT().GetConfig(j).Return(config, nil)
				req.EXPECT().NewRequest(caps, v, owner, gomock.Any(), gomock.Any()).Return(
					boruta.ReqID(0), err)
				req.EXPECT().ListRequests(nil).AnyTimes()

				h.Request(j)

				eventuallyNoti(1, false, "Failed to create request in Boruta : test error")
				eventuallyEmpty(1)

				Eventually(func() string {
					return ws.GetString()
				}).Should(ContainSubstring("Failed to create request in Boruta."))
			}
		})
	})
	Describe("With registered request", func() {
		var listRequestRet chan []boruta.ReqInfo
		states := []boruta.ReqState{
			boruta.WAIT,
			boruta.INPROGRESS,
			boruta.CANCEL,
			boruta.TIMEOUT,
			boruta.INVALID,
			boruta.DONE,
			boruta.FAILED,
		}
		ai := boruta.AccessInfo{
			Addr: &net.IPNet{
				IP:   net.IPv4(1, 2, 3, 4),
				Mask: net.IPv4Mask(5, 6, 7, 8),
			}}

		BeforeEach(func() {
			var va, dl time.Time
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(caps, priority, owner, gomock.Any(), gomock.Any()).Return(
				rid, nil).Do(
				func(c boruta.Capabilities, p boruta.Priority, ui boruta.UserInfo,
					validAfter, deadline time.Time,
				) {
					va = validAfter
					dl = deadline
				})
			listRequestRet = make(chan []boruta.ReqInfo)
			req.EXPECT().ListRequests(nil).AnyTimes().DoAndReturn(
				func(boruta.ListFilter) ([]boruta.ReqInfo, error) {
					return <-listRequestRet, nil
				})

			before := time.Now()
			h.Request(j)
			after := time.Now()

			Expect(va).To(BeTemporally(">=", before))
			Expect(va).To(BeTemporally("<=", after))
			Expect(dl).To(BeTemporally(">=", before.Add(jobTimeout)))
			Expect(dl).To(BeTemporally("<=", after.Add(jobTimeout)))

			expectRegistered(1)
		})
		It("should ignore ID of not registered request", func() {
			for _, s := range states {
				rinfo := boruta.ReqInfo{ID: boruta.ReqID(0x0BCA), State: s}
				listRequestRet <- []boruta.ReqInfo{rinfo}

				expectRegistered(1)
			}
		})
		for _, s := range states { // Every state is a separate It,
			//because objects must be reinitialized.
			It("should ignore if request's state is unchanged : "+string(s), func() {
				h.(*BoruterImpl).mutex.Lock()

				Expect(len(h.(*BoruterImpl).info)).To(Equal(1))
				info, ok := h.(*BoruterImpl).info[j]
				Expect(ok).To(BeTrue())
				info.status = s

				h.(*BoruterImpl).mutex.Unlock()

				rinfo := boruta.ReqInfo{ID: rid, State: s}
				listRequestRet <- []boruta.ReqInfo{rinfo}

				expectRegistered(1)
			})
		}
		It("should acquire Dryad if state changes to INPROGRESS", func() {
			req.EXPECT().AcquireWorker(rid).Return(ai, nil)
			jc.EXPECT().SetDryad(
				j, weles.Dryad{
					Addr:     ai.Addr,
					Key:      ai.Key,
					Username: "boruta-user",
				})

			rinfo := boruta.ReqInfo{
				ID:    rid,
				State: boruta.INPROGRESS,
				Job:   &boruta.JobInfo{Timeout: time.Now().AddDate(0, 0, 1)},
			}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyNoti(1, true, "")
		})
		It("should fail during acquire if SetDryad fails", func() {
			req.EXPECT().AcquireWorker(rid).Return(ai, nil)
			jc.EXPECT().SetDryad(
				j, weles.Dryad{
					Addr:     ai.Addr,
					Key:      ai.Key,
					Username: "boruta-user",
				}).Return(err)

			rinfo := boruta.ReqInfo{
				ID:    rid,
				State: boruta.INPROGRESS,
				Job:   &boruta.JobInfo{Timeout: time.Now().AddDate(0, 0, 1)},
			}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyNoti(1, false, "Internal Weles error while setting Dryad : test error")
			eventuallyEmpty(1)

			Eventually(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to set up the dryad."))
		})
		It("should fail during acquire if AcquireWorker fails", func() {
			req.EXPECT().AcquireWorker(rid).Return(boruta.AccessInfo{}, err)

			rinfo := boruta.ReqInfo{
				ID:    rid,
				State: boruta.INPROGRESS,
				Job:   &boruta.JobInfo{Timeout: time.Now().AddDate(0, 0, 1)},
			}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyNoti(1, false, "Cannot acquire worker from Boruta : test error")
			eventuallyEmpty(1)

			Eventually(func() string {
				return ws.GetString()
			}).Should(ContainSubstring("Failed to acquire worker from Boruta."))
		})
		It("should remove request if state changes to CANCEL", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.CANCEL}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyEmpty(1)
		})
		It("should remove request if state changes to DONE", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.DONE}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyEmpty(1)
		})
		It("should fail and remove request if state changes to TIMEOUT", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.TIMEOUT}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyNoti(1, false, "Timeout in Boruta.")
			eventuallyEmpty(1)
		})
		It("should fail and remove request if state changes to INVALID", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.INVALID}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyNoti(1, false, "No suitable device in Boruta to run test.")
			eventuallyEmpty(1)
		})
		It("should fail and remove request if state changes to FAILED", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.FAILED}
			listRequestRet <- []boruta.ReqInfo{rinfo}

			eventuallyNoti(1, false, "Boruta failed during request execution.")
			eventuallyEmpty(1)
		})
		Describe("Release", func() {
			It("should remove existing request and close it in Boruta", func() {
				req.EXPECT().CloseRequest(rid)

				h.Release(j)

				eventuallyEmpty(1)
			})
			It("should ignore not existing request", func() {
				h.Release(weles.JobID(0x0BCA))
				expectRegistered(1)

				Eventually(func() string {
					return ws.GetString()
				}).Should(SatisfyAll(
					ContainSubstring("JobID not found in BoruterImpl.info map."),
					ContainSubstring("Failed to return Dryad to Boruta's pool."),
				))
			})
		})
	})
})
