package election

import (
	"testing"
	log "github.com/sirupsen/logrus"
	"github.com/AlexandreBelling/go-boojum/network"
	"github.com/AlexandreBelling/go-boojum/aggregator"
)

func TestElection(t *testing.T) {

	log.SetLevel(log.InfoLevel)

	blockchain := BlockchainClientProtocolMock{BatchSize: 2}
	net := network.MockNetwork{Connexions: make(map[string]chan []byte)}
	agg := aggregator.MockAggregator{}

	blockchain.GenerateBatch(&agg)

	peers := []*Participant{
		NewParticipant("0"),
		NewParticipant("1"),
		NewParticipant("2"),
		NewParticipant("3"),
	}

	for _, peer := range peers {
		peer.WithNetwork(&net).WithBCInterface(&blockchain).WithAggregator(&agg)
		blockchain.Connect(peer.Blockchain)
		go peer.Run()
	}

	// Publish aggregated will send a new batch
	blockchain.NewBatch()

	for{}
}