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

//go:generate mockgen -package=manager -destination=mock_dryad_test.go git.tizen.org/tools/weles/manager/dryad SessionProvider,DeviceCommunicationProvider

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DryadJobRunner", func() {
	var (
		mockSession *MockSessionProvider
		mockDevice  *MockDeviceCommunicationProvider
		ctrl        *gomock.Controller
		djr         DryadJobRunner
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSession = NewMockSessionProvider(ctrl)
		mockDevice = NewMockDeviceCommunicationProvider(ctrl)
		djr = newDryadJobRunner(context.Background(), mockSession, mockDevice)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Deploy", func() {
		var tsCall *gomock.Call

		BeforeEach(func() {
			tsCall = mockSession.EXPECT().TS().Times(1)
		})

		It("should switch to TS", func() {
			err := djr.Deploy()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail if TS fails", func() {
			tsCall.Return(errors.New("TS failed"))

			err := djr.Deploy()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Boot", func() {
		var dutCall *gomock.Call

		BeforeEach(func() {
			dutCall = mockSession.EXPECT().DUT().Times(1)
		})

		It("should switch to DUT", func() {
			err := djr.Boot()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail if DUT fails", func() {
			dutCall.Return(errors.New("DUT failed"))

			err := djr.Boot()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Test", func() {
		var execCall *gomock.Call

		BeforeEach(func() {
			execCall = mockSession.EXPECT().Exec([]string{"echo", "healthcheck"})
		})

		It("should exec echo healthcheck", func() {
			err := djr.Test()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail if Exec fails", func() {
			execCall.Return(nil, nil, errors.New("exec failed"))

			err := djr.Test()
			Expect(err).To(HaveOccurred())
		})
	})
})
