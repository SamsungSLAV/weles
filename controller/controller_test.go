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
	"sync"
	"time"

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

var _ = Describe("JobManager", func() {
	It("should create a new object", func() {
		ctrl := gomock.NewController(GinkgoT())

		arm := mock.NewMockArtifactManager(ctrl)
		yap := mock.NewMockParser(ctrl)
		bor := cmock.NewMockRequests(ctrl)
		djm := mock.NewMockDryadJobManager(ctrl)

		bor.EXPECT().ListRequests(nil).AnyTimes()

		jm := NewJobManager(arm, yap, bor, time.Second, djm)
		Expect(jm).NotTo(BeNil())

		ctrl.Finish()
	})
})

var _ = Describe("Controller", func() {
	var (
		jc      *cmock.MockJobsController
		par     *cmock.MockParser
		dow     *cmock.MockDownloader
		bor     *cmock.MockBoruter
		dry     *cmock.MockDryader
		h       *Controller
		ctrl    *gomock.Controller
		parChan chan notifier.Notification
		dowChan chan notifier.Notification
		borChan chan notifier.Notification
		dryChan chan notifier.Notification
		done    bool
		mutex   *sync.Mutex
	)

	j := weles.JobID(0xCAFE)
	testErr := errors.New("test error")
	yaml := []byte("test yaml")
	testMsg := "test msg"
	notiOk := notifier.Notification{JobID: j, OK: true}
	notiFail := notifier.Notification{JobID: j, OK: false, Msg: testMsg}

	setDone := func(weles.JobID) {
		mutex.Lock()
		defer mutex.Unlock()
		done = true
	}
	eventuallyDone := func() {
		Eventually(func() bool {
			mutex.Lock()
			defer mutex.Unlock()
			return done
		}).Should(BeTrue())
	}

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		jc = cmock.NewMockJobsController(ctrl)
		par = cmock.NewMockParser(ctrl)
		dow = cmock.NewMockDownloader(ctrl)
		bor = cmock.NewMockBoruter(ctrl)
		dry = cmock.NewMockDryader(ctrl)

		parChan = make(chan notifier.Notification)
		dowChan = make(chan notifier.Notification)
		borChan = make(chan notifier.Notification)
		dryChan = make(chan notifier.Notification)

		par.EXPECT().Listen().AnyTimes().Return((<-chan notifier.Notification)(parChan))
		dow.EXPECT().Listen().AnyTimes().Return((<-chan notifier.Notification)(dowChan))
		bor.EXPECT().Listen().AnyTimes().Return((<-chan notifier.Notification)(borChan))
		dry.EXPECT().Listen().AnyTimes().Return((<-chan notifier.Notification)(dryChan))

		h = NewController(jc, par, dow, bor, dry)

		mutex = new(sync.Mutex)
		done = false
	})
	AfterEach(func() {
		h.Finish()
		ctrl.Finish()
	})

	Describe("NewController", func() {
		It("should create a new object", func() {
			Expect(h).NotTo(BeNil())
			Expect(h.jobs).To(Equal(jc))
			Expect(h.parser).To(Equal(par))
			Expect(h.downloader).To(Equal(dow))
			Expect(h.boruter).To(Equal(bor))
			Expect(h.dryader).To(Equal(dry))
			Expect(h.finish).NotTo(BeNil())
		})
	})
	Describe("CreateJob", func() {
		It("should create a new Job and delegate parsing", func() {
			jc.EXPECT().NewJob(yaml).Return(j, nil)
			par.EXPECT().Parse(j).Do(setDone)

			retJobID, retErr := h.CreateJob(yaml)

			Expect(retErr).NotTo(HaveOccurred())
			Expect(retJobID).To(Equal(j))
			eventuallyDone()
		})
		It("should fail if JobsController.NewJob fails", func() {
			jc.EXPECT().NewJob(yaml).Return(weles.JobID(0), testErr)

			log, logerr := testutil.WithStderrMocked(func() {
				defer GinkgoRecover()
				retJobID, retErr := h.CreateJob(yaml)

				Expect(retErr).To(Equal(testErr))
				Expect(retJobID).To(Equal(weles.JobID(0)))
			})
			Expect(logerr).NotTo(HaveOccurred())
			Expect(log).To(ContainSubstring("Failed to create new job."))
		})
	})

	Describe("CancelJob", func() {
		It("should cancel Job, stop execution on Dryad and release Dryad to Boruta", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusCANCELED, "")
			dry.EXPECT().CancelJob(j)
			bor.EXPECT().Release(j)

			retErr := h.CancelJob(j)

			Expect(retErr).To(BeNil())
		})
		It("should return error if Job fails to be cancelled", func() {
			jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusCANCELED, "").Return(testErr)

			log, logerr := testutil.WithStderrMocked(func() {
				defer GinkgoRecover()
				retErr := h.CancelJob(j)

				Expect(retErr).To(Equal(testErr))
			})
			Expect(logerr).NotTo(HaveOccurred())
			Expect(log).To(ContainSubstring("Failed to cancel job."))
		})
	})
	Describe("ListJobs", func() {
		It("should call JobsController method", func() {
			filter := weles.JobFilter{}
			sorter := weles.JobSorter{}
			paginator := weles.JobPagination{}
			list := []weles.JobInfo{
				{
					JobID: weles.JobID(3),
					Name:  "test name",
				},
			}
			info := weles.ListInfo{}
			jc.EXPECT().List(filter, sorter, paginator).Return(list, info, testErr)

			ret, retInfo, retErr := h.ListJobs(filter, sorter, paginator)

			Expect(retErr).To(Equal(testErr))
			Expect(retInfo).To(Equal(info))
			Expect(ret).To(Equal(list))
		})
	})
	Describe("Actions", func() {
		DescribeTable("Action OK",
			func(setMocks func(), cnn *chan notifier.Notification) {
				setMocks()
				*cnn <- notiOk
				eventuallyDone()
			},
			Entry("should start download when parser finished",
				func() {
					dow.EXPECT().DispatchDownloads(j).Do(setDone)
				}, &parChan),
			Entry("should request Dryad from Boruta when downloader finished",
				func() {
					bor.EXPECT().Request(j).Do(setDone)
				}, &dowChan),
			Entry("should start Job execution when Dryad is acquired from Boruta",
				func() {
					dry.EXPECT().StartJob(j).Do(setDone)
				}, &borChan),
			Entry("should complete Job after Dryad Job is done",
				func() {
					jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusCOMPLETED, "")
					bor.EXPECT().Release(j).Do(setDone)
				}, &dryChan),
		)
		DescribeTable("Action fail",
			func(cnn *chan notifier.Notification) {
				jc.EXPECT().SetStatusAndInfo(j, weles.JobStatusFAILED, testMsg)
				dry.EXPECT().CancelJob(j)
				bor.EXPECT().Release(j)
				*cnn <- notiFail
			},
			Entry("should fail when parser failed", &parChan),
			Entry("should fail when downloader failed", &dowChan),
			Entry("should fail when Boruta reports error (fail, timeout, ...)", &borChan),
			Entry("should fail when dryader fails", &dryChan),
		)
	})
})
