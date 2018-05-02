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

// File controller/notifier/notifier.go defines interface for providing notification
// channel to notify Controller about success or failure of processing Job.

package notifier

import (
	. "git.tizen.org/tools/weles"
)

// Notifier defines interface providing channel for notifying Controller.
type Notifier interface {
	// Listen returns channel which transmits notification about failure
	// or success of processed Job.
	Listen() <-chan Notification

	// SendFail notifies Controller about failure.
	SendFail(j JobID, msg string)

	// SendOK notifies Controller about success.
	SendOK(j JobID)
}
