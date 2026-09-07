package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/JMar2021/sports-data-platform/database"
	"github.com/JMar2021/sports-data-platform/internal/api"
	"github.com/JMar2021/sports-data-platform/internal/config"
	"github.com/JMar2021/sports-data-platform/internal/ingestion"
	"github.com/JMar2021/sports-data-platform/internal/jobs"
	"github.com/JMar2021/sports-data-platform/internal/mlb"
	"github.com/JMar2021/sports-data-platform/internal/repository"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Create a context that will be cancelled when the program
	// receives SIGINT (Ctrl+C) or SIGTERM.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	jobCtx := context.Background()
	config, err := config.Load()
	if err != nil {
		fmt.Printf("Config not valid")
		return
	}
	dbPool, err := database.NewPool(ctx, config.DatabaseURL)
	if err != nil {
		fmt.Printf("Failed to create database pool: %v\n", err)
		return
	}
	defer dbPool.Close()
	fmt.Println("Database connection pool created successfully.")
	repo := repository.NewRepository(dbPool)
	apiServer := api.NewServer(repo, config.HTTPAddr)
	go func() {
		if err := apiServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("API server error: %v", err)
		}
	}()
	version, err := repo.GetDatabaseVersion(ctx)
	if err != nil {
		fmt.Printf("Failed to query database version: %v\n", err)
		return
	}
	fmt.Printf("Connected to database: %s\n", version)

	// Create the job queue.
	queue := jobs.NewQueue()
	retryManager := &jobs.RetryManager{
		Queue: queue,
	}
	if err := jobs.RecoverJobs(
		ctx,
		repo,
		queue,
		logger,
	); err != nil {
		logger.Error(
			"failed to recover jobs",
			"error", err,
		)

		return
	}
	// Create the MLB client.
	mlbClient := mlb.NewClient(&http.Client{})
	ingestor := &ingestion.MLBIngestor{
		MLBClient:  mlbClient,
		Repository: repo,
	}
	var scheduleHandler jobs.JobHandler = func(ctx context.Context, job jobs.Job) error {
		return ingestor.IngestSchedule(ctx, job.Date)
	}
	var standingsHandler jobs.JobHandler = func(ctx context.Context, job jobs.Job) error {
		return ingestor.IngestStandings(ctx, job.Date)
	}
	handlers := map[jobs.JobKey]jobs.JobHandler{
		{
			Sport:     jobs.SportMLB,
			Operation: jobs.OperationGetSchedule,
		}: scheduleHandler,

		{
			Sport:     jobs.SportMLB,
			Operation: jobs.OperationGetStandings,
		}: standingsHandler,
	}
	scheduleFactory := func() (jobs.Job, error) {
		return jobs.NewScheduleJob(time.Now().Format("2006-01-02"))
	}
	standingsFactory := func() (jobs.Job, error) {
		return jobs.NewStandingsJob(time.Now().Format("2006-01-02"))
	}
	factories := []jobs.JobFactory{scheduleFactory, standingsFactory}
	executor := jobs.NewExecutor(logger, handlers)
	// Create the scheduler.
	scheduler := jobs.NewScheduler(queue, repo, 10*time.Second, factories)

	// Create and start workers.
	var wg sync.WaitGroup
	var producerWG sync.WaitGroup
	producerWG.Add(1)
	go func() {
		defer producerWG.Done()
		scheduler.Run(ctx)
	}()

	for i := 1; i <= 3; i++ {
		worker := &jobs.Worker{
			ID:           i,
			Queue:        queue,
			Executor:     executor,
			Logger:       logger,
			RetryManager: retryManager,
			Repository:   repo,
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			worker.Run(ctx, jobCtx)
		}()
	}

	fmt.Println("Sports Data Platform running. Press Ctrl+C to stop.")

	// Wait until the context is cancelled.
	<-ctx.Done()
	fmt.Println("Shutdown signal received.")
	producerWG.Wait()
	retryManager.WG.Wait()
	queue.Close()
	wg.Wait()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("API server shutdown error: %v", err)
	}

	fmt.Println("Sports Data Platform stopped.")
}
