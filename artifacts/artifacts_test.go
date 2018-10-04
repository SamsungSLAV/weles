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

// Package artifacts is responsible for Weles system's job artifact management.
package artifacts

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	const poem = `How doth the little crocodile
Improve his shining tail,
And pour the waters of the Nile
On every golden scale!

How cheerfully he seems to grin
How neatly spreads his claws,
And welcomes little fishes in,
With gently smiling jaws!

-Lewis Carroll`

	var (
		testDir string
		dbPath  string
	)

	var (
		silverKangaroo weles.ArtifactManager
		job            weles.JobID       = 58008
		validURL       weles.ArtifactURI = "validURL"
		invalidURL     weles.ArtifactURI = "invalidURL"
	)

	var (
		description = weles.ArtifactDescription{
			Alias: "alias",
			JobID: job,
			Type:  weles.ArtifactTypeIMAGE,
			URI:   "uri",
		}

		dSameJobNType = weles.ArtifactDescription{
			Alias: "other alias",
			JobID: job,
			Type:  weles.ArtifactTypeIMAGE,
			URI:   "other uri",
		}

		dSameJobOtherType = weles.ArtifactDescription{
			Alias: "another alias",
			JobID: job,
			Type:  weles.ArtifactTypeYAML,
			URI:   "another uri",
		}
	)

	BeforeEach(func() {
		var err error
		testDir, err = ioutil.TempDir("", "test-weles-")
		Expect(err).ToNot(HaveOccurred())
		dbPath = filepath.Join(testDir, "test.db")

		silverKangaroo, err = newArtifactManager(dbPath, testDir, 100, 16, 100)
		//TODO add tests against different notifier cap, queue cap and workers count.
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(testDir)
		Expect(err).ToNot(HaveOccurred())
		err = silverKangaroo.Close()
		Expect(err).ToNot(HaveOccurred())
	})

	checkPathInDb := func(path weles.ArtifactPath) bool {
		db, err := sql.Open("sqlite3", dbPath)
		Expect(err).ToNot(HaveOccurred())
		defer db.Close()
		var n int
		err = db.QueryRow("select count (*) from artifacts where path = ?", path).Scan(&n)
		Expect(err).ToNot(HaveOccurred())
		return (n > 0)
	}

	prepareServer := func(url weles.ArtifactURI) *httptest.Server {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if url == validURL {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, poem)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		return ts
	}

	It("should create new temp directory for artifacts", func() {
		var path, pathSame, pathType weles.ArtifactPath
		jobDir := filepath.Join(testDir, strconv.Itoa(int(description.JobID)))
		typeDir := filepath.Join(jobDir, string(description.Type))
		newTypeDir := filepath.Join(jobDir, string(dSameJobOtherType.Type))

		Expect(jobDir).ToNot(BeADirectory())

		var err error
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
				Expect(checkPathInDb(p)).To(BeTrue())
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
			copperPanda, err := NewArtifactManager(db, dir, 100, 16, 100)
			//TODO: add tests against different notifier cap and workers count.
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
			Entry("create database in custom directory", customDb, customDir),
		)
	})

	Describe("PushArtifact", func() {

		var (
			ch chan weles.ArtifactStatusChange

			ad weles.ArtifactDescription = weles.ArtifactDescription{
				Alias: "somealias",
				JobID: job,
				Type:  weles.ArtifactTypeIMAGE,
				URI:   validURL,
			}

			adInvalid weles.ArtifactDescription = weles.ArtifactDescription{
				Alias: "somealias",
				JobID: job,
				Type:  weles.ArtifactTypeIMAGE,
				URI:   invalidURL,
			}
		)

		BeforeEach(func() {
			ch = make(chan weles.ArtifactStatusChange, 20)
		})

		DescribeTable("Push artifact",
			func(ad weles.ArtifactDescription, finalStatus weles.ArtifactStatus) {

				ts := prepareServer(ad.URI)
				defer ts.Close()
				ad.URI = weles.ArtifactURI(ts.URL)

				path, err := silverKangaroo.PushArtifact(ad, ch)

				Expect(err).ToNot(HaveOccurred())

				Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{
					Path:      path,
					NewStatus: weles.ArtifactStatusPENDING,
				})))
				Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{
					Path:      path,
					NewStatus: weles.ArtifactStatusDOWNLOADING,
				})))
				Eventually(ch).Should(Receive(Equal(weles.ArtifactStatusChange{
					Path:      path,
					NewStatus: finalStatus,
				})))

				if finalStatus != weles.ArtifactStatusFAILED {
					By("Check if file exists and has proper content")
					content, erro := ioutil.ReadFile(string(path))

					Expect(erro).ToNot(HaveOccurred())
					Expect(string(content)).To(BeIdenticalTo(poem))

				} else {
					By("Check if file exists")
					Expect(string(path)).NotTo(BeAnExistingFile())
				}

				Eventually(func() weles.ArtifactStatus {
					ai, err := silverKangaroo.GetArtifactInfo(path)
					Expect(err).ToNot(HaveOccurred())
					return ai.Status
				}).Should(Equal(finalStatus))

				By("Check if artifact is in ArtifactDB")
				Expect(checkPathInDb(path)).To(BeTrue())
			},
			Entry("push artifact to db and download file", ad, weles.ArtifactStatusREADY),
			Entry("do not push an invalid artifact", adInvalid, weles.ArtifactStatusFAILED),
		)
	})
})
