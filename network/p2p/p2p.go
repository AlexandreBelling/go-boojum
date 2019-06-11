package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-pubsub"

	bnetwork "github.com/AlexandreBelling/go-boojum/network"
)

// Server ...
type Server struct {
	Host      host.Host
	PubSub    *pubsub.PubSub
	Bootstrap *BoostrappingRoutine
}

// Start the server
func (s *Server) Start() {
	s.Bootstrap.Start()
}

// GetTopic returns a topic
func (s *Server) GetTopic(ctx context.Context, topic string) bnetwork.Topic {

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
