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
	"github.com/go-openapi/validate"

	enums "github.com/SamsungSLAV/weles/enums"
)

// ArtifactInfoExt contains public information about single Artifact stored in ArtifactDB.
// swagger:model ArtifactInfoExt
type ArtifactInfoExt struct {
	ArtifactDescription

	// ID is unique identifier of an Artifact.
	ID int64 `json:"ID" db:",primarykey, autoincrement"`

	// Status of Artifact. For details - see documentation of
	// ArtifactStatus.
	Status enums.ArtifactStatus `json:"Status,omitempty"`

	// Timestamp is the date of creating an Artifact.
	// Format: date-time
	Timestamp strfmt.DateTime `json:"Timestamp"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *ArtifactInfoExt) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 ArtifactDescription
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.ArtifactDescription = aO0

	// now for regular properties
	var propsArtifactInfoExt struct {
		ID int64 `json:"ID"`

		Status enums.ArtifactStatus `json:"Status,omitempty"`

		Timestamp strfmt.DateTime `json:"Timestamp"`
	}
	if err := swag.ReadJSON(raw, &propsArtifactInfoExt); err != nil {
		return err
	}
	m.ID = propsArtifactInfoExt.ID

	m.Status = propsArtifactInfoExt.Status

	m.Timestamp = propsArtifactInfoExt.Timestamp

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m ArtifactInfoExt) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 1)

	aO0, err := swag.WriteJSON(m.ArtifactDescription)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	// now for regular properties
	var propsArtifactInfoExt struct {
		ID int64 `json:"ID"`

		Status enums.ArtifactStatus `json:"Status,omitempty"`

		Timestamp strfmt.DateTime `json:"Timestamp"`
	}
	propsArtifactInfoExt.ID = m.ID

	propsArtifactInfoExt.Status = m.Status

	propsArtifactInfoExt.Timestamp = m.Timestamp

	jsonDataPropsArtifactInfoExt, errArtifactInfoExt := swag.WriteJSON(propsArtifactInfoExt)
	if errArtifactInfoExt != nil {
		return nil, errArtifactInfoExt
	}
	_parts = append(_parts, jsonDataPropsArtifactInfoExt)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this artifact info ext
func (m *ArtifactInfoExt) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with ArtifactDescription
	if err := m.ArtifactDescription.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ArtifactInfoExt) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := m.Status.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Status")
		}
		return err
	}

	return nil
}

func (m *ArtifactInfoExt) validateTimestamp(formats strfmt.Registry) error {

	if swag.IsZero(m.Timestamp) { // not required
		return nil
	}

	if err := validate.FormatOf("Timestamp", "body", "date-time", m.Timestamp.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ArtifactInfoExt) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ArtifactInfoExt) UnmarshalBinary(b []byte) error {
	var res ArtifactInfoExt
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}