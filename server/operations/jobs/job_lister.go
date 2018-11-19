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
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	weles "github.com/SamsungSLAV/weles"
)

// JobListerHandlerFunc turns a function with the right signature into a job lister handler
type JobListerHandlerFunc func(JobListerParams) middleware.Responder

// Handle executing the request and returning a response
func (fn JobListerHandlerFunc) Handle(params JobListerParams) middleware.Responder {
	return fn(params)
}

// JobListerHandler interface for that can handle valid job lister params
type JobListerHandler interface {
	Handle(JobListerParams) middleware.Responder
}

// NewJobLister creates a new http.Handler for the job lister operation
func NewJobLister(ctx *middleware.Context, handler JobListerHandler) *JobLister {
	return &JobLister{Context: ctx, Handler: handler}
}

/*JobLister swagger:route POST /jobs/list jobs jobLister

List Jobs with filtering, sorting and pagination.

Returns sorted list of Jobs. If there are more records than returned
page, 206 response is returned. If the page is last - 200 response is
returned. If no Jobs satisfy passed filter, 404 response is returned.
Filling both before and after query will result in 400 error response.

Providing empty body and no query parameter will result in list with
default values - no filter, sorted in Ascending order by JobID.
Check JobFilter and JobSorter models documentation to see how to use
them.
To ease automatic pagination, URL suffixes are returned with each
2xx response.

*/
type JobLister struct {
	Context *middleware.Context
	Handler JobListerHandler
}

func (o *JobLister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewJobListerParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// JobListerBody JobFilterAndSort contains data for filtering and sorting
// Weles Jobs lists. Please refer to JobFilter and
// JobSorter documentation for more details.
// swagger:model JobListerBody
type JobListerBody struct {

	// filter
	Filter *weles.JobFilter `json:"Filter,omitempty"`

	// sorter
	Sorter *weles.JobSorter `json:"Sorter,omitempty"`
}

// Validate validates this job lister body
func (o *JobListerBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateFilter(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateSorter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *JobListerBody) validateFilter(formats strfmt.Registry) error {

	if swag.IsZero(o.Filter) { // not required
		return nil
	}

	if o.Filter != nil {
		if err := o.Filter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("jobFilterAndSort" + "." + "Filter")
			}
			return err
		}
	}

	return nil
}

func (o *JobListerBody) validateSorter(formats strfmt.Registry) error {

	if swag.IsZero(o.Sorter) { // not required
		return nil
	}

	if o.Sorter != nil {
		if err := o.Sorter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("jobFilterAndSort" + "." + "Sorter")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *JobListerBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *JobListerBody) UnmarshalBinary(b []byte) error {
	var res JobListerBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
