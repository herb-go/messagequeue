package messagequeue_test

import (
	"testing"
	"time"

	"github.com/herb-go/messagequeue"
	"github.com/herb-go/messagequeue/messagequeuetestutil"
)

var ttl = 100 * time.Microsecond

func newTestBroker() messagequeue.Driver {
	b := messagequeue.NewChanBroker()
	return b
}
func newTTLTestBroker() messagequeue.Driver {
	b := messagequeue.NewChanBroker()
	b.SetTTL(ttl)
	return b
}

func TestChanBroker(t *testing.T) {
	b := newTestBroker()
	ctx := messagequeuetestutil.NewTestContext()
	messagequeuetestutil.TestBroker(b, 5, "test", ctx, ttl, nil)
	if len(ctx.Msgs) != 3 || len(ctx.Errors) != 1 {
		t.Fatal(ctx.Msgs, ctx.Errors)
	}
}
func TestTTLChanBroker(t *testing.T) {
	b := newTTLTestBroker()
	ctx := messagequeuetestutil.NewTestContext()
	messagequeuetestutil.TestBroker(b, 5, "test", ctx, ttl, nil)
	if len(ctx.Msgs) != 2 || len(ctx.Errors) != 1 {
		t.Fatal(ctx.Msgs, ctx.Errors)
	}
}

func TestTTL(t *testing.T) {
	b := messagequeue.NewChanBroker()
	if b.TTL() != messagequeue.DefaultChanBrokerTTL {
		t.Fatal(b)
	}
	b.SetTTL(8 * time.Second)
	if b.TTL() != 8*time.Second {
		t.Fatal(b)
	}
}
