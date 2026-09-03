package jobs

import (
	"context"
	"time"
)

type Scheduler struct {
	Queue    chan Job
	Interval time.Duration
}

func NewScheduler(queue chan Job, interval time.Duration) *Scheduler {
	return &Scheduler{
		Queue:    queue,
		Interval: interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			job := Job{
				Sport:     SportMLB,
				Operation: OperationGetSchedule,
				Date:      time.Now().Format("2006-01-02"),
			}

			s.Queue <- job

		case <-ctx.Done():
			return
		}
	}
}
