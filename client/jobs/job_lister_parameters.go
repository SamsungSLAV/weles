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
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewJobListerParams creates a new JobListerParams object
// with the default values initialized.
func NewJobListerParams() *JobListerParams {
	var ()
	return &JobListerParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewJobListerParamsWithTimeout creates a new JobListerParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewJobListerParamsWithTimeout(timeout time.Duration) *JobListerParams {
	var ()
	return &JobListerParams{

		timeout: timeout,
	}
}

// NewJobListerParamsWithContext creates a new JobListerParams object
// with the default values initialized, and the ability to set a context for a request
func NewJobListerParamsWithContext(ctx context.Context) *JobListerParams {
	var ()
	return &JobListerParams{

		Context: ctx,
	}
}

// NewJobListerParamsWithHTTPClient creates a new JobListerParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewJobListerParamsWithHTTPClient(client *http.Client) *JobListerParams {
	var ()
	return &JobListerParams{
		HTTPClient: client,
	}
}

/*JobListerParams contains all the parameters to send to the API endpoint
for the job lister operation typically these are written to a http.Request
*/
type JobListerParams struct {

	/*After
	  JobID of the last element from previous page.

	*/
	After *uint64
	/*Before
	  JobID of first element from next page.

	*/
	Before *uint64
	/*JobFilterAndSort
	  Job Filter and Sort object.

	*/
	JobFilterAndSort JobListerBody
	/*Limit
	  Custom page limit. Denotes number of JobInfo structures that will be returned.

	*/
	Limit *int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the job lister params
func (o *JobListerParams) WithTimeout(timeout time.Duration) *JobListerParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the job lister params
func (o *JobListerParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the job lister params
func (o *JobListerParams) WithContext(ctx context.Context) *JobListerParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the job lister params
func (o *JobListerParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the job lister params
func (o *JobListerParams) WithHTTPClient(client *http.Client) *JobListerParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the job lister params
func (o *JobListerParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAfter adds the after to the job lister params
func (o *JobListerParams) WithAfter(after *uint64) *JobListerParams {
	o.SetAfter(after)
	return o
}

// SetAfter adds the after to the job lister params
func (o *JobListerParams) SetAfter(after *uint64) {
	o.After = after
}

// WithBefore adds the before to the job lister params
func (o *JobListerParams) WithBefore(before *uint64) *JobListerParams {
	o.SetBefore(before)
	return o
}

// SetBefore adds the before to the job lister params
func (o *JobListerParams) SetBefore(before *uint64) {
	o.Before = before
}

// WithJobFilterAndSort adds the jobFilterAndSort to the job lister params
func (o *JobListerParams) WithJobFilterAndSort(jobFilterAndSort JobListerBody) *JobListerParams {
	o.SetJobFilterAndSort(jobFilterAndSort)
	return o
}

// SetJobFilterAndSort adds the jobFilterAndSort to the job lister params
func (o *JobListerParams) SetJobFilterAndSort(jobFilterAndSort JobListerBody) {
	o.JobFilterAndSort = jobFilterAndSort
}

// WithLimit adds the limit to the job lister params
func (o *JobListerParams) WithLimit(limit *int32) *JobListerParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the job lister params
func (o *JobListerParams) SetLimit(limit *int32) {
	o.Limit = limit
}

// WriteToRequest writes these params to a swagger request
func (o *JobListerParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.After != nil {

		// query param after
		var qrAfter uint64
		if o.After != nil {
			qrAfter = *o.After
		}
		qAfter := swag.FormatUint64(qrAfter)
		if qAfter != "" {
			if err := r.SetQueryParam("after", qAfter); err != nil {
				return err
			}
		}

	}

	if o.Before != nil {

		// query param before
		var qrBefore uint64
		if o.Before != nil {
			qrBefore = *o.Before
		}
		qBefore := swag.FormatUint64(qrBefore)
		if qBefore != "" {
			if err := r.SetQueryParam("before", qBefore); err != nil {
				return err
			}
		}

	}

	if err := r.SetBodyParam(o.JobFilterAndSort); err != nil {
		return err
	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int32
		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt32(qrLimit)
		if qLimit != "" {
			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
