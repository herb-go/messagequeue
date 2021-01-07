package messagequeue

import (
	"errors"
	"sync"
	"time"
)

func testQueue(q Queue, ctx *testContext, ttl time.Duration) {
	var err error
	var locker sync.Mutex
	var failed bool
	var errhandler = func(err error) {
		locker.Lock()
		defer locker.Unlock()
		ctx.Errors = append(ctx.Errors, err)
	}

	p, err := q.NewPublisher()
	if err != nil {
		panic(err)
	}
	mh := func(m *Message) MessageStatus {
		locker.Lock()
		defer locker.Unlock()
		if m.Data[0] == 1 && !failed {
			failed = true
			return MessageStatusFail
		}
		if m.Data[0] == 2 {
			panic(errors.New("test error"))
		}
		ctx.Msgs = append(ctx.Msgs, m.Data)
		return MessageStatusSuccess
	}
	h := NewMessageHandler().
		WithHandler(mh).
		WithErrorHandler(errhandler)
	go func() {
		err := p.Publish([]byte{0})
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(200 * ttl)
	us, err := q.Subscribe(h)
	if err != nil {
		panic(err)
	}

	for i := 1; i < 5; i++ {
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
}
