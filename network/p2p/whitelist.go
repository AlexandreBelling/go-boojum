package p2p

import (
	"sync"

	net 		"github.com/libp2p/go-libp2p-net"
	peer 		"github.com/libp2p/go-libp2p-peer"
	wtore 		"github.com/libp2p/go-libp2p-peerstore"
	netCommon	"github.com/AlexandreBelling/go-boojum/network"

)

// WhiteList is a WhiteList component managing permissionning
type WhiteList struct {
	Provider			netCommon.WhiteListProvider
	newPeersChan		chan [][]byte
	peers				map[peer.ID]wtore.PeerInfo
	mut					sync.RWMutex
}

// NewWhiteList construct a new whitelist instance
func NewWhiteList(provider netCommon.WhiteListProvider) *WhiteList {
	w := &WhiteList{}
	w.SetProvider(provider)
	return w
}

// Maintain permissions the connexions by looking to a 
func (w *WhiteList) Maintain() error {

	newPeersChan, err := w.Provider.GetPeersChan()
	if err != nil {
		return err
	}
	w.newPeersChan = newPeersChan

	// Set initial white list
	err = w.GetWhiteListInit()
	if err != nil {
		return err
	}

	go w.maintainanceLoop()
	return nil
}

func (w *WhiteList) maintainanceLoop() {
	for {
		newPeers, open := <- w.newPeersChan
		if !open {
			return
		}
		w.ResetPeers(ParsePeerList(newPeers))
	}
}


// SetProvider injects a provider in the instance
func (w *WhiteList) SetProvider(provider netCommon.WhiteListProvider) {
	w.Provider = provider
}

// ResetPeers updates the whitelist with new data
func (w *WhiteList) ResetPeers(whiteList map[peer.ID]wtore.PeerInfo) {
	w.mut.Lock()
	defer w.mut.Unlock()
	w.peers = whiteList
}

// OnlyWhiteListed ensures that streams can only originate from whitelisted nodes
func (w *WhiteList) OnlyWhiteListed(conn net.Conn) {
	if !w.IsWhiteListed(conn.RemotePeer()) {
		conn.Close()
	}
	return
}

// IsWhiteListed returns true if the given peerID is whitelisted
func (w *WhiteList) IsWhiteListed(id peer.ID) bool {
	w.mut.RLock()
	defer w.mut.RUnlock()
	_, ok := w.peers[id]
	return ok
}

// GetWhiteListInit injects a blockchain client in the instance
func (w *WhiteList) GetWhiteListInit() error {
	getPeersRes, err := w.Provider.GetPeers()
	if err != nil {
		return err
	}

	// Update the whitelist
	w.ResetPeers(ParsePeerList(getPeersRes))
	return nil
}

// ParsePeerList parses the list of the peers and ignore the one malformed
func ParsePeerList(peersBytes [][]byte) map[peer.ID]wtore.PeerInfo {
	newPeers := make(map[peer.ID]wtore.PeerInfo)

	parsePeerInfoLoop:
	for _, pr := range peersBytes {
		pinfo := wtore.PeerInfo{}
		err := pinfo.UnmarshalJSON(pr)
		if err != nil {
			continue parsePeerInfoLoop
		}
		newPeers[pinfo.ID] = pinfo
	}

	return newPeers
}

