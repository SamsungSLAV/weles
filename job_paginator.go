// Code generated by go-swagger; DO NOT EDIT.

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

package weles

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// JobPaginator JobPaginator holds information neccessary to request for a single page of
// data.
// When JobID is set, and Forward is false - Controller should return a page
// of records before the supplied JobID.
// When JobID is set, and Forward is true - Controller should return page
// of record after the supplied JobID.
// In both cases, returned page should not include supplied JobID.
// Limit denotes the number of records to be returned on the page. When
// Limit is set to 0, pagination is disabled, JobID and Forward fields are
// ignored and all records are returned.
// swagger:model JobPaginator
type JobPaginator struct {

	// Forward denotes direction of pagination.
	Forward bool `json:"Forward,omitempty"`

	// JobID is the key used for pagination.
	JobID JobID `json:"JobID,omitempty"`

	// Limit the page size.
	Limit int32 `json:"Limit,omitempty"`
}

// Validate validates this job paginator
func (m *JobPaginator) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateJobID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JobPaginator) validateJobID(formats strfmt.Registry) error {

	if swag.IsZero(m.JobID) { // not required
		return nil
	}

	if err := m.JobID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("JobID")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *JobPaginator) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JobPaginator) UnmarshalBinary(b []byte) error {
	var res JobPaginator
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
