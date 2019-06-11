package network

import (
	"bytes"
	"sync"
)

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

	net.Connexions[address] = make(chan []byte, 100)
	return net.Connexions[address]
}

// MockWhiteListProvider is a backend that retrieves a list of authorized peers
type MockWhiteListProvider struct {
	mutChan   sync.Mutex
	peersChan []chan [][]byte

	mutPeers sync.RWMutex
	peers    [][]byte
}

// NewMockWhiteListProvider ....
func NewMockWhiteListProvider() *MockWhiteListProvider {
	return &MockWhiteListProvider{
		peersChan: []chan [][]byte{},
		peers:     [][]byte{},
	}
}

// GetPeersChan ....
func (m *MockWhiteListProvider) GetPeersChan() (chan [][]byte, error) {
	res := make(chan [][]byte, 32)

	m.mutChan.Lock()
	defer m.mutChan.Unlock()
	m.peersChan = append(m.peersChan, res)

	return res, nil
}

// GetPeers ......
func (m *MockWhiteListProvider) GetPeers() ([][]byte, error) {
	m.mutPeers.RLock()
	defer m.mutPeers.RUnlock()
	return m.peers, nil
}

// Add ...
func (m *MockWhiteListProvider) Add(p []byte) {
	m.mutPeers.Lock()
	defer m.mutPeers.Unlock()
	m.peers = append(m.peers, p)
}

// Rm ...
func (m *MockWhiteListProvider) Rm(p []byte) {
	m.mutPeers.Lock()
	defer m.mutPeers.Unlock()

	old := m.peers
	new := [][]byte{}

	for _, q := range old {
		if bytes.Equal(p, q) {
			continue
		}
		new = append(new, q)
	}
}
