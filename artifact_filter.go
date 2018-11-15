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
)

// ArtifactFilter is used to filter results from ArtifactDB.
// swagger:model ArtifactFilter
type ArtifactFilter struct {

	// alias
	Alias []ArtifactAlias `json:"Alias"`

	// ID
	ID []int64 `json:"ID"`

	// job ID
	JobID []JobID `json:"JobID"`

	// status
	Status []ArtifactStatus `json:"Status"`

	// type
	Type []ArtifactType `json:"Type"`
}

// Validate validates this artifact filter
func (m *ArtifactFilter) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlias(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateJobID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ArtifactFilter) validateAlias(formats strfmt.Registry) error {

	if swag.IsZero(m.Alias) { // not required
		return nil
	}

	for i := 0; i < len(m.Alias); i++ {

		if err := m.Alias[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Alias" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *ArtifactFilter) validateJobID(formats strfmt.Registry) error {

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

func (m *ArtifactFilter) validateStatus(formats strfmt.Registry) error {

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

func (m *ArtifactFilter) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	for i := 0; i < len(m.Type); i++ {

		if err := m.Type[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("Type" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ArtifactFilter) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ArtifactFilter) UnmarshalBinary(b []byte) error {
	var res ArtifactFilter
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
