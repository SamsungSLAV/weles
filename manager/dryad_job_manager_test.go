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
	"sync"

	"github.com/SamsungSLAV/slav/logger"
	. "github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("DryadJobManager", func() {
	var djm DryadJobManager
	var ws *testutil.WriterString
	jobID := JobID(666)
	artifactDBPath := "/artifact/db/path"
	log := logger.NewLogger()
	stderrLog := logger.NewLogger()

	stderrLog.AddBackend("default", logger.Backend{
		Filter:     logger.NewFilterPassAll(),
		Serializer: logger.NewSerializerText(),
		Writer:     logger.NewWriterStderr(),
	})

	BeforeEach(func() {
		djm = NewDryadJobManager(artifactDBPath)

		ws = testutil.NewWriterString()
		log.AddBackend("string", logger.Backend{
			Filter:     logger.NewFilterPassAll(),
			Serializer: logger.NewSerializerText(),
			Writer:     ws,
		})
		logger.SetDefault(log)
	})
	AfterEach(func() {
		logger.SetDefault(stderrLog)
	})

	create := func() {
		err := djm.Create(jobID, Dryad{}, Config{}, nil)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() string {
			return ws.GetString()
		}).Should(ContainSubstring("Dryad job run panicked."))
	}

	It("should work for a single job", func() {
		By("create")
		create()

		By("cancel")
		err := djm.Cancel(jobID)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail to duplicate jobs", func() {
		create()

		err := djm.Create(jobID, Dryad{}, Config{}, nil)
		Expect(err).To(Equal(ErrDuplicated))

		Eventually(func() string {
			return ws.GetString()
		}).Should(ContainSubstring("Tried to create job that already exists."))
	})

	It("should fail to cancel non-existing job", func() {
		err := djm.Cancel(jobID)
		Expect(err).To(Equal(ErrNotExist))

		Eventually(func() string {
			return ws.GetString()
		}).Should(ContainSubstring("Tried to cancel nonexistent job."))
	})

	Describe("list", func() {
		var list []DryadJobInfo

		createN := func(start, end int, status DryadJobStatus) {
			dj := djm.(*DryadJobs)
			for i := start; i <= end; i++ {
				id := JobID(i)
				info := DryadJobInfo{
					Job:    id,
					Status: status,
				}
				dj.jobs[id] = &dryadJob{
					mutex: new(sync.Mutex),
					info:  info,
				}
				list = append(list, info)
			}
		}

		BeforeEach(func() {
			list = make([]DryadJobInfo, 0, 11)
			createN(0, 2, DryadJobStatusNEW)
			createN(3, 5, DryadJobStatusDEPLOY)
			createN(6, 8, DryadJobStatusBOOT)
			createN(9, 11, DryadJobStatusTEST)
		})

		It("should list created jobs", func() {
			l, err := djm.List(nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(l).To(HaveLen(len(list)))

			Consistently(func() string {
				return ws.GetString()
			}).Should(BeEmpty())
		})

		DescribeTable("list of jobs with status",
			func(start, end int, s []DryadJobStatus) {
				l, err := djm.List(&DryadJobFilter{
					Statuses: s,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(l).To(HaveLen(end - start + 1))
				for _, j := range l {
					Expect(s).To(ContainElement(j.Status))
					Expect(j.Job).To(BeNumerically(">=", start))
					Expect(j.Job).To(BeNumerically("<=", end))
				}

				Consistently(func() string {
					return ws.GetString()
				}).Should(BeEmpty())
			},
			Entry("NEW",
				0, 2, []DryadJobStatus{DryadJobStatusNEW}),
			Entry("DEPLOY",
				3, 5, []DryadJobStatus{DryadJobStatusDEPLOY}),
			Entry("BOOT",
				6, 8, []DryadJobStatus{DryadJobStatusBOOT}),
			Entry("TEST",
				9, 11, []DryadJobStatus{DryadJobStatusTEST}),
			Entry("NEW and DEPLOY",
				0, 5, []DryadJobStatus{DryadJobStatusNEW, DryadJobStatusDEPLOY}),
		)

		DescribeTable("list of jobs with id",
			func(ids []JobID, exp []int) {
				l, err := djm.List(&DryadJobFilter{
					References: ids,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(l).To(HaveLen(len(exp)))
				expected := make([]DryadJobInfo, 0)
				for _, i := range exp {
					expected = append(expected, list[i])
				}
				for _, j := range l {
					Expect(expected).To(ContainElement(Equal(j)))
				}

				Consistently(func() string {
					return ws.GetString()
				}).Should(BeEmpty())
			},
			Entry("any - 0", []JobID{0}, []int{0}),
			Entry("any - 10", []JobID{10}, []int{10}),
			Entry("out of bounds - 128", []JobID{128}, []int{}),
			Entry("many - 1 and 8", []JobID{1, 8}, []int{1, 8}),
		)

		DescribeTable("list of jobs with status and id",
			func(filter DryadJobFilter, exp []int) {
				l, err := djm.List(&filter)
				Expect(err).ToNot(HaveOccurred())
				Expect(l).To(HaveLen(len(exp)))
				expected := make([]DryadJobInfo, 0)
				for _, i := range exp {
					expected = append(expected, list[i])
				}
				for _, j := range l {
					Expect(expected).To(ContainElement(Equal(j)))
				}

				Consistently(func() string {
					return ws.GetString()
				}).Should(BeEmpty())
			},
			Entry("NEW - 2, 3", DryadJobFilter{
				References: []JobID{2, 3},
				Statuses:   []DryadJobStatus{DryadJobStatusNEW},
			}, []int{2}),
			Entry("NEW, TEST - 0, 6, 8, 10, 11", DryadJobFilter{
				References: []JobID{0, 6, 8, 10, 11},
				Statuses:   []DryadJobStatus{DryadJobStatusNEW, DryadJobStatusTEST},
			}, []int{0, 10, 11}),
		)
	})
})
