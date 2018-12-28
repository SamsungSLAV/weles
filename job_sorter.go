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

	enums "github.com/SamsungSLAV/weles/enums"
)

// JobSorter defines the key for sorting as well as direction of sorting.
// swagger:model JobSorter
type JobSorter struct {

	// SortBy denotes key used for sorting Jobs. For more details, see JobSortBy
	// documentation.
	SortBy enums.JobSortBy `json:"SortBy,omitempty"`

	// SortOrder in which Jobs will be sorted. For more details, see SortOrder
	// documentation.
	SortOrder enums.SortOrder `json:"SortOrder,omitempty"`
}

// Validate validates this job sorter
func (m *JobSorter) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSortBy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSortOrder(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JobSorter) validateSortBy(formats strfmt.Registry) error {

	if swag.IsZero(m.SortBy) { // not required
		return nil
	}

	if err := m.SortBy.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("SortBy")
		}
		return err
	}

	return nil
}

func (m *JobSorter) validateSortOrder(formats strfmt.Registry) error {

	if swag.IsZero(m.SortOrder) { // not required
		return nil
	}

	if err := m.SortOrder.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("SortOrder")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *JobSorter) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JobSorter) UnmarshalBinary(b []byte) error {
	var res JobSorter
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
