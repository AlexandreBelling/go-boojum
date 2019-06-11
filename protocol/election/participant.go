package election

import (
	"context"

	// log "github.com/sirupsen/logrus"

	"github.com/AlexandreBelling/go-boojum/aggregator"
	"github.com/AlexandreBelling/go-boojum/blockchain"
	net "github.com/AlexandreBelling/go-boojum/network"
	"github.com/AlexandreBelling/go-boojum/protocol"
)

// Participant is the higher level protocol struct
type Participant struct {
	ctx context.Context

	ID             protocol.ID
	MemberProvider protocol.MemberProvider
	Network        net.PubSub
	Blockchain     *BCUser
	Aggregator     aggregator.Aggregator
}

// NewParticipant ..
func NewParticipant(ctx context.Context) *Participant {
	return &Participant{
		ctx: ctx,
	}
}

// Start the main routine of the participant
func (par *Participant) Start() {
	go par.background()
}

func (par *Participant) background() {
	for {
		select {
		case batch := <-par.Blockchain.NewBatch:
			round := NewRound(par.ctx, par, batch)
			round.Start()
		case <-par.ctx.Done():
			return
		}
	}
}

// SetNetwork ..
func (par *Participant) SetNetwork(network net.PubSub) {
	par.Network = network
}

// SetMemberProvider ...
func (par *Participant) SetMemberProvider(provider protocol.MemberProvider) {
	par.MemberProvider = provider
}

// SetBCInterface ..
func (par *Participant) SetBCInterface(blockchain blockchain.Client) {
	par.Blockchain = NewBCUser(blockchain)
}

// SetAggregator ..
func (par *Participant) SetAggregator(aggregator aggregator.Aggregator) {
	par.Aggregator = aggregator
}
