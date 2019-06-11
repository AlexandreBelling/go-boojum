package p2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-pubsub"
)

// A Topic is a pubsub abstraction that can be subscribed and published
type Topic struct {
	ps         *pubsub.PubSub
	ctx        context.Context
	Subs       *pubsub.Subscription
	Name       string
	cancelChan context.CancelFunc
	onceChan   sync.Once
}

// Publish a message in the topic
func (t *Topic) Publish(msg []byte) error {
	return t.ps.Publish(t.Subs.Topic(), msg)
}

// Chan get a chan for the message
// Should be called only once
func (t *Topic) Chan() (<-chan []byte, error) {
	var out chan []byte
	var err error

	t.onceChan.Do(func() {

		subs, err := t.ps.Subscribe(t.Name)
		if err != nil {
			return
		}

		t.Subs = subs
		out = make(chan []byte, 20)
		go t.background(out)
	})

	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, fmt.Errorf("Topic.Chan can be called only once")
	}

	return out, nil
}

// Close the subscription
func (t *Topic) Close() {
	t.cancelChan()
	t.Subs.Cancel()
}

func (t *Topic) background(out chan []byte) {
	defer close(out)
	ctx, can := context.WithCancel(context.Background())
	t.cancelChan = can

	for {
		fromSubscription := make(chan *pubsub.Message)
		errorChan := make(chan error)

		go func() {
			defer close(errorChan)
			defer close(fromSubscription)

			msg, err := t.Subs.Next(ctx)
			if err != nil {
				errorChan <- err
				return
			}

			fromSubscription <- msg
		}()

		select {

		case <-t.ctx.Done():
			return

		case <-errorChan:
			// It is impossible that this is triggered by a cancellation.
			// This is truly a pubsub error
			t.cancelChan() // Destroy the context to avoid leaks
			return         // Make the rest of the app, aware that there is a problem

		case msg := <-fromSubscription:
			out <- msg.GetData()
		}
	}
}
