package messagequeue

import (
	"bytes"
	"container/list"
	"testing"
	"time"
)

func newTestBroker() *Broker {
	b := NewBroker()
	c := NewOptionConfig()
	c.Driver = "chan"
	err := c.ApplyTo(b)
	if err != nil {
		panic(err)
	}
	return b
}
func testrecover() {

}
func TestBroker(t *testing.T) {
	b := newTestBroker()
	producer := Producer(b)
	consumer := Consumer(b)
	err := consumer.Listen()
	if err != nil {
		t.Fatal(err)
	}
	err = producer.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := consumer.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = producer.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()
	consumer.SetRecover(testrecover)
	testchan := make(chan *Message, 100)
	consumer.SetConsumer(NewChanConsumer(testchan))
	unreceived := list.New()
	for i := byte(0); i < 5; i++ {
		err = producer.ProduceMessage([]byte{i})
		if err != nil {
			t.Fatal(err)
		}
		unreceived.PushBack([]byte{i})
	}

	time.Sleep(time.Second)

	if len(testchan) != 5 {
		t.Fatal(len(testchan))
	}
	if unreceived.Len() != 5 {
		t.Fatal(unreceived.Len())
	}
	for i := byte(0); i < 5; i++ {
		m := <-testchan
		e := unreceived.Front()
		for e != nil {
			if bytes.Compare(e.Value.([]byte), m.Data) == 0 {
				unreceived.Remove(e)
				break
			}
			e = e.Next()
		}
	}
	if unreceived.Len() != 0 {
		t.Fatal(unreceived)
	}
}
