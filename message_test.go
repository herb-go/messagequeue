package messagequeue

import (
	"errors"
	"testing"
)

func TestNessage(t *testing.T) {
	m := NewMessage()
	if m.ID != "" || len(m.Data) != 0 {
		t.Fatal(m)
	}
	n := m.WithID("test")
	if n != m {
		t.Fatal(n, m)
	}
	if m.ID != "test" || len(m.Data) != 0 {
		t.Fatal(m)
	}
	n = m.WithData([]byte("testdata"))
	if n != m {
		t.Fatal(n, m)
	}
	if m.ID != "test" || string(m.Data) != "testdata" {
		t.Fatal(m)
	}
}

func TestNopErrorHandler(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal()
		}
	}()
	h := NewMessageHandler()
	h.OnError(errors.New("nop error"))
}
