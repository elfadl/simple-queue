package taskqueue

import (
	"sync"
	"time"
)

// HandlerFunc adalah tipe fungsi generic
type HandlerFunc[T any] func(item T)

// Queue struct utama
type Queue[T any] struct {
	dataChan  chan T
	wg        sync.WaitGroup
	handler   HandlerFunc[T]
	workerNum int
	interval  time.Duration
	isClosed  bool
	mu        sync.Mutex
}

// New membuat instance queue baru
func New[T any](bufferSize int, workerNum int, interval time.Duration, handler HandlerFunc[T]) *Queue[T] {
	return &Queue[T]{
		dataChan:  make(chan T, bufferSize),
		workerNum: workerNum,
		interval:  interval,
		handler:   handler,
	}
}

// Start menjalankan worker
func (q *Queue[T]) Start() {
	for i := 0; i < q.workerNum; i++ {
		q.wg.Add(1)
		go q.workerLoop()
	}
}

// Enqueue memasukkan data
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

// Stop mematikan queue dengan graceful shutdown
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
