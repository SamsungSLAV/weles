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

	"github.com/SamsungSLAV/perun/testutil"
	"github.com/SamsungSLAV/weles"
	cmock "github.com/SamsungSLAV/weles/controller/mock"
	"github.com/SamsungSLAV/weles/controller/notifier"
	mock "github.com/SamsungSLAV/weles/mock"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("DryaderImpl", func() {
	var r <-chan notifier.Notification
	var jc *cmock.MockJobsController
	var djm *mock.MockDryadJobManager
	var h Dryader
	var ctrl *gomock.Controller
	lock := sync.Locker(new(sync.Mutex))
	j := weles.JobID(0xCAFE)
	dryad := weles.Dryad{Addr: &net.IPNet{IP: net.IPv4(1, 2, 3, 4), Mask: net.IPv4Mask(5, 6, 7, 8)}}
	err := errors.New("test error")
	conf := weles.Config{JobName: "test123"}

	expectRegistered := func(offset int) {
		h.(*DryaderImpl).mutex.Lock()
		defer h.(*DryaderImpl).mutex.Unlock()

		ExpectWithOffset(offset, len(h.(*DryaderImpl).info)).To(Equal(1))
		info, ok := h.(*DryaderImpl).info[j]
		ExpectWithOffset(offset, ok).To(BeTrue())
		ExpectWithOffset(offset, info).To(BeTrue())
	}
	eventuallyEmpty := func(offset int) {
		EventuallyWithOffset(offset, func() int {
			h.(*DryaderImpl).mutex.Lock()
			defer h.(*DryaderImpl).mutex.Unlock()
			return len(h.(*DryaderImpl).info)
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
		lock.Lock()
		defer lock.Unlock()

		ctrl = gomock.NewController(GinkgoT())

		jc = cmock.NewMockJobsController(ctrl)
		djm = mock.NewMockDryadJobManager(ctrl)

		h = NewDryader(jc, djm)
		r = h.Listen()
	})
	AfterEach(func() {
		lock.Lock()
		defer lock.Unlock()

		h.(*DryaderImpl).Finish()
		ctrl.Finish()
	})

	Describe("NewBoruter", func() {
		It("should create a new object", func() {
			Expect(h).NotTo(BeNil())
			Expect(h.(*DryaderImpl).jobs).To(Equal(jc))
			Expect(h.(*DryaderImpl).djm).To(Equal(djm))
			Expect(h.(*DryaderImpl).info).NotTo(BeNil())
			Expect(h.(*DryaderImpl).mutex).NotTo(BeNil())
			Expect(h.(*DryaderImpl).finish).NotTo(BeNil())
		})
	})

	Describe("StartJob", func() {
		It("should register job successfully", func() {
			jc.EXPECT().GetDryad(j).Return(dryad, nil)
			jc.EXPECT().GetConfig(j).Return(conf, nil)
			djm.EXPECT().Create(j, dryad, conf, (chan<- weles.DryadJobStatusChange)(
				h.(*DryaderImpl).listener))

			h.StartJob(j)
			expectRegistered(1)
		})
		It("should fail if DryadJobManager.Create fails", func() {
			lock.Lock()
			jc.EXPECT().GetDryad(j).Return(dryad, nil)
			jc.EXPECT().GetConfig(j).Return(conf, nil)
			djm.EXPECT().Create(j, dryad, conf, (chan<- weles.DryadJobStatusChange)(
				h.(*DryaderImpl).listener)).Return(err)
			lock.Unlock()

			log, logerr := testutil.WithStderrMocked(func() {
				defer GinkgoRecover()
				lock.Lock()
				defer lock.Unlock()
				h.StartJob(j)

				eventuallyNoti(1, false, "Cannot delegate Job to Dryad : test error")
				eventuallyEmpty(1)
			})
			Expect(logerr).NotTo(HaveOccurred())
			Expect(log).To(ContainSubstring("Failed to start Job execution on Dryad."))
		})
		It("should fail if JobManager.GetDryad fails", func() {
			lock.Lock()
			jc.EXPECT().GetDryad(j).Return(weles.Dryad{}, err)
			lock.Unlock()

			log, logerr := testutil.WithStderrMocked(func() {
				defer GinkgoRecover()
				lock.Lock()
				defer lock.Unlock()
				h.StartJob(j)

				eventuallyNoti(1, false,
					"Internal Weles error while getting Dryad for Job : test error")
				eventuallyEmpty(1)
			})
			Expect(logerr).NotTo(HaveOccurred())
			Expect(log).To(ContainSubstring("Failed to get Dryad for Job."))
		})
	})

	Describe("With registered request", func() {
		updateStates := []weles.DryadJobStatus{
			weles.DryadJobStatusNEW,
			weles.DryadJobStatusDEPLOY,
			weles.DryadJobStatusBOOT,
			weles.DryadJobStatusTEST,
		}
		updateMsgs := []string{
			"Started",
			"Deploying",
			"Booting",
			"Testing",
		}
		BeforeEach(func() {
			jc.EXPECT().GetDryad(j).Return(dryad, nil)
			jc.EXPECT().GetConfig(j).Return(conf, nil)
			djm.EXPECT().Create(
				j, dryad, conf, (chan<- weles.DryadJobStatusChange)(h.(*DryaderImpl).listener))

			h.StartJob(j)

			expectRegistered(1)
		})

		It("should ignore ID of not registered request", func() {
			states := []weles.DryadJobStatus{
				weles.DryadJobStatusNEW,
				weles.DryadJobStatusDEPLOY,
				weles.DryadJobStatusBOOT,
				weles.DryadJobStatusTEST,
				weles.DryadJobStatusFAIL,
				weles.DryadJobStatusOK,
			}
			for _, s := range states {
				change := weles.DryadJobInfo{Job: weles.JobID(0x0BCA), Status: s}
				h.(*DryaderImpl).listener <- weles.DryadJobStatusChange(change)

				expectRegistered(1)
			}
		})
		It("should update status of the Job", func() {
			for i, s := range updateStates {
				change := weles.DryadJobInfo{Job: j, Status: s}
				jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusRUNNING, updateMsgs[i])

				h.(*DryaderImpl).listener <- weles.DryadJobStatusChange(change)

				expectRegistered(1)
			}
		})
		updateTableEntries := func() []TableEntry {
			var ret []TableEntry
			for i, s := range updateStates {
				ret = append(ret, Entry(string(s), s, updateMsgs[i]))
			}
			return ret
		}()
		DescribeTable("should fail if updating status of the Job fails",
			func(s weles.DryadJobStatus, msg string) {
				lock.Lock()
				change := weles.DryadJobInfo{Job: j, Status: s}
				jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusRUNNING, msg).Return(err)
				lock.Unlock()

				log, logerr := testutil.WithStderrMocked(func() {
					defer GinkgoRecover()
					lock.Lock()
					defer lock.Unlock()
					h.(*DryaderImpl).listener <- weles.DryadJobStatusChange(change)

					eventuallyNoti(1, false,
						"Internal Weles error while changing Job status : test error")
					eventuallyEmpty(1)
				})
				Expect(logerr).NotTo(HaveOccurred())
				Expect(log).To(ContainSubstring("Failed to change job state to RUNNING."))
			},
			updateTableEntries...,
		)

		It("should fail if Dryad Job fails", func() {
			change := weles.DryadJobInfo{Job: j, Status: weles.DryadJobStatusFAIL}

			h.(*DryaderImpl).listener <- weles.DryadJobStatusChange(change)

			eventuallyNoti(1, false, "Failed to execute test on Dryad.")
			eventuallyEmpty(1)
		})
		It("should notify about successfully completed Dryad Job", func() {
			change := weles.DryadJobInfo{Job: j, Status: weles.DryadJobStatusOK}

			h.(*DryaderImpl).listener <- weles.DryadJobStatusChange(change)

			eventuallyNoti(1, true, "")
			eventuallyEmpty(1)
		})

		Describe("CancelJob", func() {
			It("should remove Job and cancel it in Dryad Job Manager", func() {
				djm.EXPECT().Cancel(j)

				h.CancelJob(j)

				eventuallyEmpty(1)
			})
			It("should ignore djm's Cancel error", func() {
				djm.EXPECT().Cancel(j).Return(err)

				h.CancelJob(j)

				eventuallyEmpty(1)
			})
			It("should ignore not existing request", func() {
				h.CancelJob(weles.JobID(0x0BCA))
				expectRegistered(1)
			})
		})
	})
})
