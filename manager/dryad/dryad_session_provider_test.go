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

package dryad

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SessionProvider", func() {
	var sp SessionProvider

	BeforeEach(func() {
		if !accessInfoGiven {
			Skip("No valid access info to Dryad")
		}
		sp = NewSessionProvider(dryadInfo)
	})

	AfterEach(func() {
		sp.Close()
	})

	It("should write poem to a file", func() {
		stdout, stderr, err := sp.Exec([]string{"echo", flyingCows, " > ", flyingCowsPath})
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).To(BeEmpty())
		Expect(stderr).To(BeEmpty())
	})

	It("should read poem from a file", func() {
		stdout, stderr, err := sp.Exec([]string{"cat", flyingCowsPath})
		Expect(err).ToNot(HaveOccurred())
		Expect(strings.TrimRight(string(stdout), "\n")).To(BeIdenticalTo(flyingCows))
		Expect(stderr).To(BeEmpty())
	})

	It("should not read poem from nonexistent file", func() {
		stdout, stderr, err := sp.Exec([]string{"cat", "/Ihopethispathdoesnotexist/" + flyingCowsPath + ".txt"})
		Expect(err).To(HaveOccurred())
		Expect(stdout).To(BeEmpty())
		Expect(stderr).ToNot(BeEmpty())
	})

	It("should switch to DUT", func() {
		Expect(sp.DUT()).ToNot(HaveOccurred())
	})

	It("should tick DUT's power supply", func() {
		Expect(sp.PowerTick()).ToNot(HaveOccurred())
	})

	It("should switch to TS", func() {
		Expect(sp.TS()).ToNot(HaveOccurred())
	})
})
