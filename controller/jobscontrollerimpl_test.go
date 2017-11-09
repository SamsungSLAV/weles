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
	"fmt"
	"net"
	"time"

	"git.tizen.org/tools/weles"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobsControllerImpl", func() {
	Describe("NewJobsController", func() {
		It("should create a new object", func() {
			before := time.Now()
			jc := NewJobsController()
			after := time.Now()

			Expect(jc).NotTo(BeNil())
			Expect(jc.(*JobsControllerImpl).mutex).NotTo(BeNil())
			Expect(jc.(*JobsControllerImpl).jobs).NotTo(BeNil())
			Expect(jc.(*JobsControllerImpl).jobs).To(BeEmpty())
			Expect(jc.(*JobsControllerImpl).lastID).To(BeNumerically(">=", before.Unix()))
			Expect(jc.(*JobsControllerImpl).lastID).To(BeNumerically("<=", after.Unix()))
		})
	})
	Describe("With JobsController initialized", func() {
		var jc JobsController
		var initID, j weles.JobID
		ipAddr := &net.IPNet{IP: net.IPv4(1, 2, 3, 4), Mask: net.IPv4Mask(5, 6, 7, 8)}
		yaml := []byte("test yaml")
		var invalidID weles.JobID

		BeforeEach(func() {
			jc = NewJobsController()
			initID = jc.(*JobsControllerImpl).lastID

			var err error
			j, err = jc.NewJob(yaml)
			Expect(err).NotTo(HaveOccurred())
			Expect(j).To(Equal(initID + 1))

			invalidID = initID - 1
		})
		Describe("NewJob", func() {
			It("should create new Job structure", func() {
				before := time.Now()
				j, err := jc.NewJob(yaml)
				after := time.Now()

				Expect(err).NotTo(HaveOccurred())
				Expect(j).To(Equal(initID + 2))

				Expect(jc.(*JobsControllerImpl).lastID).To(Equal(j))
				Expect(len(jc.(*JobsControllerImpl).jobs)).To(Equal(2))

				job, ok := jc.(*JobsControllerImpl).jobs[j]
				Expect(ok).To(BeTrue())
				Expect(job.JobID).To(Equal(j))
				Expect(job.Created).To(Equal(job.Updated))
				Expect(job.Created).To(BeTemporally(">=", before))
				Expect(job.Created).To(BeTemporally("<=", after))
				Expect(job.Status).To(Equal(weles.JOB_NEW))
				Expect(job.yaml).To(Equal(yaml))
			})
		})
		Describe("GetYaml", func() {
			It("should return proper yaml for existing job", func() {
				retyaml, err := jc.GetYaml(j)
				Expect(err).NotTo(HaveOccurred())
				Expect(retyaml).To(Equal(yaml))
			})
			It("should return error for not existing job", func() {
				yaml, err := jc.GetYaml(invalidID)
				Expect(err).To(Equal(weles.ErrJobNotFound))
				Expect(yaml).To(BeZero())
			})
		})
		Describe("SetStatus", func() {
			allStatus := []weles.JobStatus{
				weles.JOB_NEW,
				weles.JOB_PARSING,
				weles.JOB_DOWNLOADING,
				weles.JOB_WAITING,
				weles.JOB_RUNNING,
				weles.JOB_FAILED,
				weles.JOB_CANCELED,
				weles.JOB_COMPLETED,
			}
			validChanges := map[weles.JobStatus](map[weles.JobStatus]bool){
				weles.JOB_NEW: map[weles.JobStatus]bool{
					weles.JOB_NEW:      true,
					weles.JOB_PARSING:  true,
					weles.JOB_FAILED:   true,
					weles.JOB_CANCELED: true,
				},
				weles.JOB_PARSING: map[weles.JobStatus]bool{
					weles.JOB_PARSING:     true,
					weles.JOB_DOWNLOADING: true,
					weles.JOB_FAILED:      true,
					weles.JOB_CANCELED:    true,
				},
				weles.JOB_DOWNLOADING: map[weles.JobStatus]bool{
					weles.JOB_DOWNLOADING: true,
					weles.JOB_WAITING:     true,
					weles.JOB_FAILED:      true,
					weles.JOB_CANCELED:    true,
				},
				weles.JOB_WAITING: map[weles.JobStatus]bool{
					weles.JOB_WAITING:  true,
					weles.JOB_RUNNING:  true,
					weles.JOB_FAILED:   true,
					weles.JOB_CANCELED: true,
				},
				weles.JOB_RUNNING: map[weles.JobStatus]bool{
					weles.JOB_RUNNING:   true,
					weles.JOB_FAILED:    true,
					weles.JOB_CANCELED:  true,
					weles.JOB_COMPLETED: true,
				},
				weles.JOB_FAILED: map[weles.JobStatus]bool{
					weles.JOB_FAILED: true,
				},
				weles.JOB_CANCELED: map[weles.JobStatus]bool{
					weles.JOB_CANCELED: true,
				},
				weles.JOB_COMPLETED: map[weles.JobStatus]bool{
					weles.JOB_COMPLETED: true,
				},
			}
			It("should return error for not existing job", func() {
				for _, status := range allStatus {
					err := jc.SetStatusAndInfo(invalidID, status, "test info")
					Expect(err).To(Equal(weles.ErrJobNotFound))
				}
			})
			It("should work to change status only for valid transitions", func() {
				job := jc.(*JobsControllerImpl).jobs[j]
				for _, oldStatus := range allStatus {
					for _, newStatus := range allStatus {
						job.Status = oldStatus
						if _, ok := validChanges[oldStatus][newStatus]; !ok {
							info := fmt.Sprintf("failing to change from '%s' to '%s'", oldStatus, newStatus)
							By(info, func() {
								oldJob := *job
								err := jc.SetStatusAndInfo(j, newStatus, info)
								Expect(err).To(Equal(weles.ErrJobStatusChangeNotAllowed))
								Expect(job).To(Equal(&oldJob))
							})
						} else {
							info := fmt.Sprintf("changing from '%s' to '%s'", oldStatus, newStatus)
							oldUpdated := job.Updated
							By(info, func() {
								err := jc.SetStatusAndInfo(j, newStatus, info)
								Expect(err).NotTo(HaveOccurred())
								Expect(job.Status).To(Equal(newStatus))
								Expect(job.Info).To(Equal(info))
								Expect(job.Updated).To(BeTemporally(">=", oldUpdated))
							})
						}
					}
				}
			})
		})
		Describe("SetConfig", func() {
			It("should set config for existing job", func() {
				config := weles.Config{JobName: "Test Job"}
				before := time.Now()
				err := jc.SetConfig(j, config)
				after := time.Now()
				Expect(err).NotTo(HaveOccurred())

				Expect(jc.(*JobsControllerImpl).jobs[j].config).To(Equal(config))
				Expect(jc.(*JobsControllerImpl).jobs[j].Updated).To(BeTemporally(">=", before))
				Expect(jc.(*JobsControllerImpl).jobs[j].Updated).To(BeTemporally("<=", after))
			})
			It("should return error for not existing job", func() {
				config := weles.Config{JobName: "Test Job"}
				err := jc.SetConfig(invalidID, config)
				Expect(err).To(Equal(weles.ErrJobNotFound))
			})
		})
		Describe("GetConfig", func() {
			It("should return proper config for existing job", func() {
				expectedConfig := weles.Config{JobName: "Test config"}
				err := jc.SetConfig(j, expectedConfig)
				Expect(err).NotTo(HaveOccurred())

				config, err := jc.GetConfig(j)
				Expect(err).NotTo(HaveOccurred())
				Expect(config).To(Equal(expectedConfig))
			})
			It("should return error for not existing job", func() {
				config, err := jc.GetConfig(invalidID)
				Expect(err).To(Equal(weles.ErrJobNotFound))
				Expect(config).To(BeZero())
			})
		})

		Describe("SetDryad", func() {
			It("should set Dryad for existing job", func() {
				dryad := weles.Dryad{Addr: ipAddr}
				err := jc.SetDryad(j, dryad)
				Expect(err).NotTo(HaveOccurred())

				Expect(jc.(*JobsControllerImpl).jobs[j].dryad).To(Equal(dryad))
			})
			It("should return error for not existing job", func() {
				dryad := weles.Dryad{Addr: ipAddr}
				err := jc.SetDryad(invalidID, dryad)
				Expect(err).To(Equal(weles.ErrJobNotFound))
			})
		})

		Describe("GetDryad", func() {
			It("should return proper Dryad structure for existing job", func() {
				expectedDryad := weles.Dryad{Addr: ipAddr}
				err := jc.SetDryad(j, expectedDryad)
				Expect(err).NotTo(HaveOccurred())

				dryad, err := jc.GetDryad(j)
				Expect(err).NotTo(HaveOccurred())
				Expect(dryad).To(Equal(expectedDryad))
			})
			It("should return error for not existing job", func() {
				dryad, err := jc.GetDryad(invalidID)
				Expect(err).To(Equal(weles.ErrJobNotFound))
				Expect(dryad).To(BeZero())
			})
		})
	})
})
