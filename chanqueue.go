package messagequeue

type ChanQueue struct {
	queue    chan *Message
	c        chan int
	consumer func(*Message) ConsumerStatus
	recover  func()
}

func (q *ChanQueue) SetRecover(r func()) {
	q.recover = r
}
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
func (q *ChanQueue) Close() error {
	close(q.c)
	return nil
}
func (q *ChanQueue) ProduceMessages(messages ...[]byte) (sent []bool, err error) {
	sent = make([]bool, len(messages))
	for k := range messages {
		q.queue <- NewMessage(messages[k])
		sent[k] = true
	}
	return sent, nil
}
func (q *ChanQueue) SetConsumer(c func(*Message) ConsumerStatus) {
	q.consumer = c
}

func NewChanQueue() *ChanQueue {
	return &ChanQueue{}
}

func ChanQueueFactory(conf Config, prefix string) (Driver, error) {
	return NewChanQueue(), nil
}

func init() {
	Register("chan", ChanQueueFactory)
}
