package main

import (
	"fmt"
	"log"

	"github.com/michaelzhan1/recent-max/internal/connection"
	"github.com/michaelzhan1/recent-max/internal/generate"
	"github.com/michaelzhan1/recent-max/internal/message"
)

func main() {
	gen := generate.NewGenerator(10.0, 0.5, 1.0)
	for i := range 10 {
		newValue := gen.Step()
		fmt.Printf("Step %d: New Value = %.2f\n", i+1, newValue)
	}

	dialer, err := connection.NewTCPDialer("server:8081")
	if err != nil {
		log.Fatal(err)
	}
	defer dialer.Close()

	msg := message.Message{
		Name:    "Generator",
		Message: fmt.Sprintf("New Value = %.2f", gen.Step()),
	}
	dialer.Send(msg)
}
