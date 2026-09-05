package jobs

import (
	"errors"
	"sync"
)

var ErrQueueClosed = errors.New("queue closed")

type JobState string

const (
	JobQueued   JobState = "queued"
	JobRunning  JobState = "running"
	JobRetrying JobState = "retrying"
)

type Queue struct {
	Jobs   []Job
	mu     sync.Mutex
	cond   *sync.Cond
	seen   map[string]JobState
	closed bool
}

func NewQueue() *Queue {
	q := &Queue{
		Jobs: make([]Job, 0),
		seen: make(map[string]JobState),
	}

	q.cond = sync.NewCond(&q.mu)

	return q
}

func (q *Queue) Enqueue(job Job) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return false
	}

	if _, exists := q.seen[job.Key]; exists {
		return false
	}

	q.Jobs = append(q.Jobs, job)
	q.seen[job.Key] = JobQueued

	q.cond.Signal()

	return true
}

func (q *Queue) Dequeue() (Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.Jobs) == 0 && !q.closed {
		q.cond.Wait()
	}

	if len(q.Jobs) == 0 && q.closed {
		return Job{}, ErrQueueClosed
	}

	job := q.Jobs[0]
	q.Jobs = q.Jobs[1:]
	q.seen[job.Key] = JobRunning

	return job, nil
}

func (q *Queue) Requeue(job Job) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return false
	}

	q.Jobs = append(q.Jobs, job)
	q.seen[job.Key] = JobRetrying

	q.cond.Signal()

	return true
}

func (q *Queue) Complete(job Job) {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.seen, job.Key)
}

func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.closed = true
	q.cond.Broadcast()
}
