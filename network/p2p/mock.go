package p2p

import (
	"fmt"
	"time"
	"context"

	// log "github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/crypto"


	bnetwork "github.com/AlexandreBelling/go-boojum/network"
)

// RandomIdentity ...
func RandomIdentity() (libp2p.Option) {
	sk, _, _ := crypto.GenerateSecp256k1Key(nil)
	return libp2p.Identity(sk)
}

// DefaultServer ...
func DefaultServer(addr string, wlp bnetwork.WhiteListProvider) (*Server, error) {

	listenAddr, err := ListenAddress(addr)
	if err != nil {
		return nil, err
	}

	hs, err := libp2p.New(context.Background(), listenAddr, RandomIdentity())
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(context.Background(), hs)
	bstr := NewBoostrappingRoutine(
		context.Background(), hs, wlp, 5, 9,
		time.Duration(1) * time.Minute,
	)

	return &Server{
		Host: 	hs,
		PubSub: ps,
		Bootstrap: bstr,
	}, nil
}

// MakeServers returns a list of started servers
func MakeServers(n int) []*Server {

	// Wait for the servers to connect
	defer time.Sleep(time.Duration(5) * time.Second)

	servers := make([]*Server, n)
	wlp := bnetwork.NewMockWhiteListProvider()

	for i := 0; i<n; i++ {

		addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", 9000 + i)
		s, _ := DefaultServer(addr, wlp)
		servers[i] = s

		pi := peer.AddrInfo{
			ID:		s.Host.ID(),
			Addrs:	s.Host.Addrs(),
		}

		marshalled, _ := pi.MarshalJSON()
		wlp.Add(marshalled)
	}

	for _, s := range servers {
		s.Start()
	}

	return servers
}