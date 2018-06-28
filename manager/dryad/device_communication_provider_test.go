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

//go:generate mockgen --package=dryad --destination=mock_session_provider_test.go git.tizen.org/tools/weles/manager/dryad SessionProvider

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeviceCommunicationProvider", func() {
	var (
		ctrl        *gomock.Controller
		mockSession *MockSessionProvider
		dcp         DeviceCommunicationProvider
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSession = NewMockSessionProvider(ctrl)
		dcp = NewDeviceCommunicationProvider(mockSession)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should call dut_login", func() {
		user := "username"
		pass := "password"
		mockSession.EXPECT().Exec("/usr/local/bin/dut_login.sh", user, pass)

		err := dcp.Login(Credentials{user, pass})
		Expect(err).ToNot(HaveOccurred())
	})

	It("should list call dut_exec", func() {
		mockSession.EXPECT().Exec("/usr/local/bin/dut_exec.sh", "ls", "-al", "/").Return([]byte("not-empty"), nil, nil)

		stdout, stderr, err := dcp.Exec("ls", "-al", "/")
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).ToNot(BeEmpty())
		Expect(stderr).To(BeEmpty())
	})

	It("should transfer files to and from DUT", func() {
		file1 := "a"
		file2 := "b"
		file3 := "c"
		files := []string{file1, file2, file3}
		target := "/tmp/dl"
		gomock.InOrder(
			mockSession.EXPECT().Exec("/usr/local/bin/dut_copyto.sh", file1, target),
			mockSession.EXPECT().Exec("/usr/local/bin/dut_copyto.sh", file2, target),
			mockSession.EXPECT().Exec("/usr/local/bin/dut_copyto.sh", file3, target),
		)

		By("Sending files to DUT")
		err := dcp.CopyFilesTo(files, target)
		Expect(err).ToNot(HaveOccurred())

		gomock.InOrder(
			mockSession.EXPECT().Exec("/usr/local/bin/dut_copyfrom.sh", file1, target),
			mockSession.EXPECT().Exec("/usr/local/bin/dut_copyfrom.sh", file2, target),
			mockSession.EXPECT().Exec("/usr/local/bin/dut_copyfrom.sh", file3, target),
		)

		By("Receiving files from DUT")
		err = dcp.CopyFilesFrom(files, target)
		Expect(err).ToNot(HaveOccurred())
	})

	//TODO: Test error paths.
})
