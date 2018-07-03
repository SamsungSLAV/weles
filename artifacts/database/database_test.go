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
				job,
				weles.AM_IMAGEFILE,
				"some alias",
				"http://example.com",
			},
			"path1",
			weles.AM_PENDING,
			time.Now().UTC(),
		}

		aImageReady = weles.ArtifactInfo{
			weles.ArtifactDescription{
				job + 1,
				weles.AM_IMAGEFILE,
				"other alias",
				"http://example.com/1",
			},
			"path2",
			weles.AM_READY,
			time.Now().UTC(),
		}

		aYamlFailed = weles.ArtifactInfo{
			weles.ArtifactDescription{
				job + 1,
				weles.AM_YAMLFILE,
				"other alias",
				"http://example.com/2",
			},
			"path3",
			weles.AM_FAILED,
			time.Now().UTC(),
		}

		aTestFailed = weles.ArtifactInfo{
			weles.ArtifactDescription{
				job + 2,
				weles.AM_TESTFILE,
				"alias",
				"http://example.com/2",
			},
			"path4",
			weles.AM_FAILED,
			time.Unix(3000, 60).UTC(),
		}

		testArtifacts = []weles.ArtifactInfo{artifact, aImageReady, aYamlFailed, aTestFailed}

		oneJobFilter  = weles.ArtifactFilter{[]weles.JobID{artifact.JobID}, nil, nil, nil}
		twoJobsFilter = weles.ArtifactFilter{[]weles.JobID{artifact.JobID, aImageReady.JobID}, nil, nil, nil}
		noJobFilter   = weles.ArtifactFilter{[]weles.JobID{invalidJob}, nil, nil, nil}

		oneTypeFilter  = weles.ArtifactFilter{nil, []weles.ArtifactType{aYamlFailed.Type}, nil, nil}
		twoTypesFilter = weles.ArtifactFilter{nil, []weles.ArtifactType{aYamlFailed.Type, aTestFailed.Type}, nil, nil}
		noTypeFilter   = weles.ArtifactFilter{nil, []weles.ArtifactType{invalidType}, nil, nil}

		oneStatusFilter = weles.ArtifactFilter{nil, nil, []weles.ArtifactStatus{artifact.Status}, nil}
		twoStatusFilter = weles.ArtifactFilter{nil, nil, []weles.ArtifactStatus{artifact.Status, aYamlFailed.Status}, nil}
		noStatusFilter  = weles.ArtifactFilter{nil, nil, []weles.ArtifactStatus{invalidStatus}, nil}

		oneAliasFilter = weles.ArtifactFilter{nil, nil, nil, []weles.ArtifactAlias{artifact.Alias}}
		twoAliasFilter = weles.ArtifactFilter{nil, nil, nil, []weles.ArtifactAlias{artifact.Alias, aImageReady.Alias}}
		noAliasFilter  = weles.ArtifactFilter{nil, nil, nil, []weles.ArtifactAlias{invalidAlias}}

		fullFilter    = weles.ArtifactFilter{twoJobsFilter.JobID, twoTypesFilter.Type, twoStatusFilter.Status, twoAliasFilter.Alias}
		noMatchFilter = weles.ArtifactFilter{oneJobFilter.JobID, oneTypeFilter.Type, nil, nil}
		emptyFilter   = weles.ArtifactFilter{}
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
			Entry("select Type", weles.AM_YAMLFILE, nil, aYamlFailed),
			Entry("select Status", weles.AM_READY, nil, aImageReady),
			Entry("select Alias", artifact.Alias, nil, artifact),
			// type bool is not supported by select.
			Entry("select unsupported value", true, ErrUnsupportedQueryType),
			// test query itsef.
			Entry("select multiple entries for JobID", aImageReady.JobID, nil, aImageReady, aYamlFailed),
			Entry("select no entries for invalid JobID", invalidJob, nil),
			Entry("select multiple entries for Type", weles.AM_IMAGEFILE, nil, artifact, aImageReady),
			Entry("select multiple entries for Alias", aImageReady.Alias, nil, aImageReady, aYamlFailed),
			Entry("select multiple entries for Status", weles.AM_FAILED, nil, aYamlFailed, aTestFailed),
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
			Entry("change status of artifact not present in ArtifactDB", weles.ArtifactStatusChange{invalidPath, weles.AM_DOWNLOADING}, sql.ErrNoRows),
			Entry("change status of artifact present in ArtifactDB", weles.ArtifactStatusChange{artifact.Path, weles.AM_DOWNLOADING}, nil),
		)
	})
})
