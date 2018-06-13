/*

  krush

  Usage: krush -h

  Execute a simple HTTP GET benchmark against the specified host.

  TODO:
   - tests
   - show latency percentiles (ext dep)

*/
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

    "github.com/ksnavely/krush/internal/controller"
    "github.com/ksnavely/krush/internal/worker"
)

type cliOpts struct {
	TargetHost  string
	Concurrency uint64
	RunDuration time.Duration
}

func parseCLI() cliOpts {
	args := cliOpts{}
	flag.StringVar(&args.TargetHost, "TargetHost", "", "The target URL to benchmark")
	flag.Uint64Var(&args.Concurrency, "Concurrency", 0, "The number of worker goroutines to use")
	flag.DurationVar(&args.RunDuration, "RunDuration", 0, "The number of worker goroutines to use")
	flag.Parse()
	fmt.Printf("\n  Arguments:\n%+v\n\n", args)
	validateCLI(args)
	return args
}

func validateCLI(args cliOpts) {
	var failed bool
	if len(args.TargetHost) == 0 {
		fmt.Fprintf(os.Stderr, "[Error] A target URL must be specified.\n")
		failed = true
	}
	if args.Concurrency == 0 {
		fmt.Printf("[Warning] Concurrency is zero, no workers will be spawned.\n")
	}
	if args.RunDuration < 0 {
		fmt.Fprintf(os.Stderr, "[Error] Benchmark duration can't be negative.\n")
		failed = true
	}
	if failed {
		fmt.Fprintf(os.Stderr, "\nkrush has encountered invalid CLI options and will exit.\n\n")
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("Welcome to Krush\n")
	args := parseCLI()

	workers := worker.NewHTTPWorkers(args.Concurrency, args.TargetHost)
	controller := controller.NewBenchmarkController(args.RunDuration, workers)
	results := controller.RunBenchmark()

	var sum float64
	for i := range results {
		sum += results[i]
	}
	fmt.Printf("\nTotal results length: %v\n", len(controller.Results))
	if len(results) > 0 {
		fmt.Printf("Average result: %v\n", sum/float64(len(results)))
	} else {
		fmt.Fprintf(os.Stderr, "krush didn't find any results")
	}
	fmt.Printf("\nKrush has finished.\n\n")
}
