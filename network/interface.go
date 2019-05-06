package network

// Network is a generic network interface used y
type Network interface{
	Send(msg []byte, to string) error
	GetConnexion(address string) chan []byte
}