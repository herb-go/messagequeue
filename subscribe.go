package messagequeue

type Subscription interface {
	Unsubscribe() error
}

type FuncSubscription func() error

func (u FuncSubscription) Unsubscribe() error {
	return u()
}
