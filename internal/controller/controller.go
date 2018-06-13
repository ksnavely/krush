/*

controller.go

This file implements the benchmarkController struct and its methods.

The benchmarkController is a key orchestrator of benchmark activity.
Notably, it holds Workers, a pointer to a slice of structs which satisfy
the WorkerRunner interface.

The method RunBenchmark triggers to overall execution of benchmarking across
the specified number of workers. The benchmark runs until RunDuration has
passed at which point the workers are terminated.

*/
package controller

import (
	"fmt"
	"sync"
	"time"

    "github.com/ksnavely/krush/internal/worker"
)

type benchmarkController struct {
	Results     []float64
	Workers     *[]worker.WorkerRunner
	StopChan    chan string
	RunDuration time.Duration
}

func (c *benchmarkController) StopBenchmark(reason string) {
	for _ = range *c.Workers {
		c.StopChan <- reason
	}
}

func (c *benchmarkController) RunBenchmark() []float64 {
	var wg sync.WaitGroup
	resultsChan := make(chan float64)
	c.StopChan = make(chan string, len(*c.Workers))

	for _, worker := range *c.Workers {
		wg.Add(1)
		thisworker := worker
		go c.ManagedRun(&thisworker, resultsChan, c.StopChan, &wg)
	}

	go func() {
		<-time.After(c.RunDuration)
		c.StopBenchmark("completed")
	}()

	fmt.Printf("\n  Waiting on worker execution...\n")
	wg.Wait()

	for _, worker := range *c.Workers {
		c.Results = append(c.Results, (*worker.Results())...)
	}
	return c.Results
}

func (c *benchmarkController) ManagedRun(workerPtr *worker.WorkerRunner, resultsChan chan<- float64, stopChan <-chan string, wg *sync.WaitGroup) {
	var result float64
	var stop bool
	worker := *workerPtr
	stopMsg := "noreason"

	fmt.Printf("Worker starting: %v\n", worker.Id())

	for {
		select {
		case stopMsg = <-stopChan:
			stop = true
		default:
		}
		if stop {
			break
		}

		result = worker.Run()
		worker.AppendResult(result)
	}

	fmt.Printf("Terminating worker with reason: %v: %v\n", worker.Id(), stopMsg)
	wg.Done()
}

func NewBenchmarkController(runDuration time.Duration, workers *[]worker.WorkerRunner) benchmarkController {
	fmt.Printf("\n  Initializing the benchmarkController...\n")
	c := benchmarkController{RunDuration: runDuration, Workers: workers}
	return c
}
