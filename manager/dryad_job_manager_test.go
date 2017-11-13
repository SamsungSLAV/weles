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

package manager_test

import (
	. "git.tizen.org/tools/weles"
	. "git.tizen.org/tools/weles/manager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DryadJobManager", func() {
	var djm DryadJobManager
	jobID := JobID(666)

	BeforeEach(func() {
		djm = NewDryadJobManager()
	})

	create := func() {
		err := djm.Create(jobID, Dryad{}, nil)
		Expect(err).ToNot(HaveOccurred())
	}

	It("should work for a single job", func() {
		By("create")
		create()

		By("cancel")
		err := djm.Cancel(jobID)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail to duplicate jobs", func() {
		create()

		err := djm.Create(jobID, Dryad{}, nil)
		Expect(err).To(Equal(ErrDuplicated))
	})

	It("should fail to cancel non-existing job", func() {
		err := djm.Cancel(jobID)
		Expect(err).To(Equal(ErrNotExist))
	})

	It("should list created jobs", func() {
		create()

		list, err := djm.List(nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(list).To(HaveLen(1))
	})
})
