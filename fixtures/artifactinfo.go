/*
 *  Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

package fixtures

import (
	"math/rand"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/tideland/golib/audit"

	"git.tizen.org/tools/weles"
)

const (
	dateLayout        = "Mon Jan 2 15:04:05 -0700 MST 2006"
	someDate          = "Tue Jan 2 15:04:05 +0100 CET 1900"
	durationIncrement = "+25h"
)

// CreateArtifactInfoSlice returns slice of ArtifactInfos of sliceLength length.
// It is filled with random data used for testing.
func CreateArtifactInfoSlice(sliceLength int) []weles.ArtifactInfo {
	// checking for errors omitted due to fixed input.
	dateTimeIter, _ := time.Parse(dateLayout, someDate)
	durationIncrement, _ := time.ParseDuration(durationIncrement)
	artifactInfo := make([]weles.ArtifactInfo, sliceLength)
	gen := audit.NewGenerator(rand.New(rand.NewSource(time.Now().UTC().UnixNano())))
	for i := range artifactInfo {
		tmp := weles.ArtifactInfo{}
		timestamp := gen.Time(time.Local, dateTimeIter, durationIncrement)
		tmp.Timestamp = strfmt.DateTime(timestamp)
		tmp.ID = int64(i + 1)
		tmp.ArtifactDescription.Alias = weles.ArtifactAlias(gen.Word())
		tmp.ArtifactDescription.JobID = weles.JobID(i + 1)
		tmp.ArtifactDescription.Type = weles.ArtifactType(gen.OneStringOf(
			string(weles.ArtifactTypeIMAGE),
			string(weles.ArtifactTypeRESULT),
			string(weles.ArtifactTypeTEST),
			string(weles.ArtifactTypeYAML)))
		tmp.ArtifactDescription.URI = weles.ArtifactURI(gen.URL())
		tmp.Path = weles.ArtifactPath(gen.URL())
		tmp.Status = weles.ArtifactStatus(gen.OneStringOf(
			string(weles.ArtifactStatusDOWNLOADING),
			string(weles.ArtifactStatusPENDING),
			string(weles.ArtifactStatusREADY),
			string(weles.ArtifactStatusFAILED)))

		dateTimeIter = dateTimeIter.Add(durationIncrement)
		artifactInfo[i] = tmp
	}
	return artifactInfo
}
