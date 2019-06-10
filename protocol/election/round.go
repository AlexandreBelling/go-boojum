package election

import (
	"context"
	"github.com/AlexandreBelling/go-boojum/protocol"
)

// Round encompasses the computation made in a single
type Round struct {
	ctx				context.Context

	ID				protocol.ID
	Batch			[][]byte
	Participant		*Participant
	Members			protocol.IDList

	TopicProvider	*TopicProvider
}

// WithTopicProvider add a topic provider in the Round object
func (r *Round) WithTopicProvider() *Round {
	r.TopicProvider = &TopicProvider{
		Network: 	r.Participant.Network,
		Round:		r,
	}
	return r
}

// Start run the Round
func (r *Round) Start() {
	if r.Participant.ID == r.GetLeaderID() {
		r.runLeader()
	}
	r.runWorker()
}

// GetLeaderID returns the ID and position of the leader
func (r *Round)	GetLeaderID() (protocol.ID) {
	_, res := r.Members.SmallestHigherThan(r.ID)
	return res
}

// run a leader instance
func (r *Round) runLeader() {
	NewLeader(r.ctx, r).Start()
}

// run a worker instance
func (r *Round) runWorker() {}





