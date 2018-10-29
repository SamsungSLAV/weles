package client

import (
	"os"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/client/artifacts"
	"github.com/SamsungSLAV/weles/client/jobs"
)

// JobSpecification is a wrapper around os.File
// It implements NamedReadCloser interface (including io.ReadCloser interface)
type JobSpecification struct {
	file *os.File
}

// NewJobSpec opens file on given path and returns JobSpecification object
// which should be used as input to NewJob method. NewJob will call Close on
// the file.
func NewJobSpec(path string) (js JobSpecification, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return JobSpecification{}, err
	}
	return JobSpecification{file: f}, nil
}

// Read implements io.Reader interface.
func (js JobSpecification) Read(b []byte) (n int, err error) {
	return js.file.Read(b)
}

func (js JobSpecification) Close() error {
	return js.file.Close()
}

func (js JobSpecification) Name() string {
	return "yamlfile"
}

// NewJob returns JobID of created request and nil error on successful creation.
// Weles should be furtherly queried on the job progress. Job will be run when
// requested device will be available in Boruta server.
func (c *Weles) NewJob(j JobSpecification) (weles.JobID, error) {
	params := jobs.NewJobCreatorParams().WithYamlfile(j)
	resp, err := c.Jobs.JobCreator(params)
	return resp.Payload, err
}

// CancelJob returns nil on success.
func (c *Weles) CancelJob(jid uint64) error {
	params := jobs.NewJobCancelerParams().WithJobID(jid)
	_, err := c.Jobs.JobCanceler(params)
	return err
}

// ListJobs returns a slice of JobInfo according to specified:
// * filter
// * sorter
// * paginator
// It is recommended to use filter to full extent rather than paginating through results.
// This is why []JobInfo is not a type with Previous/Next methods.
func (c *Weles) ListJobs(
	f weles.JobFilter, s weles.JobSorter, p weles.JobPagination) (
	jl []*weles.JobInfo, totalRecords, remainingRecords uint64, next, prev string, err error) {

	params := jobs.NewJobListerParams().WithJobFilterAndSort(
		jobs.JobListerBody{Filter: &f, Sorter: &s})

	if (p != weles.JobPagination{}) {
		if p.JobID != 0 {
			if p.Forward {
				tmp := uint64(p.JobID)
				params.SetAfter(&tmp)
			} else {
				tmp := uint64(p.JobID)
				params.SetBefore(&tmp)
			}
		}
		if p.Limit != 0 {
			params.SetLimit(&p.Limit)
		}
	}
	var respFull *jobs.JobListerOK
	var respPartial *jobs.JobListerPartialContent
	if respFull, respPartial, err = c.Jobs.JobLister(params); err != nil {
		return []*weles.JobInfo{}, 0, 0, "", "", err
	}

	if respFull != nil {
		return respFull.Payload, respFull.TotalRecords, 0, respFull.Previous, respFull.Next, nil
	}
	return respPartial.Payload, respPartial.TotalRecords, respPartial.RemainingRecords, respPartial.Previous, respPartial.Next, nil
}

// ListArtifacts returns a slice of ArtifactInfo according to specified filter,
// sorter and paginator.It is recommended to use filter to full extend, rather
// than paginating through results thus a slice is returned (instead of a type
// with Previous/Next methods.
func (c *Weles) ListArtifacts(
	f weles.ArtifactFilter, s weles.ArtifactSorter, p weles.ArtifactPagination) (
	al []*weles.ArtifactInfo, totalRecords, remainingRecords uint64, prev, next string, err error) {

	params := artifacts.NewArtifactListerParams().WithArtifactFilterAndSort(
		artifacts.ArtifactListerBody{Filter: &f, Sorter: &s})

	if (p != weles.ArtifactPagination{}) {
		if p.ID != 0 {
			if p.Forward {
				tmp := int64(p.ID)
				params.SetAfter(&tmp)
			} else {
				tmp := int64(p.ID)
				params.SetBefore(&tmp)
			}
		}
		if p.Limit != 0 {
			params.SetLimit(&p.Limit)
		}
	}
	var respFull *artifacts.ArtifactListerOK
	var respPartial *artifacts.ArtifactListerPartialContent
	if respFull, respPartial, err = c.Artifacts.ArtifactLister(params); err != nil {
		return []*weles.ArtifactInfo{}, 0, 0, "", "", err
	}
	if respFull != nil {
		return respFull.Payload, respFull.TotalRecords, 0, respFull.Previous, respFull.Next, nil
	}
	return respPartial.Payload, respPartial.TotalRecords, respPartial.RemainingRecords, respPartial.Previous, respPartial.Next, nil
}
