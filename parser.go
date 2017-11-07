/*
 *  Copyright (c) 2017 Samsung Electronics Co., Ltd All Rights Reserved
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

package weles

import (
	"time"
)

// ValidPeriod denotes timeouts.
type ValidPeriod time.Duration

// Priority denotes priority level.
type Priority string

const (
	// LOW - low priority.
	LOW Priority = "low"
	// MEDIUM - medium priority.
	MEDIUM Priority = "medium"
	// HIGH - high priority.
	HIGH Priority = "high"
)

// TestAction is either Boot/Push/Run/Pull.
// It describes test action, which can be done on already prepared DUT.
type TestAction interface{}

// Boot describes the boot part of the test.
type Boot struct {
	Login         string      `yaml:"login"`
	Password      string      `yaml:"password"`
	Prompts       []string    `yaml:"prompts"`
	FailureRetry  int         `yaml:"failure_retry"`
	Timeout       ValidPeriod `yaml:"timeout"`
	InputSequence string      `yaml:"input_sequence"`
	WaitPattern   string      `yaml:"wait_pattern"`
	WaitTime      ValidPeriod `yaml:"wait_time"`
}

// Push describes the push part of the test.
type Push struct {
	URI     string      `yaml:"uri"`
	Dest    string      `yaml:"dest"`
	Alias   string      `yaml:"alias"`
	Timeout ValidPeriod `yaml:"timeout"`

	// Path defines ArtifactDB path. It's added for Controller purposes.
	Path string `yaml:"-"`
}

// Run describes the run part of the test.
type Run struct {
	Name    string      `yaml:"name"`
	Timeout ValidPeriod `yaml:"timeout"`
}

// Pull describes the pull part of the test,
// e.g. getting the test artifacts.
type Pull struct {
	Src     string      `yaml:"src"`
	Alias   string      `yaml:"alias"`
	Timeout ValidPeriod `yaml:"timeout"`

	// Path defines ArtifactDB path. It's added for Controller purposes.
	Path string `yaml:"-"`
}

// ImageDefinition describes images required for the tests.
type ImageDefinition struct {
	URI          string `yaml:"uri"`
	ChecksumURI  string `yaml:"checksum_uri"`
	ChecksumType string `yaml:"checksum_type"`
	Compression  string `yaml:"compression"`

	// Path defines ArtifactDB path. It's added for Controller purposes.
	Path         string `yaml:"-"`
	ChecksumPath string `yaml:"-"`
}

// PartitionDefinition describes a relation of a partition to named image, its size, and type.
type PartitionDefinition struct {
	ID        int    `yaml:"id"`
	ImageName string `yaml:"image_name"`
	Size      string `yaml:"size"`
	Type      string `yaml:"type"`
}

// Deploy describes "deploy" section in YAML.
type Deploy struct {
	Timeout         ValidPeriod           `yaml:"timeout"`
	Images          []ImageDefinition     `yaml:"images"`
	PartitionLayout []PartitionDefinition `yaml:"partition_layout"`
}

// TestActions is a container for all test actions.
type TestActions []TestAction

// TestCase describes single test case.
type TestCase struct {
	CaseName    string      `yaml:"case_name"`
	TestActions TestActions `yaml:"test_actions"`
}

// Test describes "test" section in YAML.
type Test struct {
	FailureRetry int         `yaml:"failure_retry"`
	Name         string      `yaml:"name"`
	Timeout      ValidPeriod `yaml:"timeout"`
	TestCases    []TestCase  `yaml:"test_cases"`
}

// Action describes actions executed on the DUT.
// Firstly it describes how to prepare DUT for a test,
// and then the test procedure itself.
type Action struct {
	Deploy
	Boot
	Test
}

// Timeouts describes default timeouts for different actions.
type Timeouts struct {
	// JobTimeout describes default timeouts for a job.
	JobTimeout ValidPeriod `yaml:"job"`
	// ActionTimeout describes default timeouts for boot/push/run/pull.
	ActionTimeout ValidPeriod `yaml:"action"`
}

// Config contains all informtion needed for the Weles to make test.
type Config struct {
	DeviceType string   `yaml:"device_type"`
	JobName    string   `yaml:"job_name"`
	Timeouts   Timeouts `yaml:"timeouts"`
	Priority   Priority `yaml:"priority"`
	Action     Action   `yaml:"actions"`
}

// Parser defines methods of YAML parser.
type Parser interface {
	// ParseYaml converts given input to Config. YAML file format is expected.
	ParseYaml(input []byte) (*Config, error)
}
