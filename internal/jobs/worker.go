package jobs

import (
	"context"
	"fmt"
)

type Worker struct {
	ID       int
	Queue    chan Job
	Executor *Executor
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case job := <-w.Queue:
			fmt.Printf("Worker %d executing job: %+v\n", w.ID, job)

			err := w.Executor.Execute(ctx, job)
			if err != nil {
				fmt.Printf("Worker %d error executing job: %v\n", w.ID, err)
			}

		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down.\n", w.ID)
			return
		}
	}
}
