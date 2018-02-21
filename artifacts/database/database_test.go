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

// Package database is responsible for Weles system's job artifact storage.
package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-openapi/strfmt"

	"git.tizen.org/tools/weles"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArtifactDB", func() {
	var (
		job           weles.JobID          = 58008
		invalidJob    weles.JobID          = 1
		invalidPath   weles.ArtifactPath   = "invalidPath"
		invalidStatus weles.ArtifactStatus = "invalidStatus"
		invalidType   weles.ArtifactType   = "invalidType"
		invalidAlias  weles.ArtifactAlias  = "invalidAlias"
		goldenUnicorn ArtifactDB
		tmpDir        string

		artifact = weles.ArtifactInfo{
			weles.ArtifactDescription{
				"some alias",
				job,
				weles.ArtifactTypeIMAGE,
				"http://example.com",
			},
			"path1",
			weles.ArtifactStatusPENDING,
			strfmt.DateTime(time.Now().Round(time.Millisecond).UTC()),
		}

		aImageReady = weles.ArtifactInfo{
			weles.ArtifactDescription{
				"other alias",
				job + 1,
				weles.ArtifactTypeIMAGE,
				"http://example.com/1",
			},
			"path2",
			weles.ArtifactStatusREADY,
			strfmt.DateTime(time.Now().Round(time.Millisecond).UTC()),
		}

		aYamlFailed = weles.ArtifactInfo{
			weles.ArtifactDescription{
				"other alias",
				job + 1,
				weles.ArtifactTypeYAML,
				"http://example.com/2",
			},
			"path3",
			weles.ArtifactStatusFAILED,
			strfmt.DateTime(time.Now().Round(time.Millisecond).UTC()),
		}

		aTestFailed = weles.ArtifactInfo{
			weles.ArtifactDescription{
				"alias",
				job + 2,
				weles.ArtifactTypeTEST,
				"http://example.com/2",
			},
			"path4",
			weles.ArtifactStatusFAILED,
			strfmt.DateTime(time.Unix(3000, 60).Round(time.Millisecond).UTC()),
		}

		testArtifacts = []weles.ArtifactInfo{artifact, aImageReady, aYamlFailed, aTestFailed}

		oneJobFilter = weles.ArtifactFilter{
			JobID:  []weles.JobID{artifact.JobID},
			Alias:  nil,
			Status: nil,
			Type:   nil}
		twoJobsFilter = weles.ArtifactFilter{
			JobID:  []weles.JobID{artifact.JobID, aImageReady.JobID},
			Alias:  nil,
			Status: nil,
			Type:   nil}
		noJobFilter = weles.ArtifactFilter{
			JobID:  []weles.JobID{invalidJob},
			Alias:  nil,
			Status: nil,
			Type:   nil}

		oneTypeFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   []weles.ArtifactType{aYamlFailed.Type},
			Status: nil,
			Alias:  nil}

		twoTypesFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   []weles.ArtifactType{aYamlFailed.Type, aTestFailed.Type},
			Status: nil,
			Alias:  nil}

		noTypeFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   []weles.ArtifactType{invalidType},
			Status: nil,
			Alias:  nil}

		oneStatusFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: []weles.ArtifactStatus{artifact.Status},
			Alias:  nil}

		twoStatusFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: []weles.ArtifactStatus{artifact.Status, aYamlFailed.Status},
			Alias:  nil}

		noStatusFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: []weles.ArtifactStatus{invalidStatus},
			Alias:  nil}

		oneAliasFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: nil,
			Alias:  []weles.ArtifactAlias{artifact.Alias}}

		twoAliasFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: nil,
			Alias:  []weles.ArtifactAlias{artifact.Alias, aImageReady.Alias}}

		noAliasFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: nil,
			Alias:  []weles.ArtifactAlias{invalidAlias}}

		fullFilter = weles.ArtifactFilter{
			JobID:  twoJobsFilter.JobID,
			Type:   twoTypesFilter.Type,
			Status: twoStatusFilter.Status,
			Alias:  twoAliasFilter.Alias}

		noMatchFilter = weles.ArtifactFilter{
			JobID:  oneJobFilter.JobID,
			Type:   oneTypeFilter.Type,
			Status: nil,
			Alias:  nil}

		emptyFilter = weles.ArtifactFilter{}
	)

	jobsInDB := func(job weles.JobID) int64 {
		n, err := goldenUnicorn.dbmap.SelectInt(`SELECT COUNT(*)
 		FROM artifacts
 		WHERE JobID = ?`, job)
		Expect(err).ToNot(HaveOccurred())
		return n
	}

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "weles-")
		Expect(err).ToNot(HaveOccurred())
		err = goldenUnicorn.Open(filepath.Join(tmpDir, "test.db"))
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := goldenUnicorn.Close()
		Expect(err).ToNot(HaveOccurred())
		err = os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should open new database, with artifact table", func() {
		n, err := goldenUnicorn.dbmap.SelectInt(`SELECT COUNT(*)
 		FROM sqlite_master
 		WHERE name = 'artifacts'
 		AND type = 'table'`)
		Expect(err).ToNot(HaveOccurred())
		Expect(n).To(BeNumerically("==", 1))
	})

	It("should fail to open database on invalid path", func() {
		// sql.Open only validates arguments.
		// db.Ping must be called to check the connection.
		invalidDatabasePath := filepath.Join(tmpDir, "invalid", "test.db")
		err := goldenUnicorn.Open(invalidDatabasePath)
		Expect(err).To(HaveOccurred())
		Expect(invalidDatabasePath).ToNot(BeAnExistingFile())
	})

	It("should insert new artifact to database", func() {
		Expect(jobsInDB(job)).To(BeNumerically("==", 0))

		err := goldenUnicorn.InsertArtifactInfo(&artifact)
		Expect(err).ToNot(HaveOccurred())

		Expect(jobsInDB(artifact.JobID)).To(BeNumerically("==", 1))
	})

	Describe("SelectPath", func() {

		BeforeEach(func() {
			err := goldenUnicorn.InsertArtifactInfo(&artifact)
			Expect(err).ToNot(HaveOccurred())

			Expect(jobsInDB(artifact.JobID)).To(BeNumerically("==", 1))
		})

		DescribeTable("database selectpath",
			func(path weles.ArtifactPath, expectedErr error, expectedArtifact weles.ArtifactInfo) {
				result, err := goldenUnicorn.SelectPath(path)

				if expectedErr != nil {
					Expect(err).To(Equal(expectedErr))
				} else {
					Expect(err).ToNot(HaveOccurred())
				}
				Expect(result).To(Equal(expectedArtifact))
			},
			Entry("retrieve artifact based on path", artifact.Path, nil, artifact),
			Entry("retrieve artifact based on invalid path", invalidPath, sql.ErrNoRows, weles.ArtifactInfo{}),
		)
	})

	Describe("Select", func() {

		BeforeEach(func() {
			for _, a := range testArtifacts {
				err := goldenUnicorn.InsertArtifactInfo(&a)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		DescribeTable("database select",
			func(lookedFor interface{}, expectedErr error, expectedResult ...weles.ArtifactInfo) {
				result, err := goldenUnicorn.Select(lookedFor)

				if expectedErr != nil {
					Expect(err).To(Equal(expectedErr))
					Expect(result).To(BeNil())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(result).To(Equal(expectedResult))
				}
			},
			// types supported by select.
			Entry("select JobID", artifact.JobID, nil, artifact),
			Entry("select Type", weles.ArtifactTypeYAML, nil, aYamlFailed),
			Entry("select Status", weles.ArtifactStatusREADY, nil, aImageReady),
			Entry("select Alias", artifact.Alias, nil, artifact),
			// type bool is not supported by select.
			Entry("select unsupported value", true, ErrUnsupportedQueryType),
			// test query itsef.
			Entry("select multiple entries for JobID", aImageReady.JobID, nil, aImageReady, aYamlFailed),
			Entry("select no entries for invalid JobID", invalidJob, nil),
			Entry("select multiple entries for Type", weles.ArtifactTypeIMAGE, nil, artifact, aImageReady),
			Entry("select multiple entries for Alias", aImageReady.Alias, nil, aImageReady, aYamlFailed),
			Entry("select multiple entries for Status", weles.ArtifactStatusFAILED, nil, aYamlFailed, aTestFailed),
		)
	})

	Describe("List", func() {
		BeforeEach(func() {
			for _, a := range testArtifacts {
				err := goldenUnicorn.InsertArtifactInfo(&a)
				Expect(err).ToNot(HaveOccurred())
			}
		})
		DescribeTable("list artifacts matching filter",
			func(filter weles.ArtifactFilter, expected ...weles.ArtifactInfo) {
				results, err := goldenUnicorn.Filter(filter)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).To(ConsistOf(expected))
			},
			Entry("filter one JobID", oneJobFilter, artifact),
			Entry("filter more than one JobIDs", twoJobsFilter, artifact, aImageReady, aYamlFailed),
			Entry("filter JobID not in db", noJobFilter),
			Entry("filter one Type", oneTypeFilter, aYamlFailed),
			Entry("filter more than one Type", twoTypesFilter, aYamlFailed, aTestFailed),
			Entry("filter Type not in db", noTypeFilter),
			Entry("filter one Status", oneStatusFilter, artifact),
			Entry("filter more than one Status", twoStatusFilter, artifact, aTestFailed, aYamlFailed),
			Entry("filter Status not in db", noStatusFilter),
			Entry("filter one Alias", oneAliasFilter, artifact),
			Entry("filter more than one Alias", twoAliasFilter, artifact, aImageReady, aYamlFailed),
			Entry("filter Alias not in db", noAliasFilter),
			Entry("filter is completly set up", fullFilter, aYamlFailed),
			Entry("no artifact in db matches filter", noMatchFilter),
			Entry("filter is empty", emptyFilter, artifact, aImageReady, aYamlFailed, aTestFailed),
		)
	})
	Describe("SetStatus", func() {
		BeforeEach(func() {
			for _, a := range testArtifacts {
				err := goldenUnicorn.InsertArtifactInfo(&a)
				Expect(err).ToNot(HaveOccurred())
			}
		})
		DescribeTable("artifact status change",
			func(change weles.ArtifactStatusChange, expectedErr error) {

				err := goldenUnicorn.SetStatus(change)
				if expectedErr == nil {
					Expect(err).ToNot(HaveOccurred())

					a, err := goldenUnicorn.SelectPath(change.Path)
					Expect(err).ToNot(HaveOccurred())
					Expect(a.Status).To(Equal(change.NewStatus))
				} else {
					Expect(err).To(Equal(expectedErr))
					a, err := goldenUnicorn.SelectPath(change.Path)
					Expect(err).To(HaveOccurred())
					Expect(a).To(Equal(weles.ArtifactInfo{}))
				}
			},
			Entry("change status of artifact not present in ArtifactDB", weles.ArtifactStatusChange{invalidPath, weles.ArtifactStatusDOWNLOADING}, sql.ErrNoRows),
			Entry("change status of artifact present in ArtifactDB", weles.ArtifactStatusChange{artifact.Path, weles.ArtifactStatusDOWNLOADING}, nil),
		)
	})
})
