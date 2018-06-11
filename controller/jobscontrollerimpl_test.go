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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	strfmt "github.com/go-openapi/strfmt"

	"git.tizen.org/tools/weles"
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
		var initID, invalidID weles.JobID

		ipAddr := &net.IPNet{IP: net.IPv4(1, 2, 3, 4), Mask: net.IPv4Mask(5, 6, 7, 8)}
		yaml := []byte("test yaml")

		BeforeEach(func() {
			jc = NewJobsController()
			initID = jc.(*JobsControllerImpl).lastID
			invalidID = initID - 1
		})
		Describe("With Job created", func() {
			var j weles.JobID

			BeforeEach(func() {
				var err error
				j, err = jc.NewJob(yaml)
				Expect(err).NotTo(HaveOccurred())
				Expect(j).To(Equal(initID + 1))
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
					Expect(time.Time(job.Created)).To(BeTemporally(">=", before))
					Expect(time.Time(job.Created)).To(BeTemporally("<=", after))
					Expect(job.Status).To(Equal(weles.JobStatusNEW))
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
					weles.JobStatusNEW,
					weles.JobStatusPARSING,
					weles.JobStatusDOWNLOADING,
					weles.JobStatusWAITING,
					weles.JobStatusRUNNING,
					weles.JobStatusFAILED,
					weles.JobStatusCANCELED,
					weles.JobStatusCOMPLETED,
				}
				validChanges := map[weles.JobStatus](map[weles.JobStatus]bool){
					weles.JobStatusNEW: map[weles.JobStatus]bool{
						weles.JobStatusNEW:      true,
						weles.JobStatusPARSING:  true,
						weles.JobStatusFAILED:   true,
						weles.JobStatusCANCELED: true,
					},
					weles.JobStatusPARSING: map[weles.JobStatus]bool{
						weles.JobStatusPARSING:     true,
						weles.JobStatusDOWNLOADING: true,
						weles.JobStatusFAILED:      true,
						weles.JobStatusCANCELED:    true,
					},
					weles.JobStatusDOWNLOADING: map[weles.JobStatus]bool{
						weles.JobStatusDOWNLOADING: true,
						weles.JobStatusWAITING:     true,
						weles.JobStatusFAILED:      true,
						weles.JobStatusCANCELED:    true,
					},
					weles.JobStatusWAITING: map[weles.JobStatus]bool{
						weles.JobStatusWAITING:  true,
						weles.JobStatusRUNNING:  true,
						weles.JobStatusFAILED:   true,
						weles.JobStatusCANCELED: true,
					},
					weles.JobStatusRUNNING: map[weles.JobStatus]bool{
						weles.JobStatusRUNNING:   true,
						weles.JobStatusFAILED:    true,
						weles.JobStatusCANCELED:  true,
						weles.JobStatusCOMPLETED: true,
					},
					weles.JobStatusFAILED: map[weles.JobStatus]bool{
						weles.JobStatusFAILED: true,
					},
					weles.JobStatusCANCELED: map[weles.JobStatus]bool{
						weles.JobStatusCANCELED: true,
					},
					weles.JobStatusCOMPLETED: map[weles.JobStatus]bool{
						weles.JobStatusCOMPLETED: true,
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
									Expect(time.Time(job.Updated)).To(BeTemporally(">=", time.Time(oldUpdated)))
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
					Expect(time.Time(jc.(*JobsControllerImpl).jobs[j].Updated)).To(BeTemporally(">=", before))
					Expect(time.Time(jc.(*JobsControllerImpl).jobs[j].Updated)).To(BeTemporally("<=", after))
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
		Describe("List", func() {
			var elems int
			magicDate := time.Now()
			expectIDsFull := func(result []weles.JobInfo, info weles.ListInfo, expected []weles.JobID, total int, remaining int) {
				Expect(len(result)).To(Equal(len(expected)))
				for _, j := range expected {
					Expect(result).To(ContainElement(WithTransform(func(info weles.JobInfo) weles.JobID {
						return info.JobID
					}, Equal(j))))
				}
				Expect(info.TotalRecords).To(Equal(uint64(total)))
				Expect(info.RemainingRecords).To(Equal(uint64(remaining)))
			}
			expectIDs := func(result []weles.JobInfo, info weles.ListInfo, expected []weles.JobID) {
				expectIDsFull(result, info, expected, len(expected), 0)
			}
			defaultPagination := weles.JobPagination{Limit: 100}
			Describe("Filter", func() {
				jobids := []weles.JobID{}
				BeforeEach(func() {
					elems = 5
					jobids = []weles.JobID{}
					for i := 1; i <= elems; i++ {
						j, err := jc.NewJob(yaml)
						Expect(err).NotTo(HaveOccurred())
						jobids = append(jobids, j)
					}
				})
				It("should return all Jobs", func() {
					list, info, err := jc.List(weles.JobFilter{}, weles.JobSorter{}, defaultPagination)
					Expect(err).NotTo(HaveOccurred())
					expectIDs(list, info, jobids)
				})
				Describe("Created", func() {
					BeforeEach(func() {
						jc.(*JobsControllerImpl).mutex.Lock()
						defer jc.(*JobsControllerImpl).mutex.Unlock()
						for i := 0; i < elems; i++ {
							jc.(*JobsControllerImpl).jobs[jobids[i]].JobInfo.Created = strfmt.DateTime(magicDate.AddDate(i-(elems)/2, 0, 0))
						}
					})
					It("should return only jobs created after magicDate", func() {
						f := weles.JobFilter{CreatedAfter: strfmt.DateTime(magicDate), CreatedBefore: strfmt.DateTime{}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[elems/2+1:])
					})
					It("should return only jobs created before magicDate", func() {
						f := weles.JobFilter{CreatedAfter: strfmt.DateTime{}, CreatedBefore: strfmt.DateTime(magicDate)}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[:elems/2])
					})
					It("should return no jobs if created before and created after dates conflict", func() {
						f := weles.JobFilter{CreatedAfter: strfmt.DateTime(magicDate), CreatedBefore: strfmt.DateTime(magicDate)}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{})
					})
				})
				Describe("Updated", func() {
					BeforeEach(func() {
						jc.(*JobsControllerImpl).mutex.Lock()
						defer jc.(*JobsControllerImpl).mutex.Unlock()
						for i := 0; i < elems; i++ {
							jc.(*JobsControllerImpl).jobs[jobids[i]].JobInfo.Updated = strfmt.DateTime(magicDate.AddDate(i-(elems)/2, 0, 0))
						}
					})
					It("should return only jobs updated after magicDate", func() {
						f := weles.JobFilter{UpdatedAfter: strfmt.DateTime(magicDate), UpdatedBefore: strfmt.DateTime{}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[elems/2+1:])
					})
					It("should return only jobs updated before magicDate", func() {
						f := weles.JobFilter{UpdatedAfter: strfmt.DateTime{}, UpdatedBefore: strfmt.DateTime(magicDate)}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[:elems/2])
					})
					It("should return no jobs if updated before and updated after dates conflict", func() {
						f := weles.JobFilter{UpdatedAfter: strfmt.DateTime(magicDate), UpdatedBefore: strfmt.DateTime(magicDate)}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{})
					})
				})
				Describe("Info", func() {
					BeforeEach(func() {
						jc.(*JobsControllerImpl).mutex.Lock()
						defer jc.(*JobsControllerImpl).mutex.Unlock()
						jc.(*JobsControllerImpl).jobs[jobids[0]].JobInfo.Info = "Lumberjack"
						jc.(*JobsControllerImpl).jobs[jobids[1]].JobInfo.Info = "I cut down trees, I wear high heels"
						jc.(*JobsControllerImpl).jobs[jobids[2]].JobInfo.Info = "Suspenders and a bra"
						jc.(*JobsControllerImpl).jobs[jobids[3]].JobInfo.Info = "I wish I'd been a girlie"
						jc.(*JobsControllerImpl).jobs[jobids[4]].JobInfo.Info = "Just like my dear papa."
					})
					It("should return only jobs containing given substing in Info", func() {
						f := weles.JobFilter{Info: []string{"ear"}} // ear matches "wear" (line 1) and "dear" (line 4).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[4]})
					})
					It("should return only jobs containing any substing in Info", func() {
						f := weles.JobFilter{Info: []string{"ear", "I"}} // ear matches "wear" (line 1) and "dear" (line 4), "I" matches lines 1 and 3.
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[3], jobids[4]})
					})
					It("should return only jobs matching pattern", func() {
						f := weles.JobFilter{Info: []string{"a .*e"}} // matches "a girlie" (line 3).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[3]})
					})
					It("should return only jobs matching any pattern", func() {
						f := weles.JobFilter{Info: []string{"a .*e", "k$"}} // "a .*e" matches "a girlie" (line 3), "k$" matches "Lumberjack" (line 0).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[0], jobids[3]})
					})
					It("should return error if Info regexp is invalid", func() {
						f := weles.JobFilter{Info: []string{"[$$$*"}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).To(Equal(weles.ErrInvalidArgument("cannot compile regex from Info: error parsing regexp: missing closing ]: `[$$$*)`")))
						Expect(list).To(BeNil())
						Expect(info).To(BeZero())
					})
				})
				Describe("JobID", func() {
					It("should return only jobs matching JobIDs", func() {
						f := weles.JobFilter{JobID: []weles.JobID{jobids[0], jobids[2], jobids[4]}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[0], jobids[2], jobids[4]})
					})
					It("should ignore not existing JobIDs", func() {
						f := weles.JobFilter{JobID: []weles.JobID{jobids[1], invalidID}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1]})
					})
					It("should return all jobs if JobIDs slice is empty", func() {
						f := weles.JobFilter{JobID: []weles.JobID{}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids)
					})
				})
				Describe("Name", func() {
					BeforeEach(func() {
						jc.(*JobsControllerImpl).mutex.Lock()
						defer jc.(*JobsControllerImpl).mutex.Unlock()
						jc.(*JobsControllerImpl).jobs[jobids[0]].JobInfo.Name = "Lumberjack"
						jc.(*JobsControllerImpl).jobs[jobids[1]].JobInfo.Name = "I cut down trees, I wear high heels"
						jc.(*JobsControllerImpl).jobs[jobids[2]].JobInfo.Name = "Suspenders and a bra"
						jc.(*JobsControllerImpl).jobs[jobids[3]].JobInfo.Name = "I wish I'd been a girlie"
						jc.(*JobsControllerImpl).jobs[jobids[4]].JobInfo.Name = "Just like my dear papa."
					})
					It("should return only jobs containing given substing in Name", func() {
						f := weles.JobFilter{Name: []string{"ear"}} // ear matches "wear" (line 1) and "dear" (line 4).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[4]})
					})
					It("should return only jobs containing any substing in Name", func() {
						f := weles.JobFilter{Name: []string{"ear", "I"}} // ear matches "wear" (line 1) and "dear" (line 4), "I" matches lines 1 and 3.
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[3], jobids[4]})
					})
					It("should return only jobs matching pattern", func() {
						f := weles.JobFilter{Name: []string{"a .*e"}} // matches "a girlie" (line 3).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[3]})
					})
					It("should return only jobs matching any pattern", func() {
						f := weles.JobFilter{Name: []string{"a .*e", "k$"}} // "a .*e" matches "a girlie" (line 3), "k$" matches "Lumberjack" (line 0).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[0], jobids[3]})
					})
					It("should return error if Name regexp is invalid", func() {
						f := weles.JobFilter{Name: []string{"[$$$*"}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).To(Equal(weles.ErrInvalidArgument("cannot compile regex from Name: error parsing regexp: missing closing ]: `[$$$*)`")))
						Expect(list).To(BeNil())
						Expect(info).To(BeZero())
					})
				})
				Describe("Status", func() {
					BeforeEach(func() {
						jc.(*JobsControllerImpl).mutex.Lock()
						defer jc.(*JobsControllerImpl).mutex.Unlock()
						jc.(*JobsControllerImpl).jobs[jobids[0]].JobInfo.Status = weles.JobStatusNEW
						jc.(*JobsControllerImpl).jobs[jobids[1]].JobInfo.Status = weles.JobStatusPARSING
						jc.(*JobsControllerImpl).jobs[jobids[2]].JobInfo.Status = weles.JobStatusDOWNLOADING
						jc.(*JobsControllerImpl).jobs[jobids[3]].JobInfo.Status = weles.JobStatusWAITING
						jc.(*JobsControllerImpl).jobs[jobids[4]].JobInfo.Status = weles.JobStatusWAITING
					})
					It("should return all jobs if Status slice is empty", func() {
						f := weles.JobFilter{Status: []weles.JobStatus{}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids)
					})
					It("should return only jobs matching Status", func() {
						f := weles.JobFilter{Status: []weles.JobStatus{weles.JobStatusWAITING}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[3], jobids[4]})
					})
					It("should return only jobs matching any Status", func() {
						f := weles.JobFilter{Status: []weles.JobStatus{weles.JobStatusPARSING, weles.JobStatusDOWNLOADING}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[2]})
					})
					It("should ignore not set Status", func() {
						f := weles.JobFilter{Status: []weles.JobStatus{weles.JobStatusNEW, weles.JobStatus("ThereIsNoSuchStatus")}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[0]})
					})
				})
			})
		})
	})
})
