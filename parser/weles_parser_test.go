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

package parser_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/parser"
)

var expectedConfig = weles.Config{
	DeviceType: "qemu",
	JobName:    "qemu-pipeline",
	Timeouts: weles.Timeouts{
		JobTimeout:    weles.ValidPeriod(25 * time.Minute),
		ActionTimeout: weles.ValidPeriod(5 * time.Minute),
	},
	Priority: "medium",
	Action: weles.Action{
		Deploy: weles.Deploy{
			Timeout: weles.ValidPeriod(20 * time.Minute),
			Images: []weles.ImageDefinition{
				{
					URI: "https://images.validation.linaro.org/kvm/standard/" +
						"stretch-1.img.gz",
					ChecksumURI:  "https://images.validation.linaro.org/kvm/standard/stretch-1.md5",
					ChecksumType: "md5",
					Compression:  "gz",
					Path:         "",
					ChecksumPath: "",
				},
				{
					URI: "https://images.validation.linaro.org/kvm/standard/" +
						"stretch-2.img.zip",
					ChecksumURI:  "https://images.validation.linaro.org/kvm/standard/stretch-2.md5",
					ChecksumType: "md5",
					Compression:  "zip",
					Path:         "",
					ChecksumPath: "",
				},
			},
			PartitionLayout: []weles.PartitionDefinition{
				{
					ID:        1,
					ImageName: "image_name1_string",
					Size:      "12345",
					Type:      "fat",
				},
				{
					ID:        2,
					ImageName: "image_name2_string",
					Size:      "23456",
					Type:      "ext2",
				},
				{
					ID:        3,
					ImageName: "image_name3_string",
					Size:      "34567",
					Type:      "ext3",
				},
			},
		},
		Boot: weles.Boot{
			Login:         "root",
			Password:      "tizen",
			Prompts:       []string{"linaro-test", "root@debian:~#"},
			FailureRetry:  2,
			Timeout:       weles.ValidPeriod(20 * time.Minute),
			InputSequence: "input_sequence_string",
			WaitPattern:   "sample pattern 1 we wait for",
			WaitTime:      weles.ValidPeriod(4 * time.Minute),
		},
		Test: weles.Test{
			FailureRetry: 3,
			Name:         "kvm-basic-singlenode",
			Timeout:      weles.ValidPeriod(5 * time.Minute),
			TestCases: []weles.TestCase{
				{
					CaseName: "case_name1_string",
					TestActions: []weles.TestAction{
						weles.Boot{
							Login:         "root",
							Password:      "tizen",
							Prompts:       []string{"linaro-test", "root@debian:~#"},
							FailureRetry:  2,
							Timeout:       weles.ValidPeriod(20 * time.Minute),
							InputSequence: "input_sequence_string",
							WaitPattern:   "sample pattern 1 we wait for",
							WaitTime:      weles.ValidPeriod(4 * time.Minute),
						},
						weles.Push{
							URI:     "uri1_string",
							Dest:    "path1_string",
							Alias:   "alias1_string",
							Timeout: weles.ValidPeriod(6 * time.Minute),
							Path:    "",
						},
						weles.Run{
							Name:    "name1_string",
							Timeout: weles.ValidPeriod(2 * time.Minute),
						},
						weles.Pull{
							Src:     "path2_string",
							Alias:   "alias2_string",
							Timeout: weles.ValidPeriod(1 * time.Minute),
							Path:    "",
						},
					},
				},
				{
					CaseName: "case_name2_string",
					TestActions: []weles.TestAction{
						weles.Boot{
							Login:         "root",
							Password:      "tizen",
							Prompts:       []string{"linaro-test", "root@debian:~#"},
							FailureRetry:  2,
							Timeout:       weles.ValidPeriod(20 * time.Minute),
							InputSequence: "input_sequence_string",
							WaitPattern:   "sample pattern 2 we wait for",
							WaitTime:      weles.ValidPeriod(3 * time.Minute),
						},
						weles.Push{
							URI:     "uri1_string",
							Dest:    "path1_string",
							Alias:   "alias1_string",
							Timeout: weles.ValidPeriod(4 * time.Minute),
							Path:    "",
						},
						weles.Push{
							URI:     "uri1_string",
							Dest:    "path1_string",
							Alias:   "alias1_string",
							Timeout: weles.ValidPeriod(5 * time.Minute),
							Path:    "",
						},
						weles.Pull{
							Src:     "path2_string",
							Alias:   "alias2_string",
							Timeout: weles.ValidPeriod(2 * time.Minute),
							Path:    "",
						},
					},
				},
				{
					CaseName: "case_name3_string",
					TestActions: []weles.TestAction{
						weles.Pull{
							Src:     "path2_string",
							Alias:   "alias2_string",
							Timeout: weles.ValidPeriod(1 * time.Minute),
							Path:    "",
						},
					},
				},
			},
		},
	},
}

