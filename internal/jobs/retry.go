package jobs

import (
	"context"
	"sync"
	"time"
)

type RetryManager struct {
	Queue *Queue
	WG    sync.WaitGroup
}

func (r *RetryManager) Retry(ctx context.Context, job Job) {
	r.WG.Add(1)
	go func() {
		defer r.WG.Done()
		delay := RetryDelay(job.Attempts)
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-timer.C:
			if !r.Queue.Requeue(job) {
				return
			}
		case <-ctx.Done():
			return
		}
	}()

}

func RetryDelay(attempt int) time.Duration {
	delay := 1 << (attempt - 1) * time.Second
	return delay
}
