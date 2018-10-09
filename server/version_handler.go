// Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package server

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/server/operations/general"
)

// Version is Weles version information API endpoint handler.
func (a *APIDefaults) Version(params general.VersionParams) middleware.Responder {
	var v = &weles.Version{
		Server: weles.SrvVersion,
		API:    apiVersion,
		State:  apiState,
	}

	return general.NewVersionOK().WithPayload(v).
		WithWelesAPIState(v.State).WithWelesAPIVersion(v.API).WithWelesServerVersion(v.Server)
}
