package server

import (
	. "github.com/onsi/ginkgo"
	t "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/SamsungSLAV/weles"
)

var _ = Describe("responder*200", func() {
	t.DescribeTable("responderArtifact200 should not panic on empty list",
		func(ap weles.ArtifactPagination) {
			Expect(func() {
				responderArtifact200(weles.ListInfo{}, ap, []*weles.ArtifactInfo{}, 50)
			}).ShouldNot(Panic())
		},
		t.Entry("paginating backward", weles.ArtifactPagination{Forward: false}),
		t.Entry("paginating forward", weles.ArtifactPagination{Forward: true}),
	)
	t.DescribeTable("responder200 should not panic on empty list",
		func(ap weles.JobPagination) {
			Expect(func() {
				responder200(weles.ListInfo{}, ap, []*weles.JobInfo{}, 50)
			}).ShouldNot(Panic())
		},
		t.Entry("paginating backward", weles.JobPagination{Forward: false}),
		t.Entry("paginating forward", weles.JobPagination{Forward: true}),
	)

})
