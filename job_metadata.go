package cronroutine

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
)

type jobMetadata struct {
	jobConfig    *JobConfig
	history      []*History
	historyLimit int
	cron         *Cron
	running      sync.Mutex
}

func (j *jobMetadata) ID() string {
	return j.jobConfig.ID
}

func (j *jobMetadata) History() []*History {
	ret := make([]*History, len(j.history))
	if copy(ret, j.history) != len(j.history) {
		panic("unexpected history copy length")
	}

	return ret
}

func (j *jobMetadata) addResult(ranAt time.Time, err error) {
	j.history = append(j.history, &History{
		jobID: j.ID(),
		ranAt: ranAt,
		err:   err,
	})

	if len(j.history) > j.historyLimit {
		j.history = j.history[1:]
	}
}

func (j *jobMetadata) Job() *Job {
	return &Job{
		jobConfig: &JobConfig{
			ID:                  j.jobConfig.ID,
			Schedule:            j.jobConfig.Schedule,
			Timeout:             j.jobConfig.Timeout,
			StartingDeadline:    j.jobConfig.StartingDeadline,
			AllowConccurentRuns: j.jobConfig.AllowConccurentRuns,
			Func:                j.jobConfig.Func,
		},
		history: j.History(),
		cron: &Cron{
			Minute:     j.cron.Minute,
			Hour:       j.cron.Hour,
			DayOfMonth: j.cron.DayOfMonth,
			Month:      j.cron.Month,
			DayOfWeek:  j.cron.DayOfWeek,
		},
	}
}

func (j *jobMetadata) shouldRun(startTime time.Time) error {
	deadline := startTime.Add(j.jobConfig.StartingDeadline)
	if j.jobConfig.AllowConccurentRuns {
		if time.Now().UTC().After(deadline) {
			return ErrPastStartingDeadline{}
		}
		return nil
	}

	lockAquired := make(chan struct{})
	go func() {
		j.running.Lock()
		lockAquired <- struct{}{}
	}()

	startingDeadlineExceeded := make(chan struct{})
	go func() {
		time.Sleep(time.Until(deadline))
		startingDeadlineExceeded <- struct{}{}
	}()

	select {
	case <-lockAquired:
		return nil
	case <-startingDeadlineExceeded:
		return ErrJobRunning{}
	}
}

func (j *jobMetadata) run(logger logr.Logger, startTime time.Time) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		err := j.shouldRun(startTime)
		if err != nil {
			logger.Error(err, "job could not start")
			j.addResult(startTime, err)
			return err
		}

		logger.Info("job started", "start_time", startTime.UTC())
		ctx, cancel := context.WithTimeout(ctx, j.jobConfig.Timeout)
		defer cancel()

		errChan := make(chan error)
		go func() {
			err := j.jobConfig.Func(ctx)
			if err != nil {
				logger.Error(err, "job failed")
			} else {
				logger.Info("job finished successfully", "end_time", time.Now().UTC())
			}
			errChan <- nil
		}()

		select {
		case <-ctx.Done():
			cancel()
			err := ErrJobTimeout{}
			logger.Error(err, "job execution timed out")
			j.addResult(startTime, err)
			return err
		case err := <-errChan:
			j.addResult(startTime, err)
			return err
		}
	}
}
