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

// File controller/jobscontrollerimpl.go contains JobsController interface
// implementation.

package controller

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/SamsungSLAV/weles"
)

// JobsControllerImpl structure stores Weles' Jobs data. It controls
// collision-free JobID creation. It stores state of Jobs' execution and saves
// data to DB. It implements JobsController interface.
type JobsControllerImpl struct {
	JobsController
	// mutex protects JobsControllerImpl structure.
	mutex *sync.RWMutex
	// lastID is the last used ID for the Job.
	lastID weles.JobID
	// jobs stores information about Weles' Jobs.
	jobs map[weles.JobID]*Job
}

// setupLastID initializes last used ID. Value is read from DB meta data.
func (js *JobsControllerImpl) setupLastID() {
	// TODO initialize with meta data read from DB.
	// Current implementation starts with seconds from Epoch to avoid problems with
	// artifacts database.

	js.lastID = weles.JobID(time.Now().Unix())
}

// NewJobsController creates and initializes a new instance of Jobs structure.
// It is the only valid way of creating it.
func NewJobsController() JobsController {
	js := &JobsControllerImpl{
		mutex: new(sync.RWMutex),
		jobs:  make(map[weles.JobID]*Job),
	}

	js.setupLastID()

	// TODO load Jobs data from DB.

	return js
}

// nextID generates and returns ID assigned to a new Job.
// It also updates lastID and saves the information in DB meta data.
func (js *JobsControllerImpl) nextID() weles.JobID {
	js.lastID++

	// TODO save new lastID in DB.

	return js.lastID
}

// NewJob creates and initializes a new Job.
func (js *JobsControllerImpl) NewJob(yaml []byte) (weles.JobID, error) {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	j := js.nextID()

	now := strfmt.DateTime(time.Now())
	js.jobs[j] = &Job{
		JobInfo: weles.JobInfo{
			JobID:   j,
			Created: now,
			Updated: now,
			Status:  weles.JobStatusNEW,
		},
		yaml: yaml,
	}

	// TODO save struct in DB

	return j, nil
}

// GetYaml returns yaml Job description.
func (js *JobsControllerImpl) GetYaml(j weles.JobID) ([]byte, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	job, ok := js.jobs[j]
	if !ok {
		return nil, weles.ErrJobNotFound
	}

	return job.yaml, nil
}

// SetConfig stores config in Jobs structure.
func (js *JobsControllerImpl) SetConfig(j weles.JobID, conf weles.Config) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.ErrJobNotFound
	}

	job.config = conf
	job.Updated = strfmt.DateTime(time.Now())
	return nil
}

// isStatusChangeValid verifies if Job's status change is valid.
// It is a helper function for SetStatusAndInfo.
func isStatusChangeValid(oldStatus, newStatus weles.JobStatus) bool {
	if oldStatus == newStatus {
		return true
	}
	switch oldStatus {
	case weles.JobStatusNEW:
		switch newStatus {
		case weles.JobStatusPARSING, weles.JobStatusCANCELED, weles.JobStatusFAILED:
			return true
		}
	case weles.JobStatusPARSING:
		switch newStatus {
		case weles.JobStatusDOWNLOADING, weles.JobStatusCANCELED, weles.JobStatusFAILED:
			return true
		}
	case weles.JobStatusDOWNLOADING:
		switch newStatus {
		case weles.JobStatusWAITING, weles.JobStatusCANCELED, weles.JobStatusFAILED:
			return true
		}
	case weles.JobStatusWAITING:
		switch newStatus {
		case weles.JobStatusRUNNING, weles.JobStatusCANCELED, weles.JobStatusFAILED:
			return true
		}
	case weles.JobStatusRUNNING:
		switch newStatus {
		case weles.JobStatusCOMPLETED, weles.JobStatusCANCELED, weles.JobStatusFAILED:
			return true
		}
	}
	return false
}

