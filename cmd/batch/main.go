package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JMar2021/sports-data-platform/database"
	"github.com/JMar2021/sports-data-platform/internal/config"
	"github.com/JMar2021/sports-data-platform/internal/domain"
	"github.com/JMar2021/sports-data-platform/internal/ingestion"
	"github.com/JMar2021/sports-data-platform/internal/jobs"
	"github.com/JMar2021/sports-data-platform/internal/mlb"
	"github.com/JMar2021/sports-data-platform/internal/nfl"
	"github.com/JMar2021/sports-data-platform/internal/repository"
)

const workerCount = 3

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	config, err := config.Load()
	if err != nil {
		log.Printf("Config not valid: %v", err)
		os.Exit(1)
	}

	dbPool, err := database.NewPool(ctx, config.DatabaseURL)
	if err != nil {
		log.Printf("Failed to create database pool: %v", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	repo := repository.NewRepository(dbPool)

	queue := jobs.NewQueue()

	retryManager := &jobs.RetryManager{
		Queue: queue,
	}

	mlbClient := mlb.NewClient(&http.Client{})
	nflClient := nfl.NewClient(&http.Client{})

	ingestor := &ingestion.MLBIngestor{
		MLBClient:  mlbClient,
		Repository: repo,
	}

	nflIngestor := &ingestion.NFLIngestor{
		NFLClient:  nflClient,
		Repository: repo,
	}

	var scheduleHandler jobs.JobHandler = func(
		ctx context.Context,
		job jobs.Job,
	) error {
		return ingestor.IngestSchedule(ctx, job.Date)
	}

	var nflScheduleHandler jobs.JobHandler = func(
		ctx context.Context,
		job jobs.Job,
	) error {
		return nflIngestor.IngestSchedule(ctx, job.Date)
	}

	var standingsHandler jobs.JobHandler = func(
		ctx context.Context,
		job jobs.Job,
	) error {
		return ingestor.IngestStandings(ctx, job.Date)
	}

	handlers := map[jobs.JobKey]jobs.JobHandler{
		{
			Sport:     domain.SportMLB,
			Operation: jobs.OperationGetSchedule,
		}: scheduleHandler,
		{
			Sport:     domain.SportNFL,
			Operation: jobs.OperationGetSchedule,
		}: nflScheduleHandler,
		{
			Sport:     domain.SportMLB,
			Operation: jobs.OperationGetStandings,
		}: standingsHandler,
	}

	factories := makeScheduleFactories(
		domain.SportMLB,
		-1,
		9,
	)

	factories = append(
		factories,
		makeScheduleFactories(
			domain.SportNFL,
			-1,
			9,
		)...,
	)

	factories = append(
		factories,
		func() (jobs.Job, error) {
			return jobs.NewStandingsJob(
				time.Now().Format("2006-01-02"),
			)
		},
	)

	executor := jobs.NewExecutor(logger, handlers)

	batchRunner := jobs.NewBatchRunner(
		queue,
		executor,
		repo,
		retryManager,
		logger,
		factories,
		workerCount,
	)

	logger.Info(
		"starting batch",
		"workers", workerCount,
	)

	if err := batchRunner.Run(ctx); err != nil {
		logger.Error(
			"batch failed",
			"error", err,
		)

		retryManager.WG.Wait()
		os.Exit(1)
	}

	retryManager.WG.Wait()

	logger.Info("batch completed successfully")
}

func makeScheduleFactories(
	sport domain.Sport,
	startOffset int,
	numberOfDays int,
) []jobs.JobFactory {
	factories := make([]jobs.JobFactory, 0, numberOfDays)

	for offset := 0; offset < numberOfDays; offset++ {
		offset := offset

		factories = append(
			factories,
			func() (jobs.Job, error) {
				date := time.Now().
					AddDate(0, 0, startOffset+offset).
					Format("2006-01-02")

				return jobs.NewScheduleJob(sport, date)
			},
		)
	}

	return factories
}
