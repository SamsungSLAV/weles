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
//

package enums

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// JobSortBy denotes key for sorting Jobs list.
//
// * ID - default sort key.
//
// * CreatedDate - sorting by date of creation of the weles Job.
//
// * UpdatedDate - sorting by date of update of the weles Job.
//
// * JobStatus - sorting by the Job Status. Descending order will sort in
// the order JobStatuses are listed in the docs (from NEW at the start to
// CANCELED at the end). Ascending will reverse this order.
//
// When sorting is applied, and there are many Jobs with the same
// date/status, they will be sorted by JobID (Ascending)
//
// swagger:model JobSortBy
type JobSortBy string

const (

	// JobSortByID captures enum value "ID"
	JobSortByID JobSortBy = "ID"

	// JobSortByCreatedDate captures enum value "CreatedDate"
	JobSortByCreatedDate JobSortBy = "CreatedDate"

	// JobSortByUpdatedDate captures enum value "UpdatedDate"
	JobSortByUpdatedDate JobSortBy = "UpdatedDate"

	// JobSortByJobStatus captures enum value "JobStatus"
	JobSortByJobStatus JobSortBy = "JobStatus"
)

// for schema
var jobSortByEnum []interface{}

func init() {
	var res []JobSortBy
	if err := json.Unmarshal([]byte(`["ID","CreatedDate","UpdatedDate","JobStatus,"""]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		jobSortByEnum = append(jobSortByEnum, v)
	}
}

func (m JobSortBy) validateJobSortByEnum(path, location string, value JobSortBy) error {
	if err := validate.Enum(path, location, value, jobSortByEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this job sort by
func (m JobSortBy) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateJobSortByEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
