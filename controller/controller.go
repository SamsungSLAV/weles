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

// Package controller provides Controller implementation. Controller binds all
// major components of Weles and provides logic layer for running proper methods
// of these components.
//
// Controller implements also JobManager interface providing API for controlling
// Weles' Jobs. This interface should be used by HTTP API server.
package controller

import (
	"sync"
	"time"

	"git.tizen.org/tools/boruta"

	"git.tizen.org/tools/weles"
)

// Controller binds all major components of Weles and provides logic layer
// for running proper methods of these components.
//
// It implements JobManager interface providing API for controlling
// Weles' Jobs. This interface should be used by HTTP API server.
// Controller should be created with NewController function.
type Controller struct {
	weles.JobManager

	// jobs references module implementing Jobs management.
	jobs JobsController
	// parser controls parsing yaml file and creation of Job's config.
	parser Parser
	// downloader controls downloading artifacts required for the Job.
	downloader Downloader
	// boruter acquires, releases and monitors Dryads from Boruta.
	boruter Boruter
	// dryader delegates Jobs execution to DryadJobManager and monitors progress.
	dryader Dryader
	// finish is channel for stopping internal goroutine.
	finish chan int
	// looper waits for internal goroutine running loop to finish.
	looper sync.WaitGroup
}

// NewJobManager creates and initializes a new instance of Controller with
// internal submodules and returns JobManager interface.
// It is the only valid way to get JobManager interface.
func NewJobManager(arm weles.ArtifactManager, yap weles.Parser, bor boruta.Requests,
	borutaRefreshPeriod time.Duration, djm weles.DryadJobManager) weles.JobManager {

	js := NewJobsController()
	pa := NewParser(js, arm, yap)
	do := NewDownloader(js, arm)
	bo := NewBoruter(js, bor, borutaRefreshPeriod)
	dr := NewDryader(js, djm)

	return NewController(js, pa, do, bo, dr)
}

// NewController creates and initializes a new instance of Controller.
// It requires internal Controller's submodules.
func NewController(js JobsController, pa Parser, do Downloader, bo Boruter, dr Dryader,
) *Controller {
	c := &Controller{
		jobs:       js,
		parser:     pa,
		downloader: do,
		boruter:    bo,
		dryader:    dr,
		finish:     make(chan int),
	}
	c.looper.Add(1)
	go c.loop()
	return c
}

// Finish internal goroutine.
func (c *Controller) Finish() {
	c.finish <- 1
	c.looper.Wait()
}

// CreateJob creates a new Job in Weles using recipe passed in YAML format.
// It is a part of JobManager implementation.
func (c *Controller) CreateJob(yaml []byte) (weles.JobID, error) {
	j, err := c.jobs.NewJob(yaml)
	if err != nil {
		return weles.JobID(0), err
	}

	go c.parser.Parse(j)

	return j, nil
}

// CancelJob cancels Job identified by argument. Job execution is stopped.
// It is a part of JobManager implementation.
func (c *Controller) CancelJob(j weles.JobID) error {
	err := c.jobs.SetStatusAndInfo(j, weles.JobStatusCANCELED, "")
	if err != nil {
		return err
	}
	c.dryader.CancelJob(j)
	c.boruter.Release(j)
	return nil
}

// ListJobs returns information on Jobs.
// It is a part of JobManager implementation.
func (c *Controller) ListJobs(filter weles.JobFilter, sorter weles.JobSorter,
	paginator weles.JobPagination) ([]weles.JobInfo, weles.ListInfo, error) {
	return c.jobs.List(filter, sorter, paginator)
}

// loop implements main loop of the Controller reacting to different events
// related to processed Jobs.
func (c *Controller) loop() {
	defer c.looper.Done()
	for {
		select {
		case <-c.finish:
			return
		case noti := <-c.parser.Listen():
			if !noti.OK {
				c.fail(noti.JobID, noti.Msg)
				continue
			}
			c.downloader.DispatchDownloads(noti.JobID)
		case noti := <-c.downloader.Listen():
			if !noti.OK {
				c.fail(noti.JobID, noti.Msg)
				continue
			}
			c.boruter.Request(noti.JobID)
		case noti := <-c.boruter.Listen():
			if !noti.OK {
				c.fail(noti.JobID, noti.Msg)
				continue
			}
			c.dryader.StartJob(noti.JobID)
		case noti := <-c.dryader.Listen():
			if !noti.OK {
				c.fail(noti.JobID, noti.Msg)
				continue
			}
			c.succeed(noti.JobID)
		}
	}
}

// fail sets Job in FAILED state and if needed stops Job's execution on Dryad
// and releases Dryad to Boruta.
func (c *Controller) fail(j weles.JobID, msg string) {
	c.jobs.SetStatusAndInfo(j, weles.JobStatusFAILED, msg)
	c.dryader.CancelJob(j)
	c.boruter.Release(j)
}

// succeed sets Job in COMPLETED state.
func (c *Controller) succeed(j weles.JobID) {
	c.jobs.SetStatusAndInfo(j, weles.JobStatusCOMPLETED, "")
	c.boruter.Release(j)
}
