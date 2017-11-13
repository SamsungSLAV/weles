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
	. "git.tizen.org/tools/weles"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("dryadJob", func() {
	var changes chan DryadJobStatusChange
	var jobID JobID
	var dj *dryadJob

	BeforeEach(func() {
		jobID = 666
		changes = make(chan DryadJobStatusChange, 6)
		dj = newDryadJob(jobID, Dryad{}, changes)
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
