package jobs

import (
	"context"

	"github.com/JMar2021/sports-data-platform/internal/domain"
)

type Operation string

const (
	OperationGetSchedule  Operation = "get_schedule"
	OperationGetStandings Operation = "get_standings"
)

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
	StatusRetrying  JobStatus = "retrying"
)

type Job struct {
	ID        string
	Key       string
	Sport     domain.Sport
	Operation Operation
	Date      string
	Attempts  int
	Status    JobStatus
}

type JobHandler func(context.Context, Job) error

type JobKey struct {
	Sport     domain.Sport
	Operation Operation
}
