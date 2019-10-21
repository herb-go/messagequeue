package messagequeue

import (
	"testing"
	"time"
)

func TestChans(t *testing.T) {
	var err error
	time.Sleep(time.Millisecond)
	if len(chans) != 0 {
		t.Fatal(chans)
	}
	conf := NewOptionConfigMap()
	conf.Config.Set("Name", "name1")
	conf.Driver = "chan"
	conf2 := NewOptionConfigMap()
	conf2.Config.Set("Name", "name2")
	conf2.Driver = "chan"
	c := NewBroker()
	err = conf.ApplyTo(c)
	if err != nil {
		t.Fatal(err)
	}
	c2 := NewBroker()
	err = conf2.ApplyTo(c2)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Listen()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if len(chans) != 1 {
		t.Fatal(chans)
	}
	err = c2.Listen()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if len(chans) != 2 {
		t.Fatal(chans)
	}
	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if len(chans) != 1 {
		t.Fatal(chans)
	}
	err = c2.Close()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if len(chans) != 0 {
		t.Fatal(chans)
	}
}
