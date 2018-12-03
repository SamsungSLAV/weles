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
	"context"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/manager/dryad"
	dmock "github.com/SamsungSLAV/weles/manager/dryad/mock"
	"github.com/SamsungSLAV/weles/manager/mock"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DryadJobRunner", func() {
	var (
		mockSession *dmock.MockSessionProvider
		mockDevice  *mock.MockDeviceCommunicationProvider
		ctrl        *gomock.Controller
		djr         DryadJobRunner
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSession = dmock.NewMockSessionProvider(ctrl)
		mockDevice = mock.NewMockDeviceCommunicationProvider(ctrl)
		djr = newDryadJobRunner(context.Background(), mockSession, mockDevice, weles.Config{}, "")
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should execute the basic weles job definition", func() {
		djr = newDryadJobRunner(context.Background(), mockSession, mockDevice, basicConfig, "")
		By("Deploy")
		gomock.InOrder(
			mockSession.EXPECT().TS(),
			mockSession.EXPECT().Exec("echo", "'{\"image name_1\":\"1\",\"image_name 2\":\"2\"}'",
				">", fotaFilePath),
			mockSession.EXPECT().Exec(newFotaCmd(fotaSDCardPath, fotaFilePath,
				[]string{basicConfig.Action.Deploy.Images[0].Path,
					basicConfig.Action.Deploy.Images[1].Path}).GetCmd()),
		)

		Expect(djr.Deploy()).To(Succeed())

		By("Boot")
		gomock.InOrder(
			mockDevice.EXPECT().Boot(),
			mockDevice.EXPECT().Login(
				dryad.Credentials{
					Username: basicConfig.Action.Boot.Login,
					Password: basicConfig.Action.Boot.Password,
				}),
		)

		Expect(djr.Boot()).To(Succeed())

		By("Test")
		gomock.InOrder(
			mockDevice.EXPECT().CopyFilesTo(
				[]string{basicConfig.Action.Test.TestCases[0].TestActions[0].(weles.Push).Path},
				basicConfig.Action.Test.TestCases[0].TestActions[0].(weles.Push).Dest),
			mockDevice.EXPECT().Exec("command to be run"),
			mockDevice.EXPECT().CopyFilesFrom(
				[]string{basicConfig.Action.Test.TestCases[0].TestActions[2].(weles.Pull).Src},
				basicConfig.Action.Test.TestCases[0].TestActions[2].(weles.Pull).Path),
		)

		Expect(djr.Test()).To(Succeed())
	})
})
