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

// File manager/dryad_job_runner_fota.go provides wrapper of CLI to fota tool available on Dryad.

package manager

import (
	"encoding/json"
	"strconv"
)

const (
	fotaCmdPath    = "/usr/local/bin/fota"
	fotaSDCardPath = "/dev/sda"
	fotaFilePath   = "/tmp/fota.json"
)

type fotaCmd struct {
	sdcard  string
	mapping string
	md5sums string
	URLs    []string
}

// newFotaCmd creates new fotaCmd instance.
// Currently it always receives same params thus nolint directive.
// Those params should be based on job submission file.
// nolint:unparam
func newFotaCmd(sdcard, mapping string, urls []string) *fotaCmd {
	return &fotaCmd{
		sdcard:  sdcard,
		mapping: mapping,
		URLs:    urls,
	}
}

func (f *fotaCmd) SetMD5(url string) {
	f.md5sums = url
}

func (f *fotaCmd) GetCmd() (cmd []string) {
	cmd = []string{fotaCmdPath,
		"-map", f.mapping,
		"-card", f.sdcard}
	if f.md5sums != "" {
		cmd = append(cmd, "-md5")
		cmd = append(cmd, f.md5sums)
	}
	cmd = append(cmd, f.URLs...)
	return
}

type fotaMap struct {
	name string
	part int
}

func newMapping(fms []fotaMap) []byte {
	fotaMapping := make(map[string]string)
	for _, fm := range fms {
		fotaMapping[fm.name] = strconv.Itoa(fm.part)
	}
	ret, err := json.Marshal(fotaMapping)
	if err != nil {
		panic(err)
	}
	return ret
}
