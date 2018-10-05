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

// Package notifier defines structures and constants used by controller
// package and provides Notifier interface with implementation
// for communication between submodules and Controller.
package notifier

import (
	"github.com/SamsungSLAV/weles"
)

// BuffSize is the default channel buffer size for channels passing
// notification messages to Controller.
const BuffSize = 32

// Notification describes single notification message for Controller.
// It is passed to Controller's listening channels.
type Notification struct {
	// JobID identifies Job.
	weles.JobID
	// Ok reports if Job processing stage has ended with success.
	OK bool
	// Msg contains additional information for the final user.
	Msg string
}
