
package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-pubsub"
)

// A Topic is a pubsub abstraction that can be subscribed and published
type Topic struct {
	ps		*pubsub.PubSub
	ctx		context.Context
	Subs 	pubsub.Subscription
}

// Publish a message in the topic
func (t *Topic) Publish(msg []byte) error {
	return t.ps.Publish(t.Subs.Topic(), msg)
}

// Chan get a chan for the message
// Should be called only once
func (t *Topic) Chan() <-chan []byte {
	out := make(chan []byte, 20)
	go t.background(out)
	return out
}

func (t *Topic) background(out chan []byte) {

	for {
		ctx, can := context.WithCancel(context.Background())
		defer can()
		
		fromSubscription := make(chan *pubsub.Message)
		errorChan		 := make(chan error)

		go func(){
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

		case <- t.ctx.Done():
			return

		case <- errorChan:
			// It is impossible that this is triggered by a cancellation.
			// This is truly a pubsub error
			return // Make the rest of the app, aware that there is a problem

		case msg := <- fromSubscription:
			out <- msg.GetData()
		}
	}
}
