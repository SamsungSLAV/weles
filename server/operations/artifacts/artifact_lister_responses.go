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
	/*URI to request next page of data. Please note that the same body must be used as in initial request.


	 */
	Next string `json:"Next"`
	/*URI to request next page of data. Please note that the same body must be used as in initial request.


	 */
	Previous string `json:"Previous"`
	/*count of records currently fulfilling the requested ArtifactFilter. Please note that this value may change when requesting for the same data at a different moment in time.


	 */
	TotalRecords uint64 `json:"TotalRecords"`

	/*
	  In: Body
	*/
	Payload []*weles.ArtifactInfo `json:"body,omitempty"`
}

// NewArtifactListerOK creates ArtifactListerOK with default headers values
func NewArtifactListerOK() *ArtifactListerOK {

	return &ArtifactListerOK{}
}

// WithNext adds the next to the artifact lister o k response
func (o *ArtifactListerOK) WithNext(next string) *ArtifactListerOK {
	o.Next = next
	return o
}

// SetNext sets the next to the artifact lister o k response
func (o *ArtifactListerOK) SetNext(next string) {
	o.Next = next
}

// WithPrevious adds the previous to the artifact lister o k response
func (o *ArtifactListerOK) WithPrevious(previous string) *ArtifactListerOK {
	o.Previous = previous
	return o
}

// SetPrevious sets the previous to the artifact lister o k response
func (o *ArtifactListerOK) SetPrevious(previous string) {
	o.Previous = previous
}

// WithTotalRecords adds the totalRecords to the artifact lister o k response
func (o *ArtifactListerOK) WithTotalRecords(totalRecords uint64) *ArtifactListerOK {
	o.TotalRecords = totalRecords
	return o
}

// SetTotalRecords sets the totalRecords to the artifact lister o k response
func (o *ArtifactListerOK) SetTotalRecords(totalRecords uint64) {
	o.TotalRecords = totalRecords
}

// WithPayload adds the payload to the artifact lister o k response
func (o *ArtifactListerOK) WithPayload(payload []*weles.ArtifactInfo) *ArtifactListerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the artifact lister o k response
func (o *ArtifactListerOK) SetPayload(payload []*weles.ArtifactInfo) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ArtifactListerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Next

	next := o.Next
	if next != "" {
		rw.Header().Set("Next", next)
	}

	// response header Previous

	previous := o.Previous
	if previous != "" {
		rw.Header().Set("Previous", previous)
	}

	// response header TotalRecords

	totalRecords := swag.FormatUint64(o.TotalRecords)
	if totalRecords != "" {
		rw.Header().Set("TotalRecords", totalRecords)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*weles.ArtifactInfo, 0, 50)
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
	/*URI to request next page of data. Please note that the same body must be used as in initial request.


	 */
	Next string `json:"Next"`
	/*URI to request next page of data. Please note that the same body must be used as in initial request.


	 */
	Previous string `json:"Previous"`
	/*number of records after current page. Please note that this value may change when requesting for the same data at a different moment in time.


	 */
	RemainingRecords uint64 `json:"RemainingRecords"`
	/*count of records currently fulfilling the requested ArtifactFilter. Please note that this value may change when requesting for the same data at a different moment in time.


	 */
	TotalRecords uint64 `json:"TotalRecords"`

	/*
	  In: Body
	*/
	Payload []*weles.ArtifactInfo `json:"body,omitempty"`
}

// NewArtifactListerPartialContent creates ArtifactListerPartialContent with default headers values
func NewArtifactListerPartialContent() *ArtifactListerPartialContent {

	return &ArtifactListerPartialContent{}
}

// WithNext adds the next to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithNext(next string) *ArtifactListerPartialContent {
	o.Next = next
	return o
}

// SetNext sets the next to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetNext(next string) {
	o.Next = next
}

// WithPrevious adds the previous to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithPrevious(previous string) *ArtifactListerPartialContent {
	o.Previous = previous
	return o
}

// SetPrevious sets the previous to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetPrevious(previous string) {
	o.Previous = previous
}

// WithRemainingRecords adds the remainingRecords to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithRemainingRecords(remainingRecords uint64) *ArtifactListerPartialContent {
	o.RemainingRecords = remainingRecords
	return o
}

// SetRemainingRecords sets the remainingRecords to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetRemainingRecords(remainingRecords uint64) {
	o.RemainingRecords = remainingRecords
}

// WithTotalRecords adds the totalRecords to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithTotalRecords(totalRecords uint64) *ArtifactListerPartialContent {
	o.TotalRecords = totalRecords
	return o
}

// SetTotalRecords sets the totalRecords to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetTotalRecords(totalRecords uint64) {
	o.TotalRecords = totalRecords
}

// WithPayload adds the payload to the artifact lister partial content response
func (o *ArtifactListerPartialContent) WithPayload(payload []*weles.ArtifactInfo) *ArtifactListerPartialContent {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the artifact lister partial content response
func (o *ArtifactListerPartialContent) SetPayload(payload []*weles.ArtifactInfo) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ArtifactListerPartialContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Next

	next := o.Next
	if next != "" {
		rw.Header().Set("Next", next)
	}

	// response header Previous

	previous := o.Previous
	if previous != "" {
		rw.Header().Set("Previous", previous)
	}

	// response header RemainingRecords

	remainingRecords := swag.FormatUint64(o.RemainingRecords)
	if remainingRecords != "" {
		rw.Header().Set("RemainingRecords", remainingRecords)
	}

	// response header TotalRecords

	totalRecords := swag.FormatUint64(o.TotalRecords)
	if totalRecords != "" {
		rw.Header().Set("TotalRecords", totalRecords)
	}

	rw.WriteHeader(206)
	payload := o.Payload
	if payload == nil {
		payload = make([]*weles.ArtifactInfo, 0, 50)
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
