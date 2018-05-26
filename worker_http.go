/*

worker_http.go

Defined in this file is the httpWorkerState struct which satisfies
the workerRunner interface.

httpWorkerState's Run method performs a simple HTTP GET and returns the
time needed to execute the request.

At the time of writing it's the only composer of workerState, howerver
this is meant to signal a direction towards multiple benchmark worker types.

*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// TODO later extend this to different types of HTTP requests
type httpWorkerState struct {
	url string
	workerState
}

func (h *httpWorkerState) Run() float64 {
	return timed(func() { httpGet(h.url) })
}

func timed(fn func()) float64 {
	start_time := time.Now()
	fn()
	end_time := time.Now()
	elapsed_time := end_time.Sub(start_time).Seconds()
	return elapsed_time
}

func httpGet(url string) {
	_, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

func NewHTTPWorkers(num_workers uint64, url string) *[]workerRunner {
	workers := []workerRunner{}
	for i := uint64(0); i < num_workers; i++ {
		fmt.Printf("Creating worker with id: %v\n", i)
		var worker workerRunner
		worker = &httpWorkerState{url: url, workerState: workerState{id: i, results: new([]float64)}}
		workers = append(workers, worker)
	}
	return &workers
}
