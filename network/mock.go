package network

// MockNetwork ..
type MockNetwork struct {
	Connexions map[string]chan []byte
}

// Send push the message in the indicated channel
func (net *MockNetwork) Send(msg []byte, to string) error {
	net.Connexions[to] <- msg
	return nil
}

// GetConnexion returns a connexion channel if it does not exist, it creates one
func (net *MockNetwork) GetConnexion(address string) chan []byte {
	if con, ok := net.Connexions[address]; ok {
		return con
	}

	net.Connexions[address] = make(chan []byte)
	return net.Connexions[address]
}