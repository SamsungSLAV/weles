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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"flag"
	"io/ioutil"
	"net"
	"testing"

	"crypto/x509"

	"encoding/pem"

	. "git.tizen.org/tools/weles"
)

func TestManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dryad Suite")
}

var (
	dryadAddress string
	port         int
	userName     string
	keyFile      string

	dryadInfo       Dryad
	accessInfoGiven bool
)

func init() {
	flag.StringVar(&dryadAddress, "address", "", "IP address to dryad")
	flag.IntVar(&port, "port", 22, "SSH port to use for connection")
	flag.StringVar(&userName, "userName", "", "user name")
	flag.StringVar(&keyFile, "keyFile", "", "path to file containing private part of ssh key")
}

var _ = BeforeSuite(func() {
	if dryadAddress != "" && userName != "" && keyFile != "" {
		accessInfoGiven = true
		strkey, err := ioutil.ReadFile(keyFile)
		if err != nil {
			Skip("Error reading key file: " + err.Error())
		}

		block, _ := pem.Decode(strkey)
		if block == nil {
			Skip("Error decoding PEM block from key file contents")
		}

		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			Skip("Error parsing key file: " + err.Error())
		}

		dryadInfo = Dryad{
			Addr:     &net.TCPAddr{IP: net.ParseIP(dryadAddress), Port: port, Zone: ""},
			Username: userName,
			Key:      *key,
		}
	}
})
