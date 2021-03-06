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

// ArtifactType denotes type and function of an artifact.
//
// * IMAGE - image file.
//
// * RESULT - all outputs, files built during tests, etc.
//
// * TEST - additional files uploaded by user for conducting test.
//
// * YAML - yaml file describing Weles Job.
//
// swagger:model ArtifactType
type ArtifactType string

const (

	// ArtifactTypeIMAGE captures enum value "IMAGE"
	ArtifactTypeIMAGE ArtifactType = "IMAGE"

	// ArtifactTypeRESULT captures enum value "RESULT"
	ArtifactTypeRESULT ArtifactType = "RESULT"

	// ArtifactTypeTEST captures enum value "TEST"
	ArtifactTypeTEST ArtifactType = "TEST"

	// ArtifactTypeYAML captures enum value "YAML"
	ArtifactTypeYAML ArtifactType = "YAML"
)

// for schema
var artifactTypeEnum []interface{}

func init() {
	var res []ArtifactType
	if err := json.Unmarshal([]byte(`["IMAGE","RESULT","TEST","YAML",""]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		artifactTypeEnum = append(artifactTypeEnum, v)
	}
}

func (m ArtifactType) validateArtifactTypeEnum(path, location string, value ArtifactType) error {
	if err := validate.Enum(path, location, value, artifactTypeEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this artifact type
func (m ArtifactType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateArtifactTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
