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

package dryad

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
)

var _ = Describe("reverseSSHFS", func() {
	var (
		client *ssh.Client
		dir    string

		sshfs *reverseSSHFS
	)

	BeforeEach(func() {
		if !accessInfoGiven {
			Skip("No valid access info to Dryad")
		}
		var err error
		client, err = ssh.Dial(
			"tcp", dryadInfo.Addr.String(), prepareSSHConfig(dryadInfo.Username, dryadInfo.Key))
		Expect(err).ToNot(HaveOccurred())
		dir, err = ioutil.TempDir("", "sshfs-test")
		Expect(err).ToNot(HaveOccurred())

		sshfs = newReverseSSHFS(context.Background(), dir, dir)
	})

	AfterEach(func() {
		err := sshfs.close()
		Expect(err).ToNot(HaveOccurred())
		client.Close()
		os.RemoveAll(dir)
	})

	getSession := func() *ssh.Session {
		session, err := client.NewSession()
		Expect(err).ToNot(HaveOccurred())
		return session
	}

	It("should close even if open was never called", func() {
	})

	It("should open and check", func() {
		err := sshfs.open(getSession())
		Expect(err).ToNot(HaveOccurred())

		err = sshfs.check(getSession())
		Expect(err).ToNot(HaveOccurred())
	})

	It("should have files accessible on both ends", func() {
		err := sshfs.open(getSession())
		Expect(err).ToNot(HaveOccurred())

		err = sshfs.check(getSession())
		Expect(err).ToNot(HaveOccurred())

		By("create on local and access from remote", func() {
			filename := filepath.Join(dir, "testfile1")

			f, err := os.Create(filename)
			Expect(err).ToNot(HaveOccurred())
			err = f.Close()
			Expect(err).ToNot(HaveOccurred())

			session := getSession()
			err = session.Run("stat " + filename)
			Expect(err).ToNot(HaveOccurred())
		})

		By("create on remote and access from local", func() {
			filename := filepath.Join(dir, "testfile2")

			session := getSession()
			err = session.Run("touch " + filename)
			Expect(err).ToNot(HaveOccurred())

			_, err = os.Stat(filename)
			Expect(err).ToNot(HaveOccurred())
		})

		err = sshfs.check(getSession())
		Expect(err).ToNot(HaveOccurred())
	})

	It("should not be accessible after close", func() {
		err := sshfs.open(getSession())
		Expect(err).ToNot(HaveOccurred())

		err = sshfs.check(getSession())
		Expect(err).ToNot(HaveOccurred())

		err = sshfs.close()
		Expect(err).ToNot(HaveOccurred())

		filename := filepath.Join(dir, "testfile2")

		session := getSession()
		err = session.Run("touch " + filename)
		Expect(err).ToNot(HaveOccurred())

		_, err = os.Stat(filename)
		Expect(err).To(HaveOccurred())

		// Make AfterEach work.
		sshfs = newReverseSSHFS(context.Background(), dir, dir)
	})
})
