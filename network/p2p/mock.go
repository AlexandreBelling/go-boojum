package p2p

import (
	"time"
	"context"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-pubsub"

	bnetwork "github.com/AlexandreBelling/go-boojum/network"
)

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