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

package jobs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
)

// JobCancelerOptionsOKCode is the HTTP code returned for type JobCancelerOptionsOK
const JobCancelerOptionsOKCode int = 200

/*JobCancelerOptionsOK 200 OK

swagger:response jobCancelerOptionsOK
*/
type JobCancelerOptionsOK struct {
	/*

	 */
	AccessControlAllowHeaders []string `json:"Access-Control-Allow-Headers"`
	/*

	 */
	AccessControlAllowMethods []string `json:"Access-Control-Allow-Methods"`
	/*

	 */
	AccessControlAllowOrigin string `json:"Access-Control-Allow-Origin"`
}

// NewJobCancelerOptionsOK creates JobCancelerOptionsOK with default headers values
func NewJobCancelerOptionsOK() *JobCancelerOptionsOK {

	return &JobCancelerOptionsOK{}
}

// WithAccessControlAllowHeaders adds the accessControlAllowHeaders to the job canceler options o k response
func (o *JobCancelerOptionsOK) WithAccessControlAllowHeaders(accessControlAllowHeaders []string) *JobCancelerOptionsOK {
	o.AccessControlAllowHeaders = accessControlAllowHeaders
	return o
}

// SetAccessControlAllowHeaders sets the accessControlAllowHeaders to the job canceler options o k response
func (o *JobCancelerOptionsOK) SetAccessControlAllowHeaders(accessControlAllowHeaders []string) {
	o.AccessControlAllowHeaders = accessControlAllowHeaders
}

// WithAccessControlAllowMethods adds the accessControlAllowMethods to the job canceler options o k response
func (o *JobCancelerOptionsOK) WithAccessControlAllowMethods(accessControlAllowMethods []string) *JobCancelerOptionsOK {
	o.AccessControlAllowMethods = accessControlAllowMethods
	return o
}

// SetAccessControlAllowMethods sets the accessControlAllowMethods to the job canceler options o k response
func (o *JobCancelerOptionsOK) SetAccessControlAllowMethods(accessControlAllowMethods []string) {
	o.AccessControlAllowMethods = accessControlAllowMethods
}

// WithAccessControlAllowOrigin adds the accessControlAllowOrigin to the job canceler options o k response
func (o *JobCancelerOptionsOK) WithAccessControlAllowOrigin(accessControlAllowOrigin string) *JobCancelerOptionsOK {
	o.AccessControlAllowOrigin = accessControlAllowOrigin
	return o
}

// SetAccessControlAllowOrigin sets the accessControlAllowOrigin to the job canceler options o k response
func (o *JobCancelerOptionsOK) SetAccessControlAllowOrigin(accessControlAllowOrigin string) {
	o.AccessControlAllowOrigin = accessControlAllowOrigin
}

// WriteResponse to the client
func (o *JobCancelerOptionsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Access-Control-Allow-Headers

	var accessControlAllowHeadersIR []string
	for _, accessControlAllowHeadersI := range o.AccessControlAllowHeaders {
		accessControlAllowHeadersIS := accessControlAllowHeadersI
		if accessControlAllowHeadersIS != "" {
			accessControlAllowHeadersIR = append(accessControlAllowHeadersIR, accessControlAllowHeadersIS)
		}
	}
	accessControlAllowHeaders := swag.JoinByFormat(accessControlAllowHeadersIR, "csv")
	if len(accessControlAllowHeaders) > 0 {
		hv := accessControlAllowHeaders[0]
		if hv != "" {
			rw.Header().Set("Access-Control-Allow-Headers", hv)
		}
	}

	// response header Access-Control-Allow-Methods

	var accessControlAllowMethodsIR []string
	for _, accessControlAllowMethodsI := range o.AccessControlAllowMethods {
		accessControlAllowMethodsIS := accessControlAllowMethodsI
		if accessControlAllowMethodsIS != "" {
			accessControlAllowMethodsIR = append(accessControlAllowMethodsIR, accessControlAllowMethodsIS)
		}
	}
	accessControlAllowMethods := swag.JoinByFormat(accessControlAllowMethodsIR, "csv")
	if len(accessControlAllowMethods) > 0 {
		hv := accessControlAllowMethods[0]
		if hv != "" {
			rw.Header().Set("Access-Control-Allow-Methods", hv)
		}
	}

	// response header Access-Control-Allow-Origin

	accessControlAllowOrigin := o.AccessControlAllowOrigin
	if accessControlAllowOrigin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", accessControlAllowOrigin)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}