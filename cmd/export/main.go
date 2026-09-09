package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/JMar2021/sports-data-platform/database"
	"github.com/JMar2021/sports-data-platform/internal/config"
	exporter "github.com/JMar2021/sports-data-platform/internal/export"
	"github.com/JMar2021/sports-data-platform/internal/repository"
)

func main() {
	ctx := context.Background()

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := database.NewPool(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	repo := repository.NewRepository(dbPool)

	now := time.Now()
	location := now.Location()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		location,
	)

	start := today.AddDate(0, 0, -1)
	end := today.AddDate(0, 0, 8)

	games, err := repo.GetGames(ctx, start, end)
	if err != nil {
		log.Fatal(err)
	}

	outputPath := os.Getenv("EXPORT_OUTPUT_PATH")
	if outputPath == "" {
		outputPath = "site/data/games.json"
	}

	if err := exporter.ExportGames(outputPath, games); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Exported %d games to %s\n", len(games), outputPath)
}
