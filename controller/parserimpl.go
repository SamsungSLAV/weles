/*
 *  Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
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

// File controller/parserimpl.go implements Parser.

package controller

import (
	"fmt"
	"io/ioutil"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/controller/notifier"
)

// ParserImpl implements Parser for Controller.
type ParserImpl struct {
	// Notifier provides channel for communication with Controller.
	notifier.Notifier
	// jobs references Controller's submodule responsible for Jobs management.
	jobs JobsController
	// artifacts manages ArtifactsDB.
	artifacts weles.ArtifactManager
	// parser creates Job's recipe from yaml.
	parser weles.Parser
}

// NewParser creates a new ParserImpl structure setting up references
// to Weles' modules.
func NewParser(j JobsController, a weles.ArtifactManager, p weles.Parser) Parser {
	return &ParserImpl{
		Notifier:  notifier.NewNotifier(),
		jobs:      j,
		artifacts: a,
		parser:    p,
	}
}

// Parse prepares new Job to be processed by saving yaml file in ArtifactDB,
// parsing yaml and preparing Job's configuration.
func (h *ParserImpl) Parse(j weles.JobID) {
	err := h.jobs.SetStatusAndInfo(j, weles.JobStatusPARSING, "")
	if err != nil {
		h.SendFail(j, fmt.Sprintf("Internal Weles error while changing Job status : %s",
			err.Error()))
		return
	}

	yaml, err := h.jobs.GetYaml(j)
	if err != nil {
		h.SendFail(j, fmt.Sprintf("Internal Weles error while getting yaml description : %s",
			err.Error()))
		return
	}

	path, err := h.artifacts.CreateArtifact(weles.ArtifactDescription{
		JobID: j,
		Type:  weles.ArtifactTypeYAML,
	})
	if err != nil {
		h.SendFail(j, fmt.Sprintf(
			"Internal Weles error while creating file path in ArtifactDB : %s",
			err.Error()))
		return
	}

	err = ioutil.WriteFile(string(path), yaml, 0644)
	if err != nil {
		h.SendFail(
			j, fmt.Sprintf("Internal Weles error while saving file in ArtifactDB : %s",
				err.Error()))
		return
	}

	conf, err := h.parser.ParseYaml(yaml)
	if err != nil {
		h.SendFail(j, fmt.Sprintf("Error parsing yaml file : %s",
			err.Error()))
		return
	}

	err = h.jobs.SetConfig(j, *conf)
	if err != nil {
		h.SendFail(j, fmt.Sprintf("Internal Weles error while setting config : %s",
			err.Error()))
		return
	}

	h.SendOK(j)
}
