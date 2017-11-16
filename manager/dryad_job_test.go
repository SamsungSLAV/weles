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

package manager

import (
	"context"

	. "git.tizen.org/tools/weles"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
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
			dJob.run(ctx)
		}()
		return dJob, dJobSync
	}

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDryadJobRunner = NewMockDryadJobRunner(ctrl)
		mockOfDryadJobRunner := mockDryadJobRunner.(*MockDryadJobRunner)
		gomock.InOrder(
			mockOfDryadJobRunner.EXPECT().Deploy(),
			mockOfDryadJobRunner.EXPECT().Boot(),
			mockOfDryadJobRunner.EXPECT().Test(),
		)

		jobID = 666
		changes = make(chan DryadJobStatusChange, 6)
		dj, djSync = newMockDryadJob(jobID)
	})

	AfterEach(func() {
		Eventually(djSync).Should(BeClosed())
		ctrl.Finish()
	})

	It("should go through proper states", func() {
		states := []DryadJobStatus{DJ_NEW, DJ_DEPLOY, DJ_BOOT, DJ_TEST, DJ_OK}
		for _, state := range states {
			change := DryadJobStatusChange{jobID, state}
			Eventually(changes).Should(Receive(Equal(change)))
		}
	})

	It("should return DryadJobInfo", func() {
		info := dj.GetJobInfo()
		Expect(info.Job).To(Equal(jobID))
	})
})
