package generate

import "math/rand/v2"

// Generator represents a random walk generator
type Generator struct {
	value     float64 // current value
	incChance float64 // 0 to 1, chance to increment value
	incAmt    float64 // amount to increment value by
}

// NewGenerator returns a new Generator instance with given parameters
func NewGenerator(initialValue, incChance, incAmt float64) *Generator {
	return &Generator{
		value:     initialValue,
		incChance: incChance,
		incAmt:    incAmt,
	}
}

// Step performs a single step in the random walk, with a floor at 0.
// It returns the new value after the step.
func (g *Generator) Step() float64 {
	if rand.Float64() < g.incChance {
		g.value += g.incAmt
	} else {
		g.value -= g.incAmt
	}
	if g.value < 0 {
		g.value = 0
	}
	return g.value
}
