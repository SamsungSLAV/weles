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

package artifacts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	weles "github.com/SamsungSLAV/weles"
)

// ArtifactListerOKCode is the HTTP code returned for type ArtifactListerOK
const ArtifactListerOKCode int = 200

/*ArtifactListerOK OK

swagger:response artifactListerOK
*/
type ArtifactListerOK struct {
	/*Weles-List-Batch-Size is the count of records returned on
	current page.

	*/
	WelesListBatchSize int32 `json:"Weles-List-Batch-Size"`
	/*Weles-List-Total is count of records currently fulfilling the
	requested ArtifactFilter. Please note that this value may
	change when requesting for the same data at a different moment
	in time.

	*/
	WelesListTotal uint64 `json:"Weles-List-Total"`
	/*Weles-Next-Page is the URL suffix to request next page of data.
	Please note that the same body must be used as in initial
	request.

	*/
	WelesNextPage string `json:"Weles-Next-Page"`
	/*Weles-Previous-Page is the URL suffix to request next page of
	data.  Please note that the same body must be used as in
	initial request.

	*/
	WelesPreviousPage string `json:"Weles-Previous-Page"`

	/*
	  In: Body
	*/
	Payload []*weles.ArtifactInfoExt `json:"body,omitempty"`
}

// NewArtifactListerOK creates ArtifactListerOK with default headers values
func NewArtifactListerOK() *ArtifactListerOK {

	return &ArtifactListerOK{}
}

// WithWelesListBatchSize adds the welesListBatchSize to the artifact lister o k response
func (o *ArtifactListerOK) WithWelesListBatchSize(welesListBatchSize int32) *ArtifactListerOK {
	o.WelesListBatchSize = welesListBatchSize
	return o
}

// SetWelesListBatchSize sets the welesListBatchSize to the artifact lister o k response
func (o *ArtifactListerOK) SetWelesListBatchSize(welesListBatchSize int32) {
	o.WelesListBatchSize = welesListBatchSize
}

// WithWelesListTotal adds the welesListTotal to the artifact lister o k response
func (o *ArtifactListerOK) WithWelesListTotal(welesListTotal uint64) *ArtifactListerOK {
	o.WelesListTotal = welesListTotal
	return o
}

// SetWelesListTotal sets the welesListTotal to the artifact lister o k response
func (o *ArtifactListerOK) SetWelesListTotal(welesListTotal uint64) {
	o.WelesListTotal = welesListTotal
}

// WithWelesNextPage adds the welesNextPage to the artifact lister o k response
func (o *ArtifactListerOK) WithWelesNextPage(welesNextPage string) *ArtifactListerOK {
	o.WelesNextPage = welesNextPage
	return o
}

// SetWelesNextPage sets the welesNextPage to the artifact lister o k response
func (o *ArtifactListerOK) SetWelesNextPage(welesNextPage string) {
	o.WelesNextPage = welesNextPage
}

// WithWelesPreviousPage adds the welesPreviousPage to the artifact lister o k response
func (o *ArtifactListerOK) WithWelesPreviousPage(welesPreviousPage string) *ArtifactListerOK {
	o.WelesPreviousPage = welesPreviousPage
	return o
}

// SetWelesPreviousPage sets the welesPreviousPage to the artifact lister o k response
func (o *ArtifactListerOK) SetWelesPreviousPage(welesPreviousPage string) {
	o.WelesPreviousPage = welesPreviousPage
}

