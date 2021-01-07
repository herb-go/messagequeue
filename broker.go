package messagequeue

type Driver interface {
	SubscribeTopic(string, MessageHandler) (Unsubscriber, error)
	NewTopicPublisher(string) (Publisher, error)
	Close() error
}
