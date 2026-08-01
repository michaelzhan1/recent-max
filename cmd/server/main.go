package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/generate"
)

func runTCPServer(ln *connection.TCPListener) {
	err := ln.AcceptAndHandle(func(dec *json.Decoder) error {
		for {
			var msg generate.Message
			if err := dec.Decode(&msg); err != nil {
				return err
			}
			log.Printf("Value: %.2f, Timestamp: %s\n", msg.Value, msg.Timestamp.Format("2006-01-02 15:04:05"))
		}
	})
	if err != nil {
		log.Println("Error handling connection:", err)
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
		ln.Close() // needed to break out of the AcceptAndHandle loop
	}()

	// waitgroups for connection handling
	var wg sync.WaitGroup
	wg.Go(func() {
		log.Println("TCP server is running on port 8081.")
		runTCPServer(ln)
	})

	<-ctx.Done()

	wg.Wait()
	log.Println("Server stopped.")
}
