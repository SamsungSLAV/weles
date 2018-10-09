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
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Version defines version of Weles API (and its state) and server.
//
// swagger:model Version
type Version struct {

	// Version of Weles API.
	API string `json:"API,omitempty"`

	// Version of Weles server.
	Server string `json:"Server,omitempty"`

	// State of Weles API.
	// Enum: [devel stable deprecated]
	State string `json:"State,omitempty"`
}

// Validate validates this version
func (m *Version) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateState(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var versionTypeStatePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["devel","stable","deprecated"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		versionTypeStatePropEnum = append(versionTypeStatePropEnum, v)
	}
}

const (

	// VersionStateDevel captures enum value "devel"
	VersionStateDevel string = "devel"

	// VersionStateStable captures enum value "stable"
	VersionStateStable string = "stable"

	// VersionStateDeprecated captures enum value "deprecated"
	VersionStateDeprecated string = "deprecated"
)

// prop value enum
func (m *Version) validateStateEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, versionTypeStatePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Version) validateState(formats strfmt.Registry) error {

	if swag.IsZero(m.State) { // not required
		return nil
	}

	// value enum
	if err := m.validateStateEnum("State", "body", m.State); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Version) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Version) UnmarshalBinary(b []byte) error {
	var res Version
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
