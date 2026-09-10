package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/repository"
)

type BatchRunner struct {
	Queue        *Queue
	Executor     *Executor
	Repository   *repository.Repository
	RetryManager *RetryManager
	Logger       *slog.Logger
	Factories    []JobFactory
	WorkerCount  int
}

func NewBatchRunner(
	queue *Queue,
	executor *Executor,
	repository *repository.Repository,
	retryManager *RetryManager,
	logger *slog.Logger,
	factories []JobFactory,
	workerCount int,
) *BatchRunner {
	return &BatchRunner{
		Queue:        queue,
		Executor:     executor,
		Repository:   repository,
		RetryManager: retryManager,
		Logger:       logger,
		Factories:    factories,
		WorkerCount:  workerCount,
	}
}

func (b *BatchRunner) Run(ctx context.Context) error {
	jobs, err := b.createJobs(ctx)
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		b.Logger.Info("batch contains no new jobs")
		return nil
	}

	tracker := NewBatchTracker()

	var workerWG sync.WaitGroup

	for i := 0; i < b.WorkerCount; i++ {
		worker := &Worker{
			ID:            i + 1,
			Queue:         b.Queue,
			Executor:      b.Executor,
			Logger:        b.Logger,
			RetryManager:  b.RetryManager,
			Repository:    b.Repository,
			OnJobComplete: tracker.Done,
		}

		workerWG.Add(1)

		go func(w *Worker) {
			defer workerWG.Done()
			w.Run(ctx, ctx)
		}(worker)
	}

	for _, job := range jobs {
		tracker.Add()

		if !b.Queue.Enqueue(job) {
			tracker.Done(job)

			b.Logger.Error(
				"failed to enqueue job",
				"job_id", job.ID,
			)

			b.Queue.Close()
			workerWG.Wait()

			return fmt.Errorf("failed to enqueue job %s", job.ID)
		}
	}

	b.Logger.Info(
		"batch started",
		"jobs", len(jobs),
		"workers", b.WorkerCount,
	)

	tracker.Wait()

	b.Logger.Info(
		"all batch jobs reached terminal state",
		"failed_jobs", tracker.FailedJobs(),
	)

	b.Queue.Close()
	workerWG.Wait()

	if tracker.FailedJobs() > 0 {
		return fmt.Errorf(
			"batch completed with %d failed jobs",
			tracker.FailedJobs(),
		)
	}

	return nil
}

func (b *BatchRunner) createJobs(ctx context.Context) ([]Job, error) {
	jobs := make([]Job, 0, len(b.Factories))

	for _, factory := range b.Factories {
		job, err := factory()
		if err != nil {
			return nil, fmt.Errorf("create job: %w", err)
		}

		jobDate, err := time.Parse("2006-01-02", job.Date)
		if err != nil {
			return nil, fmt.Errorf(
				"parse job date %q: %w",
				job.Date,
				err,
			)
		}

		created, err := b.Repository.CreateJob(
			ctx,
			job.ID,
			job.Key,
			string(job.Sport),
			string(job.Operation),
			jobDate,
			string(job.Status),
			job.Attempts,
			time.Now(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"persist job %s: %w",
				job.ID,
				err,
			)
		}

		if !created {
			b.Logger.Info(
				"job already exists, skipping",
				"job_id", job.ID,
				"job_key", job.Key,
			)
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
