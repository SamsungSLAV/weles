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

// ArtifactDescription contains information needed to create new artifact in ArtifactDB.
// swagger:model ArtifactDescription
type ArtifactDescription struct {

	// alias
	Alias ArtifactAlias `json:"Alias,omitempty"`

	// specifies  Job for which artifact was created.
	JobID JobID `json:"JobID,omitempty"`

	// type
	Type ArtifactType `json:"Type,omitempty"`

	// URI
	// Format: uri
	URI ArtifactURI `json:"URI,omitempty"`
}

// Validate validates this artifact description
func (m *ArtifactDescription) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlias(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateJobID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateURI(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ArtifactDescription) validateAlias(formats strfmt.Registry) error {

	if swag.IsZero(m.Alias) { // not required
		return nil
	}

	if err := m.Alias.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Alias")
		}
		return err
	}

	return nil
}

func (m *ArtifactDescription) validateJobID(formats strfmt.Registry) error {

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

func (m *ArtifactDescription) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	if err := m.Type.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Type")
		}
		return err
	}

	return nil
}

func (m *ArtifactDescription) validateURI(formats strfmt.Registry) error {

	if swag.IsZero(m.URI) { // not required
		return nil
	}

	if err := m.URI.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("URI")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ArtifactDescription) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ArtifactDescription) UnmarshalBinary(b []byte) error {
	var res ArtifactDescription
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
