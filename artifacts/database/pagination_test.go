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

package database

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"git.tizen.org/tools/weles"
)

//TODO: add some consts and limit 'magic' to multipliers

var _ = Describe("ArtifactDB pagination", func() {

	emptyF := weles.ArtifactFilter{}
	//below struct is not empty but filled with default values passed from server.
	descS := weles.ArtifactSorter{
		SortOrder: weles.SortOrderDescending,
		SortBy:    weles.ArtifactSortByID,
	}
	ascS := weles.ArtifactSorter{
		SortOrder: weles.SortOrderAscending,
		SortBy:    weles.ArtifactSortByID,
	}
	Context("Database is filled with 100 records", func() {
		DescribeTable("paginating through artifacts",
			func(
				filter weles.ArtifactFilter,
				paginator weles.ArtifactPagination,
				sorter weles.ArtifactSorter,
				expectedResponseLength, expectedTotalRecords, expectedRemainingRecords int) {
				result, list, err := silverHoneybadger.Filter(filter, sorter, paginator)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(result)).To(BeEquivalentTo(expectedResponseLength))
				Expect(list.TotalRecords).To(BeEquivalentTo(expectedTotalRecords))
				Expect(list.RemainingRecords).To(BeEquivalentTo(expectedRemainingRecords))
			},
			Entry("p1/1, no limit, forward, sort desc", emptyF,
				weles.ArtifactPagination{Forward: true},
				descS, 100, 100, 0),
			Entry("p1/4, limit 30, forward, sort desc", emptyF,
				weles.ArtifactPagination{Limit: 30, Forward: true},
				descS, 30, 100, 70),
			Entry("p2/4, limit 30, forward, sort desc", emptyF,
				weles.ArtifactPagination{ID: 71, Limit: 30, Forward: true},
				descS, 30, 100, 40),
			Entry("p3/4, limit 30, forward, sort desc", emptyF,
				weles.ArtifactPagination{ID: 41, Limit: 30, Forward: true},
				descS, 30, 100, 10),
			Entry("p4/4, limit 30, forward, sort desc", emptyF,
				weles.ArtifactPagination{ID: 11, Limit: 30, Forward: true},
				descS, 10, 100, 0),
			Entry("p3/4, limit 30, backward, sort desc", emptyF,
				weles.ArtifactPagination{ID: 10, Limit: 30, Forward: false},
				descS, 30, 100, 60),
			Entry("p2/4, limit 30, backward, sort desc", emptyF,
				weles.ArtifactPagination{ID: 40, Limit: 30, Forward: false},
				descS, 30, 100, 30),
			Entry("p1/4, limit 30, backward, sort desc", emptyF,
				weles.ArtifactPagination{ID: 70, Limit: 30, Forward: false},
				descS, 30, 100, 0),
			Entry("p1/1, no limit, forward, sort asc", emptyF,
				weles.ArtifactPagination{Forward: true},
				ascS, 100, 100, 0),
			Entry("1/4, no limit, forward, sort asc", emptyF,
				weles.ArtifactPagination{Limit: 30, Forward: true},
				ascS, 30, 100, 70),
			Entry("p2/4, no limit, forward, sort asc", emptyF,
				weles.ArtifactPagination{ID: 30, Limit: 30, Forward: true},
				ascS, 30, 100, 40),
			Entry("p3/4, no limit, forward, sort asc", emptyF,
				weles.ArtifactPagination{ID: 60, Limit: 30, Forward: true},
				ascS, 30, 100, 10),
			Entry("p4/4, no limit, forward, sort asc", emptyF,
				weles.ArtifactPagination{ID: 90, Limit: 30, Forward: true},
				ascS, 10, 100, 0),
		)
	})

})
