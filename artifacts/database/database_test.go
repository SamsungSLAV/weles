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

// Package database is responsible for Weles system's job artifact storage.
package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/SamsungSLAV/weles"
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

		artifact = weles.ArtifactInfo{
			ArtifactDescription: weles.ArtifactDescription{
				Alias: "some alias",
				JobID: job,
				Type:  weles.ArtifactTypeIMAGE,
				URI:   "http://example.com",
			},
			Path:      "path1",
			Status:    weles.ArtifactStatusPENDING,
			Timestamp: strfmt.DateTime(time.Now().Round(time.Millisecond).UTC()),
		}
		aImageReady = weles.ArtifactInfo{
			ArtifactDescription: weles.ArtifactDescription{
				Alias: "other alias",
				JobID: job + 1,
				Type:  weles.ArtifactTypeIMAGE,
				URI:   "http://example.com/1",
			},
			Path:      "path2",
			Status:    weles.ArtifactStatusREADY,
			Timestamp: strfmt.DateTime(time.Now().Round(time.Millisecond).UTC()),
		}
		aYamlFailed = weles.ArtifactInfo{
			ArtifactDescription: weles.ArtifactDescription{
				Alias: "other alias",
				JobID: job + 1,
				Type:  weles.ArtifactTypeYAML,
				URI:   "http://example.com/2",
			},
			Path:      "path3",
			Status:    weles.ArtifactStatusFAILED,
			Timestamp: strfmt.DateTime(time.Now().Round(time.Millisecond).UTC()),
		}
		aTestFailed = weles.ArtifactInfo{
			ArtifactDescription: weles.ArtifactDescription{
				Alias: "alias",
				JobID: job + 2,
				Type:  weles.ArtifactTypeTEST,
				URI:   "http://example.com/3",
			},
			Path:      "path4",
			Status:    weles.ArtifactStatusFAILED,
			Timestamp: strfmt.DateTime(time.Unix(3000, 60).Round(time.Millisecond).UTC()),
		}
		testArtifacts = []weles.ArtifactInfo{artifact, aImageReady, aYamlFailed, aTestFailed}

		oneJobFilter = weles.ArtifactFilter{
			JobID:  []weles.JobID{artifact.JobID},
			Alias:  nil,
			Status: nil,
			Type:   nil,
		}
		twoJobsFilter = weles.ArtifactFilter{
			JobID:  []weles.JobID{artifact.JobID, aImageReady.JobID},
			Alias:  nil,
			Status: nil,
			Type:   nil,
		}
		noJobFilter = weles.ArtifactFilter{
			JobID:  []weles.JobID{invalidJob},
			Alias:  nil,
			Status: nil,
			Type:   nil,
		}
		oneTypeFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   []weles.ArtifactType{aYamlFailed.Type},
			Status: nil,
			Alias:  nil,
		}
		twoTypesFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   []weles.ArtifactType{aYamlFailed.Type, aTestFailed.Type},
			Status: nil,
			Alias:  nil,
		}
		noTypeFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   []weles.ArtifactType{invalidType},
			Status: nil,
			Alias:  nil,
		}
		oneStatusFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: []weles.ArtifactStatus{artifact.Status},
			Alias:  nil,
		}
		twoStatusFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: []weles.ArtifactStatus{artifact.Status, aYamlFailed.Status},
			Alias:  nil,
		}
		noStatusFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: []weles.ArtifactStatus{invalidStatus},
			Alias:  nil,
		}
		oneAliasFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: nil,
			Alias:  []weles.ArtifactAlias{artifact.Alias},
		}
		twoAliasFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: nil,
			Alias:  []weles.ArtifactAlias{artifact.Alias, aImageReady.Alias},
		}
		noAliasFilter = weles.ArtifactFilter{
			JobID:  nil,
			Type:   nil,
			Status: nil,
			Alias:  []weles.ArtifactAlias{invalidAlias},
		}
		fullFilter = weles.ArtifactFilter{
			JobID:  []weles.JobID{artifact.JobID, aImageReady.JobID, aYamlFailed.JobID},
			Type:   twoTypesFilter.Type,
			Status: twoStatusFilter.Status,
			Alias:  twoAliasFilter.Alias,
		}
		noMatchFilter = weles.ArtifactFilter{
			JobID:  oneJobFilter.JobID,
			Type:   oneTypeFilter.Type,
			Status: nil,
			Alias:  nil,
		}
		emptyFilter = weles.ArtifactFilter{}

		emptyPaginator = weles.ArtifactPagination{}

		descendingSorter = weles.ArtifactSorter{
			SortOrder: weles.SortOrderDescending,
			SortBy:    weles.ArtifactSortByID,
		}

		defaultSorter = descendingSorter

		ascendingSorter = weles.ArtifactSorter{
			SortOrder: weles.SortOrderAscending,
			SortBy:    weles.ArtifactSortByID,
		}
	)
	jobInDB := func(job weles.JobID, db ArtifactDB) bool {
		n, err := db.dbmap.SelectInt(
			`SELECT COUNT(*)
 			FROM artifacts
 			WHERE JobID = ?`, job)
		Expect(err).ToNot(HaveOccurred())
		return bool(n > 0)
	}

	Describe("Not pagination", func() {

		var (
			goldenUnicorn ArtifactDB
			tmpDir        string
		)
		BeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir("", tmpDirPrefix)
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
			n, err := goldenUnicorn.dbmap.SelectInt(
				`SELECT COUNT(*)
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
			Expect(jobInDB(job, goldenUnicorn)).To(BeFalse())

			err := goldenUnicorn.InsertArtifactInfo(&artifact)
			Expect(err).ToNot(HaveOccurred())

			Expect(jobInDB(artifact.JobID, goldenUnicorn)).To(BeTrue())
		})

		Describe("SetStatus", func() {
			BeforeEach(func() {
				trans, err := goldenUnicorn.dbmap.Begin()
				Expect(err).ToNot(HaveOccurred())
				defer trans.Commit()
				for _, a := range testArtifacts {
					err := trans.Insert(&a)
					Expect(err).ToNot(HaveOccurred())
				}
			})
			DescribeTable("artifact status change",
				func(change weles.ArtifactStatusChange, expectedErr error) {
					var a weles.ArtifactInfo
					err := goldenUnicorn.SetStatus(change)
					if expectedErr == nil {
						Expect(err).ToNot(HaveOccurred())
						a, err = goldenUnicorn.SelectPath(change.Path)
						Expect(err).ToNot(HaveOccurred())
						Expect(a.Status).To(Equal(change.NewStatus))
					} else {
						Expect(err).To(Equal(expectedErr))
						a, err = goldenUnicorn.SelectPath(change.Path)
						Expect(err).To(HaveOccurred())
						Expect(a).To(Equal(weles.ArtifactInfo{}))
					}
				},
				Entry("change status of artifact not present in ArtifactDB",
					weles.ArtifactStatusChange{
						Path:      invalidPath,
						NewStatus: weles.ArtifactStatusDOWNLOADING,
					},
					sql.ErrNoRows),
				Entry("change status of artifact present in ArtifactDB",
					weles.ArtifactStatusChange{
						Path:      artifact.Path,
						NewStatus: weles.ArtifactStatusDOWNLOADING,
					},
					nil),
			)
		})

		Describe("SelectPath", func() {

			BeforeEach(func() {
				err := goldenUnicorn.InsertArtifactInfo(&artifact)
				Expect(err).ToNot(HaveOccurred())

				Expect(jobInDB(artifact.JobID, goldenUnicorn)).To(BeTrue())
			})

			DescribeTable("database selectpath",
				func(path weles.ArtifactPath, expectedErr error,
					expectedArtifact weles.ArtifactInfo) {
					result, err := goldenUnicorn.SelectPath(path)

					if expectedErr != nil {
						Expect(err).To(Equal(expectedErr))
					} else {
						Expect(err).ToNot(HaveOccurred())
					}
					expectedArtifact.ID = result.ID
					Expect(result).To(Equal(expectedArtifact))
				},
				Entry("retrieve artifact based on path", artifact.Path, nil, artifact),
				Entry("retrieve artifact based on invalid path", invalidPath, sql.ErrNoRows,
					weles.ArtifactInfo{}),
			)
		})

		Describe("List", func() {
			BeforeEach(func() {
				trans, err := goldenUnicorn.dbmap.Begin()
				Expect(err).ToNot(HaveOccurred())
				defer trans.Commit()
				for _, a := range testArtifacts {
					err := trans.Insert(&a)
					Expect(err).ToNot(HaveOccurred())
				}
			})
			DescribeTable("list artifacts matching filter",
				func(filter weles.ArtifactFilter, expected ...weles.ArtifactInfo) {
					results, _, err := goldenUnicorn.Filter(filter, defaultSorter, emptyPaginator)
					Expect(err).ToNot(HaveOccurred())
					//TODO: match all fields except ID.
					for i := range results {
						for j := range expected {
							if results[i].JobID == expected[j].JobID {
								if results[i].URI == expected[j].URI {
									expected[j].ID = results[i].ID
								}
							}
						}
					}
					Expect(results).To(ConsistOf(expected))
				},
				Entry("filter one JobID", oneJobFilter, artifact),
				Entry("filter more than one JobIDs", twoJobsFilter, artifact, aImageReady,
					aYamlFailed),
				Entry("filter one Type", oneTypeFilter, aYamlFailed),
				Entry("filter more than one Type", twoTypesFilter, aYamlFailed, aTestFailed),
				Entry("filter one Status", oneStatusFilter, artifact),
				Entry("filter more than one Status", twoStatusFilter, artifact, aTestFailed,
					aYamlFailed),
				Entry("filter one Alias", oneAliasFilter, artifact),
				Entry("filter more than one Alias", twoAliasFilter, artifact, aImageReady,
					aYamlFailed),
				Entry("filter is completly set up", fullFilter, aYamlFailed),
				Entry("filter is empty", emptyFilter, artifact, aImageReady, aYamlFailed,
					aTestFailed),
			)

			DescribeTable("return artifact not found error",
				func(filter weles.ArtifactFilter, expected ...weles.ArtifactInfo) {
					_, _, err := goldenUnicorn.Filter(filter, defaultSorter, emptyPaginator)
					Expect(err).To(Equal(weles.ErrArtifactNotFound))
				},
				Entry("filter JobID not in db", noJobFilter),
				Entry("filter Type not in db", noTypeFilter),
				Entry("filter Status not in db", noStatusFilter),
				Entry("filter Alias not in db", noAliasFilter),
				Entry("no artifact in db matches filter", noMatchFilter),
			)
		})
		Describe("Sorting", func() {
			BeforeEach(func() {
				trans, err := goldenUnicorn.dbmap.Begin()
				Expect(err).ToNot(HaveOccurred())
				defer trans.Commit()
				for _, a := range testArtifacts {
					err := trans.Insert(&a)
					Expect(err).ToNot(HaveOccurred())
				}
			})
			DescribeTable("Should return correctly sorted artifacts",
				func(sorter weles.ArtifactSorter) {
					result, _, err := goldenUnicorn.Filter(emptyFilter, sorter, emptyPaginator)
					Expect(err).ToNot(HaveOccurred())
					var currID int
					if sorter.SortOrder == weles.SortOrderAscending {
						for _, a := range result {
							if currID == 0 {
								currID = int(a.ID)
								continue
							}
							Expect(a.ID).To(BeEquivalentTo(currID + 1))
							currID = int(a.ID)
						}
					} else {
						for _, a := range result {
							if currID == 0 {
								currID = int(a.ID)
								continue
							}
							Expect(a.ID).To(BeEquivalentTo(currID - 1))
							currID = int(a.ID)
						}
					}
				},
				Entry("By ID, Ascending", ascendingSorter),
				Entry("By ID, Descending", descendingSorter),
			)

		})
	})
	Describe("Pagination", func() {
		Context("Database is filled with generatedRecordsCount records", func() {
			DescribeTable("paginating through artifacts",
				func(paginator weles.ArtifactPagination,
					expectedResponseLength, expectedRemainingRecords int) {

					result, list, err := silverHoneybadger.Filter(emptyFilter, defaultSorter, paginator)
					Expect(err).ToNot(HaveOccurred())

					Expect(len(result)).To(BeEquivalentTo(expectedResponseLength))
					Expect(list.TotalRecords).To(BeEquivalentTo(generatedRecordsCount))
					Expect(list.RemainingRecords).To(BeEquivalentTo(expectedRemainingRecords))
				},
				// please keep in mind that data is sorted in descending order.
				Entry("first and last page (limit is 0)",
					weles.ArtifactPagination{ID: 0, Limit: 0, Forward: true},
					generatedRecordsCount, 0),

				Entry("first page, paginating forward",
					weles.ArtifactPagination{ID: 0, Limit: pageLimit, Forward: true},
					pageLimit, (generatedRecordsCount-pageLimit)),

				Entry("second page, paginating forward",
					weles.ArtifactPagination{
						ID:      (generatedRecordsCount - pageLimit + 1),
						Limit:   pageLimit,
						Forward: true,
					},
					pageLimit,
					(generatedRecordsCount-(2*pageLimit))),

				Entry("last page, paginating forward",
					weles.ArtifactPagination{
						ID: int64(generatedRecordsCount -
							(pageLimit * int(generatedRecordsCount/pageLimit)) + 1),
						Limit:   pageLimit,
						Forward: true,
					},
					generatedRecordsCount-(pageLimit*(pageCount-1)),
					0),

				Entry("second to last page, paginating backward",
					weles.ArtifactPagination{
						ID: int64(generatedRecordsCount -
							(pageLimit * int(generatedRecordsCount/pageLimit))),
						Limit:   pageLimit,
						Forward: false,
					},
					pageLimit,
					pageLimit*(pageCount-2)),

				Entry("first page, paginating backward",
					weles.ArtifactPagination{
						ID:      int64(generatedRecordsCount - pageLimit),
						Limit:   pageLimit,
						Forward: false,
					},
					pageLimit,
					0),
			)
		})
	})
})
