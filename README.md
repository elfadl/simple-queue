# Simple Queue

![Version](https://img.shields.io/badge/version-0.0.2-blue)
[![Go Reference](https://pkg.go.dev/badge/github.com/elfadl/simple-queue.svg)](https://pkg.go.dev/github.com/elfadl/simple-queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/elfadl/simple-queue)](https://goreportcard.com/report/github.com/elfadl/simple-queue)

`simple-queue` is a simple Go library for managing task queues using goroutines and channels. It supports generic types, multiple workers, and graceful shutdown.

## Features

- **Generic Support**: Can be used with any data type.
- **Concurrent Workers**: Run multiple workers simultaneously to process data.
- **Configurable Interval**: Provide a delay between task processing.
- **Graceful Shutdown**: Ensures all queued tasks are processed before the application stops.

## Installation

Use `go get` to add this library to your project:

```bash
go get github.com/elfadl/simple-queue
```

## Usage

Here are the steps to use `simple-queue`:

### 1. Define a Handler
Create a function that will process each item in the queue.

```go
myHandler := func(msg string) {
    fmt.Printf("Processing: %s\n", msg)
}
```

### 2. Initialize the Queue
Use `taskqueue.New` to create a new instance. You need to specify the buffer size, number of workers, interval between tasks, and the handler.

```go
// Buffer: 10, Workers: 2, Interval: 1 second
q := taskqueue.New(10, 2, 1*time.Second, myHandler)
```

### 3. Start the Queue
Call the `Start()` method to begin running workers in the background.

```go
q.Start()
```

### 4. Enqueue Data
Use `Enqueue()` to send data into the queue.

```go
q.Enqueue("Data 1")
q.Enqueue("Data 2")
```

### 5. Stop Gracefully
Call `Stop()` to close the queue and wait for all workers to finish their tasks.

```go
q.Stop()
```

## Full Example

```go
package main

import (
	"fmt"
	"time"
	"github.com/elfadl/simple-queue/pkg/taskqueue"
)

func main() {
	handler := func(item int) {
		fmt.Printf("Processing number: %d\n", item)
		time.Sleep(500 * time.Millisecond)
	}

	// Initialize queue for int data type
	q := taskqueue.New(5, 3, 0, handler)
	q.Start()

	// Send data
	for i := 1; i <= 10; i++ {
		q.Enqueue(i)
	}

	// Stop after completion
	time.Sleep(2 * time.Second)
	q.Stop()
	fmt.Println("Finished.")
}
```

## License

MIT
