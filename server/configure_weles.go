// This file is safe to edit. Once it exists it will not be overwritten

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

package server

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/server/operations"
	"github.com/SamsungSLAV/weles/server/operations/artifacts"
	"github.com/SamsungSLAV/weles/server/operations/general"
	"github.com/SamsungSLAV/weles/server/operations/jobs"
)

const (
	apiVersion = "0.1.0"
	apiState   = weles.VersionStateDevel
)

func configureFlags(api *operations.WelesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

// WelesConfigureAPI configures the API and handlers.
func (s *Server) WelesConfigureAPI(a *APIDefaults) {
	if s.api != nil {
		s.handler = welesConfigureAPI(s.api, a)
	}
}

func welesConfigureAPI(api *operations.WelesAPI, a *APIDefaults) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	api.SetDefaultProduces("application/json")
	api.SetDefaultConsumes("application/json")

	api.JobsJobCreatorHandler = jobs.JobCreatorHandlerFunc(a.Managers.JobCreator)
	api.JobsJobCancelerHandler = jobs.JobCancelerHandlerFunc(a.Managers.JobCanceller)
	api.JobsJobListerHandler = jobs.JobListerHandlerFunc(a.JobLister)

	api.ArtifactsArtifactListerHandler = artifacts.ArtifactListerHandlerFunc(a.ArtifactLister)

	api.GeneralVersionHandler = general.VersionHandlerFunc(a.Version)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later,
// this is the place. This function can be called multiple times, depending on the number of serving
// schemes. Scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the
// swagger.json document. The middleware executes after routing but before authentication,
// binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to
// serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func configureAPI(api *operations.WelesAPI) http.Handler {
	// WARNING
	// as go-swagger generated code (server.go) includes calls to this function its definition
	// must be present. This function should not be called anywhere.

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}
