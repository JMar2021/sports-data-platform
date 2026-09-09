package jobs

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/domain"
	"github.com/JMar2021/sports-data-platform/internal/repository"
)

type JobFactory func() (Job, error)

type Scheduler struct {
	Queue      *Queue
	Repository *repository.Repository
	Interval   time.Duration
	Factories  []JobFactory
}

func NewScheduler(queue *Queue, repository *repository.Repository, interval time.Duration, factories []JobFactory) *Scheduler {
	return &Scheduler{
		Queue:      queue,
		Repository: repository,
		Interval:   interval,
		Factories:  factories,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, factory := range s.Factories {
				job, err := factory()
				if err != nil {
					fmt.Println("Error creating job:", err)
					continue
				}

				s.createAndEnqueueJob(ctx, job)
			}

		case <-ctx.Done():
			return
		}
	}
}

func GenerateID() (string, error) {
	id := make([]byte, 16)

	_, err := rand.Read(id)
	if err != nil {
		fmt.Println("Error generating random bytes")
		return "", err
	}

	return hex.EncodeToString(id), nil
}

func NewScheduleJob(sport domain.Sport, date string) (Job, error) {
	id, err := GenerateID()
	if err != nil {
		return Job{}, err
	}

	return Job{
		ID:        id,
		Key:       string(sport) + ":get_schedule:" + date,
		Sport:     sport,
		Operation: OperationGetSchedule,
		Date:      date,
		Status:    StatusPending,
	}, nil
}

func NewStandingsJob(date string) (Job, error) {
	id, err := GenerateID()
	if err != nil {
		return Job{}, err
	}

	return Job{
		ID:        id,
		Key:       "mlb:get_standings:" + date,
		Sport:     domain.SportMLB,
		Operation: OperationGetStandings,
		Date:      date,
		Status:    StatusPending,
	}, nil
}

func (s *Scheduler) createAndEnqueueJob(
	ctx context.Context,
	job Job,
) {
	jobDate, err := time.Parse("2006-01-02", job.Date)
	if err != nil {
		fmt.Println("Error parsing job date:", err)
		return
	}

	created, err := s.Repository.CreateJob(
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
		fmt.Println("Error persisting job:", err)
		return
	}

	// An active job with this key already exists.
	if !created {
		return
	}

	if s.Queue.Enqueue(job) {
		return
	}

	// The job was persisted but could not be queued.
	// Mark it failed so it doesn't remain stuck as pending.
	err = s.Repository.FailJob(
		ctx,
		job.ID,
		job.Attempts,
		"failed to enqueue job",
		time.Now(),
	)
	if err != nil {
		fmt.Println("Error marking job as failed:", err)
	}
}
