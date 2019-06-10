package election

import (
	"context"

	"github.com/AlexandreBelling/go-boojum/protocol"
	"github.com/AlexandreBelling/go-boojum/aggregator"
	"github.com/AlexandreBelling/go-boojum/blockchain"
	net "github.com/AlexandreBelling/go-boojum/network"
)

// Participant is the higher level protocol struct
type Participant struct {
	ctx						context.Context

	ID						protocol.ID
	MemberProvider			protocol.MemberProvider
	Network 				net.PubSub
	Blockchain				*BCUser
	Aggregator				aggregator.Aggregator
}

// NewParticipant ..
func NewParticipant(ctx context.Context) *Participant {
	return &Participant{
		ctx: 	ctx,
	}
}

// Run starts the main routine of the participant
func (par *Participant) Start() {}

// SetNetwork ..
func (par *Participant) SetNetwork(network net.PubSub) {
	par.Network = network
}

// SetBCInterface ..
func (par *Participant) SetBCInterface(blockchain blockchain.Client) {
	par.Blockchain = NewBCUser(blockchain)
}

// SetAggregator ..
func (par *Participant) SetAggregator(aggregator aggregator.Aggregator) {
	par.Aggregator = aggregator
}