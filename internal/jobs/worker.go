package jobs

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/repository"
)

type Worker struct {
	ID            int
	Queue         *Queue
	Executor      *Executor
	Logger        *slog.Logger
	RetryManager  *RetryManager
	Repository    *repository.Repository
	OnJobComplete func(Job)
}

const MaxAttempts = 3

func (w *Worker) jobComplete(job Job) {
	if w.OnJobComplete != nil {
		w.OnJobComplete(job)
	}
}

func (w *Worker) Run(shutdownCtx context.Context, jobCtx context.Context) {
	for {
		job, err := w.Queue.Dequeue()
		if err != nil {
			if errors.Is(err, ErrQueueClosed) {
				w.Logger.Info("queue closed", "worker", w.ID)
				return
			}

			w.Logger.Error(
				"queue error",
				"worker", w.ID,
				"error", err,
			)
			return
		}

		w.executeJob(jobCtx, shutdownCtx, job)
	}
}

func (w *Worker) executeJob(
	jobCtx context.Context,
	shutdownCtx context.Context,
	job Job,
) {
	job.Attempts++
	job.Status = StatusRunning

	if err := w.Repository.StartJob(
		jobCtx,
		job.ID,
		time.Now(),
	); err != nil {
		w.Logger.Error(
			"failed to start job",
			"worker", w.ID,
			"job_id", job.ID,
			"error", err,
		)

		w.Queue.Complete(job)
		w.jobComplete(job)

		return
	}

	w.Logger.Info(
		"executing job",
		"worker", w.ID,
		"job_id", job.ID,
		"sport", job.Sport,
		"operation", job.Operation,
		"attempt", job.Attempts,
	)

	err := w.Executor.Execute(jobCtx, job)
	if err == nil {
		job.Status = StatusCompleted

		if err := w.Repository.CompleteJob(
			jobCtx,
			job.ID,
			time.Now(),
		); err != nil {
			w.Logger.Error(
				"failed to complete job",
				"job_id", job.ID,
				"error", err,
			)
		}

		w.Queue.Complete(job)
		w.jobComplete(job)

		w.Logger.Info(
			"job complete",
			"worker", w.ID,
			"job_id", job.ID,
		)

		return
	}

	w.Logger.Error(
		"job execution failed",
		"worker", w.ID,
		"job_id", job.ID,
		"attempt", job.Attempts,
		"error", err,
	)

	if job.Attempts < MaxAttempts && shutdownCtx.Err() == nil {
		job.Status = StatusRetrying

		if repoErr := w.Repository.RetryJob(
			jobCtx,
			job.ID,
			job.Attempts,
			err.Error(),
		); repoErr != nil {
			w.Logger.Error(
				"failed to mark job for retry",
				"job_id", job.ID,
				"error", repoErr,
			)

			w.Queue.Complete(job)
			w.jobComplete(job)

			return
		}

		w.Logger.Info(
			"retrying job",
			"worker", w.ID,
			"job_id", job.ID,
			"attempt", job.Attempts,
		)

		w.RetryManager.Retry(shutdownCtx, job)

		return
	}

	job.Status = StatusFailed

	if repoErr := w.Repository.FailJob(
		jobCtx,
		job.ID,
		job.Attempts,
		err.Error(),
		time.Now(),
	); repoErr != nil {
		w.Logger.Error(
			"failed to mark job as failed",
			"job_id", job.ID,
			"error", repoErr,
		)
	}

	w.Logger.Error(
		"job abandoned",
		"worker", w.ID,
		"job_id", job.ID,
		"attempt", job.Attempts,
	)

	w.Queue.Complete(job)
	w.jobComplete(job)
}
