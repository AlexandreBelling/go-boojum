package p2p

import (
	"context"
	pubsub 		"github.com/libp2p/go-libp2p-pubsub"
)

// Topic manages communication around a topic in a pubsub manner
type Topic struct {
	ctx				context.Context
	cancelCtx		context.CancelFunc
	topic			string
	subscription 	*pubsub.Subscription
	close			chan bool
	callback		chan []byte
}

// NewTopic is the basic topic constructor
func NewTopic(subscription *pubsub.Subscription) *Topic {
	ctx, cancel := context.WithCancel(context.Background())
	topic := Topic{
		ctx:			ctx,
		cancelCtx:		cancel,
		topic:			subscription.Topic(),
		subscription:	subscription,
		callback:		make(chan []byte, 32),
	}
	go topic.flowCallback()
	return &topic
}

func (top *Topic) flowCallback() {
	for {
		msg, err := top.subscription.Next(top.ctx)
		if err != nil {
			return
		}
		top.callback <- msg.Data
	}
}

// Chan return the consumable topic channel
func (top *Topic) Chan() <-chan []byte {
	return top.callback
}

// Close terminate the topic subscription
func (top *Topic) Close() {
	top.cancelCtx()
	top.subscription.Cancel()
	close(top.callback)
}