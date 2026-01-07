package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Import path sesuai nama module di go.mod + path foldernya
	"github.com/elfadl/simple-queue/pkg/taskqueue"
)

func main() {
	// Definisikan handler: apa yang dilakukan dengan data?
	myHandler := func(msg string) {
		log.Printf("Processing: %s", msg)
	}

	// Buat Queue: Buffer 10, 2 Worker, Delay 1 detik per job
	q := taskqueue.New(10, 2, 1*time.Second, myHandler)
	q.Start()

	fmt.Println("System running... (Ctrl+C to stop)")

	// Simulasi kirim data
	go func() {
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("Message #%d", i)
			if q.Enqueue(msg) {
				fmt.Printf(" -> Enqueued: %s\n", msg)
			}
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Tunggu sinyal Stop
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	fmt.Println("\nStopping system...")
	q.Stop()
	fmt.Println("System stopped.")
}
