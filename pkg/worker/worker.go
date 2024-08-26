package worker

import (
	"sync"
	"time"
)

type Message[T any] struct {
	WorkerId int
	Duration time.Duration
	Result   T
}

type Task[T any] func() T

type IWorker[T any] interface {
	Work()
	Stop()
	AddTask(t Task[T])
	TasksCompleted() int
	NumberOfWorkers() int
}

func NewWorker[T any](numberOfWorkers int, messageChan chan Message[T]) IWorker[T] {
	return &Worker[T]{
		numberOfWorkers: numberOfWorkers,
		messageChan:     messageChan,
		taskChan:        make(chan Task[T]),
		stopWorkerChan:  make(chan bool, numberOfWorkers),
	}
}

var _ IWorker[any] = (*Worker[any])(nil)

type Worker[T any] struct {
	numberOfWorkers    int
	taskChan           chan Task[T]
	messageChan        chan Message[T]
	stopWorkerChan     chan bool
	jobsProcessed      int
	jobsProcessedMutex sync.Mutex
}

func (w *Worker[T]) Work() {
	w.jobsProcessed = 0
	for i := 1; i <= w.numberOfWorkers; i += 1 {
		go func(workerId int) {
			doWork := true
			for doWork {
				select {
				case t := <-w.taskChan:
					st := time.Now()
					r := t()
					w.incrementJobsProcessed()
					w.messageChan <- Message[T]{WorkerId: workerId, Duration: time.Now().Sub(st), Result: r}
				case <-w.stopWorkerChan:
					doWork = false
				}
			}
		}(i)
	}
}

func (w *Worker[T]) Stop() {
	for i := 0; i < w.numberOfWorkers; i += 1 {
		w.stopWorkerChan <- true
	}
}

func (w *Worker[T]) NumberOfWorkers() int {
	return w.numberOfWorkers
}

func (w *Worker[T]) AddTask(t Task[T]) {
	w.taskChan <- t
}

func (w *Worker[T]) TasksCompleted() int {
	return w.jobsProcessed
}

func (w *Worker[T]) incrementJobsProcessed() {
	w.jobsProcessedMutex.Lock()
	defer w.jobsProcessedMutex.Unlock()
	w.jobsProcessed += 1
}
