package network

// MessageBroker is a generic network interface used y
type MessageBroker interface{
	Send(msg []byte, to string) error
	GetConnexion(address string) chan []byte
}

// WhiteListProvider is a backend that retrieves a list of authorized peers
type WhiteListProvider interface{
	GetPeersChan() (chan [][]byte, error)
	GetPeers() ([][]byte, error)

}