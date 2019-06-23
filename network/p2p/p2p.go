package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-pubsub"

	"github.com/AlexandreBelling/go-boojum/identity"
	"github.com/AlexandreBelling/go-boojum/network"
	"github.com/libp2p/go-libp2p"
)

// Server ...
type Server struct {
	Host      host.Host
	PubSub    *pubsub.PubSub
	Bootstrap *BoostrappingRoutine
}

// NewServerWithID returns a server object with
func NewServerWithID(wlp network.WhiteListProvider, priv *identity.PrivKey, addr string) (*Server, error) {
	listenAddr, err := ListenAddress(addr)
	if err != nil {
		return nil, err
	}

	privP2P, err := Identity(priv.Libp2p())
	if err != nil {
		return nil, err
	}

	hs, err := libp2p.New(context.Background(), listenAddr, privP2P)
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(context.Background(), hs)
	bstr := NewBoostrappingRoutine(
		context.Background(), hs, wlp, 5, 9,
		time.Duration(2)*time.Second,
	)

	return &Server{
		Host:      hs,
		PubSub:    ps,
		Bootstrap: bstr,
	}, nil
}

// Start the server
func (s *Server) Start() {
	s.Bootstrap.Start()
}

// GetPeerInfo returns the marshalled peer info of the server
func (s *Server) GetPeerInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    s.Host.ID(),
		Addrs: s.Host.Addrs(),
	}
}

// GetTopic returns a topic
func (s *Server) GetTopic(ctx context.Context, topic string) network.Topic {
	res := Topic{
		ctx: ctx,

		ps:   s.PubSub,
		Name: topic,
	}
	return &res
}

// Publish sends a message in a topic
func (s *Server) Publish(topic string, msg []byte) error {
	return s.PubSub.Publish(topic, msg)
}
