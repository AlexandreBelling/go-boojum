package network

// Network is a generic network interface used y
type Network interface{
	Send(msg interface{}, to string) error
	GetConnexion(address string) chan interface{}
}