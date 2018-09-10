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
	"github.com/go-openapi/validate"
)

// ArtifactStatus describes artifact status and availability.
//
// * DOWNLOADING - artifact is currently being downloaded.
//
// * READY - artifact has been downloaded and is ready to use.
//
// * FAILED - file is not available for use (e.g. download failed).
//
// * PENDING - artifact download has not started yet.
//
// swagger:model ArtifactStatus
type ArtifactStatus string

const (

	// ArtifactStatusDOWNLOADING captures enum value "DOWNLOADING"
	ArtifactStatusDOWNLOADING ArtifactStatus = "DOWNLOADING"

	// ArtifactStatusREADY captures enum value "READY"
	ArtifactStatusREADY ArtifactStatus = "READY"

	// ArtifactStatusFAILED captures enum value "FAILED"
	ArtifactStatusFAILED ArtifactStatus = "FAILED"

	// ArtifactStatusPENDING captures enum value "PENDING"
	ArtifactStatusPENDING ArtifactStatus = "PENDING"
)

// for schema
var artifactStatusEnum []interface{}

func init() {
	var res []ArtifactStatus
	if err := json.Unmarshal([]byte(`["DOWNLOADING","READY","FAILED","PENDING",""]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		artifactStatusEnum = append(artifactStatusEnum, v)
	}
}

func (m ArtifactStatus) validateArtifactStatusEnum(path, location string, value ArtifactStatus) error {
	if err := validate.Enum(path, location, value, artifactStatusEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this artifact status
func (m ArtifactStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateArtifactStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
