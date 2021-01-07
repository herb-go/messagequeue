package messagequeue

//Message message struct
type Message struct {
	// Data message data
	Data []byte
	//ID message id
	ID string
}

//WithID set id to message
//Return message self
func (m *Message) WithID(id string) *Message {
	m.ID = id
	return m
}

//WithData set data to message
//Return message self
func (m *Message) WithData(data []byte) *Message {
	m.Data = data
	return m
}

//NewMessage create message
func NewMessage() *Message {
	return &Message{}
}

func HandleMesage(h MessageHandler, m *Message) {
	defer func() {
		r := recover()
		if r != nil {
			err := r.(error)
			go func() {
				h.OnError(err)
			}()
		}
	}()
	h.MustHandleMessage(m)
}

type MessageHandler interface {
	MustHandleMessage(*Message) MessageStatus
	OnError(err error)
}

type PlainMessageHandler struct {
	Handler      func(*Message) MessageStatus
	ErrorHandler func(err error)
}

func (h *PlainMessageHandler) WithHandler(f func(*Message) MessageStatus) *PlainMessageHandler {
	h.Handler = f
	return h
}
func (h *PlainMessageHandler) WithErrorHandler(f func(err error)) *PlainMessageHandler {
	h.ErrorHandler = f
	return h
}
func (h *PlainMessageHandler) MustHandleMessage(m *Message) MessageStatus {
	return h.Handler(m)
}
func (h *PlainMessageHandler) OnError(err error) {
	h.ErrorHandler(err)
}

var NopErrorHander = func(err error) {
	panic(err)
}

func NewMessageHandler() *PlainMessageHandler {
	return &PlainMessageHandler{
		ErrorHandler: NopErrorHander,
	}
}

//MessageStatus message status type
type MessageStatus int

//MessageStatusSuccess message status success
const MessageStatusSuccess = MessageStatus(0)

//MessageStatusFail message status fail
const MessageStatusFail = MessageStatus(-1)
