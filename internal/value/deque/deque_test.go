package deque

import (
	"testing"
	"time"

	"github.com/michaelzhan1/recent-max/internal/value"
)

type testOp string

var (
	PushOp testOp = "push"
	PeekOp testOp = "peek"
)

func TestValueDeque(t *testing.T) {
	staticTime := func() time.Time {
		return time.Unix(0, 0)
	}

	d := &ValueDeque{
		values: []value.Message{},
		now:    staticTime,
	}

	ops := []struct {
		op         testOp
		v          value.Message
		expPeekVal float64
	}{
		{
			op: PushOp,
			v: value.Message{
				Value:     2,
				Timestamp: staticTime(),
			},
		},
		{
			op:         PeekOp,
			expPeekVal: 2,
		},
		{
			op: PushOp,
			v: value.Message{
				Value:     1,
				Timestamp: staticTime(),
			},
		},
		{
			op:         PeekOp,
			expPeekVal: 2,
		},
		{
			op: PushOp,
			v: value.Message{
				Value:     3,
				Timestamp: staticTime(),
			},
		},
		{
			op:         PeekOp,
			expPeekVal: 3,
		},
	}

	for _, op := range ops {
		switch op.op {
		case PushOp:
			d.Push(op.v)
		case PeekOp:
			v, ok := d.Peek()
			if !ok {
				t.Fatalf("expected value, got none")
			}
			if v.Value != op.expPeekVal {
				t.Fatalf("expected %v, got %v", op.expPeekVal, v.Value)
			}
		default:
			t.Fatalf("unknown operation: %v", op.op)
		}
	}
}

func TestValueDequePrune(t *testing.T) {
	expiredTime := time.Unix(0, 0)
	nonExpTime := time.Unix(10, 0)

	d := ValueDeque{
		values: []value.Message{},
		maxAge: 5 * time.Second,
		now: func() time.Time {
			return nonExpTime
		},
	}

	ops := []struct {
		op         testOp
		v          value.Message
		expPeekVal float64
	}{
		{
			op: PushOp,
			v: value.Message{
				Value:     2,
				Timestamp: expiredTime,
			},
		},
		{
			op: PushOp,
			v: value.Message{
				Value:     1,
				Timestamp: nonExpTime,
			},
		},
		{
			op:         PeekOp,
			expPeekVal: 1,
		},
	}

	for _, op := range ops {
		switch op.op {
		case PushOp:
			d.Push(op.v)
		case PeekOp:
			v, ok := d.Peek()
			if !ok {
				t.Fatalf("expected value, got none")
			}
			if v.Value != op.expPeekVal {
				t.Fatalf("expected %v, got %v", op.expPeekVal, v.Value)
			}
		default:
			t.Fatalf("unknown operation: %v", op.op)
		}
	}
}
