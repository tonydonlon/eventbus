package eventbus_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tonydonlon/eventbus"
)

func TestEventBus(t *testing.T) {
	assert := assert.New(t)
	expectedMessages := []string{"one", "two", "three", "four", "five"}

	b := eventbus.New("test")
	done := make(chan bool, 1)
	subscriber1 := make(chan []byte)
	b.Start(done)

	b.Subscribe <- subscriber1
	assert.Equal("test", b.Name, "bus has a name")
	assert.Equal(1, b.ConnectionCount(), "should have one subscriber")

	count := 0
	messagesReceived := []string{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case c := <-subscriber1:
				count++
				messagesReceived = append(messagesReceived, string(c))
				if count == len(expectedMessages) {
					b.Unsubscribe <- subscriber1
					done <- true
					wg.Done()
					break
				}
			}
		}
	}()

	setupTimerProducer(b, expectedMessages, time.Millisecond)
	wg.Wait()

	assert.Equal(len(expectedMessages), len(messagesReceived), "received correct number of messages")

	for i := range messagesReceived {
		assert.Equal(expectedMessages[i], messagesReceived[i], "messages match")
	}
	// TODO fix race condition
	// assert.Equal(t, 0, b.ConnectionCount(), "should have no subscribers")
}

func setupTimerProducer(bus *eventbus.EventBus, expectedMessages []string, dur time.Duration) {
	go func() {
		for _, v := range expectedMessages {
			bus.Messages <- []byte(v)
			time.Sleep(dur)
		}
	}()
}
