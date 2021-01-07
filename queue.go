package messagequeue

type Queue interface {
	PushChan() chan []byte
	PopChan() chan []byte
	ShouldClose()
}
