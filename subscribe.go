package messagequeue

type Unsubscriber interface {
	Unsubscribe() error
}

type FuncUnsubscriber func() error

func (u FuncUnsubscriber) Unsubscribe() error {
	return u()
}
