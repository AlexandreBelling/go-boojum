package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"

	bnetwork "github.com/AlexandreBelling/go-boojum/network"
)

// BoostrappingRoutine embedded all the bootstrapping logic
type BoostrappingRoutine struct {
	ctx  context.Context
	Host host.Host

	mut         sync.RWMutex
	Wlp         bnetwork.WhiteListProvider
	peerChan    chan [][]byte
	sortedPeers []peer.ID
	whitelist   map[peer.ID][]ma.Multiaddr

	minConns int
	target   int
	maxConns int

	bootstrappingPeriod time.Duration
}

// NewBoostrappingRoutine construct a new BootstrappingRoutine object
func NewBoostrappingRoutine(
	ctx context.Context,
	host host.Host,

	Wlp bnetwork.WhiteListProvider,

	minConns int,
	maxConns int,

	bootstrappingPeriod time.Duration,

) *BoostrappingRoutine {

	pchan, _ := Wlp.GetPeersChan()
	return &BoostrappingRoutine{
		ctx:                 ctx,
		Host:                host,
		Wlp:                 Wlp,
		peerChan:            pchan,
		minConns:            minConns,
		target:              (minConns + maxConns) / 2,
		maxConns:            maxConns,
		bootstrappingPeriod: bootstrappingPeriod,
	}
}

// Start the bootsrapping routine
func (br *BoostrappingRoutine) Start() error {
	initialPeers, err := br.Wlp.GetPeers()
	if err != nil {
		return err
	}

	br.SetNewWhiteList(initialPeers)
	go br.background()
	return nil
}

func (br *BoostrappingRoutine) background() error {

	err := br.RunBootstrap()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(br.bootstrappingPeriod)
	for {
		select {

		case <-br.ctx.Done():
			// Graceful shutdown, nothing to do yet

		case <-ticker.C:
			err := br.RunBootstrap()
			if err != nil {
				return err
			}

		case newPeers := <-br.peerChan:
			err := br.SetNewWhiteList(newPeers)
			if err != nil {
				return err
			}
			// Remove every non-listed connexions
			br.TrimUnlistedPeers()
		}
	}
}

// RunBootstrap will adjust the number of connected peers
func (br *BoostrappingRoutine) RunBootstrap() (err error) {

	nConns := len(br.Host.Network().Peers())
	if nConns < br.minConns {
		err = br.AddConnections()
	}

	if nConns > br.maxConns {
		err = br.PruneConnections()
	}

	return err
}

// AddConnections add peers until the target is reached
func (br *BoostrappingRoutine) AddConnections() error {

	br.mut.RLock()
	defer br.mut.RUnlock()

	self := br.Host.ID()

	nConns := len(br.Host.Network().Peers())
ForEachPeerConnectLoop:
	for _, id := range br.sortedPeers {

		connectedNess := br.Host.Network().Connectedness(id)
		if connectedNess == network.CanConnect || connectedNess == network.NotConnected {

			if id == self {
				continue ForEachPeerConnectLoop
			}

			addrInfo := peer.AddrInfo{
				ID:    id,
				Addrs: br.whitelist[id],
			}

			if err := br.Host.Connect(context.Background(), addrInfo); err != nil {
				continue ForEachPeerConnectLoop
			}

			log.Debugf("Found new peer : %v, from : %v", id, br.Host.ID())
			nConns++
			if nConns == br.target {
				return nil
			}
		}
	}

	return nil
}

// PruneConnections remove peers until the target is reached
func (br *BoostrappingRoutine) PruneConnections() error {

	br.mut.RLock()
	defer br.mut.RUnlock()

	nConns := len(br.Host.Network().Peers())
ForEachPeerDisconnectLoop:
	for i := len(br.sortedPeers) - 1; i > 0; i-- {

		id := br.sortedPeers[i]
		connectedNess := br.Host.Network().Connectedness(id)
		if connectedNess == network.Connected {

			err := br.Host.Network().ClosePeer(id)
			if err != nil {
				continue ForEachPeerDisconnectLoop
			}

			nConns--
			if nConns == br.target {
				return nil
			}
		}
	}

	return nil
}

// SetNewWhiteList replace the old whitelist with the new one
func (br *BoostrappingRoutine) SetNewWhiteList(newPeers [][]byte) error {

	newWhiteList := make(map[peer.ID][]ma.Multiaddr)
	for _, marshalled := range newPeers {

		var unmarshalled peer.AddrInfo
		err := unmarshalled.UnmarshalJSON(marshalled)
		if err != nil {
			return err
		}
		newWhiteList[unmarshalled.ID] = unmarshalled.Addrs
	}

	bxd := byXORDistance{
		slice:     make([]peer.ID, len(newPeers)),
		reference: br.Host.ID(),
	}

	index := 0 // Population the slice with the keys (peer.IDs) of newWhiteList
	for id := range newWhiteList {
		bxd.slice[index] = id
		index++
	}
	bxd.Sort()

	br.mut.Lock()
	br.whitelist = newWhiteList
	br.sortedPeers = bxd.slice
	br.mut.Unlock()

	return nil
}

// TrimUnlistedPeers remove peers that are not listed anymore
func (br *BoostrappingRoutine) TrimUnlistedPeers() {
	for _, connectedPeer := range br.Host.Network().Peers() {
		if !br.authorized(connectedPeer) {
			br.Host.Network().ClosePeer(connectedPeer)
		}
	}
}

// Permissionize blocks remote connexions from unauthorized peers
func (br *BoostrappingRoutine) Permissionize() {
	onlyWhiteListed := func(conn network.Conn) {
		if !br.authorized(conn.RemotePeer()) {
			conn.Close()
		}
	}
	br.Host.Network().SetConnHandler(onlyWhiteListed)
}

func (br *BoostrappingRoutine) authorized(p peer.ID) bool {
	_, listed := br.whitelist[p]
	return listed
}
