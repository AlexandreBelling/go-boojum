package network

import "context"

// PubSub is a generic network interface used y
type PubSub interface {
	GetTopic(ctx context.Context, topic string) Topic
	Publish(topic string, msg []byte) error
}

// WhiteListProvider is a backend that retrieves a list of authorized peers
type WhiteListProvider interface {
	GetPeersChan() (chan [][]byte, error)
	GetPeers() ([][]byte, error)
}

// A Topic is a pubsub abstraction that can be subscribed and published
type Topic interface {
	Close()
	Publish(msg []byte) error
	Chan() (<-chan []byte, error)
}
