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

// File job_status_extra.go enhances generated JobStatus type.

package weles

// ToInt converts JobStatus to int.
func (status JobStatus) ToInt() int {
	switch status {
	case JobStatusNEW:
		return 1
	case JobStatusPARSING:
		return 2
	case JobStatusDOWNLOADING:
		return 3
	case JobStatusWAITING:
		return 4
	case JobStatusRUNNING:
		return 5
	case JobStatusCOMPLETED:
		return 6
	case JobStatusFAILED:
		return 7
	case JobStatusCANCELED:
		return 8
	default:
		return -1
	}
}
