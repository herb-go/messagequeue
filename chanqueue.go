package messagequeue

//ChanQueue chan queue driver
type ChanQueue struct {
	queue    chan *Message
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
	q.queue = make(chan *Message)
	q.c = make(chan int)
	go func() {
		for {
			select {
			case m := <-q.queue:
				go q.consumer(m)
			case <-q.c:
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
	sent = make([]bool, len(messages))
	for k := range messages {
		q.queue <- NewMessage(messages[k])
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
