package value

import "time"

type Message struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
