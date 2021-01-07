package messagequeue

import (
	"sync"
	"time"
)

var DefaultChanBrokerTTL = 24 * time.Hour

type ChanBroker struct {
	locker sync.Mutex
	topics map[string]chan ([]byte)
	ttl    time.Duration
}

func (b *ChanBroker) unsafeGetChanByName(name string) chan ([]byte) {
	c, ok := b.topics[name]
	if !ok {
		c = make(chan ([]byte))
		b.topics[name] = c
	}
	return c
}
func (b *ChanBroker) SetTTL(d time.Duration) {
	b.locker.Lock()
	defer b.locker.Unlock()
	b.ttl = d
}
func (b *ChanBroker) TTL() time.Duration {
	b.locker.Lock()
	defer b.locker.Unlock()
	return b.ttl
}
func (b *ChanBroker) NewTopicPublisher(topic string) (Publisher, error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	p := &ChanPublisher{
		c:   b.unsafeGetChanByName(topic),
		ttl: b.ttl,
	}
	return p, nil
}
func (b *ChanBroker) SubscribeTopic(topic string, h MessageHandler) (Subscription, error) {
	b.locker.Lock()
	defer b.locker.Unlock()
	quitchan := make(chan int)
	u := FuncSubscription(func() error {
		close(quitchan)
		return nil
	})
	c := b.unsafeGetChanByName(topic)
	go func() {
		for {
			select {
			case bs := <-c:
				HandleMesage(h, NewMessage().WithData(bs))
			case <-quitchan:
				return
			}
		}
	}()
	return u, nil
}

func (b *ChanBroker) Close() error {
	b.locker.Lock()
	defer b.locker.Unlock()
	b.topics = map[string]chan ([]byte){}
	return nil
}

func NewChanBroker() *ChanBroker {
	return &ChanBroker{
		topics: map[string]chan ([]byte){},
		ttl:    DefaultChanBrokerTTL,
	}
}

type ChanPublisher struct {
	c   chan []byte
	ttl time.Duration
}

func (p *ChanPublisher) Publish(bs []byte) error {
	go func() {
		select {
		case p.c <- bs:
			return
		case <-time.After(p.ttl):
			return
		}
	}()
	return nil
}
func (p *ChanPublisher) Close() error {
	return nil
}

type ChanBrokerConfig struct {
	TTLDuration string
}

//ChanBrokerFactory chan broker factory
//Create driver with given loader
func ChanBrokerFactory(loader func(interface{}) error) (Driver, error) {
	c := NewChanBroker()
	conf := &ChanBrokerConfig{}
	err := loader(conf)
	if err != nil {
		return nil, err
	}
	if conf.TTLDuration != "" {
		d, err := time.ParseDuration(conf.TTLDuration)
		if err != nil {
			return nil, err
		}
		c.SetTTL(d)
	}
	return c, nil
}

func init() {
	Register("chan", ChanBrokerFactory)
}
