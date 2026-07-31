package main

import (
	"fmt"

	"github.com/michaelzhan1/recent-max/internal/generate"
)

func main() {
	gen := generate.NewGenerator(10.0, 0.5, 1.0)
	for i := range 10 {
		newValue := gen.Step()
		fmt.Printf("Step %d: New Value = %.2f\n", i+1, newValue)
	}
}
