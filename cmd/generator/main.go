package main

import (
	"log"
	"time"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/generate"
)

// sendValue wraps the generator's current value and timestamp and sends it through the dialer
func sendValue(dialer *connection.TCPDialer, g *generate.Generator) error {
	msg := generate.Message{
		Value:     g.Value(),
		Timestamp: time.Now(),
	}
	return dialer.Send(msg)
}

func main() {
	dialer, err := connection.NewTCPDialer("server:8081")
	if err != nil {
		log.Fatal(err)
	}
	defer dialer.Close()

	gen := generate.NewGenerator(10.0, 0.6, 1.0)
	for {
		sendValue(dialer, gen)
		time.Sleep(1 * time.Second)
		gen.Step()
	}
}
