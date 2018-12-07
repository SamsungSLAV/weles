/*
 *  Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
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

package testutil

import (
	"strings"
	"sync"

	"github.com/SamsungSLAV/slav/logger"
)

// WriterString is a simple writer that stores logs in a string using strings.Builder.
// It synchronizes reads and writes with mutex.
// It implements logger.Writer interface.
type WriterString struct {
	b     strings.Builder
	mutex sync.Locker
}

// NewWriterString creates a new WriterString object.
func NewWriterString() *WriterString {
	return &WriterString{
		mutex: new(sync.Mutex),
	}
}

// Write writes to string using strings.Builder.
// It implements logger.Writer interface in WriterString.
func (w *WriterString) Write(_ logger.Level, p []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.b.Write(append(p, '\n'))
}

// GetString returns contents stored in built string.
func (w *WriterString) GetString() string {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.b.String()
}
