package messagequeue

type Publisher interface {
	Publish([]byte) error
	Close() error
}
