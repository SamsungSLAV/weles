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

package dryad

import (
	"io/ioutil"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	flyingCows = `Cows called Daisy
Are often lazy.
But cows called Brian
They be flyin'
Up in the air
And out into space
Because of the grass
And the gasses it makes!`

	flyingCowsPath = "flyingCow.txt"
)

var _ = Describe("SessionProvider", func() {
	var (
		sp      SessionProvider
		testDir string
	)

	BeforeEach(func() {
		if !accessInfoGiven {
			Skip("No valid access info to Dryad")
		}
		var err error
		testDir, err = ioutil.TempDir("", "test")
		Expect(err).ToNot(HaveOccurred())

		sp = NewSessionProvider(dryadInfo, testDir)
	})

	AfterEach(func() {
		sp.Close()

		err := os.RemoveAll(testDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should write poem to a file and read from it", func() {
		stdout, stderr, err := sp.Exec("echo", "\""+flyingCows+"\"", " > ", flyingCowsPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).To(BeEmpty())
		Expect(stderr).To(BeEmpty())

		stdout, stderr, err = sp.Exec("cat", flyingCowsPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(strings.TrimRight(string(stdout), "\n")).To(BeIdenticalTo(flyingCows))
		Expect(stderr).To(BeEmpty())

		_, _, err = sp.Exec("rm", flyingCowsPath)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should read local file from remote", func() {
		content := []byte("test file contents")
		tmpfile, err := ioutil.TempFile(testDir, "testfile")
		Expect(err).ToNot(HaveOccurred())
		_, err = tmpfile.Write(content)
		Expect(err).ToNot(HaveOccurred())
		tmpfile.Close()

		stdout, stderr, err := sp.Exec("cat", tmpfile.Name())
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).To(Equal(content))
		Expect(stderr).To(BeEmpty())
	})

	It("should not read poem from nonexistent file", func() {
		stdout, stderr, err := sp.Exec("cat", "/Ihopethispathdoesnotexist/"+flyingCowsPath+".txt")
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