// WithPayload adds the payload to the artifact lister o k response
func (o *ArtifactListerOK) WithPayload(payload []*weles.ArtifactInfoExt) *ArtifactListerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the artifact lister o k response
func (o *ArtifactListerOK) SetPayload(payload []*weles.ArtifactInfoExt) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ArtifactListerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Weles-List-Batch-Size

	welesListBatchSize := swag.FormatInt32(o.WelesListBatchSize)
	if welesListBatchSize != "" {
		rw.Header().Set("Weles-List-Batch-Size", welesListBatchSize)
	}

	// response header Weles-List-Total

	welesListTotal := swag.FormatUint64(o.WelesListTotal)
	if welesListTotal != "" {
		rw.Header().Set("Weles-List-Total", welesListTotal)
	}

	// response header Weles-Next-Page

	welesNextPage := o.WelesNextPage
	if welesNextPage != "" {
		rw.Header().Set("Weles-Next-Page", welesNextPage)
	}

	// response header Weles-Previous-Page

	welesPreviousPage := o.WelesPreviousPage
	if welesPreviousPage != "" {
		rw.Header().Set("Weles-Previous-Page", welesPreviousPage)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*weles.ArtifactInfoExt, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// ArtifactListerPartialContentCode is the HTTP code returned for type ArtifactListerPartialContent
const ArtifactListerPartialContentCode int = 206

/*ArtifactListerPartialContent Partial Content

swagger:response artifactListerPartialContent
*/
type ArtifactListerPartialContent struct {
	/*Weles-List-Batch-Size is the count of records returned on
	current page.

	*/
	WelesListBatchSize int32 `json:"Weles-List-Batch-Size"`
	/*Weles-List-Remaining is number of records after current page.
	Please note that this value may change when requesting for the
	same data at a different moment in time.

	*/
	WelesListRemaining uint64 `json:"Weles-List-Remaining"`
	/*Weles-List-Total is count of records currently fulfilling the
	requested ArtifactFilter. Please note that this value may
	change when requesting for the same data at a different moment
	in time.

	*/
	WelesListTotal uint64 `json:"Weles-List-Total"`
	/*Weles-Next-Page is URL to request next page of data. Please
	note that the same body must be used as in initial request.

	*/
	WelesNextPage string `json:"Weles-Next-Page"`
	/*Weles-Previous-Page is URL suffix to request next page of data.
	Please note that the same body must be used as in initial
	request.

	*/
	WelesPreviousPage string `json:"Weles-Previous-Page"`

	/*
	  In: Body
	*/
	Payload []*weles.ArtifactInfoExt `json:"body,omitempty"`
}

// NewArtifactListerPartialContent creates ArtifactListerPartialContent with default headers values
func NewArtifactListerPartialContent() *ArtifactListerPartialContent {

	return &ArtifactListerPartialContent{}
}

// WithWelesListBatchSize adds the welesListBatchSize to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithWelesListBatchSize(welesListBatchSize int32) *ArtifactListerPartialContent {
	o.WelesListBatchSize = welesListBatchSize
	return o
}

// SetWelesListBatchSize sets the welesListBatchSize to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetWelesListBatchSize(welesListBatchSize int32) {
	o.WelesListBatchSize = welesListBatchSize
}

// WithWelesListRemaining adds the welesListRemaining to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithWelesListRemaining(welesListRemaining uint64) *ArtifactListerPartialContent {
	o.WelesListRemaining = welesListRemaining
	return o
}

// SetWelesListRemaining sets the welesListRemaining to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetWelesListRemaining(welesListRemaining uint64) {
	o.WelesListRemaining = welesListRemaining
}

// WithWelesListTotal adds the welesListTotal to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithWelesListTotal(welesListTotal uint64) *ArtifactListerPartialContent {
	o.WelesListTotal = welesListTotal
	return o
}

// SetWelesListTotal sets the welesListTotal to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetWelesListTotal(welesListTotal uint64) {
	o.WelesListTotal = welesListTotal
}

// WithWelesNextPage adds the welesNextPage to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithWelesNextPage(welesNextPage string) *ArtifactListerPartialContent {
	o.WelesNextPage = welesNextPage
	return o
}

// SetWelesNextPage sets the welesNextPage to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetWelesNextPage(welesNextPage string) {
	o.WelesNextPage = welesNextPage
}

// WithWelesPreviousPage adds the welesPreviousPage to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithWelesPreviousPage(welesPreviousPage string) *ArtifactListerPartialContent {
	o.WelesPreviousPage = welesPreviousPage
	return o
}

// SetWelesPreviousPage sets the welesPreviousPage to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetWelesPreviousPage(welesPreviousPage string) {
	o.WelesPreviousPage = welesPreviousPage
}

