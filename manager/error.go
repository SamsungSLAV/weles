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

// File manager/error.go contains definitions of errors.

package manager

import "errors"

var (
	// ErrDuplicated is returned when a creation of dryadJob with known JobID has been requested.
	ErrDuplicated = errors.New("job with given ID already exists")
	// ErrNotExist is returned when an argument did not match any known dryadJob.
	ErrNotExist = errors.New("job with given ID does not exist")
)
