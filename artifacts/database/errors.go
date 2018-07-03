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

// File errors.go provides definitions of errors for Weles' database package.

package database

import (
	"errors"
)

var (
	// ErrUnsupportedQueryType is returned when wrong type of argument is passed to
	// ArtifactDB's Select().
	ErrUnsupportedQueryType = errors.New("unsupported argument type")
)
