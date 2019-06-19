package election

import (
	"context"

	// log "github.com/sirupsen/logrus"

	"github.com/AlexandreBelling/go-boojum/identity"
	"github.com/AlexandreBelling/go-boojum/aggregator"
	// "github.com/AlexandreBelling/go-boojum/blockchain"
	net "github.com/AlexandreBelling/go-boojum/network"
	"github.com/AlexandreBelling/go-boojum/protocol"
)

// Participant is the higher level protocol struct
type Participant struct {
	ctx context.Context

	ID             identity.ID
	MemberProvider protocol.MemberProvider

	Network        net.PubSub
	BatchPubSub    BatchPubSub
	Aggregator     aggregator.Aggregator
}

// NewParticipant ..
func NewParticipant(
	ctx 				context.Context,
	id					identity.ID,
	aggregator	 		aggregator.Aggregator,
	memberProvider 		protocol.MemberProvider,
	batchPubSub			BatchPubSub,
	network				net.PubSub,

	) *Participant {
	return &Participant{
		ctx: 				ctx,

		MemberProvider:		memberProvider,
		ID:					id,

		Network:			network,		
		BatchPubSub:		batchPubSub,
		Aggregator:			aggregator,
	}
}

// Start the main routine of the participant
func (par *Participant) Start() {
	go par.background()
}

func (par *Participant) background() {
	for {
		batch := par.BatchPubSub.NextNewBatch(par.ctx)
		NewRound(par.ctx, par, batch).Start()

		select{
		case <-par.ctx.Done():
			return
		}
	}
}
