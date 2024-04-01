package cronroutine

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/cilium/workerpool"
	"github.com/go-logr/logr"
)

type Scheduler struct {
	jobsLock sync.RWMutex
	jobs     map[string]*jobMetadata

	logger       logr.Logger
	historyLimit int
	workerpool   *workerpool.WorkerPool
}

type SchedulerConfig struct {
	Logger       logr.Logger
	HistoryLimit int
	WorkerCount  int
}

func DefaultSchedulerConfig() *SchedulerConfig {
	return &SchedulerConfig{
		Logger:       logr.Discard(),
		HistoryLimit: 10,
		WorkerCount:  runtime.NumCPU(),
	}
}

func StartNewScheduler(cfg *SchedulerConfig) *Scheduler {
	if cfg == nil {
		cfg = DefaultSchedulerConfig()
	}

	s := &Scheduler{
		jobsLock: sync.RWMutex{},
		jobs:     make(map[string]*jobMetadata),

		workerpool:   workerpool.New(cfg.WorkerCount),
		logger:       cfg.Logger,
		historyLimit: cfg.HistoryLimit,
	}

	s.start()
	return s
}

func StopScheduler(s *Scheduler) error {
	if err := s.workerpool.Close(); err != nil {
		return fmt.Errorf("failed to close worker pool: %w", err)
	}

	// The workerpool keeps a history of all the tasks. We don't need to keep
	// the history, so this library may not be suitable.
	if _, err := s.workerpool.Drain(); err != nil {
		return fmt.Errorf("failed to drain worker pool: %w", err)
	}

	return nil
}

func (s *Scheduler) AddJob(job JobConfig) error {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()

	if _, existed := s.jobs[job.ID]; existed {
		return fmt.Errorf("job with ID %s already exists", job.ID)
	}

	cron, err := ParseCron(job.Schedule)
	if err != nil {
		return fmt.Errorf("failed to parse cron schedule: %w", err)
	}

	s.jobs[job.ID] = &jobMetadata{
		jobConfig:    &job,
		history:      make([]*History, 0, s.historyLimit),
		historyLimit: s.historyLimit,
		cron:         cron,
		running:      sync.Mutex{},
	}

	return nil
}

func (s *Scheduler) ListJobs() []*Job {
	s.jobsLock.RLock()
	defer s.jobsLock.RUnlock()

	jobs := make([]*Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job.Job())
	}

	return jobs
}

func (s *Scheduler) GetJob(jobID string) (*Job, error) {
	s.jobsLock.RLock()
	defer s.jobsLock.RUnlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return nil, fmt.Errorf("job with ID %s does not exist", jobID)
	}

	return job.Job(), nil
}

func (s *Scheduler) RemoveJob(jobID string) error {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()

	if _, existed := s.jobs[jobID]; !existed {
		return fmt.Errorf("job with ID %s does not exist", jobID)
	}

	delete(s.jobs, jobID)

	return nil
}

func (s *Scheduler) start() {
	type scheduledJob struct {
		job       *jobMetadata
		startTime time.Time
	}

	workQueue := make(chan *scheduledJob, 100)
	go func() {
		wLog := s.logger.WithName("worker-loop")
		wLog.Info("worker loop started")
		for job := range workQueue {
			time.Sleep(time.Until(job.startTime))
			wLog.Info("job started", "job_id", job.job.ID(), "time", job.startTime)
			err := s.workerpool.Submit(job.job.ID(), job.job.run(s.logger.WithName(job.job.ID()), job.startTime))
			if err != nil {
				wLog.Error(err, "failed to submit job to worker pool", "job_id", job.job.ID())
			}
		}
	}()

	go func() {
		qLog := s.logger.WithName("queue-loop")
		qLog.Info("queue loop started")
		for {
			scheduledJobs := make([]*scheduledJob, 0, 100)

			s.jobsLock.RLock()
			for _, job := range s.jobs {
				schedule := job.cron.NextFor(3 * time.Minute)
				for _, t := range schedule {
					scheduledJobs = append(scheduledJobs, &scheduledJob{
						job:       job,
						startTime: t,
					})
				}
			}
			s.jobsLock.RUnlock()

			sort.SliceStable(scheduledJobs, func(i, j int) bool {
				return scheduledJobs[i].startTime.Before(scheduledJobs[j].startTime)
			})

			lastJobTime := time.Now().UTC()
			if len(scheduledJobs) > 0 {
				lastJobTime = scheduledJobs[len(scheduledJobs)-1].startTime
			}

			for _, job := range scheduledJobs {
				qLog.Info("job scheduled", "job_id", job.job.ID(), "time", job.startTime)
				workQueue <- job
			}

			qLog.Info("sleeping for until last job", "time", lastJobTime)
			time.Sleep(time.Until(lastJobTime.Add(10 * time.Millisecond)))
		}
	}()
}
