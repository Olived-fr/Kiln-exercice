package worker

import "context"

type Task func(ctx context.Context, input any) (any, error)

type Worker struct {
	id         int
	taskQueue  <-chan any
	resultChan chan<- Result
	fn         Task
	ctx        context.Context
}

func (w *Worker) Start() {
	go func() {
		for in := range w.taskQueue {
			output, err := w.fn(w.ctx, in)
			w.resultChan <- Result{workerID: w.id, Output: output, Err: err}
		}
	}()
}

type WorkerPool struct {
	ctx         context.Context
	taskQueue   chan any
	resultChan  chan Result
	workerCount int
}

type Result struct {
	workerID int
	Output   any
	Err      error
}

func NewWorkerPool(ctx context.Context, workerCount int) *WorkerPool {
	return &WorkerPool{
		taskQueue:   make(chan any),
		resultChan:  make(chan Result),
		workerCount: workerCount,
		ctx:         ctx,
	}
}

func (wp *WorkerPool) Start(fn Task) {
	for i := 0; i < wp.workerCount; i++ {
		worker := Worker{id: i, taskQueue: wp.taskQueue, resultChan: wp.resultChan, fn: fn, ctx: wp.ctx}
		worker.Start()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
	close(wp.resultChan)
}

func (wp *WorkerPool) Submit(input any) {
	wp.taskQueue <- input
}

func (wp *WorkerPool) GetResult() Result {
	return <-wp.resultChan
}
