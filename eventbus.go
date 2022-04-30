package eventbus

import "sync/atomic"

// EventBus is a simple pub sub event bus
type EventBus struct {
	Name            string
	Messages        chan []byte
	Subscribe       chan chan []byte
	Unsubscribe     chan chan []byte
	subscriptions   map[chan []byte]bool
	halt            <-chan bool
	connectionCount int64
}

// New creates a new EventBus
func New(name string) *EventBus {
	bus := &EventBus{
		Messages:      make(chan []byte, 1),
		Subscribe:     make(chan chan []byte),
		Unsubscribe:   make(chan chan []byte),
		subscriptions: make(map[chan []byte]bool),
		Name:          name,
		halt:          make(<-chan bool),
	}
	return bus
}

// Start makes the EventBus start listening
func (b *EventBus) Start(done <-chan bool) {
	b.halt = done
	go b.listen()
}

// SubscriberCount returns the number of active client connections
func (b *EventBus) SubscriberCount() int64 {
	return atomic.LoadInt64(&b.connectionCount)
}

func (b *EventBus) listen() {
	for {
		select {
		case chnl := <-b.Subscribe:
			b.subscriptions[chnl] = true
			atomic.AddInt64(&b.connectionCount, 1)
		case chnl := <-b.Unsubscribe:
			atomic.AddInt64(&b.connectionCount, ^int64(0))
			delete(b.subscriptions, chnl)
		case event := <-b.Messages:
			for subscriberChannel := range b.subscriptions {
				subscriberChannel <- event
			}
		case <-b.halt:
			return
		}
	}
}
