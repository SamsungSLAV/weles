/*
 *  Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
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

package testutil

import (
	"github.com/SamsungSLAV/slav/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WriterString", func() {
	const (
		testMsg  = "testMessage"
		anyLevel = logger.EmergLevel
	)

	var ws *WriterString

	BeforeEach(func() {
		ws = NewWriterString()
	})
	Describe("NewWriterString", func() {
		It("should create a new empty object", func() {
			Expect(ws).NotTo(BeNil())
		})
	})
	Describe("Write", func() {
		It("should write to string", func() {
			n, err := ws.Write(anyLevel, []byte(testMsg))
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(SatisfyAll(Equal(len(testMsg)+1), Equal(ws.b.Len())))
			Expect(ws.b.String()).To(Equal(testMsg + "\n"))
		})
	})
	Describe("GetString", func() {
		It("should return no content before string is built", func() {
			s := ws.GetString()
			Expect(s).To(BeEmpty())
		})
		It("should return contents stored in built string", func() {
			_, err := ws.Write(anyLevel, []byte(testMsg))
			Expect(err).NotTo(HaveOccurred())

			s := ws.GetString()
			Expect(s).To(Equal(testMsg + "\n"))
		})
	})
})
