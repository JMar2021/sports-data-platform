package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/jobs"
	"github.com/JMar2021/sports-data-platform/internal/mlb"
)

func main() {
	httpClient := &http.Client{}
	client := mlb.NewClient(httpClient)

	executor := &jobs.Executor{
		Client: client,
	}
	job := jobs.Job{
		Sport:     "mlb",
		Operation: "get_schedule",
		Date:      time.Now().Format("2006-01-02"),
	}
	err := executor.Execute(job)
	if err != nil {
		fmt.Println("Error executing job:", err)
	}
}
