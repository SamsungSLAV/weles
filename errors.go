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

// File errors.go provides definitions of errors common to Weles components.

package weles

import (
	"errors"
)

var (
	// ErrNotImplemented is returned when function is not implemented yet.
	ErrNotImplemented = errors.New("function not implemented")
	// ErrJobNotFound is returned when Job is not found.
	ErrJobNotFound = errors.New("job not found")
	// ErrJobStatusChangeNotAllowed is returned when Job status change is not
	// possible. It suggests internal Weles logic error.
	ErrJobStatusChangeNotAllowed = errors.New("job status change not allowed")
)
