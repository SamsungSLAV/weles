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

package controller

import (
	"errors"

	"github.com/SamsungSLAV/weles"
	cmock "github.com/SamsungSLAV/weles/controller/mock"
	"github.com/SamsungSLAV/weles/controller/notifier"
	"github.com/SamsungSLAV/weles/enums"
	mock "github.com/SamsungSLAV/weles/mock"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParserImpl", func() {
	var r <-chan notifier.Notification
	var jc *cmock.MockJobsController
	var am *mock.MockArtifactManager
	var yp *mock.MockParser
	var h Parser
	var ctrl *gomock.Controller
	j := weles.JobID(0xCAFE)
	goodpath := weles.ArtifactPath("/tmp/weles_test")
	badpath := weles.ArtifactPath("/such/path/does/not/exist")
	config := weles.Config{JobName: "Test name"}
	yaml := []byte("test yaml")
	err := errors.New("test error")

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		jc = cmock.NewMockJobsController(ctrl)
		am = mock.NewMockArtifactManager(ctrl)
		yp = mock.NewMockParser(ctrl)

		h = NewParser(jc, am, yp)
		r = h.Listen()
	})
	AfterEach(func() {
		ctrl.Finish()
	})
	Describe("NewParser", func() {
		It("should create a new object", func() {
			Expect(h).NotTo(BeNil())
			Expect(h.(*ParserImpl).jobs).To(Equal(jc))
			Expect(h.(*ParserImpl).artifacts).To(Equal(am))
			Expect(h.(*ParserImpl).parser).To(Equal(yp))
		})
	})
	Describe("Parse", func() {
		It("should handle job successfully", func() {
			gomock.InOrder(
				jc.EXPECT().SetStatusAndInfo(j, enums.JobStatusPARSING, ""),
				jc.EXPECT().GetYaml(j).Return(yaml, nil),
				am.EXPECT().CreateArtifact(
					weles.ArtifactDescription{
						JobID: j,
						Type:  enums.ArtifactTypeYAML,
					}).Return(goodpath, nil),
				yp.EXPECT().ParseYaml(yaml).Return(&config, nil),
				jc.EXPECT().SetConfig(j, config),
			)

			h.Parse(j)

			Eventually(r).Should(Receive(Equal(notifier.Notification{
				JobID: j,
				OK:    true,
			})))
		})
		It("should fail when unable to set config", func() {
			gomock.InOrder(
				jc.EXPECT().SetStatusAndInfo(j, enums.JobStatusPARSING, ""),
				jc.EXPECT().GetYaml(j).Return(yaml, nil),
				am.EXPECT().CreateArtifact(
					weles.ArtifactDescription{
						JobID: j,
						Type:  enums.ArtifactTypeYAML,
					}).Return(goodpath, nil),
				yp.EXPECT().ParseYaml(yaml).Return(&config, nil),
				jc.EXPECT().SetConfig(j, config).Return(err),
			)

			h.Parse(j)

			expectedNotification := notifier.Notification{
				JobID: j,
				OK:    false,
				Msg:   "Internal Weles error while setting config : " + err.Error(),
			}
			Eventually(r).Should(Receive(Equal(expectedNotification)))
		})
		It("should fail when unable to parse yaml", func() {
			gomock.InOrder(
				jc.EXPECT().SetStatusAndInfo(j, enums.JobStatusPARSING, ""),
				jc.EXPECT().GetYaml(j).Return(yaml, nil),
				am.EXPECT().CreateArtifact(
					weles.ArtifactDescription{
						JobID: j,
						Type:  enums.ArtifactTypeYAML,
					}).Return(goodpath, nil),
				yp.EXPECT().ParseYaml(yaml).Return(&weles.Config{}, err),
			)

			h.Parse(j)

			expectedNotification := notifier.Notification{
				JobID: j,
				OK:    false,
				Msg:   "Error parsing yaml file : " + err.Error(),
			}
			Eventually(r).Should(Receive(Equal(expectedNotification)))
		})
		It("should fail when unable to write yaml file", func() {
			gomock.InOrder(
				jc.EXPECT().SetStatusAndInfo(j, enums.JobStatusPARSING, ""),
				jc.EXPECT().GetYaml(j).Return(yaml, nil),
				am.EXPECT().CreateArtifact(
					weles.ArtifactDescription{
						JobID: j,
						Type:  enums.ArtifactTypeYAML,
					}).Return(badpath, nil),
			)

			h.Parse(j)

			expectedNotification := notifier.Notification{
				JobID: j,
				OK:    false,
				Msg: "Internal Weles error while saving file in ArtifactDB : " +
					"open " + string(badpath) + ": no such file or directory",
			}
			Eventually(r).Should(Receive(Equal(expectedNotification)))
		})
		It("should fail when unable to create path in ArtifactDB", func() {
			gomock.InOrder(
				jc.EXPECT().SetStatusAndInfo(j, enums.JobStatusPARSING, ""),
				jc.EXPECT().GetYaml(j).Return(yaml, nil),
				am.EXPECT().CreateArtifact(
					weles.ArtifactDescription{
						JobID: j,
						Type:  enums.ArtifactTypeYAML,
					}).Return(weles.ArtifactPath(""), err),
			)

			h.Parse(j)

			expectedNotification := notifier.Notification{
				JobID: j,
				OK:    false,
				Msg: "Internal Weles error while creating file path in ArtifactDB : " +
					err.Error(),
			}
			Eventually(r).Should(Receive(Equal(expectedNotification)))
		})
		It("should fail when unable to get yaml", func() {
			gomock.InOrder(
				jc.EXPECT().SetStatusAndInfo(j, enums.JobStatusPARSING, ""),
				jc.EXPECT().GetYaml(j).Return([]byte{}, err),
			)

			h.Parse(j)

			expectedNotification := notifier.Notification{
				JobID: j,
				OK:    false,
				Msg:   "Internal Weles error while getting yaml description : " + err.Error(),
			}
			Eventually(r).Should(Receive(Equal(expectedNotification)))
		})
		It("should fail when unable to change job status", func() {
			jc.EXPECT().SetStatusAndInfo(j, enums.JobStatusPARSING, "").Return(err)

			h.Parse(j)

			expectedNotification := notifier.Notification{
				JobID: j,
				OK:    false,
				Msg:   "Internal Weles error while changing Job status : " + err.Error(),
			}
			Eventually(r).Should(Receive(Equal(expectedNotification)))
		})
	})
})
