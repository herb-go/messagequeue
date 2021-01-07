package messagequeue

import (
	"sync"
	"time"
)

type ChanBroker struct {
	locker sync.Mutex
	topics map[string]chan ([]byte)
}

func (b *ChanBroker) getChanByName(name string) chan ([]byte) {
	b.locker.Lock()
	defer b.locker.Unlock()
	c, ok := b.topics[name]
	if !ok {
		c = make(chan ([]byte))
		b.topics[name] = c
	}
	return c
}

func NewChanBroker() *ChanBroker {
	return &ChanBroker{
		topics: map[string]chan ([]byte){},
	}
}

var DefaultChanBroker = NewChanBroker()

type ChanQueue struct {
	broker *ChanBroker
	Name   string
	TTL    time.Duration
}

func NewChanQueue() *ChanQueue {
	return &ChanQueue{
		broker: DefaultChanBroker,
	}
}
func (q *ChanQueue) Publish(bs []byte) error {
	go func() {
		select {
		case q.broker.getChanByName(q.Name) <- bs:
			return
		case <-time.After(q.TTL):
			return
		}
	}()
	return nil
}
func (q *ChanQueue) Close() error {
	return nil
}

func (q *ChanQueue) Subscribe(h MessageHandler) (Unsubscriber, error) {
	quitchan := make(chan int)
	u := FuncUnsubscriber(func() error {
		close(quitchan)
		return nil
	})
	go func() {
		for {
			select {
			case bs := <-q.broker.getChanByName(q.Name):
				HandleMesage(h, NewMessage().WithData(bs))
			case <-quitchan:
				return
			}
		}
	}()
	return u, nil
}
func (q *ChanQueue) NewPublisher() (Publisher, error) {
	return q, nil
}
