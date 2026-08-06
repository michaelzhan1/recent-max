package generate

import (
	"math"
	"math/rand/v2"
)

// Generator represents a random walk generator
type Generator struct {
	value float64 // current value
	mu    float64 // expected annual return (drift)
	sigma float64 // annualized volatility (standard deviation)
}

// NewGenerator returns a new Generator instance with given parameters
func NewGenerator(initialValue, mu, sigma float64) *Generator {
	return &Generator{
		value: initialValue,
		mu:    mu,
		sigma: sigma,
	}
}

// Step performs a single step in geometric Brownian motion and updates the current value
func (g *Generator) Step() float64 {
	dt := 1.0 / 252.0
	exp := (g.mu-0.5*g.sigma*g.sigma)*dt + g.sigma*math.Sqrt(dt)*rand.NormFloat64()
	g.value *= math.Exp(exp)
	return g.value
}
