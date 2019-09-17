package messagequeue

import (
	"bytes"
	"container/list"
	"testing"
	"time"
)

func newTestBroker() *Broker {
	b := NewBroker()
	c := NewOptionConfigMap()
	c.Driver = "chan"
	c.ApplyTo(b)
	return b
}
func testrecover() {

}
func TestBroker(t *testing.T) {
	b := newTestBroker()
	err := b.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := b.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	b.SetRecover(testrecover)
	testchan := make(chan *Message, 100)
	b.SetConsumer(NewChanConsumer(testchan))
	messages := [][]byte{}
	unreceived := list.New()
	for i := byte(0); i < 5; i++ {
		messages = append(messages, []byte{i})
		unreceived.PushBack([]byte{i})
	}
	sent, err := b.ProduceMessages(messages...)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	for k := range sent {
		if sent[k] == false {
			t.Fatal(k)
		}
	}
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
