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

	"git.tizen.org/tools/weles/fixtures"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var (
	silverHoneybadger ArtifactDB
	tmpDirPath        string
)

var _ = BeforeSuite(func() {

	var err error
	tmpDirPath, err = ioutil.TempDir("", "weles-")
	Expect(err).ToNot(HaveOccurred())
	err = silverHoneybadger.Open(filepath.Join(tmpDirPath, "test_pagination.db"))
	Expect(err).ToNot(HaveOccurred())
	artifacts := fixtures.CreateArtifactInfoSlice(100)
	trans, err := silverHoneybadger.dbmap.Begin()
	Expect(err).ToNot(HaveOccurred())
	for _, artifact := range artifacts {
		err = silverHoneybadger.InsertArtifactInfo(&artifact)
		Expect(err).ToNot(HaveOccurred())
	}
	trans.Commit()

})

var _ = AfterSuite(func() {

	err := silverHoneybadger.Close()
	Expect(err).ToNot(HaveOccurred())
	err = os.RemoveAll(tmpDirPath)
	Expect(err).ToNot(HaveOccurred())
})
