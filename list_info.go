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

// ListInfo is a struct holding information to be returned by JobManager and ArtifactManager with
// lists of data.
// TotalRecords is returned to the user by API.
// RemainingRecords is used by API internally to return correct HTTP status codes (206 partial
// content and 200 ok).
type ListInfo struct {
	TotalRecords     uint64
	RemainingRecords uint64
}
