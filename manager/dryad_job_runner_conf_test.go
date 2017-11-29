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

package manager

import (
	"git.tizen.org/tools/weles"
)

// basicConfig is incomplete, but contains necessary definitions
// for some of the tests. It can be used as a base for more complex
// definitions.
var basicConfig = weles.Config{
	Action: weles.Action{
		Deploy: weles.Deploy{
			Images: []weles.ImageDefinition{
				{Path: "artifact/path/image1"},
				{Path: "artifact/path/image2"},
			},
			PartitionLayout: []weles.PartitionDefinition{
				{
					ID:        1,
					ImageName: "image name_1",
					Size:      "4096",
				},
				{
					ID:        2,
					ImageName: "image_name 2",
					Size:      "1234",
				},
			},
		},
		Boot: weles.Boot{
			Login:         "test_login",
			Password:      "test_password",
			Prompts:       []string{"prompt ~1"},
			InputSequence: "\n",
			WaitPattern:   "device login:",
		},
		Test: weles.Test{
			TestCases: []weles.TestCase{
				{
					TestActions: []weles.TestAction{
						weles.Push{
							Path: "artifact/path/test_push1",
							Dest: "dest of test_push1",
						},
						weles.Run{
							Name: "command to be run",
						},
						weles.Pull{
							Src:  "src of test_pull1",
							Path: "artifact/path/test_push1",
						},
					},
				},
			},
		},
	},
}
