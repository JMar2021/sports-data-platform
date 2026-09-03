package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/JMar2021/sports-data-platform/database"
	"github.com/JMar2021/sports-data-platform/internal/api"
	"github.com/JMar2021/sports-data-platform/internal/ingestion"
	"github.com/JMar2021/sports-data-platform/internal/jobs"
	"github.com/JMar2021/sports-data-platform/internal/mlb"
	"github.com/JMar2021/sports-data-platform/internal/repository"
)

func main() {
	// Create a context that will be cancelled when the program
	// receives SIGINT (Ctrl+C) or SIGTERM.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	connString := os.Getenv("SPORTS_DATABASE_URL")
	dbPool, err := database.NewPool(ctx, connString)
	if err != nil {
		fmt.Printf("Failed to create database pool: %v\n", err)
		return
	}
	defer dbPool.Close()
	fmt.Println("Database connection pool created successfully.")
	repo := repository.NewRepository(dbPool)
	apiServer := api.NewServer(repo)
	go func() {
		if err := apiServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("API server error: %v", err)
		}
	}()
	games, err := repo.GetGames(context.Background())
	if err != nil {
		fmt.Errorf("Failed to get games: %v", err)
	}

	for _, game := range games {
		fmt.Printf(
			"%s %d @ %s %d — %s — %s\n",
			game.AwayTeam,
			game.AwayScore,
			game.HomeTeam,
			game.HomeScore,
			game.Status,
			game.Venue,
		)
	}
	teamID, err := repo.UpsertTeam(ctx, 111, "Test Team")
	if err != nil {
		fmt.Println("Failed to upsert team:", err)
		return
	}

	fmt.Println("Team database ID:", teamID)
	version, err := repo.GetDatabaseVersion(ctx)
	if err != nil {
		fmt.Printf("Failed to query database version: %v\n", err)
		return
	}
	fmt.Printf("Connected to database: %s\n", version)

	// Create the job queue.
	queue := make(chan jobs.Job, 10)

	// Create the MLB client.
	mlbClient := mlb.NewClient(&http.Client{})
	ingestor := &ingestion.MLBIngestor{
		MLBClient:  mlbClient,
		Repository: repo,
	}
	// Create the executor.
	executor := &jobs.Executor{
		Ingestor: ingestor,
	}

	// Create the scheduler.
	scheduler := jobs.NewScheduler(queue, 10*time.Second)

	// Start the scheduler.
	go scheduler.Run(ctx)

	// Create and start workers.
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		worker := &jobs.Worker{
			ID:       i,
			Queue:    queue,
			Executor: executor,
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			worker.Run(ctx)
		}()
	}

	fmt.Println("Sports Data Platform running. Press Ctrl+C to stop.")

	// Wait until the context is cancelled.
	<-ctx.Done()

	fmt.Println("Shutdown signal received.")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("API server shutdown error: %v", err)
	}

	// Wait for workers to finish.
	wg.Wait()

	fmt.Println("Sports Data Platform stopped.")
}
