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

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/generate"
)

func runTCPServer(ctx context.Context, ln *connection.TCPListener) {
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
			var msg generate.Message
			if err := dec.Decode(&msg); err != nil {
				return err
			}
			log.Printf("Value: %.2f, Timestamp: %s\n", msg.Value, msg.Timestamp.Format("2006-01-02 15:04:05"))
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

	// waitgroups for connection handling
	var wg sync.WaitGroup
	wg.Go(func() {
		log.Println("TCP server is running on port 8081.")
		runTCPServer(ctx, ln)
	})

	<-ctx.Done()

	wg.Wait()
	log.Println("Server stopped.")
}
