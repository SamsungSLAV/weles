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

package notifier

import (
	"fmt"

	. "github.com/SamsungSLAV/weles"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Impl", func() {
	var h Notifier

	BeforeEach(func() {
		h = NewNotifier()
	})
	Describe("NewNotifier", func() {
		It("should create a new object", func() {
			Expect(h).NotTo(BeNil())
			Expect(h.(*Impl).channel).NotTo(BeNil())
		})
	})
	Describe("Listen", func() {
		It("should return read only channel", func() {
			c := h.Listen()

			Expect(c).To(Equal((<-chan Notification)(h.(*Impl).channel)))
		})
	})
	Describe("SendFail", func() {
		It("should send proper notifications to channel", func() {
			for i := 1; i <= 5; i++ {
				h.SendFail(JobID(i), fmt.Sprintf("message number %d", i))
			}

			c := h.Listen()
			for i := 1; i <= 5; i++ {
				Eventually(c).Should(Receive(Equal(Notification{
					JobID: JobID(i),
					OK:    false,
					Msg:   fmt.Sprintf("message number %d", i),
				})))
			}
		})
	})
	Describe("SendOK", func() {
		It("should send proper notifications to channel", func() {
			for i := 1; i <= 5; i++ {
				h.SendOK(JobID(i))
			}

			c := h.Listen()
			for i := 1; i <= 5; i++ {
				Eventually(c).Should(Receive(Equal(Notification{
					JobID: JobID(i),
					OK:    true,
				})))
			}
		})
	})
})
