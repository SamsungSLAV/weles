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

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewJobListerParams creates a new JobListerParams object
// no default values defined in spec.
func NewJobListerParams() JobListerParams {

	return JobListerParams{}
}

// JobListerParams contains all the bound params for the job lister operation
// typically these are obtained from a http.Request
//
// swagger:parameters JobLister
type JobListerParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*After should be filled with JobID of the last element from current
	page to receive next one.
	  In: query
	*/
	After *uint64
	/*Before should be filled with JobID of the first element from
	current page to receive previous one.
	  In: query
	*/
	Before *uint64
	/*Job Filter and Sort object.
	  In: body
	*/
	JobFilterAndSort JobListerBody
	/*Limit is the number of records to return. Overrides default server
	page limit.
	  In: query
	*/
	Limit *int32
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewJobListerParams() beforehand.
func (o *JobListerParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qAfter, qhkAfter, _ := qs.GetOK("after")
	if err := o.bindAfter(qAfter, qhkAfter, route.Formats); err != nil {
		res = append(res, err)
	}

	qBefore, qhkBefore, _ := qs.GetOK("before")
	if err := o.bindBefore(qBefore, qhkBefore, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body JobListerBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("jobFilterAndSort", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.JobFilterAndSort = body
			}
		}
	}
	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAfter binds and validates parameter After from query.
func (o *JobListerParams) bindAfter(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertUint64(raw)
	if err != nil {
		return errors.InvalidType("after", "query", "uint64", raw)
	}
	o.After = &value

	return nil
}

// bindBefore binds and validates parameter Before from query.
func (o *JobListerParams) bindBefore(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertUint64(raw)
	if err != nil {
		return errors.InvalidType("before", "query", "uint64", raw)
	}
	o.Before = &value

	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *JobListerParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int32", raw)
	}
	o.Limit = &value

	return nil
}
