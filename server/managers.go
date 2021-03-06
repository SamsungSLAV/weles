// Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
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

package server

import (
	"github.com/SamsungSLAV/weles"
)

// Managers provide implementation of JobManager and ArtifactManager interfaces.
type Managers struct {
	JM weles.JobManager
	AM weles.ArtifactManager
}

// APIDefaults contains interface implementations (Managers) and default values
// (set via CLI flags) for the API.
type APIDefaults struct {
	Managers  *Managers
	PageLimit int32
}

// NewManagers creates managers struct and assigns JobManager and ArtifactManager implementation
// to it.
func NewManagers(jm weles.JobManager, am weles.ArtifactManager) (m *Managers) {
	return &Managers{JM: jm, AM: am}
}
