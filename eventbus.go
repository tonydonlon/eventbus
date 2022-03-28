package eventbus

// EventBus is a simple pub sub event bus
type EventBus struct {
	Name          string
	Messages      chan []byte
	Subscribe     chan chan []byte
	Unsubscribe   chan chan []byte
	subscriptions map[chan []byte]bool
	halt          <-chan bool
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

// ConnectionCount returns the number of active client connections
func (b *EventBus) ConnectionCount() int {
	return len(b.subscriptions)
}

func (b *EventBus) listen() {
	for {
		select {
		case chnl := <-b.Subscribe:
			b.subscriptions[chnl] = true
		case chnl := <-b.Unsubscribe:
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
