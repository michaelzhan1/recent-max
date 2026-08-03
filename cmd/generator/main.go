package main

import (
	"log"
	"time"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/value"
	"github.com/michaelzhan1/recent-max/internal/value/generate"
)

// sendValue wraps the generator's current value and timestamp and sends it through the dialer
func sendValue(dialer *connection.TCPDialer, g *generate.Generator) error {
	msg := value.Message{
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

	gen := generate.NewGenerator(100, 0.08, 0.25)
	for {
		sendValue(dialer, gen)
		time.Sleep(100 * time.Millisecond)
		val := gen.Step()
		log.Println("Generated value:", val)
	}
}
