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

package general

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	weles "github.com/SamsungSLAV/weles"
)

// VersionOKCode is the HTTP code returned for type VersionOK
const VersionOKCode int = 200

/*VersionOK OK

swagger:response versionOK
*/
type VersionOK struct {
	/*State of Weles API.

	 */
	APIState string `json:"API-State"`
	/*Version of Weles API.

	 */
	APIVersion string `json:"API-Version"`
	/*Version of Weles server.

	 */
	ServerVersion string `json:"Server-Version"`

	/*
	  In: Body
	*/
	Payload *weles.Version `json:"body,omitempty"`
}

// NewVersionOK creates VersionOK with default headers values
func NewVersionOK() *VersionOK {

	return &VersionOK{}
}

// WithAPIState adds the apiState to the version o k response
func (o *VersionOK) WithAPIState(aPIState string) *VersionOK {
	o.APIState = aPIState
	return o
}

// SetAPIState sets the apiState to the version o k response
func (o *VersionOK) SetAPIState(aPIState string) {
	o.APIState = aPIState
}

// WithAPIVersion adds the apiVersion to the version o k response
func (o *VersionOK) WithAPIVersion(aPIVersion string) *VersionOK {
	o.APIVersion = aPIVersion
	return o
}

// SetAPIVersion sets the apiVersion to the version o k response
func (o *VersionOK) SetAPIVersion(aPIVersion string) {
	o.APIVersion = aPIVersion
}

// WithServerVersion adds the serverVersion to the version o k response
func (o *VersionOK) WithServerVersion(serverVersion string) *VersionOK {
	o.ServerVersion = serverVersion
	return o
}

// SetServerVersion sets the serverVersion to the version o k response
func (o *VersionOK) SetServerVersion(serverVersion string) {
	o.ServerVersion = serverVersion
}

// WithPayload adds the payload to the version o k response
func (o *VersionOK) WithPayload(payload *weles.Version) *VersionOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the version o k response
func (o *VersionOK) SetPayload(payload *weles.Version) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *VersionOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header API-State

	aPIState := o.APIState
	if aPIState != "" {
		rw.Header().Set("API-State", aPIState)
	}

	// response header API-Version

	aPIVersion := o.APIVersion
	if aPIVersion != "" {
		rw.Header().Set("API-Version", aPIVersion)
	}

	// response header Server-Version

	serverVersion := o.ServerVersion
	if serverVersion != "" {
		rw.Header().Set("Server-Version", serverVersion)
	}

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// VersionInternalServerErrorCode is the HTTP code returned for type VersionInternalServerError
const VersionInternalServerErrorCode int = 500

/*VersionInternalServerError Internal Server error

swagger:response versionInternalServerError
*/
type VersionInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *weles.ErrResponse `json:"body,omitempty"`
}

// NewVersionInternalServerError creates VersionInternalServerError with default headers values
func NewVersionInternalServerError() *VersionInternalServerError {

	return &VersionInternalServerError{}
}

// WithPayload adds the payload to the version internal server error response
func (o *VersionInternalServerError) WithPayload(payload *weles.ErrResponse) *VersionInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the version internal server error response
func (o *VersionInternalServerError) SetPayload(payload *weles.ErrResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *VersionInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
