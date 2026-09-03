package jobs

import (
	"sync"
)

type Queue struct {
	Jobs []Job
	mu   sync.Mutex
	cond *sync.Cond
}

func NewQueue() *Queue {
	q := &Queue{
		Jobs: []Job{},
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Enqueue(job Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Jobs = append(q.Jobs, job)
	q.cond.Signal()
}

func (q *Queue) Dequeue() (Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.Jobs) == 0 {
		q.cond.Wait()
	}

	job := q.Jobs[0]
	q.Jobs = q.Jobs[1:]
	return job, nil
}
