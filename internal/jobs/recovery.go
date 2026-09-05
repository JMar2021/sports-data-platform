package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/repository"
)

const StaleJobAge = 5 * time.Minute

func RecoverJobs(
	ctx context.Context,
	repository *repository.Repository,
	queue *Queue,
	logger *slog.Logger,
) error {
	staleBefore := time.Now().Add(-StaleJobAge)

	jobs, err := repository.RecoverStaleJobs(ctx, staleBefore)
	if err != nil {
		return fmt.Errorf("recover stale jobs: %w", err)
	}

	for _, staleJob := range jobs {
		job := Job{
			ID:        staleJob.ID,
			Key:       staleJob.Key,
			Sport:     Sport(staleJob.Sport),
			Operation: Operation(staleJob.Operation),
			Date:      staleJob.Date.Format("2006-01-02"),
			Attempts:  staleJob.Attempts,
			Status:    StatusRetrying,
		}

		if !queue.Enqueue(job) {
			logger.Error(
				"failed to requeue recovered job",
				"job_id", job.ID,
			)

			continue
		}

		logger.Info(
			"recovered stale job",
			"job_id", job.ID,
			"attempt", job.Attempts,
		)
	}

	return nil
}
