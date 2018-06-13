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

// File controller/notifier/notifierimpl.go contains Notifier interface
// implementation. The Impl structure creates channel for communication
// between Controller and its internal submodules. Channel is used
// for notifying Controller about either failure or success of Job processing.

package notifier

import (
	"git.tizen.org/tools/weles"
)

// Impl implements Notifier interface
type Impl struct {
	Notifier
	channel chan Notification
}

// NewNotifier creates a new Impl structure setting up channel.
func NewNotifier() Notifier {
	return &Impl{
		channel: make(chan Notification, BuffSize),
	}
}

// Listen returns channel which transmits notification about failure
// or success of processing Job.
func (h *Impl) Listen() <-chan Notification {
	return h.channel
}

// SendFail notifies Controller about failure.
func (h *Impl) SendFail(j weles.JobID, msg string) {
	h.channel <- Notification{
		JobID: j,
		OK:    false,
		Msg:   msg,
	}
}

// SendOK notifies Controller about success.
func (h *Impl) SendOK(j weles.JobID) {
	h.channel <- Notification{
		JobID: j,
		OK:    true,
	}
}
