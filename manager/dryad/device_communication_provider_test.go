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
	"io/ioutil"
	"strconv"
	"time"

	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeviceCommunicationProvider", func() {
	var dcp DeviceCommunicationProvider
	timeout := time.Second * 30
	onelineCows := `Cows called Daisy Are often lazy. But cows called Brian They be flyin Up in the air And out into space Because of the grass And the gasses it makes!`

	var testFiles []string
	for i := 1; i < 6; i++ {
		buff := make([]byte, i*1024*1024)
		fileName := "/tmp/weles_test_file_" + strconv.FormatUint(uint64(i), 10) + ".bin"
		ioutil.WriteFile(fileName, buff, 0644)
		testFiles = append(testFiles, fileName)
	}

	BeforeEach(func() {
		if !accessInfoGiven {
			Skip("No valid access info to Dryad")
		}
		sp := NewSessionProvider(dryadInfo)
		dcp = NewDeviceCommunicationProvider(sp)
		sp.DUT()
		time.Sleep(2 * time.Second)
		for t := 0; t < 3; t++ { // Try to boot DUT 3 times. For some reason Odroid U3 won't boot every time.
			sp.PowerTick()
			for i := 0; i < 10; i++ {
				time.Sleep(10 * time.Second)
				if dcp.Login(Credentials{"", ""}) == nil {
					return
				}
			}
		}

		Skip("Target device (DUT) not available.")
	})

	AfterEach(func() {
		dcp.Close()
	})

	It("should 'login' to DUT", func() {
		err := dcp.Login(Credentials{"", ""})
		Expect(err).ToNot(HaveOccurred())
	})

	It("should list / dir", func() {
		stdout, stderr, err := dcp.Exec([]string{"ls", "-al", "/"}, timeout)
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).ToNot(BeEmpty())
		Expect(stderr).To(BeEmpty())
	})

	It("should write and read poem from a file", func() {
		By("Writing poem to a file")
		stdout, stderr, err := dcp.Exec([]string{"echo", onelineCows, " \">\" ", flyingCowsPath}, timeout)
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).To(BeEmpty())
		Expect(stderr).To(BeEmpty())

		By("Reading poem from a file")
		stdout, stderr, err = dcp.Exec([]string{"cat", flyingCowsPath}, timeout)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(stdout[:len(stdout)-2])).To(BeIdenticalTo(onelineCows))
		Expect(stderr).To(BeEmpty())
	})

	It("should not read poem from nonexistent file", func() {
		stdout, stderr, err := dcp.Exec([]string{"cat", flyingCowsPath + ".txt"}, timeout)
		Expect(err).ToNot(HaveOccurred())                   // When sdb is used no error is returned.
		Expect(stdout).To(ContainSubstring("No such file")) // And stderr are redirected to stdout.
		Expect(stderr).To(BeEmpty())
	})

	It("should transfer files to and from DUT", func() {
		os.Mkdir("/tmp/dl", 0755)

		By("Sending files to DUT")
		err := dcp.CopyFilesTo(testFiles, "/tmp/")
		Expect(err).ToNot(HaveOccurred())

		By("Receiving files from DUT")
		err = dcp.CopyFilesFrom(testFiles, "/tmp/dl/")
		Expect(err).ToNot(HaveOccurred())
		for _, path := range testFiles {
			fl, _ := os.Open(path)
			ls, _ := fl.Stat()
			fr, _ := os.Open("/tmp/dl/" + filepath.Base(path))
			rs, _ := fr.Stat()
			Expect(ls.Size()).To(BeIdenticalTo(rs.Size()))
			fl.Close()
			fr.Close()
		}
	})

	It("should not transfer files to Dryad's nonexistent directory", func() {
		err := dcp.CopyFilesTo(testFiles, "/nonexistent_dir/")
		Expect(err).To(HaveOccurred())
	})

	It("should not transfer nonexistent file from Dryad's", func() {
		err := dcp.CopyFilesFrom([]string{"/nonexistent_dir/nonexistent_file"}, "/tmp/dl/")
		Expect(err).To(HaveOccurred())
	})
})
