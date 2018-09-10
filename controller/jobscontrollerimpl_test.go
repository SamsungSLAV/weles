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
	. "github.com/onsi/ginkgo/extensions/table"
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
		testYaml := []byte("test yaml")

		BeforeEach(func() {
			jc = NewJobsController()
			initID = jc.(*JobsControllerImpl).lastID
			invalidID = initID - 1
		})
		Describe("With Job created", func() {
			var j weles.JobID

			BeforeEach(func() {
				var err error
				j, err = jc.NewJob(testYaml)
				Expect(err).NotTo(HaveOccurred())
				Expect(j).To(Equal(initID + 1))
			})
			Describe("NewJob", func() {
				It("should create new Job structure", func() {
					var err error
					before := time.Now()
					j, err = jc.NewJob(testYaml)
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
					Expect(job.yaml).To(Equal(testYaml))
				})
			})
			Describe("GetYaml", func() {
				It("should return proper yaml for existing job", func() {
					yaml, err := jc.GetYaml(j)
					Expect(err).NotTo(HaveOccurred())
					Expect(yaml).To(Equal(testYaml))
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
								info := fmt.Sprintf("failing to change from '%s' to '%s'",
									oldStatus, newStatus)
								By(info, func() {
									oldJob := *job
									err := jc.SetStatusAndInfo(j, newStatus, info)
									Expect(err).To(Equal(weles.ErrJobStatusChangeNotAllowed))
									Expect(job).To(Equal(&oldJob))
								})
							} else {
								info := fmt.Sprintf("changing from '%s' to '%s'",
									oldStatus, newStatus)
								oldUpdated := job.Updated
								By(info, func() {
									err := jc.SetStatusAndInfo(j, newStatus, info)
									Expect(err).NotTo(HaveOccurred())
									Expect(job.Status).To(Equal(newStatus))
									Expect(job.Info).To(Equal(info))
									Expect(time.Time(job.Updated)).To(BeTemporally(">=",
										time.Time(oldUpdated)))
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
					Expect(time.Time(jc.(*JobsControllerImpl).jobs[j].Updated)).To(
						BeTemporally(">=", before))
					Expect(time.Time(jc.(*JobsControllerImpl).jobs[j].Updated)).To(
						BeTemporally("<=", after))
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
			expectIDsFull := func(result []weles.JobInfo, info weles.ListInfo,
				expected []weles.JobID, total int, remaining int) {
				Expect(len(result)).To(Equal(len(expected)))
				for _, j := range expected {
					Expect(result).To(ContainElement(WithTransform(func(info weles.JobInfo,
					) weles.JobID {
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
						j, err := jc.NewJob(testYaml)
						Expect(err).NotTo(HaveOccurred())
						jobids = append(jobids, j)
					}
				})
				It("should return all Jobs", func() {
					list, info, err := jc.List(weles.JobFilter{}, weles.JobSorter{},
						defaultPagination)
					Expect(err).NotTo(HaveOccurred())
					expectIDs(list, info, jobids)
				})
				Describe("Created", func() {
					BeforeEach(func() {
						jc.(*JobsControllerImpl).mutex.Lock()
						defer jc.(*JobsControllerImpl).mutex.Unlock()
						for i := 0; i < elems; i++ {
							jc.(*JobsControllerImpl).jobs[jobids[i]].JobInfo.Created =
								strfmt.DateTime(magicDate.AddDate(i-(elems)/2, 0, 0))
						}
					})
					It("should return only jobs created after magicDate", func() {
						f := weles.JobFilter{
							CreatedAfter:  strfmt.DateTime(magicDate),
							CreatedBefore: strfmt.DateTime{},
						}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[elems/2+1:])
					})
					It("should return only jobs created before magicDate", func() {
						f := weles.JobFilter{
							CreatedAfter:  strfmt.DateTime{},
							CreatedBefore: strfmt.DateTime(magicDate),
						}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[:elems/2])
					})
					It("should return no jobs if created before and created after dates conflict",
						func() {
							f := weles.JobFilter{
								CreatedAfter:  strfmt.DateTime(magicDate),
								CreatedBefore: strfmt.DateTime(magicDate),
							}
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
							jc.(*JobsControllerImpl).jobs[jobids[i]].JobInfo.Updated =
								strfmt.DateTime(magicDate.AddDate(i-(elems)/2, 0, 0))
						}
					})
					It("should return only jobs updated after magicDate", func() {
						f := weles.JobFilter{
							UpdatedAfter:  strfmt.DateTime(magicDate),
							UpdatedBefore: strfmt.DateTime{},
						}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[elems/2+1:])
					})
					It("should return only jobs updated before magicDate", func() {
						f := weles.JobFilter{
							UpdatedAfter:  strfmt.DateTime{},
							UpdatedBefore: strfmt.DateTime(magicDate),
						}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, jobids[:elems/2])
					})
					It("should return no jobs if updated before and updated after dates conflict",
						func() {
							f := weles.JobFilter{
								UpdatedAfter:  strfmt.DateTime(magicDate),
								UpdatedBefore: strfmt.DateTime(magicDate),
							}
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
						jc.(*JobsControllerImpl).jobs[jobids[1]].JobInfo.Info =
							"I cut down trees, I wear high heels"
						jc.(*JobsControllerImpl).jobs[jobids[2]].JobInfo.Info =
							"Suspenders and a bra"
						jc.(*JobsControllerImpl).jobs[jobids[3]].JobInfo.Info =
							"I wish I'd been a girlie"
						jc.(*JobsControllerImpl).jobs[jobids[4]].JobInfo.Info =
							"Just like my dear papa."
					})
					It("should return only jobs containing given substing in Info", func() {
						f := weles.JobFilter{Info: []string{"ear"}}
						// ear matches "wear" (line 1) and "dear" (line 4).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[4]})
					})
					It("should return only jobs containing any substing in Info", func() {
						f := weles.JobFilter{Info: []string{"ear", "I"}}
						// ear matches "wear" (line 1) and "dear" (line 4),
						// "I" matches lines 1 and 3.
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[3], jobids[4]})
					})
					It("should return only jobs matching pattern", func() {
						f := weles.JobFilter{Info: []string{"a .*e"}}
						// matches "a girlie" (line 3).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[3]})
					})
					It("should return only jobs matching any pattern", func() {
						f := weles.JobFilter{Info: []string{"a .*e", "k$"}}
						// "a .*e" matches "a girlie" (line 3), "k$" matches "Lumberjack" (line 0).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[0], jobids[3]})
					})
					It("should return error if Info regexp is invalid", func() {
						f := weles.JobFilter{Info: []string{"[$$$*"}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).To(Equal(weles.ErrInvalidArgument(
							"cannot compile regex from Info: error parsing regexp: " +
								"missing closing ]: `[$$$*)`")))
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
						jc.(*JobsControllerImpl).jobs[jobids[1]].JobInfo.Name =
							"I cut down trees, I wear high heels"
						jc.(*JobsControllerImpl).jobs[jobids[2]].JobInfo.Name =
							"Suspenders and a bra"
						jc.(*JobsControllerImpl).jobs[jobids[3]].JobInfo.Name =
							"I wish I'd been a girlie"
						jc.(*JobsControllerImpl).jobs[jobids[4]].JobInfo.Name =
							"Just like my dear papa."
					})
					It("should return only jobs containing given substing in Name", func() {
						f := weles.JobFilter{Name: []string{"ear"}}
						// ear matches "wear" (line 1) and "dear" (line 4).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[4]})
					})
					It("should return only jobs containing any substing in Name", func() {
						f := weles.JobFilter{Name: []string{"ear", "I"}}
						// ear matches "wear" (line 1) and "dear" (line 4),
						// "I" matches lines 1 and 3.
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[3], jobids[4]})
					})
					It("should return only jobs matching pattern", func() {
						f := weles.JobFilter{Name: []string{"a .*e"}}
						// matches "a girlie" (line 3).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[3]})
					})
					It("should return only jobs matching any pattern", func() {
						f := weles.JobFilter{Name: []string{"a .*e", "k$"}}
						// "a .*e" matches "a girlie" (line 3), "k$" matches "Lumberjack" (line 0).
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[0], jobids[3]})
					})
					It("should return error if Name regexp is invalid", func() {
						f := weles.JobFilter{Name: []string{"[$$$*"}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).To(Equal(weles.ErrInvalidArgument(
							"cannot compile regex from Name: error parsing regexp: " +
								"missing closing ]: `[$$$*)`")))
						Expect(list).To(BeNil())
						Expect(info).To(BeZero())
					})
				})
				Describe("Status", func() {
					BeforeEach(func() {
						jc.(*JobsControllerImpl).mutex.Lock()
						defer jc.(*JobsControllerImpl).mutex.Unlock()
						jc.(*JobsControllerImpl).jobs[jobids[0]].JobInfo.Status = weles.JobStatusNEW
						jc.(*JobsControllerImpl).jobs[jobids[1]].JobInfo.Status =
							weles.JobStatusPARSING
						jc.(*JobsControllerImpl).jobs[jobids[2]].JobInfo.Status =
							weles.JobStatusDOWNLOADING
						jc.(*JobsControllerImpl).jobs[jobids[3]].JobInfo.Status =
							weles.JobStatusWAITING
						jc.(*JobsControllerImpl).jobs[jobids[4]].JobInfo.Status =
							weles.JobStatusWAITING
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
						f := weles.JobFilter{
							Status: []weles.JobStatus{weles.JobStatusPARSING,
								weles.JobStatusDOWNLOADING}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[1], jobids[2]})
					})
					It("should ignore not set Status", func() {
						f := weles.JobFilter{
							Status: []weles.JobStatus{
								weles.JobStatusNEW, weles.JobStatus("ThereIsNoSuchStatus")}}
						list, info, err := jc.List(f, weles.JobSorter{}, defaultPagination)
						Expect(err).NotTo(HaveOccurred())
						expectIDs(list, info, []weles.JobID{jobids[0]})
					})
				})
			})
			Describe("Sorter", func() {
				jobids := []weles.JobID{}
				checkOrder := func(
					result []weles.JobInfo, info weles.ListInfo, expected []weles.JobID,
					order []int) {
					Expect(len(result)).To(Equal(len(expected)))
					for i := range result {
						ExpectWithOffset(1, result[i].JobID).To(Equal(expected[order[i]]))
					}
					ExpectWithOffset(1, info.TotalRecords).To(Equal(uint64(len(expected))))
					ExpectWithOffset(1, info.RemainingRecords).To(BeZero())
				}
				BeforeEach(func() {
					jobids = []weles.JobID{}
					elems = 10
					for i := 1; i <= elems; i++ {
						j, err := jc.NewJob(testYaml)
						Expect(err).NotTo(HaveOccurred())
						jobids = append(jobids, j)
					}
					// Manipulate data so the sort effect will be visible.
					jc.(*JobsControllerImpl).mutex.Lock()
					defer jc.(*JobsControllerImpl).mutex.Unlock()
					jc.(*JobsControllerImpl).jobs[jobids[0]].JobInfo = weles.JobInfo{
						JobID:   jobids[0],
						Created: strfmt.DateTime(magicDate.AddDate(5, 0, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(3, 0, 0)),
						Status:  weles.JobStatusNEW,
					}
					jc.(*JobsControllerImpl).jobs[jobids[1]].JobInfo = weles.JobInfo{
						JobID:   jobids[1],
						Created: strfmt.DateTime(magicDate.AddDate(4, 0, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(1, 0, 0)),
						Status:  weles.JobStatusWAITING,
					}
					jc.(*JobsControllerImpl).jobs[jobids[2]].JobInfo = weles.JobInfo{
						JobID:   jobids[2],
						Created: strfmt.DateTime(magicDate.AddDate(2, 0, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(2, 0, 0)),
						Status:  weles.JobStatusCANCELED,
					}
					jc.(*JobsControllerImpl).jobs[jobids[3]].JobInfo = weles.JobInfo{
						JobID:   jobids[3],
						Created: strfmt.DateTime(magicDate.AddDate(3, 0, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(4, 0, 0)),
						Status:  weles.JobStatusPARSING,
					}
					jc.(*JobsControllerImpl).jobs[jobids[4]].JobInfo = weles.JobInfo{
						JobID:   jobids[4],
						Created: strfmt.DateTime(magicDate.AddDate(1, 0, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(5, 0, 0)),
						Status:  weles.JobStatusDOWNLOADING,
					}

					jc.(*JobsControllerImpl).jobs[jobids[5]].JobInfo = weles.JobInfo{
						JobID:   jobids[5],
						Created: strfmt.DateTime(magicDate.AddDate(6, 0, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(6, 0, 0)),
						Status:  weles.JobStatusRUNNING,
					}
					jc.(*JobsControllerImpl).jobs[jobids[6]].JobInfo = weles.JobInfo{
						JobID:   jobids[6],
						Created: strfmt.DateTime(magicDate.AddDate(6, 1, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(6, 1, 0)),
						Status:  weles.JobStatusFAILED,
					}
					jc.(*JobsControllerImpl).jobs[jobids[7]].JobInfo = weles.JobInfo{
						JobID:   jobids[7],
						Created: strfmt.DateTime(magicDate.AddDate(6, 2, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(6, 2, 0)),
						Status:  weles.JobStatusCOMPLETED,
					}
					jc.(*JobsControllerImpl).jobs[jobids[8]].JobInfo = weles.JobInfo{
						JobID:   jobids[8],
						Created: strfmt.DateTime(magicDate.AddDate(6, 3, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(6, 3, 0)),
						Status:  weles.JobStatus("InvalidJobStatus"),
					}
					jc.(*JobsControllerImpl).jobs[jobids[9]].JobInfo = weles.JobInfo{
						JobID:   jobids[9],
						Created: strfmt.DateTime(magicDate.AddDate(1, 0, 0)),
						Updated: strfmt.DateTime(magicDate.AddDate(5, 0, 0)),
						Status:  weles.JobStatusDOWNLOADING,
					}
				})
				DescribeTable("sorter",
					func(s weles.JobSorter, order []int) {
						list, info, err := jc.List(weles.JobFilter{}, s, defaultPagination)

						Expect(err).NotTo(HaveOccurred())
						checkOrder(list, info, jobids, order)
					},
					Entry("should sort by JobID if Sorter is empty",
						weles.JobSorter{},
						[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
					Entry("should sort by JobID if SortBy is invalid",
						weles.JobSorter{SortBy: weles.JobSortBy("InvalidSortBy")},
						[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
					Entry("should sort by JobID if SortBy is CreatedDate and SortOrder is invalid ",
						weles.JobSorter{
							SortBy:    weles.JobSortByCreatedDate,
							SortOrder: weles.SortOrder("InvalidSortOrder"),
						},
						[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
					Entry("should sort by JobID if SortBy is UpdatedDate and SortOrder is invalid ",
						weles.JobSorter{
							SortBy:    weles.JobSortByUpdatedDate,
							SortOrder: weles.SortOrder("InvalidSortOrder"),
						},
						[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
					Entry("should sort by JobID if SortBy is Status and SortOrder is invalid ",
						weles.JobSorter{
							SortBy:    weles.JobSortByJobStatus,
							SortOrder: weles.SortOrder("InvalidSortOrder"),
						},
						[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
					Entry("should sort by CreatedAsc",
						weles.JobSorter{
							SortBy:    weles.JobSortByCreatedDate,
							SortOrder: weles.SortOrderAscending,
						},
						[]int{4, 9, 2, 3, 1, 0, 5, 6, 7, 8}),
					Entry("should sort by CreatedDesc",
						weles.JobSorter{
							SortBy:    weles.JobSortByCreatedDate,
							SortOrder: weles.SortOrderDescending,
						},
						[]int{8, 7, 6, 5, 0, 1, 3, 2, 4, 9}),
					Entry("should sort by UpdatedAsc",
						weles.JobSorter{
							SortBy:    weles.JobSortByUpdatedDate,
							SortOrder: weles.SortOrderAscending,
						},
						[]int{1, 2, 0, 3, 4, 9, 5, 6, 7, 8}),
					Entry("should sort by UpdatesDesc",
						weles.JobSorter{
							SortBy:    weles.JobSortByUpdatedDate,
							SortOrder: weles.SortOrderDescending,
						},
						[]int{8, 7, 6, 5, 4, 9, 3, 0, 2, 1}),
					Entry("should sort by StatusAsc",
						weles.JobSorter{
							SortBy:    weles.JobSortByJobStatus,
							SortOrder: weles.SortOrderAscending,
						},
						[]int{8, 0, 3, 4, 9, 1, 5, 7, 6, 2}),
					Entry("should sort by StatusDesc",
						weles.JobSorter{
							SortBy:    weles.JobSortByJobStatus,
							SortOrder: weles.SortOrderDescending,
						},
						[]int{2, 6, 7, 5, 1, 4, 9, 3, 0, 8}),
				)
			})
			Describe("Paginator", func() {
				jobids := []weles.JobID{}
				evenjobids := []weles.JobID{}
				elems = 10
				for i := 0; i < elems; i++ {
					j := weles.JobID(i + 100)
					jobids = append(jobids, j)
					if i%2 == 0 {
						evenjobids = append(evenjobids, j)
					}
				}
				evenFilter := weles.JobFilter{JobID: evenjobids}
				singleFilter := weles.JobFilter{JobID: []weles.JobID{jobids[3]}}
				emptyFilter := weles.JobFilter{JobID: []weles.JobID{invalidID}}

				BeforeEach(func() {
					jc.(*JobsControllerImpl).mutex.Lock()
					defer jc.(*JobsControllerImpl).mutex.Unlock()
					for i := 0; i < elems; i++ {
						j := jobids[i]
						jc.(*JobsControllerImpl).jobs[j] = &Job{
							JobInfo: weles.JobInfo{
								JobID: j,
							},
						}
					}
				})
				DescribeTable("paginator",
					func(f weles.JobFilter, p weles.JobPagination, expected []weles.JobID,
						total, remaining int) {
						list, info, err := jc.List(f, weles.JobSorter{}, p)
						Expect(err).NotTo(HaveOccurred())
						expectIDsFull(list, info, expected, total, remaining)
					},
					Entry("should return all records if limit is 0 (pagination disabled)",
						weles.JobFilter{},
						weles.JobPagination{Limit: 0},
						jobids, 10, 0),
					Entry("should return slice of records if page is too small",
						weles.JobFilter{},
						weles.JobPagination{Limit: 3},
						jobids[:3], 10, 7),
					Entry("should return all records if page fits exactly",
						weles.JobFilter{},
						weles.JobPagination{Limit: 10},
						jobids, 10, 0),
					Entry("should return all records if page is big",
						weles.JobFilter{},
						weles.JobPagination{Limit: 100},
						jobids, 10, 0),

					Entry("when iterating forward should return all records if limit is 0"+
						" (pagination disabled)",
						weles.JobFilter{},
						weles.JobPagination{Limit: 0, JobID: jobids[3], Forward: true},
						jobids, 10, 0),
					Entry("when iterating forward should return slice of records "+
						"if page is too small",
						weles.JobFilter{},
						weles.JobPagination{Limit: 3, JobID: jobids[3], Forward: true},
						jobids[4:7], 10, 3),
					Entry("when iterating forward should return all records if page fits exactly",
						weles.JobFilter{},
						weles.JobPagination{Limit: 6, JobID: jobids[3], Forward: true},
						jobids[4:], 10, 0),
					Entry("when iterating forward should return all records if page is big",
						weles.JobFilter{},
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: true},
						jobids[4:], 10, 0),

					Entry("when iterating backwards should return all records if limit is 0 "+
						"(pagination disabled)",
						weles.JobFilter{},
						weles.JobPagination{Limit: 0, JobID: jobids[3], Forward: false},
						jobids, 10, 0),
					Entry("when iterating backwards should return slice of records if "+
						"page is too small",
						weles.JobFilter{},
						weles.JobPagination{Limit: 1, JobID: jobids[3], Forward: false},
						jobids[2:3], 10, 2),
					Entry("when iterating backwards should return all records if page fits exactly",
						weles.JobFilter{},
						weles.JobPagination{Limit: 3, JobID: jobids[3], Forward: false},
						jobids[:3], 10, 0),
					Entry("when iterating backwards should return all records if page is big",
						weles.JobFilter{},
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: false},
						jobids[:3], 10, 0),

					Entry("should return all matching records if limit is 0 (pagination disabled) "+
						"with filter",
						evenFilter,
						weles.JobPagination{Limit: 0},
						evenjobids, 5, 0),
					Entry("should return slice of records if page is too small with filter",
						evenFilter,
						weles.JobPagination{Limit: 3},
						evenjobids[:3], 5, 2),
					Entry("should return all records if page fits exactly with filter",
						evenFilter,
						weles.JobPagination{Limit: 5},
						evenjobids, 5, 0),
					Entry("should return all records if page is big with filter",
						evenFilter,
						weles.JobPagination{Limit: 100},
						evenjobids, 5, 0),

					Entry("when iterating forward should return all matching records if "+
						"limit is 0 (pagination disabled) with filter",
						evenFilter,
						weles.JobPagination{Limit: 0, JobID: jobids[3], Forward: true},
						evenjobids, 5, 0),
					Entry("when iterating forward should return slice of records if page "+
						"is too small with filter",
						evenFilter,
						weles.JobPagination{Limit: 2, JobID: jobids[3], Forward: true},
						[]weles.JobID{jobids[4], jobids[6]}, 5, 1),
					Entry("when iterating forward should return all records if page "+
						"fits exactly with filter",
						evenFilter,
						weles.JobPagination{Limit: 3, JobID: jobids[3], Forward: true},
						[]weles.JobID{jobids[4], jobids[6], jobids[8]}, 5, 0),
					Entry("when iterating forward should return all records if page "+
						"is big with filter",
						evenFilter,
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: true},
						[]weles.JobID{jobids[4], jobids[6], jobids[8]}, 5, 0),

					Entry("when iterating backwards should return all matching records if limit "+
						"is 0 (pagination disabled) with filter",
						evenFilter,
						weles.JobPagination{Limit: 0, JobID: jobids[3], Forward: false},
						evenjobids, 5, 0),
					Entry("when iterating backwards should return slice of records if page "+
						"is too small with filter",
						evenFilter,
						weles.JobPagination{Limit: 1, JobID: jobids[3], Forward: false},
						[]weles.JobID{jobids[2]}, 5, 1),
					Entry("when iterating backwards should return all records if page "+
						"fits exactly with filter",
						evenFilter,
						weles.JobPagination{Limit: 2, JobID: jobids[3], Forward: false},
						[]weles.JobID{jobids[0], jobids[2]}, 5, 0),
					Entry("when iterating backwards should return all records if page "+
						"is big with filter",
						evenFilter,
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: false},
						[]weles.JobID{jobids[0], jobids[2]}, 5, 0),

					Entry("when iterating forward should return no records with single filter",
						singleFilter,
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: true},
						[]weles.JobID{}, 1, 0),
					Entry("when iterating forward should return no records with empty filter",
						emptyFilter,
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: true},
						[]weles.JobID{}, 0, 0),
					Entry("when iterating backwards should return no records with single filter",
						singleFilter,
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: false},
						[]weles.JobID{}, 1, 0),
					Entry("when iterating backwards should return no records with empty filter",
						emptyFilter,
						weles.JobPagination{Limit: 100, JobID: jobids[3], Forward: false},
						[]weles.JobID{}, 0, 0),
				)
				It("when iterating forward should return error when JobID does not exist",
					func() {
						list, info, err := jc.List(weles.JobFilter{}, weles.JobSorter{},
							weles.JobPagination{Limit: 100, JobID: invalidID, Forward: true})
						Expect(err).To(Equal(weles.ErrInvalidArgument(
							fmt.Sprintf("JobID: %d not found", invalidID))))
						Expect(list).To(BeNil())
						Expect(info).To(BeZero())
					})
				It("when iterating backwards should return error when JobID does not exist",
					func() {
						list, info, err := jc.List(weles.JobFilter{}, weles.JobSorter{},
							weles.JobPagination{Limit: 100, JobID: invalidID, Forward: false})
						Expect(err).To(Equal(weles.ErrInvalidArgument(
							fmt.Sprintf("JobID: %d not found", invalidID))))
						Expect(list).To(BeNil())
						Expect(info).To(BeZero())
					})
			})
		})
	})
})
