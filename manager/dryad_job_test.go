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

package manager

import (
	"context"
	"errors"

	. "git.tizen.org/tools/weles"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("dryadJob", func() {
	var (
		changes            chan DryadJobStatusChange
		jobID              JobID
		dj                 *dryadJob
		djSync             chan struct{}
		ctrl               *gomock.Controller
		mockDryadJobRunner DryadJobRunner
		deploy, boot, test *gomock.Call
		cancel             context.CancelFunc
	)

	newMockDryadJob := func(job JobID) (*dryadJob, chan struct{}) {
		dJobSync := make(chan struct{})
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		dJob := newDryadJobWithCancel(job, changes, mockDryadJobRunner, cancel)
		go func() {
			defer close(dJobSync)
			defer GinkgoRecover()
			<-dJobSync
			dJob.run(ctx)
		}()
		return dJob, dJobSync
	}

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDryadJobRunner = NewMockDryadJobRunner(ctrl)
		mockOfDryadJobRunner := mockDryadJobRunner.(*MockDryadJobRunner)
		deploy = mockOfDryadJobRunner.EXPECT().Deploy().Times(1)
		boot = mockOfDryadJobRunner.EXPECT().Boot().Times(1).After(deploy)
		test = mockOfDryadJobRunner.EXPECT().Test().Times(1).After(boot)

		jobID = 666
		changes = make(chan DryadJobStatusChange, 6)
		dj, djSync = newMockDryadJob(jobID)
	})

	AfterEach(func() {
		Eventually(djSync).Should(BeClosed())
		ctrl.Finish()
	})

	It("should go through proper states", func() {
		djSync <- struct{}{}
		states := []DryadJobStatus{DJ_NEW, DJ_DEPLOY, DJ_BOOT, DJ_TEST, DJ_OK}
		for _, state := range states {
			change := DryadJobStatusChange{Job: jobID, Status: state}
			Eventually(changes).Should(Receive(Equal(change)))
		}
	})

	registerPhaseErr := func(c *gomock.Call, err error, times int) {
		c.Return(err).Times(times)
	}
	registerErr := func(deployErr, bootErr, testErr error) []DryadJobStatus {
		ret := []DryadJobStatus{DJ_NEW, DJ_DEPLOY}
		switch {
		case deployErr != nil:
			registerPhaseErr(deploy, deployErr, 1)
			registerPhaseErr(boot, bootErr, 0)
			registerPhaseErr(test, testErr, 0)
		case bootErr != nil:
			registerPhaseErr(deploy, deployErr, 1)
			registerPhaseErr(boot, bootErr, 1)
			ret = append(ret, DJ_BOOT)
			registerPhaseErr(test, testErr, 0)
		case testErr != nil:
			registerPhaseErr(deploy, deployErr, 1)
			registerPhaseErr(boot, bootErr, 1)
			ret = append(ret, DJ_BOOT)
			registerPhaseErr(test, testErr, 1)
			ret = append(ret, DJ_TEST)
		}
		ret = append(ret, DJ_FAIL)
		djSync <- struct{}{}
		return ret
	}
	DescribeTable("fail when one of the stages does",
		func(f func() []DryadJobStatus) {
			states := f()
			for _, state := range states {
				change := DryadJobStatusChange{Job: jobID, Status: state}
				Eventually(changes).Should(Receive(Equal(change)))
			}
		},
		Entry("after deploy", func() []DryadJobStatus {
			return registerErr(errors.New("deploy failed"), nil, nil)
		}),
		Entry("after boot", func() []DryadJobStatus {
			return registerErr(nil, errors.New("boot failed"), nil)
		}),
		Entry("after test", func() []DryadJobStatus {
			return registerErr(nil, nil, errors.New("test failed"))
		}),
	)

	DescribeTable("should recover a panic and go to failed state",
		func(f func()) {
			f()
			djSync <- struct{}{}
			fail := DryadJobStatusChange{Job: jobID, Status: DJ_FAIL}
			Eventually(changes).Should(Receive(Equal(fail)))
		},
		Entry("deploy", func() {
			deploy.Do(func() { panic("deploy") })
			boot.Times(0)
			test.Times(0)
		}),
		Entry("boot", func() {
			boot.Do(func() { panic("boot") })
			test.Times(0)
		}),
		Entry("test", func() {
			test.Do(func() { panic("test") })
		}),
	)

	It("should return DryadJobInfo", func() {
		djSync <- struct{}{}
		info := dj.GetJobInfo()
		Expect(info.Job).To(Equal(jobID))
	})
})
