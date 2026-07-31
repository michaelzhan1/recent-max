package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/message"
)

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

	for {
		err := ln.AcceptAndHandle(func(dec *json.Decoder) error {
			var msg message.Message
			if err := dec.Decode(&msg); err != nil {
				return err
			}
			log.Printf("Received message: %+v\n", msg)
			return nil
		})
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}

			log.Println("Error handling connection:", err)
		}
	}

	log.Println("Server stopped.")
}
