package cronroutine

import (
	"context"
	"time"
)

type JobConfig struct {
	// ID is the unique identifier of the job.
	ID string

	// Schedule is the cron schedule that determines when the job will run.
	Schedule string

	// Timeout is the amount of time each instance of the job is allowed to
	// run before it is killed.
	Timeout time.Duration

	// StartingDeadline is the maximum time the job can be delayed. If the
	// job is delayed more than this, it will be skipped.
	StartingDeadline time.Duration

	// AllowConccurentRuns determines whether the next job will start
	// if it is currently running.
	AllowConccurentRuns bool

	// This function will be run when the job is executed.
	Func func(context.Context) error
}

type Job struct {
	jobConfig *JobConfig
	history   []*History
	cron      *Cron
}

func (j *Job) ID() string                          { return j.jobConfig.ID }
func (j *Job) Schedule() string                    { return j.jobConfig.Schedule }
func (j *Job) Timeout() time.Duration              { return j.jobConfig.Timeout }
func (j *Job) StartingDeadline() time.Duration     { return j.jobConfig.StartingDeadline }
func (j *Job) AllowConccurentRuns() bool           { return j.jobConfig.AllowConccurentRuns }
func (j *Job) NextRun() time.Time                  { return j.cron.Next() }
func (j *Job) NextFor(t time.Duration) []time.Time { return j.cron.NextFor(t) }

func (j *Job) History() []*History {
	ret := make([]*History, len(j.history))
	if copy(ret, j.history) != len(j.history) {
		panic("unexpected history copy length")
	}

	return ret
}
