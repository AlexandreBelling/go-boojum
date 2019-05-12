package network

// PubSub is a generic network interface used y
type PubSub interface{
	Subscribe(topic string) (*Topic, error)
}

// WhiteListProvider is a backend that retrieves a list of authorized peers
type WhiteListProvider interface{
	GetPeersChan() (chan [][]byte, error)
	GetPeers() ([][]byte, error)

}

// A Topic is a pubsub abstraction that can be subscribed and published
type Topic interface {
	Publish(msg []byte)
	Chan() <-chan []byte
}