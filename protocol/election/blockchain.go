package election

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// BatchPubSub is the blockchain interface of participant 
type BatchPubSub interface {
	// PublishAggregated publish the proof aggregate on the blockchain
	// It should be no fail
	PublishAggregated(context.Context, []byte) error
	// NextNewBatch waits and returns the next new batch. 
	NextNewBatch(context.Context) ([][]byte, error)
	// NextBatchDone block and waits untils the current aggregation step is over.
	NextBatchDone(context.Context) error
}

// MockBlockchain simulate a blockchain for tests
type MockBlockchain struct {
	agents			[]*MockBatchPS
	batch			[][]byte
}

// MockBatchPS is a mock for blockchain agent
type MockBatchPS struct {
	blockchain			*MockBlockchain
	newBatches			chan [][]byte
	batchesDone			chan struct{}
}

// NewMockBlockchain returns a mock of a blockchain
func NewMockBlockchain(batch [][]byte) *MockBlockchain {
	return &MockBlockchain{ batch: batch }
}

// CreateAgent builds a new agent and link it to the others
func (m *MockBlockchain) CreateAgent() *MockBatchPS {
	agent := &MockBatchPS{
		blockchain: 		m,
		newBatches: 		make(chan [][]byte, 32),
		batchesDone:		make(chan struct{}, 32),
	}
	m.agents = append(m.agents, agent)
	return agent
}

// NewBatch builds a new agent and link it to the others
func (m *MockBlockchain) NewBatch() {
	for _, agent := range m.agents {
		agent.newBatches <- m.batch
	}
}

// PublishAggregated builds a new agent and link it to the others
func (m *MockBlockchain) PublishAggregated() error {
	for _, agent := range m.agents {
		log.Infof("Publishing a new proof")
		agent.batchesDone <- struct{}{}
	}
	m.NewBatch()
	return nil
}

// NextBatchDone blocks until any agent triggers publish aggregated
func (b *MockBatchPS) NextBatchDone(ctx context.Context) error {
	select{
	case <- ctx.Done():
		return ctx.Err()
	case <- b.batchesDone:
		return nil
	}
}

// NextNewBatch blocks until a new batch is available
func (b *MockBatchPS) NextNewBatch(ctx context.Context) ([][]byte, error) {
	select{
	case <- ctx.Done():
		return [][]byte{}, ctx.Err()
	case batch := <- b.newBatches:
		return batch, nil
	}
}

// PublishAggregated publish that a new batch have been aggregated
func (b *MockBatchPS) PublishAggregated(ctx context.Context, aggregated []byte) error {
	b.blockchain.PublishAggregated()
	return nil
}
