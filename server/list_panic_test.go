// Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package server

import (
	. "github.com/onsi/ginkgo"
	t "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SamsungSLAV/weles"
)

var _ = Describe("responder*200", func() {
	t.DescribeTable("responderArtifact200 should not panic on empty list",
		func(ap weles.ArtifactPaginator) {
			Expect(func() {
				responderArtifact200(weles.ListInfo{}, ap, []*weles.ArtifactInfo{}, 50)
			}).ShouldNot(Panic())
		},
		t.Entry("paginating backward", weles.ArtifactPaginator{Forward: false}),
		t.Entry("paginating forward", weles.ArtifactPaginator{Forward: true}),
	)
	t.DescribeTable("responder200 should not panic on empty list",
		func(ap weles.JobPaginator) {
			Expect(func() {
				responder200(weles.ListInfo{}, ap, []*weles.JobInfo{}, 50)
			}).ShouldNot(Panic())
		},
		t.Entry("paginating backward", weles.JobPaginator{Forward: false}),
		t.Entry("paginating forward", weles.JobPaginator{Forward: true}),
	)
})
