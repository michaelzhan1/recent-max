package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/michaelzhan1/recent-max/internal/handler"
	"github.com/michaelzhan1/recent-max/internal/middleware"
	"github.com/michaelzhan1/recent-max/internal/value"
	"github.com/michaelzhan1/recent-max/internal/value/generate"
)

// runHTTPServer handles incoming messages over HTTP
func runHTTPServer(ctx context.Context, dataChan chan value.Message, pauseChan chan bool) {
	mux := http.NewServeMux()

	mux.HandleFunc("/stream/data", handler.DataStreamHandlerFactory(dataChan))
	mux.HandleFunc("/pause", handler.PauseHandlerFactory(pauseChan))
	mux.HandleFunc("/resume", handler.ResumeHandlerFactory(pauseChan))
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

func runGenerator(ctx context.Context, gen *generate.Generator, dataChan chan value.Message, pauseChan chan bool) {
	ticker := time.NewTicker(100 * time.Millisecond) // generate data every 100ms
	defer ticker.Stop()

	paused := false

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping generator...")
			return
		case p := <-pauseChan:
			paused = p
			log.Printf("Generator pause state: %v\n", paused)

		case <-ticker.C:
			if paused {
				continue
			}

			newValue := gen.Step()
			msg := value.Message{
				Timestamp: time.Now(),
				Value:     newValue,
			}

			select {
			case <-ctx.Done():
				log.Println("Stopping generator...")
				return
			case dataChan <- msg:
			}
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// graceful shutdown
	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
	}()

	// generator
	gen := generate.NewGenerator(100.0, 0.05, 0.5) // initial value, mu, sigma

	// data channel logic
	dataChan := make(chan value.Message)
	pauseChan := make(chan bool)

	// waitgroups for servers
	var wg sync.WaitGroup
	wg.Go(func() {
		log.Println("Generator is running.")
		runGenerator(ctx, gen, dataChan, pauseChan)
	})
	wg.Go(func() {
		log.Println("HTTP server is running on port 8080.")
		runHTTPServer(ctx, dataChan, pauseChan)
	})

	wg.Wait()
	log.Println("Server stopped.")
}
