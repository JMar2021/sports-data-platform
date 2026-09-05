package jobs

import (
	"context"
	"fmt"
	"log/slog"
)

type Executor struct {
	Logger   *slog.Logger
	Handlers map[JobKey]JobHandler
}

func NewExecutor(logger *slog.Logger, handlers map[JobKey]JobHandler) *Executor {
	return &Executor{Logger: logger, Handlers: handlers}
}
func (e *Executor) Execute(ctx context.Context, job Job) error {
	key := JobKey{
		Sport:     job.Sport,
		Operation: job.Operation,
	}

	handler, ok := e.Handlers[key]
	if !ok {
		e.Logger.Error("unsupported job",
			"job_id", job.ID,
			"sport", job.Sport,
			"operation", job.Operation,
		)

		return fmt.Errorf(
			"unsupported job: sport=%s operation=%s",
			job.Sport,
			job.Operation,
		)
	}
	err := handler(ctx, job)
	if err != nil {
		e.Logger.Error("job failed",
			"job_id", job.ID,
			"error", err)
		return err
	}
	e.Logger.Info("job completed", "job_id", job.ID)
	return nil
}
