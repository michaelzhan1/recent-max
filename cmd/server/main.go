package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/value"
	"github.com/michaelzhan1/recent-max/internal/value/deque"
)

// runTCPServer handles handles incoming messages over TCP
func runTCPServer(ctx context.Context, ln *connection.TCPListener, dq *deque.ValueDeque) {
	conn, err := ln.Accept()
	if err != nil {
		log.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close() // conn.Close will help dec.Decode close properly
	}()

	err = ln.Handle(conn, func(dec *json.Decoder) error {
		for {
			var msg value.Message
			if err := dec.Decode(&msg); err != nil {
				return err
			}
			dq.Push(msg)
		}
	})
	if err != nil {
		if err == io.EOF {
			log.Println("Connection closed by client.")
			return
		}
		log.Println("Error handling connection:", err)
	}
}

func logMaxValue(ctx context.Context, dq *deque.ValueDeque) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			maxValue, ok := dq.Peek()
			if !ok {
				log.Println("No values in deque.")
				continue
			}
			log.Printf("Current max value in the last 3 seconds: %.2f\n", maxValue.Value)
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// open TCP listener on 8081 for the generator to send data in
	ln, err := connection.NewTCPListener("8081")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// graceful shutdown
	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		ln.Close() // break any waiting Accept()
	}()

	// deque logic
	dq := deque.NewValueDeque(3 * time.Second)

	// waitgroups for connection handling and logging
	var wg sync.WaitGroup
	wg.Go(func() {
		log.Println("TCP server is running on port 8081.")
		runTCPServer(ctx, ln, dq)
	})
	wg.Go(func() {
		log.Println("Logging max value every 3 seconds.")
		logMaxValue(ctx, dq)
	})

	<-ctx.Done()

	wg.Wait()
	log.Println("Server stopped.")
}
