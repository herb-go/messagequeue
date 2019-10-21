package messagequeue

import (
	"sync"
)

var chans = map[string]chan *Message{}

var chansLock = sync.Mutex{}

func getChanByName(name string) chan *Message {
	chansLock.Lock()
	defer chansLock.Unlock()
	c, ok := chans[name]
	if ok == false {
		c = make(chan *Message)
		chans[name] = c
		return c
	}
	return c
}

func closeChanByName(name string) {
	chansLock.Lock()
	defer chansLock.Unlock()
	c, ok := chans[name]
	if ok == false {
		return
	}
	close(c)
}

//ChanQueue chan queue driver
type ChanQueue struct {
	name     string
	c        chan int
	consumer func(*Message) ConsumerStatus
	recover  func()
}

//SetRecover set recover
func (q *ChanQueue) SetRecover(r func()) {
	q.recover = r
}

// Start start queue
//Return any error if raised
func (q *ChanQueue) Start() error {
	var queue = getChanByName(q.name)
	q.c = make(chan int)
	go func() {
		for {
			select {
			case m := <-queue:
				go q.consumer(m)
			case <-q.c:
				closeChanByName(q.name)
				return
			}
		}
	}()
	return nil
}

//Close close queue
//Return any error if raised
func (q *ChanQueue) Close() error {
	close(q.c)
	return nil
}

// ProduceMessages produce messages to broke
//Return sent result and any error if raised
func (q *ChanQueue) ProduceMessages(messages ...[]byte) (sent []bool, err error) {
	var queue = getChanByName(q.name)
	sent = make([]bool, len(messages))
	for k := range messages {
		queue <- NewMessage(messages[k])
		sent[k] = true
	}
	return sent, nil
}

//SetConsumer set message consumer
func (q *ChanQueue) SetConsumer(c func(*Message) ConsumerStatus) {
	q.consumer = c
}

//NewChanQueue create new chan queue
func NewChanQueue() *ChanQueue {
	return &ChanQueue{}
}

//ChanQueueFactory chan queue factory
//Create driver with given config and prefix
func ChanQueueFactory(conf Config, prefix string) (Driver, error) {
	return NewChanQueue(), nil
}

func init() {
	Register("chan", ChanQueueFactory)
}
