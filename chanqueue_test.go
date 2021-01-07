package messagequeue

import (
	"testing"
	"time"
)

var ttl = 100 * time.Microsecond

func newTestQueue() Queue {
	q := NewChanQueue()
	q.Name = "test"
	q.TTL = time.Hour
	return q
}
func newTTLTestQueue() Queue {
	q := NewChanQueue()
	q.Name = "test"
	q.TTL = ttl
	return q
}

type testContext struct {
	Errors []error
	Msgs   [][]byte
}

func newTestContext() *testContext {
	return &testContext{}
}

func TestQueue(t *testing.T) {
	q := newTestQueue()
	ctx := newTestContext()
	testQueue(q, ctx, ttl)
	if len(ctx.Msgs) != 3 || len(ctx.Errors) != 1 {
		t.Fatal(ctx.Msgs, ctx.Errors)
	}
}
func TestTTLQueue(t *testing.T) {
	q := newTTLTestQueue()
	ctx := newTestContext()
	testQueue(q, ctx, ttl)
	if len(ctx.Msgs) != 2 || len(ctx.Errors) != 1 {
		t.Fatal(ctx.Msgs, ctx.Errors)
	}
}
