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

// File manager/dryad/error.go contains definitions of errors.

package dryad

import "errors"

var (
	// ErrConnectionClosed is returned when caller tries to close already closed connection
	// to Dryad.
	ErrConnectionClosed = errors.New("attempt to close already closed connection")
	// ErrNotMounted is returned when the check for sshfs mount fails.
	ErrNotMounted = errors.New("filesystem not mounted")
)
