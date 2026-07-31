package main

import (
	"log"
	"net"

	"github.com/michaelzhan1/recent-max/internal/connection"
)

func main() {
	// open TCP listener on 8081 for the generator to send data in
	ln, err := connection.NewTCPListener("8081")
	if err != nil {
		log.Fatal(err)
	}

	for {
		err := ln.AcceptAndHandle(func(conn net.Conn) error {
			// handle the connection here
			// for example, read data from conn and process it
			// or send data to conn
			return nil
		})
		if err != nil {
			log.Println("Error handling connection:", err)
		}
	}
}
