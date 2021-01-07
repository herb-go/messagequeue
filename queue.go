package messagequeue

type Queue interface {
	Subscribe(MessageHandler) (Unsubscriber, error)
	NewPublisher() (Publisher, error)
}
