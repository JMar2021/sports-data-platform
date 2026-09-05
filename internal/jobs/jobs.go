package jobs

import "context"

type Sport string

const (
	SportMLB Sport = "mlb"
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
	Sport     Sport
	Operation Operation
	Date      string
	Attempts  int
	Status    JobStatus
}

type JobHandler func(context.Context, Job) error

type JobKey struct {
	Sport     Sport
	Operation Operation
}