// SetStatusAndInfo changes status of the Job and updates info. Only valid
// changes are allowed.
// There are 3 terminal statuses: JobStatusFAILED, JobStatusCANCELED, JobStatusCOMPLETED;
// and 5 non-terminal statuses: JobStatusNEW, JobStatusPARSING, JobStatusDOWNLOADING,
// JobStatusWAITING, JobStatusRUNNING.
// Only below changes of statuses are allowed:
// * JobStatusNEW --> {JobStatusPARSING, JobStatusCANCELED, JobStatusFAILED}
// * JobStatusPARSING --> {JobStatusDOWNLOADING, JobStatusCANCELED, JobStatusFAILED}
// * JobStatusDOWNLOADING --> {JobStatusWAITING, JobStatusCANCELED, JobStatusFAILED}
// * JobStatusWAITING --> {JobStatusRUNNING, JobStatusCANCELED, JobStatusFAILED}
// * JobStatusRUNNING --> {JobStatusCOMPLETED, JobStatusCANCELED, JobStatusFAILED}
func (js *JobsControllerImpl) SetStatusAndInfo(j weles.JobID, newStatus weles.JobStatus, msg string,
) error {

	js.mutex.Lock()
	defer js.mutex.Unlock()

	job, ok := js.jobs[j]
	if !ok {
		log.Println(weles.ErrJobNotFound.Error(), "JobID:", j)
		return weles.ErrJobNotFound
	}

	if !isStatusChangeValid(job.Status, newStatus) {
		log.Println(weles.ErrJobStatusChangeNotAllowed.Error(), "from:", job.Status, "to:",
			newStatus)
		return weles.ErrJobStatusChangeNotAllowed
	}

	job.Status = newStatus
	job.Info = msg
	job.Updated = strfmt.DateTime(time.Now())
	return nil
}

// GetConfig returns Job's config.
func (js *JobsControllerImpl) GetConfig(j weles.JobID) (weles.Config, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.Config{}, weles.ErrJobNotFound
	}

	return job.config, nil
}

// SetDryad saves access info for acquired Dryad.
func (js *JobsControllerImpl) SetDryad(j weles.JobID, d weles.Dryad) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.ErrJobNotFound
	}

	job.dryad = d
	return nil
}

// GetDryad returns Dryad acquired for the Job.
func (js *JobsControllerImpl) GetDryad(j weles.JobID) (weles.Dryad, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()

	job, ok := js.jobs[j]
	if !ok {
		return weles.Dryad{}, weles.ErrJobNotFound
	}

	return job.dryad, nil
}

func (js *JobsControllerImpl) filter(filter weles.JobFilter, paginator weles.JobPagination) (
	[]weles.JobInfo, bool, error) {
	// extra defines if the returned collection of JobInfo contain additionally pagination JobID.
	var extra = false

	// Prepare filter.
	f, err := prepareFilter(&filter)
	if err != nil {
		return nil, extra, err
	}

	// Filter jobs.
	ret := make([]weles.JobInfo, 0, len(js.jobs))
	for _, job := range js.jobs {
		if job.passesFilter(f) {
			ret = append(ret, job.JobInfo)
		}
	}

	// If pagination is enabled and there's a valid pagination JobID, it must be added
	// to the returned collection of JobInfos to allow finding proper scope of jobs (page).
	// It should be added only if it does not pass filter, otherwise it is already added.
	if paginator.Limit != 0 && paginator.JobID != weles.JobID(0) {
		job, present := js.jobs[paginator.JobID]
		if !present {
			return nil, extra, weles.ErrInvalidArgument(fmt.Sprintf("JobID: %d not found",
				paginator.JobID))
		}
		if !job.passesFilter(f) {
			ret = append(ret, job.JobInfo)
			extra = true
		}
	}
	return ret, extra, nil
}

func (js *JobsControllerImpl) sort(ret []weles.JobInfo, sorter weles.JobSorter) []weles.JobInfo {
	// Sort jobs.
	ps := &jobSorter{
		jobs: ret,
		by:   byJobIDAsc,
	}
	switch sorter.By {
	case weles.JobSortByCreatedDate:
		ps.setByFunction(sorter.Order, byCreatedAsc, byCreatedDesc)
	case weles.JobSortByUpdatedDate:
		ps.setByFunction(sorter.Order, byUpdatedAsc, byUpdatedDesc)
	case weles.JobSortByJobStatus:
		ps.setByFunction(sorter.Order, byStatusAsc, byStatusDesc)
	}
	sort.Sort(ps)
	return ps.jobs
}

