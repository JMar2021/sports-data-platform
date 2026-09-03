package jobs

import (
	"context"
	"fmt"

	"github.com/JMar2021/sports-data-platform/internal/ingestion"
)

type Executor struct {
	Ingestor *ingestion.MLBIngestor
}

func (e *Executor) Execute(ctx context.Context, job Job) error {
	switch job.Sport {
	case SportMLB:
		switch job.Operation {
		case OperationGetSchedule:
			return e.Ingestor.IngestSchedule(ctx, job.Date)
		default:
			return fmt.Errorf("unsupported operation: %s", job.Operation)
		}
	default:
		return fmt.Errorf("unsupported sport: %s", job.Sport)
	}
}
