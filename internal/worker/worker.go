/*

worker.go

This file defines two key concepts:
  - workerState: This struct type serves as a basic element
    for composition in other struct types which fulfill the
    WorkerRunner interface.
  - WorkerRunner: This interface demands the existence of the
    Run method, which is the argument-less entrypoint for execution
    of a single benchmark worker's task.

Note that while workerState meets most of the specification of the
WorkerRunner interface, it's missing the Run method which should be defined
on concrete composers of the workerState struct.

*/
package worker

type WorkerRunner interface {
	Id() uint64
	Run() float64
	Results() *[]float64
	AppendResult(float64)
}

type workerState struct {
	id      uint64
	results *[]float64
}

func (w *workerState) Id() uint64 {
	return w.id
}

func (w *workerState) AppendResult(result float64) {
	*(w.results) = append(*(w.results), result)
}

func (w *workerState) Results() *[]float64 {
	return w.results
}
