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
	reqfilter "github.com/SamsungSLAV/boruta/filter"
	"github.com/SamsungSLAV/weles"
	cmock "github.com/SamsungSLAV/weles/controller/mock"
	"github.com/SamsungSLAV/weles/controller/notifier"
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
	j := weles.JobID(0xCAFE)
	rid := boruta.ReqID(0xD0DA)
	period := 50 * time.Millisecond
	jobTimeout := time.Hour
	owner := boruta.UserInfo{}
	err := errors.New("test error")

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
	})
	AfterEach(func() {
		h.(*BoruterImpl).Finish()
		ctrl.Finish()
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
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes().Return(
				[]boruta.ReqInfo{}, &boruta.ListInfo{}, err).Do(
				func(boruta.ListFilter, *boruta.SortInfo, *boruta.RequestsPaginator) {
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
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes()

			before := time.Now()
			h.Request(j)
			after := time.Now()

			Expect(va).To(BeTemporally(">=", before))
			Expect(va).To(BeTemporally("<=", after))
			Expect(dl).To(BeTemporally(">=", before.Add(jobTimeout)))
			Expect(dl).To(BeTemporally("<=", after.Add(jobTimeout)))

			expectRegistered(1)
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
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes()

			before := time.Now()
			h.Request(j)
			after := time.Now()

			Expect(va).To(BeTemporally(">=", before))
			Expect(va).To(BeTemporally("<=", after))
			Expect(dl).To(BeTemporally(">=", before.Add(defaultDelay)))
			Expect(dl).To(BeTemporally("<=", after.Add(defaultDelay)))

			expectRegistered(1)
		})
		It("should fail if NewRequest fails", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(caps, priority, owner, gomock.Any(), gomock.Any()).Return(
				boruta.ReqID(0), err)
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Failed to create request in Boruta : test error")
			eventuallyEmpty(1)
		})
		It("should fail if GetConfig fails", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(weles.Config{}, err)
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Internal Weles error while getting Job config : test error")
			eventuallyEmpty(1)
		})
		It("should fail if SetStatusAndInfo fails", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "").Return(err)
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Internal Weles error while changing Job status : test error")
			eventuallyEmpty(1)
		})
		It("should call NewRequest with empty caps if no device type provided", func() {
			config.DeviceType = ""
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusWAITING, "")
			jc.EXPECT().GetConfig(j).Return(config, nil)
			req.EXPECT().NewRequest(boruta.Capabilities{}, priority, owner, gomock.Any(),
				gomock.Any()).Return(boruta.ReqID(0), err)
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes()

			h.Request(j)

			eventuallyNoti(1, false, "Failed to create request in Boruta : test error")
			eventuallyEmpty(1)
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
				req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes()

				h.Request(j)

				eventuallyNoti(1, false, "Failed to create request in Boruta : test error")
				eventuallyEmpty(1)
			}
		})
	})
	Describe("With registered request", func() {
		var listRequestRet chan []boruta.ReqInfo
		var remainingRequests chan uint64
		var sinceID chan boruta.ReqID
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
			listRequestRet = make(chan []boruta.ReqInfo, 10)
			remainingRequests = make(chan uint64, 10)
			sinceID = make(chan boruta.ReqID, 10)
			req.EXPECT().ListRequests(gomock.Any(), nil, gomock.Any()).AnyTimes().DoAndReturn(
				func(f boruta.ListFilter, s *boruta.SortInfo, p *boruta.RequestsPaginator) (
					[]boruta.ReqInfo, *boruta.ListInfo, error) {

					elems := <-listRequestRet
					remaining := <-remainingRequests
					lastID := <-sinceID

					Expect(f).NotTo(BeNil())
					fr := f.(*reqfilter.Requests)
					Expect(fr).NotTo(BeNil())
					Expect(fr.IDs).To(ConsistOf(rid))
					Expect(fr.Priorities).To(BeNil())
					Expect(fr.States).To(ConsistOf(
						boruta.INPROGRESS,
						boruta.CANCEL,
						boruta.DONE,
						boruta.TIMEOUT,
						boruta.INVALID,
						boruta.FAILED))
					Expect(s).To(BeNil())
					Expect(p).NotTo(BeNil())
					Expect(p.ID).To(Equal(lastID))
					Expect(p.Direction).To(Equal(boruta.DirectionForward))
					Expect(p.Limit).To(Equal(boruta.MaxPageLimit))

					return elems, &boruta.ListInfo{
						RemainingItems: remaining,
						TotalItems:     uint64(len(elems)),
					}, nil
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
				remainingRequests <- 0
				sinceID <- boruta.ReqID(0)

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
				remainingRequests <- 0
				sinceID <- boruta.ReqID(0)

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
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

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
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

			eventuallyNoti(1, false, "Internal Weles error while setting Dryad : test error")
			eventuallyEmpty(1)
		})
		It("should fail during acquire if AcquireWorker fails", func() {
			req.EXPECT().AcquireWorker(rid).Return(boruta.AccessInfo{}, err)

			rinfo := boruta.ReqInfo{
				ID:    rid,
				State: boruta.INPROGRESS,
				Job:   &boruta.JobInfo{Timeout: time.Now().AddDate(0, 0, 1)},
			}
			listRequestRet <- []boruta.ReqInfo{rinfo}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

			eventuallyNoti(1, false, "Cannot acquire worker from Boruta : test error")
			eventuallyEmpty(1)
		})
		It("should remove request if state changes to CANCEL", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.CANCEL}
			listRequestRet <- []boruta.ReqInfo{rinfo}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

			eventuallyEmpty(1)
		})
		It("should remove request if state changes to DONE", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.DONE}
			listRequestRet <- []boruta.ReqInfo{rinfo}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

			eventuallyEmpty(1)
		})
		It("should fail and remove request if state changes to TIMEOUT", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.TIMEOUT}
			listRequestRet <- []boruta.ReqInfo{rinfo}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

			eventuallyNoti(1, false, "Timeout in Boruta.")
			eventuallyEmpty(1)
		})
		It("should fail and remove request if state changes to INVALID", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.INVALID}
			listRequestRet <- []boruta.ReqInfo{rinfo}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

			eventuallyNoti(1, false, "No suitable device in Boruta to run test.")
			eventuallyEmpty(1)
		})
		It("should fail and remove request if state changes to FAILED", func() {
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.FAILED}
			listRequestRet <- []boruta.ReqInfo{rinfo}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

			eventuallyNoti(1, false, "Boruta failed during request execution.")
			eventuallyEmpty(1)
		})
		It("iterate through many pages of boruta's response", func() {
			infoF1 := boruta.ReqInfo{ID: boruta.ReqID(rid + 1), State: boruta.WAIT}
			infoF2 := boruta.ReqInfo{ID: boruta.ReqID(rid + 2), State: boruta.WAIT}
			infoF3 := boruta.ReqInfo{ID: boruta.ReqID(rid + 3), State: boruta.WAIT}
			infoF4 := boruta.ReqInfo{ID: boruta.ReqID(rid + 4), State: boruta.WAIT}
			infoF5 := boruta.ReqInfo{ID: boruta.ReqID(rid + 5), State: boruta.WAIT}
			infoF6 := boruta.ReqInfo{ID: boruta.ReqID(rid + 6), State: boruta.WAIT}
			infoF7 := boruta.ReqInfo{ID: boruta.ReqID(rid + 7), State: boruta.WAIT}
			infoF8 := boruta.ReqInfo{ID: boruta.ReqID(rid + 8), State: boruta.WAIT}

			// send first 3 items
			listRequestRet <- []boruta.ReqInfo{infoF1, infoF2, infoF3}
			remainingRequests <- 5
			sinceID <- boruta.ReqID(0)

			// send next 2 items
			listRequestRet <- []boruta.ReqInfo{infoF4, infoF5}
			remainingRequests <- 3
			sinceID <- boruta.ReqID(rid + 3)

			// send last 3 items
			listRequestRet <- []boruta.ReqInfo{infoF6, infoF7, infoF8}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(rid + 5)

			// make a new ListRequests call
			rinfo := boruta.ReqInfo{ID: rid, State: boruta.DONE}
			listRequestRet <- []boruta.ReqInfo{rinfo}
			remainingRequests <- 0
			sinceID <- boruta.ReqID(0)

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
			})
		})
	})
})
