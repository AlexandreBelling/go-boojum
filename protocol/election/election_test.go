package election

import (
	"testing"

	"context"
	log "github.com/sirupsen/logrus"

	"github.com/AlexandreBelling/go-boojum/aggregator"
	"github.com/AlexandreBelling/go-boojum/network/p2p"
	"github.com/AlexandreBelling/go-boojum/protocol"
)

func TestElection(t *testing.T) {

	const n = 10
	const batchSize = 25

	networks := p2p.MakeServers(n)
	boojum := &aggregator.MockAggregator{} // It is stateless therefore safe to copy
	blockchain := MakeBCClientMock(batchSize)
	blockchain.GenerateBatch(boojum)
	memberProvider := &protocol.DefaultMembersProvider{
		WLP: networks[0].Bootstrap.Wlp,
	}

	participants := make([]*Participant, n)
	for index := range participants {
		participants[index] = NewParticipant(context.Background())

		participants[index].SetAggregator(boojum)
		participants[index].SetBCInterface(blockchain)
		participants[index].SetNetwork(networks[index])
		participants[index].SetMemberProvider(memberProvider)

		blockchain.Connect(participants[index].Blockchain)
		id, _ := networks[index].Host.ID().Marshal()
		copy(participants[index].ID[:32], id)
		participants[index].Start()
	}

	blockchain.NewBatch()
	log.Infof("Started a new batch")

	<-context.Background().Done()
}
