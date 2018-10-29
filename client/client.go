package client

import (
	"bytes"

	"github.com/SamsungSLAV/weles"
	"github.com/SamsungSLAV/weles/client/artifacts"
	"github.com/SamsungSLAV/weles/client/jobs"
)

// JobSpecification is a wrapper around os.File
// It implements NamedReadCloser interface (including io.ReadCloser interface)
type jobSpecification struct {
	yamlfile *bytes.Buffer
}

func newJobSpec(yaml []byte) jobSpecification {
	return jobSpecification{yamlfile: bytes.NewBuffer(yaml)}
}

func (js jobSpecification) Read(b []byte) (n int, err error) {
	return js.yamlfile.Read(b)
}

func (js jobSpecification) Close() (err error) {
	// no need to close buffer
	return
}

func (js jobSpecification) Name() string {
	return "yamlfile"
}

// NewJob returns JobID of created request and nil error on successful creation.
// Weles should be furtherly queried on the job progress. Job will be run when
// requested device will be available in Boruta server.
func (c *Weles) CreateJob(yaml []byte) (weles.JobID, error) {
	js := newJobSpec(yaml)
	params := jobs.NewJobCreatorParams().WithYamlfile(js)
	resp, err := c.Jobs.JobCreator(params)
	return resp.Payload, err
}

// CancelJob returns nil on success.
func (c *Weles) CancelJob(jid weles.JobID) error {
	params := jobs.NewJobCancelerParams().WithJobID(uint64(jid))
	_, err := c.Jobs.JobCanceler(params)
	return err
}

// ListJobs returns a slice of JobInfo according to specified:
// * filter
// * sorter
// * paginator
// It is recommended to use filter to full extent rather than paginating through results.
// This is why []JobInfo is not a type with Previous/Next methods.
func (c *Weles) ListJobs(f weles.JobFilter, s weles.JobSorter, p weles.JobPagination) (
	[]weles.JobInfo, weles.ListInfo, error) {

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
	var err error
	if respFull, respPartial, err = c.Jobs.JobLister(params); err != nil {
		return []weles.JobInfo{}, weles.ListInfo{}, err
	}

	if respFull != nil {
		tmp := dereferenceSliceElems(respFull)
		ji, ok := tmp.([]weles.JobInfo)
		if !ok {
			return []weles.JobInfo{}, weles.ListInfo{}, ErrDereferencing
		}
		return ji,
			weles.ListInfo{
				TotalRecords:     respFull.TotalRecords,
				RemainingRecords: 0,
			},
			nil
	}
	tmp := dereferenceSliceElems(respPartial)
	ji, ok := tmp.([]weles.JobInfo)
	if !ok {
		return []weles.JobInfo{}, weles.ListInfo{}, ErrDereferencing
	}
	return ji,
		weles.ListInfo{
			TotalRecords:     respPartial.TotalRecords,
			RemainingRecords: respPartial.RemainingRecords,
		},
		nil
}

// ListArtifacts returns a slice of ArtifactInfo according to specified filter,
// sorter and paginator.It is recommended to use filter to full extend, rather
// than paginating through results thus a slice is returned (instead of a type
// with Previous/Next methods.
func (c *Weles) ListArtifacts(
	f weles.ArtifactFilter, s weles.ArtifactSorter, p weles.ArtifactPagination) (
	[]weles.ArtifactInfo, weles.ListInfo, error) {

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
	var err error
	if respFull, respPartial, err = c.Artifacts.ArtifactLister(params); err != nil {
		return []weles.ArtifactInfo{}, weles.ListInfo{}, err
	}
	if respFull != nil {
		tmp := dereferenceSliceElems(respFull)
		ji, ok := tmp.([]weles.ArtifactInfo)
		if !ok {
			return []weles.ArtifactInfo{}, weles.ListInfo{}, ErrDereferencing
		}
		return ji,
			weles.ListInfo{
				TotalRecords:     respFull.TotalRecords,
				RemainingRecords: 0,
			},
			nil
	}
	tmp := dereferenceSliceElems(respPartial)
	ji, ok := tmp.([]weles.ArtifactInfo)
	if !ok {
		return []weles.ArtifactInfo{}, weles.ListInfo{}, ErrDereferencing
	}
	return ji,
		weles.ListInfo{
			TotalRecords:     respPartial.TotalRecords,
			RemainingRecords: respPartial.RemainingRecords,
		},
		nil
}

// helper function for dereferencing slice of pointers to slice of elements
// works only for JobInfo and ArtifactInfo so should probably be renamed
// i dont know if i can make it more generic without reflections (which are heavy)
func dereferenceSliceElems(i interface{}) interface{} {
	if sliceOfPointers, ok := i.([]*weles.JobInfo); ok {
		sliceOfElems := make([]weles.JobInfo, len(sliceOfPointers))
		for i := range sliceOfPointers {
			sliceOfElems[i] = *sliceOfPointers[i]
		}
		return sliceOfElems
	}
	if sliceOfPointers, ok := i.([]*weles.ArtifactInfo); ok {
		sliceOfElems := make([]weles.ArtifactInfo, len(sliceOfPointers))
		for i := range sliceOfPointers {
			sliceOfElems[i] = *sliceOfPointers[i]
		}
		return sliceOfElems
	}
	return nil
}
