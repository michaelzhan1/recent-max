package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/handler"
	"github.com/michaelzhan1/recent-max/internal/middleware"
	"github.com/michaelzhan1/recent-max/internal/value"
	"github.com/michaelzhan1/recent-max/internal/value/deque"
)

// runHTTPServer handles incoming messages over HTTP
func runHTTPServer(ctx context.Context, dq *deque.ValueDeque, dataChan chan value.Message) {
	mux := http.NewServeMux()

	mux.HandleFunc("/stream/data", handler.DataStreamHandlerFactory(dataChan))
	mux.HandleFunc("/stream/stats", handler.StatStreamHandlerFactory(dq))
	corsMux := middleware.EnableCORS(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	go func() {
		<-ctx.Done()

		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // shutdown timeout
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Println("Error shutting down HTTP server:", err)
		}
	}()

	err := server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("Error starting HTTP server:", err)
	}
}

// runTCPServer handles handles incoming messages over TCP
func runTCPServer(ctx context.Context, ln *connection.TCPListener, dq *deque.ValueDeque, dataChan chan value.Message) {
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
			dataChan <- msg
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
	dq := deque.NewValueDeque(5 * time.Second)

	// data channel logic
	dataChan := make(chan value.Message)

	// waitgroups for servers
	var wg sync.WaitGroup
	wg.Go(func() {
		log.Println("TCP server is running on port 8081.")
		runTCPServer(ctx, ln, dq, dataChan)
	})
	wg.Go(func() {
		log.Println("HTTP server is running on port 8080.")
		runHTTPServer(ctx, dq, dataChan)
	})

	<-ctx.Done()

	wg.Wait()
	log.Println("Server stopped.")
}
