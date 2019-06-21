package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/AlexandreBelling/go-boojum/aggregator"
	"github.com/AlexandreBelling/go-boojum/identity"
	"github.com/AlexandreBelling/go-boojum/network"
	"github.com/AlexandreBelling/go-boojum/network/p2p"
	"github.com/AlexandreBelling/go-boojum/protocol"
	"github.com/AlexandreBelling/go-boojum/protocol/election"
)

func makeBatch(agg aggregator.Aggregator, batchSize int) [][]byte {
	e := agg.MakeExample()
	res := make([][]byte, batchSize)
	for i := 0; i < batchSize; i++ {
		res[i] = e
	}
	return res
}

func main() {

	const n = 7
	const batchSize = 25

	boojum := &aggregator.MockAggregator{} // The mock is stateless therefore safe to copy
	batch := makeBatch(boojum, batchSize)
	blockchain := election.NewMockBlockchain(batch)
	wlp := network.NewMockWhiteListProvider()
	memberProvider := &protocol.DefaultMembersProvider{WLP: wlp}

	participants := make([]*election.Participant, n)
	p2pServers := make([]*p2p.Server, n)
	for index := range participants {
		id := identity.Generate()

		// Setting up the p2p server
		addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", 9000+index)
		s, _ := p2p.NewServerWithID(wlp, id, addr)

		pi := s.GetPeerInfo()
		marshalled, _ := pi.MarshalJSON()
		wlp.Add(marshalled)

		p2pServers[index] = s
		participants[index] = election.NewParticipant(context.Background(),
			id.ID(),
			boojum,
			memberProvider,
			blockchain.CreateAgent(),
			s,
		)
	}

	// Start the servers only once they have all been added to
	for _, s := range p2pServers {
		s.Start()
	}

	time.Sleep(time.Duration(3) * time.Second)
	for _, p := range participants {
		p.Start()
	}

	blockchain.NewBatch()
	log.Infof("Started a new batch")

	<-context.Background().Done()
}
