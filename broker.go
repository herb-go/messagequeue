package messagequeue

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

//ConsumerStatus consumer status type
type ConsumerStatus int

//ConsumerStatusSuccess consumer status success
const ConsumerStatusSuccess = ConsumerStatus(0)

//ConsumerStatusFail consumer status fail
const ConsumerStatusFail = ConsumerStatus(-1)

//Broker message queue broker struct
type Broker struct {
	Driver
}

//SendBytes send bytes to brokcer.
//Return any error if raised
func (b *Broker) SendBytes(bs []byte) error {
	_, err := b.ProduceMessages(bs)
	return err
}

//NewBroker create new message queue broker
func NewBroker() *Broker {
	return &Broker{}
}

//NewChanConsumer create new chan consumer with given message chan
func NewChanConsumer(c chan *Message) func(*Message) ConsumerStatus {
	return func(message *Message) ConsumerStatus {
		go func() {
			c <- message
		}()
		return ConsumerStatusSuccess
	}
}

//Driver message queue driver interface
type Driver interface {
	// Start start queue
	//Return any error if raised
	Start() error
	//Close close queue
	//Return any error if raised
	Close() error
	//SetRecover set recover
	SetRecover(func())
	// ProduceMessages produce messages to broke
	//Return sent result and any error if raised
	ProduceMessages(...[]byte) (sent []bool, err error)
	//SetConsumer set message consumer
	SetConsumer(func(*Message) ConsumerStatus)
}

// Factory unique id generator driver create factory.
type Factory func(conf Config, prefix string) (Driver, error)

var (
	factorysMu sync.RWMutex
	factories  = make(map[string]Factory)
)

// Register makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, f Factory) {
	factorysMu.Lock()
	defer factorysMu.Unlock()
	if f == nil {
		panic(errors.New("messagequeue: Register messagequeue factory is nil"))
	}
	if _, dup := factories[name]; dup {
		panic(errors.New("messagequeue: Register called twice for factory " + name))
	}
	factories[name] = f
}

//UnregisterAll unregister all driver
func UnregisterAll() {
	factorysMu.Lock()
	defer factorysMu.Unlock()
	// For tests.
	factories = make(map[string]Factory)
}

//Factories returns a sorted list of the names of the registered factories.
func Factories() []string {
	factorysMu.RLock()
	defer factorysMu.RUnlock()
	var list []string
	for name := range factories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewDriver create new driver with given name,config and prefix.
//Reutrn driver created and any error if raised.
func NewDriver(name string, conf Config, prefix string) (Driver, error) {
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("messagequeue: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
