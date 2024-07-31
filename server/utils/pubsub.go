package utils

type IPubSub[EventType any] interface {
	Subscribe() *chan EventType
	Emit(event EventType)
}

type PubSub[EventType any] struct {
	subscribers []chan EventType
}

func NewPubSub[EventType any]() *PubSub[EventType] {
	return &PubSub[EventType]{
		subscribers: make([]chan EventType, 0),
	}
}

func (p *PubSub[EventType]) Subscribe() *chan EventType {
	subscriber := make(chan EventType)
	p.subscribers = append(p.subscribers, subscriber)
	return &subscriber
}

func (p *PubSub[EventType]) Emit(event EventType) {
	for _, subscriber := range p.subscribers {
		go func(sub chan EventType) {
			sub <- event
		}(subscriber)
	}
}
