package main

import (
	"context"
	"fmt"
	"time"

	"github.com/drewgonzales360/cronroutine"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func main() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log := zapr.NewLogger(zapLog)

	scheduler := cronroutine.StartNewScheduler(&cronroutine.SchedulerConfig{
		Logger:       log,
		HistoryLimit: 10,
		WorkerCount:  1,
	})

	done := make(chan struct{})

	err = scheduler.AddJob(cronroutine.JobConfig{
		ID:                  "test",
		Schedule:            "* * * * *",
		Timeout:             2 * time.Second,
		StartingDeadline:    100 * time.Millisecond,
		AllowConccurentRuns: false,
		Func: func(ctx context.Context) error {
			fmt.Println("test job finished at", time.Now().Format(time.RFC3339))
			done <- struct{}{}
			return nil
		},
	})

	<-done
}
