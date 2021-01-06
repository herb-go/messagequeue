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
	delete(chans, name)
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

//Connect to brocker as producer
func (q *ChanQueue) Connect() error {
	return nil
}

//Disconnect stop producing and disconnect
func (q *ChanQueue) Disconnect() error {
	return nil
}

// Listen listen queue
//Return any error if raised
func (q *ChanQueue) Listen() error {
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

// ProduceMessage produce messages to broke
//Return any error if raised
func (q *ChanQueue) ProduceMessage(message []byte) error {
	var queue = getChanByName(q.name)
	queue <- NewMessage(message)
	return nil
}

//SetConsumer set message consumer
func (q *ChanQueue) SetConsumer(c func(*Message) ConsumerStatus) {
	q.consumer = c
}

//NewChanQueue create new chan queue
func NewChanQueue() *ChanQueue {
	return &ChanQueue{}
}

type ChanQueueConfig struct {
	Name string
}

//ChanQueueFactory chan queue factory
//Create driver with given loader
func ChanQueueFactory(loader func(interface{}) error) (Driver, error) {
	c := NewChanQueue()
	conf := &ChanQueueConfig{}
	err := loader(conf)
	if err != nil {
		return nil, err
	}
	c.name = conf.Name
	return c, nil
}

func init() {
	Register("chan", ChanQueueFactory)
}