func (js *JobsControllerImpl) paginate(ret []weles.JobInfo, paginator weles.JobPagination) (
	total, elems, left, index int) {
	// Pagination.
	total = len(ret)

	if paginator.Limit == 0 {
		// Pagination is disabled. Return all records.
		elems = total
		return
	}
	if paginator.JobID == weles.JobID(0) {
		// Starting pagination - 1st page.
		elems = min(int(paginator.Limit), total)
		left = total - elems
		return
	}
	// Find index of pagination JobID.
	for i, job := range ret {
		if job.JobID == paginator.JobID {
			index = i
			break
		}
	}

	if paginator.Forward {
		index++
		elems = min(int(paginator.Limit), total-index)
		left = total - index - elems
	} else {
		elems = min(int(paginator.Limit), index)
		left = index - elems
		index -= elems
	}
	return
}

// List returns information on Jobs. It takes 3 arguments:
// - JobFilter containing filters
// - JobSorter containing sorting key and sorting direction
// - JobPagination containing element after/before which a page should be returned. It also
// contains information about direction of listing and the size of the returned page which
// must always be set.
func (js *JobsControllerImpl) List(filter weles.JobFilter, sorter weles.JobSorter,
	paginator weles.JobPagination) ([]weles.JobInfo, weles.ListInfo, error) {

	js.mutex.RLock()
	defer js.mutex.RUnlock()

	// Filter jobs.
	ret, extra, err := js.filter(filter, paginator)
	if err != nil {
		return nil, weles.ListInfo{}, err
	}

	// Sort jobs.
	ret = js.sort(ret, sorter)

	// Pagination.
	total, elems, left, index := js.paginate(ret, paginator)

	// Don't count pagination JobID if added extra.
	if extra {
		total--
	}

	info := weles.ListInfo{TotalRecords: uint64(total), RemainingRecords: uint64(left)}
	return ret[index : index+elems], info, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type filter struct {
	CreatedAfter  time.Time
	CreatedBefore time.Time
	UpdatedAfter  time.Time
	UpdatedBefore time.Time
	Info          *regexp.Regexp
	JobID         map[weles.JobID]interface{}
	Name          *regexp.Regexp
	Status        map[weles.JobStatus]interface{}
}

func prepareFilterRegexp(arr []string) (*regexp.Regexp, error) {
	if len(arr) == 0 {
		return nil, nil
	}

	var size int
	for _, s := range arr {
		size += 3 + len(s)
	}

	var str strings.Builder
	str.Grow(size)
	for _, s := range arr {
		_, _ = str.WriteString("|(" + s + ")")
	}

	return regexp.Compile(str.String()[1:])
}

func prepareFilter(in *weles.JobFilter) (out *filter, err error) {
	var regErr error

	out = new(filter)

	out.CreatedAfter = time.Time(in.CreatedAfter)
	out.CreatedBefore = time.Time(in.CreatedBefore)
	out.Info, regErr = prepareFilterRegexp(in.Info)
	if regErr != nil {
		return nil, weles.ErrInvalidArgument("cannot compile regex from Info: " + regErr.Error())
	}
	if len(in.JobID) > 0 {
		out.JobID = make(map[weles.JobID]interface{})
		for _, x := range in.JobID {
			out.JobID[x] = nil
		}
	}
	out.Name, regErr = prepareFilterRegexp(in.Name)
	if regErr != nil {
		return nil, weles.ErrInvalidArgument("cannot compile regex from Name: " + regErr.Error())
	}
	if len(in.Status) > 0 {
		out.Status = make(map[weles.JobStatus]interface{})
		for _, x := range in.Status {
			out.Status[x] = nil
		}
	}
	out.UpdatedAfter = time.Time(in.UpdatedAfter)
	out.UpdatedBefore = time.Time(in.UpdatedBefore)

	return out, nil
}

func (job *Job) passesCreatedAfterFilter(f *filter) bool {
	return f.CreatedAfter.IsZero() || time.Time(job.JobInfo.Created).After(f.CreatedAfter)
}

func (job *Job) passesCreatedBeforeFilter(f *filter) bool {
	return f.CreatedBefore.IsZero() || time.Time(job.JobInfo.Created).Before(f.CreatedBefore)
}

func (job *Job) passesUpdatedAfterFilter(f *filter) bool {
	return f.UpdatedAfter.IsZero() || time.Time(job.JobInfo.Updated).After(f.UpdatedAfter)
}

func (job *Job) passesUpdatedBeforeFilter(f *filter) bool {
	return f.UpdatedBefore.IsZero() || time.Time(job.JobInfo.Updated).Before(f.UpdatedBefore)
}

func (job *Job) passesInfoFilter(f *filter) bool {
	return f.Info == nil || f.Info.MatchString(job.JobInfo.Info)
}

func (job *Job) passesJobIDFilter(f *filter) bool {
	if f.JobID == nil {
		return true
	}
	_, present := f.JobID[job.JobInfo.JobID]
	return present
}

func (job *Job) passesNameFilter(f *filter) bool {
	return f.Name == nil || f.Name.MatchString(job.JobInfo.Name)
}

func (job *Job) passesStatusFilter(f *filter) bool {
	if f.Status == nil {
		return true
	}
	_, present := f.Status[job.JobInfo.Status]
	return present
}

func (job *Job) passesFilter(f *filter) bool {
	return job.passesCreatedAfterFilter(f) &&
		job.passesCreatedBeforeFilter(f) &&
		job.passesUpdatedAfterFilter(f) &&
		job.passesUpdatedBeforeFilter(f) &&
		job.passesInfoFilter(f) &&
		job.passesJobIDFilter(f) &&
		job.passesNameFilter(f) &&
		job.passesStatusFilter(f)
}

func byCreatedAsc(i1, i2 *weles.JobInfo) bool {
	if time.Time(i1.Created).Equal(time.Time(i2.Created)) {
		return byJobIDAsc(i1, i2)
	}
	return time.Time(i1.Created).Before(time.Time(i2.Created))
}

func byCreatedDesc(i1, i2 *weles.JobInfo) bool {
	if time.Time(i1.Created).Equal(time.Time(i2.Created)) {
		return byJobIDAsc(i1, i2)
	}
	return time.Time(i1.Created).After(time.Time(i2.Created))
}

func byUpdatedAsc(i1, i2 *weles.JobInfo) bool {
	if time.Time(i1.Updated).Equal(time.Time(i2.Updated)) {
		return byJobIDAsc(i1, i2)
	}
	return time.Time(i1.Updated).Before(time.Time(i2.Updated))
}

func byUpdatedDesc(i1, i2 *weles.JobInfo) bool {
	if time.Time(i1.Updated).Equal(time.Time(i2.Updated)) {
		return byJobIDAsc(i1, i2)
	}
	return time.Time(i1.Updated).After(time.Time(i2.Updated))
}

func byStatusAsc(i1, i2 *weles.JobInfo) bool {
	if i1.Status.ToInt() == i2.Status.ToInt() {
		return byJobIDAsc(i1, i2)
	}
	return i1.Status.ToInt() < i2.Status.ToInt()
}

func byStatusDesc(i1, i2 *weles.JobInfo) bool {
	if i1.Status.ToInt() == i2.Status.ToInt() {
		return byJobIDAsc(i1, i2)
	}
	return i1.Status.ToInt() > i2.Status.ToInt()
}

func byJobIDAsc(i1, i2 *weles.JobInfo) bool {
	return i1.JobID < i2.JobID
}

type jobSorter struct {
	jobs []weles.JobInfo
	by   func(i1, i2 *weles.JobInfo) bool
}

// Len is part of sort.Interface.
func (s *jobSorter) Len() int {
	return len(s.jobs)
}

// Swap is part of sort.Interface.
func (s *jobSorter) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *jobSorter) Less(i, j int) bool {
	return s.by(&s.jobs[i], &s.jobs[j])
}

// by is the type of a "less" function that defines the ordering of its JobInfo arguments.
type by func(p1, p2 *weles.JobInfo) bool

func (s *jobSorter) setByFunction(order weles.SortOrder, asc, desc by) {
	switch order {
	case weles.SortOrderAscending:
		s.by = asc
	case weles.SortOrderDescending:
		s.by = desc
	}
}
