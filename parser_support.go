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
	"errors"
	"reflect"
	"time"
)

var validTimeouts = map[string]time.Duration{
	"seconds": time.Second,
	"minutes": time.Minute,
	"hours":   time.Hour,
	"days":    24 * time.Hour,
}

// UnmarshalYAML unmarshals ValidPeriod type.
func (t *ValidPeriod) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var localTime map[string]int
	err := unmarshal(&localTime)
	if err != nil {
		return err
	}

	for k, v := range validTimeouts {
		if _v, ok := localTime[k]; ok {
			*t = ValidPeriod(v) * ValidPeriod(_v)
			return nil
		}
	}

	return errors.New("invalid timeout")
}

// LocalTestActionContainer contains fields for all types of test cases.
type LocalTestActionContainer struct {
	Boot
	Push
	Run
	Pull
}

// TestActionTab contains all possible test cases.
type TestActionTab []LocalTestActionContainer

// UnmarshalYAML is an unmarshalling function for TestActions type.
func (t *TestActions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var testAct TestActionTab
	err := unmarshal(&testAct)
	if err != nil {
		return err
	}

	for _, e := range testAct {
		if !reflect.DeepEqual(e.Boot, Boot{}) {
			b := e.Boot
			*t = append(*t, b)
		}
		if e.Push != (Push{}) {
			p := e.Push
			*t = append(*t, p)
		}
		if e.Run != (Run{}) {
			r := e.Run
			*t = append(*t, r)
		}
		if e.Pull != (Pull{}) {
			p := e.Pull
			*t = append(*t, p)
		}
	}

	return nil
}

// LocalActionContainer contains fields for all types of actions.
type LocalActionContainer Action

// ActionTab contains all possible actions.
type ActionTab []LocalActionContainer

// UnmarshalYAML is an unmarshalling function for Action type.
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var act ActionTab
	err := unmarshal(&act)
	if err != nil {
		return err
	}

	for _, e := range act {
		if !reflect.DeepEqual(e.Deploy, Deploy{}) {
			a.Deploy = e.Deploy
		}
		if !reflect.DeepEqual(e.Boot, Boot{}) {
			a.Boot = e.Boot
		}
		if !reflect.DeepEqual(e.Test, Test{}) {
			a.Test = e.Test
		}
	}

	return nil
}
