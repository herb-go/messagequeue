package messagequeue

type Message struct {
	Data []byte
	ID   string
}

func (m *Message) SetID(id string) *Message {
	m.ID = id
	return m
}

func NewMessage(data []byte) *Message {
	return &Message{
		Data: data,
	}
}
