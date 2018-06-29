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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DryadJobRunnerFota", func() {
	var (
		sdcard  = "/dev/path to sdcard"
		mapping = "path to mapping"
		md5     = "url to md5sums"
		urls    = []string{"https://some secure server", "http://some not so secure one"}
		f       *fotaCmd
	)

	BeforeEach(func() {
		f = newFotaCmd(sdcard, mapping, urls)
	})

	checkOrder := func(cmd []string) {
		for i, part := range cmd {
			switch part {
			case "fota":
				Expect(i).To(Equal(0))
			case "-map":
				Expect(cmd[i+1]).To(Equal(mapping))
			case "-md5":
				Expect(cmd[i+1]).To(Equal(md5))
			case "-card":
				Expect(cmd[i+1]).To(Equal(sdcard))
			}
		}
	}

	It("should work for some arguments", func() {
		cmd := f.GetCmd()
		Expect(cmd).To(ConsistOf(
			fotaCmdPath,
			"-card", sdcard,
			"-map", mapping,
			urls[0], urls[1],
		))
		checkOrder(cmd)
	})

	It("should add md5sum argument", func() {
		f.SetMD5(md5)
		cmd := f.GetCmd()
		Expect(cmd).To(ConsistOf(
			fotaCmdPath,
			"-card", sdcard,
			"-map", mapping,
			"-md5", md5,
			urls[0], urls[1],
		))
		checkOrder(cmd)
	})
})