// WithPayload adds the payload to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithPayload(payload []*weles.ArtifactInfoExt) *ArtifactListerPartialContent {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetPayload(payload []*weles.ArtifactInfoExt) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ArtifactListerPartialContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Weles-List-Batch-Size

	welesListBatchSize := swag.FormatInt32(o.WelesListBatchSize)
	if welesListBatchSize != "" {
		rw.Header().Set("Weles-List-Batch-Size", welesListBatchSize)
	}

	// response header Weles-List-Remaining

	welesListRemaining := swag.FormatUint64(o.WelesListRemaining)
	if welesListRemaining != "" {
		rw.Header().Set("Weles-List-Remaining", welesListRemaining)
	}

	// response header Weles-List-Total

	welesListTotal := swag.FormatUint64(o.WelesListTotal)
	if welesListTotal != "" {
		rw.Header().Set("Weles-List-Total", welesListTotal)
	}

	// response header Weles-Next-Page

	welesNextPage := o.WelesNextPage
	if welesNextPage != "" {
		rw.Header().Set("Weles-Next-Page", welesNextPage)
	}

	// response header Weles-Previous-Page

	welesPreviousPage := o.WelesPreviousPage
	if welesPreviousPage != "" {
		rw.Header().Set("Weles-Previous-Page", welesPreviousPage)
	}

	rw.WriteHeader(206)
	payload := o.Payload
	if payload == nil {
		payload = make([]*weles.ArtifactInfoExt, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// ArtifactListerBadRequestCode is the HTTP code returned for type ArtifactListerBadRequest
const ArtifactListerBadRequestCode int = 400

/*ArtifactListerBadRequest Bad Request

swagger:response artifactListerBadRequest
*/
type ArtifactListerBadRequest struct {

	/*
	  In: Body
	*/
	Payload *weles.ErrResponse `json:"body,omitempty"`
}

// NewArtifactListerBadRequest creates ArtifactListerBadRequest with default headers values
func NewArtifactListerBadRequest() *ArtifactListerBadRequest {

	return &ArtifactListerBadRequest{}
}

// WithPayload adds the payload to the artifact lister bad request response
func (o *ArtifactListerBadRequest) WithPayload(payload *weles.ErrResponse) *ArtifactListerBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the artifact lister bad request response
func (o *ArtifactListerBadRequest) SetPayload(payload *weles.ErrResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ArtifactListerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ArtifactListerNotFoundCode is the HTTP code returned for type ArtifactListerNotFound
const ArtifactListerNotFoundCode int = 404

/*ArtifactListerNotFound Not Found

swagger:response artifactListerNotFound
*/
type ArtifactListerNotFound struct {

	/*
	  In: Body
	*/
	Payload *weles.ErrResponse `json:"body,omitempty"`
}

// NewArtifactListerNotFound creates ArtifactListerNotFound with default headers values
func NewArtifactListerNotFound() *ArtifactListerNotFound {

	return &ArtifactListerNotFound{}
}

// WithPayload adds the payload to the artifact lister not found response
func (o *ArtifactListerNotFound) WithPayload(payload *weles.ErrResponse) *ArtifactListerNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the artifact lister not found response
func (o *ArtifactListerNotFound) SetPayload(payload *weles.ErrResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ArtifactListerNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ArtifactListerInternalServerErrorCode is the HTTP code returned for type ArtifactListerInternalServerError
const ArtifactListerInternalServerErrorCode int = 500

/*ArtifactListerInternalServerError Internal Server error

swagger:response artifactListerInternalServerError
*/
type ArtifactListerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *weles.ErrResponse `json:"body,omitempty"`
}

// NewArtifactListerInternalServerError creates ArtifactListerInternalServerError with default headers values
func NewArtifactListerInternalServerError() *ArtifactListerInternalServerError {

	return &ArtifactListerInternalServerError{}
}

// WithPayload adds the payload to the artifact lister internal server error response
func (o *ArtifactListerInternalServerError) WithPayload(payload *weles.ErrResponse) *ArtifactListerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the artifact lister internal server error response
func (o *ArtifactListerInternalServerError) SetPayload(payload *weles.ErrResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ArtifactListerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
