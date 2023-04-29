package framework

type EventHandlerInterface interface {
	GetExpectedMessageKey() string
	GetExpectedEventType() any
	Handle(event any) error
	Commit() error
}
