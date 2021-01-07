package messagequeuetestutil

import (
	"errors"
	"sync"
	"time"

	"github.com/herb-go/messagequeue"
)

type TestOpt struct {
}
type TestContext struct {
	Errors []error
	Msgs   [][]byte
}

func NewTestContext() *TestContext {
	return &TestContext{}
}
func TestBroker(d messagequeue.Driver, count int, topic string, ctx *TestContext, ttl time.Duration, opt *TestOpt) {
	var err error
	var locker sync.Mutex
	var failed bool
	var errhandler = func(err error) {
		locker.Lock()
		defer locker.Unlock()
		ctx.Errors = append(ctx.Errors, err)
	}

	p, err := d.NewTopicPublisher(topic)
	if err != nil {
		panic(err)
	}
	mh := func(m *messagequeue.Message) messagequeue.MessageStatus {
		locker.Lock()
		defer locker.Unlock()
		if m.Data[0] == 1 && !failed {
			failed = true
			return messagequeue.MessageStatusFail
		}
		if m.Data[0] == 2 {
			panic(errors.New("test error"))
		}
		ctx.Msgs = append(ctx.Msgs, m.Data)
		return messagequeue.MessageStatusSuccess
	}
	h := messagequeue.NewMessageHandler().
		WithHandler(mh).
		WithErrorHandler(errhandler)
	go func() {
		err := p.Publish([]byte{0})
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(200 * ttl)
	us, err := d.SubscribeTopic(topic, h)
	if err != nil {
		panic(err)
	}

	for i := 1; i < count; i++ {
		data := []byte{byte(i)}
		go func() {
			err := p.Publish(data)
			if err != nil {
				panic(err)
			}
		}()
	}
	time.Sleep(time.Second)

	err = p.Close()
	if err != nil {
		panic(err)
	}
	err = us.Unsubscribe()
	if err != nil {
		panic(err)
	}
	err = d.Close()
	if err != nil {
		panic(err)
	}
}
