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

	middleware "github.com/go-openapi/runtime/middleware"
)

// JobCancelerHandlerFunc turns a function with the right signature into a job canceler handler
type JobCancelerHandlerFunc func(JobCancelerParams) middleware.Responder

// Handle executing the request and returning a response
func (fn JobCancelerHandlerFunc) Handle(params JobCancelerParams) middleware.Responder {
	return fn(params)
}

// JobCancelerHandler interface for that can handle valid job canceler params
type JobCancelerHandler interface {
	Handle(JobCancelerParams) middleware.Responder
}

// NewJobCanceler creates a new http.Handler for the job canceler operation
func NewJobCanceler(ctx *middleware.Context, handler JobCancelerHandler) *JobCanceler {
	return &JobCanceler{Context: ctx, Handler: handler}
}

/*JobCanceler swagger:route POST /jobs/{JobID}/cancel jobs jobCanceler

Cancel a Job

Stop execution of Job identified by JobID. Returns 204 on success. If
Job does not exist, 404 response will be returned. If the Job is
already in final state, 403 response will be returned.


*/
type JobCanceler struct {
	Context *middleware.Context
	Handler JobCancelerHandler
}

func (o *JobCanceler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewJobCancelerParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
