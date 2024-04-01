package cronroutine

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newTestScheduler(parallelization int) *Scheduler {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log := zapr.NewLogger(zapLog)

	return StartNewScheduler(&SchedulerConfig{
		Logger:       log,
		HistoryLimit: 10,
		WorkerCount:  parallelization,
	})
}

func TestScheduler_limitParallelization(t *testing.T) {
	t.Parallel()
	numCPU := runtime.NumCPU()
	scheduler := newTestScheduler(numCPU)
	var jobsRun atomic.Uint64

	newJobConfig := func(id string) JobConfig {
		return JobConfig{
			ID:                  id,
			Schedule:            "0/1 * * * *",
			Timeout:             2 * time.Second,
			StartingDeadline:    100 * time.Millisecond,
			AllowConccurentRuns: false,
			Func: func(ctx context.Context) error {
				fmt.Println("test job running!")
				jobsRun.Add(1)
				time.Sleep(1 * time.Second)
				return nil
			},
		}
	}

	for i := 0; i < numCPU+5; i++ {
		err := scheduler.AddJob(newJobConfig(fmt.Sprintf("test-%d", i)))
		assert.NoError(t, err)
	}

	jobOne, err := scheduler.GetJob("test-0")
	assert.NoError(t, err)
	time.Sleep(time.Until(jobOne.NextRun()) + 2*time.Second)

	assert.Equal(t, numCPU, int(jobsRun.Load()))
}

func TestScheduler_jobsRespectTimeLimit(t *testing.T) {
	t.Parallel()
	scheduler := newTestScheduler(1)
	testID := "test-0"

	err := scheduler.AddJob(JobConfig{
		ID:                  testID,
		Schedule:            "* * * * *",
		Timeout:             100 * time.Millisecond,
		StartingDeadline:    100 * time.Millisecond,
		AllowConccurentRuns: false,
		Func: func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		},
	})

	jobOne, err := scheduler.GetJob(testID)
	assert.NoError(t, err)
	time.Sleep(time.Until(jobOne.NextRun()) + 2*time.Second)

	jobOne, err = scheduler.GetJob(testID)
	assert.NoError(t, err)

	history := jobOne.History()
	assert.Equal(t, 1, len(history))
	for _, h := range history {
		t.Log(h.Error(), h.RanAt())
	}

	assert.Error(t, history[0].Error())
}

func TestScheduler_dontStartConcurrentRuns(t *testing.T) {
	t.Parallel()
	scheduler := newTestScheduler(1)
	testID := "test-0"

	err := scheduler.AddJob(JobConfig{
		ID:                  testID,
		Schedule:            "* * * * *",
		Timeout:             2 * time.Minute,
		StartingDeadline:    100 * time.Millisecond,
		AllowConccurentRuns: false,
		Func: func(ctx context.Context) error {
			time.Sleep(90 * time.Second)
			return nil
		},
	})

	jobOne, err := scheduler.GetJob(testID)
	assert.NoError(t, err)
	scheduledJobs := jobOne.NextFor(3 * time.Minute)
	furthestJobInFuture := scheduledJobs[len(scheduledJobs)-1]

	time.Sleep(time.Until(furthestJobInFuture))

	jobOne, err = scheduler.GetJob(testID)
	assert.NoError(t, err)

	history := jobOne.History()
	for _, h := range history {
		t.Logf("Ran at: %s / Error: %s", h.RanAt(), h.Error())
	}

	assert.NoError(t, history[0].Error())

	expectedErr := ErrJobRunning{}
	assert.EqualError(t, history[1].Error(), expectedErr.Error())
}
