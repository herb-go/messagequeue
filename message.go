package messagequeue

//Message message struct
type Message struct {
	// Data message data
	Data []byte
	//ID message id
	ID string
}

//SetID set id to message
//Return message self
func (m *Message) SetID(id string) *Message {
	m.ID = id
	return m
}

//NewMessage create message with given data
func NewMessage(data []byte) *Message {
	return &Message{
		Data: data,
	}
}
