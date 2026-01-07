// Package taskqueue provides a simple, generic task queue implementation
// with support for multiple workers and graceful shutdown.
package taskqueue

import (
	"sync"
	"time"
)

// HandlerFunc is a generic function type that processes items of type T.
type HandlerFunc[T any] func(item T)

// Queue is the main structure that manages the task channel and workers.
type Queue[T any] struct {
	dataChan  chan T
	wg        sync.WaitGroup
	handler   HandlerFunc[T]
	workerNum int
	interval  time.Duration
	isClosed  bool
	mu        sync.Mutex
}

// New creates a new Queue instance.
// bufferSize: the capacity of the internal channel.
// workerNum: the number of concurrent workers to process tasks.
// interval: the delay between processing each task (0 for no delay).
// handler: the function that will be called for each item.
func New[T any](bufferSize int, workerNum int, interval time.Duration, handler HandlerFunc[T]) *Queue[T] {
	return &Queue[T]{
		dataChan:  make(chan T, bufferSize),
		workerNum: workerNum,
		interval:  interval,
		handler:   handler,
	}
}

// Start begins the worker loops in the background.
func (q *Queue[T]) Start() {
	for i := 0; i < q.workerNum; i++ {
		q.wg.Add(1)
		go q.workerLoop()
	}
}

// Enqueue adds an item to the queue.
// Returns true if the item was successfully added, false if the queue is closed.
func (q *Queue[T]) Enqueue(item T) bool {
	q.mu.Lock()
	if q.isClosed {
		q.mu.Unlock()
		return false
	}
	q.mu.Unlock()

	q.dataChan <- item
	return true
}

// Stop closes the queue and waits for all workers to finish processing
// the remaining items in the buffer (graceful shutdown).
func (q *Queue[T]) Stop() {
	q.mu.Lock()
	if !q.isClosed {
		close(q.dataChan)
		q.isClosed = true
	}
	q.mu.Unlock()
	q.wg.Wait()
}

func (q *Queue[T]) workerLoop() {
	defer q.wg.Done()
	for item := range q.dataChan {
		q.handler(item)
		if q.interval > 0 {
			time.Sleep(q.interval)
		}
	}
}
