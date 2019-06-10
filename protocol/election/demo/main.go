package main

import (
	"context"

	"github.com/AlexandreBelling/go-boojum/protocol"
	"github.com/AlexandreBelling/go-boojum/aggregator"
	"github.com/AlexandreBelling/go-boojum/network/p2p"
	"github.com/AlexandreBelling/go-boojum/protocol/election"
)



func main() {

	const n = 2
	const batchSize = 32

	networks := p2p.MakeServers(n)
	boojum := &aggregator.MockAggregator{} // It is stateless therefore safe to copy
	blockchain := election.MakeBCClientMock(batchSize)
	blockchain.GenerateBatch(boojum)
	memberProvider := &protocol.DefaultMembersProvider{ 
		WLP: networks[0].Bootstrap.Wlp,
	}

	participants := make([]*election.Participant, n)
	for index := range participants {
		participants[index] = election.NewParticipant(context.Background())

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
	<- context.Background().Done()

}