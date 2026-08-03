package deque

import (
	"time"

	"github.com/michaelzhan1/recent-max/internal/value"
)

// ValueDeque is a monotonically decreasing deque that stores value.Message objects
type ValueDeque struct {
	values []value.Message
	maxAge time.Duration
	now    func() time.Time
}

// NewValueDeque creates a new ValueDeque with the specified maximum age for values
func NewValueDeque(maxAge time.Duration) *ValueDeque {
	return &ValueDeque{
		values: []value.Message{},
		maxAge: maxAge,
		now:    time.Now,
	}
}

// Len returns the number of elements in the deque
func (s *ValueDeque) Len() int {
	return len(s.values)
}

// Empty checks if the deque is empty
func (s *ValueDeque) Empty() bool {
	return len(s.values) == 0
}

// Push adds a new value.Message to the deque, maintaining the monotonically decreasing order
func (s *ValueDeque) Push(v value.Message) {
	for !s.Empty() && s.values[len(s.values)-1].Value < v.Value {
		s.values = s.values[:len(s.values)-1]
	}
	s.values = append(s.values, v)
}

// Peek returns the oldest value.Message in the deque that is within the maxAge, or false if there are no valid values.
// It prunes any values that are older than maxAge before returning the oldest valid value.
func (s *ValueDeque) Peek() (value.Message, bool) {
	for !s.Empty() && s.now().Sub(s.values[0].Timestamp) > s.maxAge {
		s.values = s.values[1:]
	}

	if s.Empty() {
		return value.Message{}, false
	}

	return s.values[0], true
}