var input = []byte(`
device_type: qemu
job_name: qemu-pipeline
timeouts:
  job:
    minutes: 25     # timeout for the whole job
  action:
    minutes: 5      # default timeout applied for each action; can be overriden in the action itself
priority: medium

actions:

  - deploy:
      timeout:
        minutes: 20
      images:       # list of images
         - uri: https://images.validation.linaro.org/kvm/standard/stretch-1.img.gz
           checksum_uri: https://images.validation.linaro.org/kvm/standard/stretch-1.md5
           checksum_type: md5
           compression: gz
         - uri: https://images.validation.linaro.org/kvm/standard/stretch-2.img.zip
           checksum_uri: https://images.validation.linaro.org/kvm/standard/stretch-2.md5
           checksum_type: md5
           compression: zip
      partition_layout:     # list of partitions structures
         - id: 1
           device_name: device_name1_string
           image_name: image_name1_string
           size: 12345
           type: fat
         - id: 2
           device_name: device_name2_string
           image_name: image_name2_string
           size: 23456
           type: ext2
         - id: 3
           device_name: device_name3_string
           image_name: image_name3_string
           size: 34567
           type: ext3
  - boot:
      login: root
      password: tizen
      prompts:
        - 'linaro-test'
        - 'root@debian:~#'
      failure_retry: 2
      timeout:
        minutes: 20
      input_sequence: input_sequence_string
      wait_pattern: 'sample pattern 1 we wait for'
      wait_time:
        minutes: 4
  - test:
      failure_retry: 3
      name: kvm-basic-singlenode
      timeout:
        minutes: 5
      test_cases:
        - case_name: case_name1_string
          test_actions:
            - boot:
                login: root
                password: tizen
                prompts:
                  - 'linaro-test'
                  - 'root@debian:~#'
                failure_retry: 2
                timeout:
                  minutes: 20
                input_sequence: input_sequence_string
                wait_pattern: 'sample pattern 1 we wait for'
                wait_time:
                  minutes: 4
            - push:
                uri: uri1_string
                dest: path1_string
                alias: alias1_string
                timeout:
                  minutes: 6
            - run:
                name: name1_string
                timeout:
                  minutes: 2
            - pull:
                src: path2_string
                alias: alias2_string
                timeout:
                  minutes: 1
        - case_name: case_name2_string
          test_actions:
            - boot:
                login: root
                password: tizen
                prompts:
                  - 'linaro-test'
                  - 'root@debian:~#'
                failure_retry: 2
                timeout:
                  minutes: 20
                input_sequence: input_sequence_string
                wait_pattern: 'sample pattern 2 we wait for'
                wait_time:
                  minutes: 3
            - push:
                uri: uri1_string
                dest: path1_string
                alias: alias1_string
                timeout:
                  minutes: 4
            - push:
                uri: uri1_string
                dest: path1_string
                alias: alias1_string
                timeout:
                  minutes: 5
            - pull:
                src: path2_string
                alias: alias2_string
                timeout:
                  minutes: 2
        - case_name: case_name3_string
          test_actions:
             - pull:
                src: path2_string
                alias: alias2_string
                timeout:
                  minutes: 1
`)

var _ = Describe("WelesParser", func() {
	It("should parse input bytes into proper config structure", func() {
		var p parser.Parser
		conf, err := p.ParseYaml(input)
		Expect(err).ToNot(HaveOccurred())
		Expect(conf).To(Equal(&expectedConfig))
		expectedConfig = weles.Config{}
	})
})
