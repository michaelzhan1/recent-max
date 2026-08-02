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

func NewValueDeque(maxAge time.Duration) *ValueDeque {
	return &ValueDeque{
		values: []value.Message{},
		maxAge: maxAge,
		now:    time.Now,
	}
}

func (s *ValueDeque) Len() int {
	return len(s.values)
}

func (s *ValueDeque) Empty() bool {
	return len(s.values) == 0
}

func (s *ValueDeque) Push(v value.Message) {
	for !s.Empty() && s.values[len(s.values)-1].Value < v.Value {
		s.values = s.values[:len(s.values)-1]
	}
	s.values = append(s.values, v)
}

func (s *ValueDeque) Peek() (value.Message, bool) {
	for !s.Empty() && s.now().Sub(s.values[0].Timestamp) > s.maxAge {
		s.values = s.values[1:]
	}

	if s.Empty() {
		return value.Message{}, false
	}

	return s.values[0], true
}
