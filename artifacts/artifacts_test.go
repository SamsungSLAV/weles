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

// Package artifacts is responsible for Weles system's job artifact management.
package artifacts

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"git.tizen.org/tools/weles"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArtifactManager", func() {

	var (
		silverKangaroo weles.ArtifactManager
		testDir        string
		dbPath         string
		err            error
		job            weles.JobID = 58008
	)

	var (
		description = weles.ArtifactDescription{
			job,
			weles.AM_IMAGEFILE,
			"alias",
			"uri",
		}

		dSameJobNType = weles.ArtifactDescription{
			job,
			weles.AM_IMAGEFILE,
			"other alias",
			"other uri",
		}

		dSameJobOtherType = weles.ArtifactDescription{
			job,
			weles.AM_YAMLFILE,
			"another alias",
			"another uri",
		}
	)

	BeforeEach(func() {
		testDir, err = ioutil.TempDir("", "test-weles-")
		Expect(err).ToNot(HaveOccurred())
		dbPath = filepath.Join(testDir, "test.db")

		silverKangaroo, err = newArtifactManager(dbPath, testDir)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(testDir)
		Expect(err).ToNot(HaveOccurred())
		err = silverKangaroo.Close()
		Expect(err).ToNot(HaveOccurred())
	})

	It("should create new temp directory for artifacts", func() {
		var path, pathSame, pathType weles.ArtifactPath

		jobDir := filepath.Join(testDir, strconv.Itoa(int(description.JobID)))
		typeDir := filepath.Join(jobDir, string(description.Type))
		newTypeDir := filepath.Join(jobDir, string(dSameJobOtherType.Type))

		Expect(jobDir).ToNot(BeADirectory())

		By("CreateArtifact", func() {
			path, err = silverKangaroo.CreateArtifact(description)
			Expect(err).ToNot(HaveOccurred())
			Expect(path).NotTo(BeNil())
		})

		By("Check if all subdirs, and new file exists", func() {
			Expect(jobDir).To(BeADirectory())
			Expect(typeDir).To(BeADirectory())
			Expect(string(path)).To(BeAnExistingFile())
			Expect(string(path)).To(ContainSubstring(string(description.Alias)))
		})

		By("Add new artifact for the same JobID", func() {
			pathSame, err = silverKangaroo.CreateArtifact(dSameJobNType)
			Expect(err).ToNot(HaveOccurred())

			Expect(jobDir).To(BeADirectory())
			Expect(typeDir).To(BeADirectory())

			Expect(string(pathSame)).To(BeAnExistingFile())
			Expect(string(pathSame)).To(ContainSubstring(string(dSameJobNType.Alias)))
		})

		By("Add artifact with other type for the same JobID", func() {
			pathType, err = silverKangaroo.CreateArtifact(dSameJobOtherType)

			Expect(err).ToNot(HaveOccurred())
			Expect(jobDir).To(BeADirectory())
			Expect(newTypeDir).To(BeADirectory())

			Expect(string(pathType)).To(BeAnExistingFile())
			Expect(string(pathType)).To(ContainSubstring(string(dSameJobOtherType.Alias)))
		})

		paths := []weles.ArtifactPath{path, pathSame, pathType}
		By("Check if artifact with path is in ArtifactDB", func() {
			db, err := sql.Open("sqlite3", dbPath)
			Expect(err).ToNot(HaveOccurred())
			var n int
			for _, p := range paths {
				err = db.QueryRow("select count (*) from artifacts where path = ?", p).Scan(&n)
				Expect(err).ToNot(HaveOccurred())
				Expect(n).NotTo(BeZero())
			}
		})

		By("Check if it's possible to GetFileInfo", func() {
			for _, p := range paths {
				ai, err := silverKangaroo.GetArtifactInfo(p)
				Expect(err).ToNot(HaveOccurred())
				Expect(ai.Path).To(Equal(p))
			}
		})
	})

	Describe("Public initializer", func() {
		var (
			defaultDb  = "weles.db"
			defaultDir = "/tmp/weles/"
			customDb   = "nawia.db"
			customDir  = "/tmp/weles-custom/"
		)

		DescribeTable("NewArtifactManager()", func(db, dir string) {
			copperPanda, err := NewArtifactManager(db, dir)
			Expect(err).ToNot(HaveOccurred())

			if db == "" {
				db = defaultDb
			}
			if dir == "" {
				dir = defaultDir
			}

			Expect(dir).To(BeADirectory())
			Expect(filepath.Join(dir, db)).To(BeAnExistingFile())

			err = copperPanda.Close()
			Expect(err).ToNot(HaveOccurred())

			err = os.RemoveAll(dir)
			Expect(err).ToNot(HaveOccurred())
		},
			Entry("create database in default directory", defaultDb, defaultDir),
			Entry("create database in default directory, when arguments are empty", "", ""),
			Entry("create database in custom directory", customDb, customDir),
		)
	})
})
