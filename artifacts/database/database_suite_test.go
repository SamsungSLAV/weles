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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SamsungSLAV/weles/fixtures"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var (
	silverHoneybadger ArtifactDB
	tmpDirPath        string
)

const (
	dbToRead              = "test-db-pagination.db"
	tmpDirPrefix          = "weles-"
	generatedRecordsCount = 100
	pageLimit             = 30
	pageCount             = generatedRecordsCount/pageLimit + 1
)

var _ = BeforeSuite(func() {
	var (
		err error
		db  ArtifactDB
	)

	tmpDirPath, err = ioutil.TempDir("", tmpDirPrefix)
	Expect(err).ToNot(HaveOccurred())

	err = db.Open(filepath.Join(tmpDirPath, dbToRead))
	Expect(err).ToNot(HaveOccurred())
	silverHoneybadgerArtifacts := fixtures.CreateArtifactInfoSlice(generatedRecordsCount)

	trans, err := db.dbmap.Begin()
	Expect(err).ToNot(HaveOccurred())

	for _, artifact := range silverHoneybadgerArtifacts {
		err = trans.Insert(&artifact)
		Expect(err).ToNot(HaveOccurred())
	}
	trans.Commit()
	db.Close()

	err = silverHoneybadger.Open(filepath.Join(tmpDirPath, dbToRead))
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := silverHoneybadger.Close()
	Expect(err).ToNot(HaveOccurred())
	err = os.RemoveAll(tmpDirPath)
	Expect(err).ToNot(HaveOccurred())
})
