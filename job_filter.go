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
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	enums "github.com/SamsungSLAV/weles/enums"
)

// JobFilter is used to filter Weles Jobs.i
// Filling more than one struct member (e.g.  JobID, Name) will result in
// searching for a Job with filled JobID and Name.
// Filling more than one member of an array/slice (in specific struct
// member) will result in searching for all members of array.
// Both aforementioned behaviours may occur simultainously. Some JobFilter
// fields support regular expressions (see fields documentation).
// swagger:model JobFilter
type JobFilter struct {

	// CreatedAfter is used to omit records created before supplied date.
	// Format: date-time
	CreatedAfter strfmt.DateTime `json:"CreatedAfter,omitempty"`

	// CreatedBefore is used to omit records created after supplied date.
	// Format: date-time
	CreatedBefore strfmt.DateTime `json:"CreatedBefore,omitempty"`

	// Info is used to filter by Job info (detailed information from Weles
	// about Job execution).
	// Allows usage of regular expressions.
	Info []string `json:"Info"`

	// JobID is used to filter Jobs by it's ID. Most commonly used filter.
	JobID []JobID `json:"JobID"`

	// Name is used to filter using name acquired form Job Submission file
	// (yaml format, job_name key's value).
	// Allows usage of regular expressions.
	Name []string `json:"Name"`

	// Status is used to receive only Jobs in specific status. When filled
	// with more than one element, returned jobs will only be in those
	// statuses.
	Status []enums.JobStatus `json:"Status"`

	// UpdatedAfter is used to omit records updated before supplied date.
	// Format: date-time
	UpdatedAfter strfmt.DateTime `json:"UpdatedAfter,omitempty"`

	// UpdatedBefore is used to omit records updated after supplied date.
	// Format: date-time
	UpdatedBefore strfmt.DateTime `json:"UpdatedBefore,omitempty"`
}

// Validate validates this job filter
func (m *JobFilter) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAfter(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedBefore(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateJobID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAfter(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedBefore(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JobFilter) validateCreatedAfter(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAfter) { // not required
		return nil
	}

	if err := validate.FormatOf("CreatedAfter", "body", "date-time", m.CreatedAfter.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *JobFilter) validateCreatedBefore(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedBefore) { // not required
		return nil
	}

	if err := validate.FormatOf("CreatedBefore", "body", "date-time", m.CreatedBefore.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *JobFilter) validateJobID(formats strfmt.Registry) error {

	if swag.IsZero(m.JobID) { // not required
		return nil
	}

	for i := 0; i < len(m.JobID); i++ {

		if err := m.JobID[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("JobID" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *JobFilter) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	for i := 0; i < len(m.Status); i++ {

		if err := m.Status[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Status" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *JobFilter) validateUpdatedAfter(formats strfmt.Registry) error {

	if swag.IsZero(m.UpdatedAfter) { // not required
		return nil
	}

	if err := validate.FormatOf("UpdatedAfter", "body", "date-time", m.UpdatedAfter.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *JobFilter) validateUpdatedBefore(formats strfmt.Registry) error {

	if swag.IsZero(m.UpdatedBefore) { // not required
		return nil
	}

	if err := validate.FormatOf("UpdatedBefore", "body", "date-time", m.UpdatedBefore.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *JobFilter) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JobFilter) UnmarshalBinary(b []byte) error {
	var res JobFilter
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
