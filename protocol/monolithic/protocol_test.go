package monolithic

import (
	"testing"
	"github.com/AlexandreBelling/go-boojum/aggregator"
)

func TestProtocol(t *testing.T) {

	boo := aggregator.NewBoojum().WithDir("./../aggregator/setup").RunGenerators()
	batch := make([][]byte, 8)

	// Initialize the backlog
	for i:=0; i<8; i++ {
		batch[i] = boo.MakeExample()
	}
	
	worker := &Worker{boo: boo}
	round := NewRound(batch)
	done := make(chan bool, 1)

	go worker.StartConsuming(round.pendings, done)
	go round.Launch()

	// This waits for the round to be complete
	valid := round.Verify(boo)
	done <- true // Frees the worker

	if !valid {
		t.FailNow()
	}
}