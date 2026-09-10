package jobs

import "sync"

type BatchTracker struct {
	wg         sync.WaitGroup
	mu         sync.Mutex
	failedJobs int
}

func NewBatchTracker() *BatchTracker {
	return &BatchTracker{}
}

func (b *BatchTracker) Add() {
	b.wg.Add(1)
}

func (b *BatchTracker) Done(job Job) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if job.Status == StatusFailed {
		b.failedJobs++
	}

	b.wg.Done()
}

func (b *BatchTracker) Wait() {
	b.wg.Wait()
}

func (b *BatchTracker) FailedJobs() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.failedJobs
}
