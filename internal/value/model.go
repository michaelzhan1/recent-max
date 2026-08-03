package value

import "time"

// Message is the data structure of values sent from the generator to the server
type Message struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
