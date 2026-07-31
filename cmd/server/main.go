package main

import (
	"encoding/json"
	"log"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/message"
)

func main() {
	// open TCP listener on 8081 for the generator to send data in
	ln, err := connection.NewTCPListener("8081")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

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
			log.Println("Error handling connection:", err)
		}
	}
}
