// Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package weles

// JobPagination holds information neccessary to request for a single page of data.
// When JobID is set, and Forward is false - Controller should return a page of records before the
// supplied JobID.
// When JobID is set, and Forward is true -  Controller should return page of record after the
// supplied JobID.
// In both cases, returned page should not include supplied JobID.
// Limit denotes the number of records to be returned on the page.
// When Limit is set to 0, pagination is disabled, JobID and Forward fields are ignored
// and all records are returned.
type JobPagination struct {
	JobID   JobID
	Forward bool
	Limit   int32
}
