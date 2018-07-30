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

// JobCreatorHandlerFunc turns a function with the right signature into a job creator handler
type JobCreatorHandlerFunc func(JobCreatorParams) middleware.Responder

// Handle executing the request and returning a response
func (fn JobCreatorHandlerFunc) Handle(params JobCreatorParams) middleware.Responder {
	return fn(params)
}

// JobCreatorHandler interface for that can handle valid job creator params
type JobCreatorHandler interface {
	Handle(JobCreatorParams) middleware.Responder
}

// NewJobCreator creates a new http.Handler for the job creator operation
func NewJobCreator(ctx *middleware.Context, handler JobCreatorHandler) *JobCreator {
	return &JobCreator{Context: ctx, Handler: handler}
}

/*JobCreator swagger:route POST /jobs jobs jobCreator

Add new job

adds new Job in Weles using recipe passed in YAML format.

*/
type JobCreator struct {
	Context *middleware.Context
	Handler JobCreatorHandler
}

func (o *JobCreator) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewJobCreatorParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
