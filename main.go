package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Success bool
	Duration time.Duration
	Error    error
}

func worker(url string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		results <- Result{Success: false, Duration: time.Since(start), Error: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		results <- Result{Success: true, Duration: time.Since(start), Error: nil}
	} else {
		results <- Result{Success: false, Duration: time.Since(start), Error: fmt.Errorf("status code: %d", resp.StatusCode)}
	}
}

func main() {
	url := "http://localhost:3000" 
	numRequests := 1000           // Total # of requests
	concurrent := 10              // # of concurrent workers (10)

	results := make(chan Result, numRequests)
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go worker(url, results, &wg)

		// Control concurrency
		if (i+1)%concurrent == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	close(results)

	totalDuration := time.Since(start)
	successCount := 0
	failureCount := 0
	var totalLatency time.Duration

	for result := range results {
		if result.Success {
			successCount++
			totalLatency += result.Duration
		} else {
			failureCount++
			fmt.Printf("Error: %v\n", result.Error)
		}
	}

	averageLatency := totalLatency / time.Duration(successCount)

	fmt.Printf("\nResults:\n")
	fmt.Printf("Total Requests: %d\n", numRequests)
	fmt.Printf("Successful Requests: %d\n", successCount)
	fmt.Printf("Failed Requests: %d\n", failureCount)
	fmt.Printf("Total Duration: %v\n", totalDuration)
	fmt.Printf("Average Latency: %v\n", averageLatency)
}
