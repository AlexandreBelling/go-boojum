package monolithic
// Round keep tracks of aggregation rounds

import(
	"math"
	"github.com/AlexandreBelling/go-boojum/aggregator"
)

type Round struct{
	BackLog []Tree
	Root *Tree
	pendings chan Tree
}

// NewRound construct a new round from payload
func NewRound(payloads []aggregator.Tree) (rou *Round) {

	height := int(
		math.Log2(
			float64(
				len(payloads),		
	))) + 1

	root := NewTree(height)
	backlog := root.GetLeaves()

	rou = &Round{
		BackLog: backlog,
		Root: root,
		pendings: make(chan Tree, len(payloads)),
	}

	for i := 0; i < len(payloads); i++ {
		rou.BackLog[i].payloadChan <- payloads[i] 
	}

	return rou
}

// Launch the scheduler must 
// be run as a go routine or it will deadlock
func (rou *Round) Launch() {

	rou.Root.Schedule(rou.pendings)
}

// Verify the round has been correctly executed
func (rou *Round) Verify(boo *aggregator.Boojum) (bool) {

	rootPayload := <- rou.Root.payloadChan
	rou.Root.payload = &rootPayload
	return boo.Verify(rou.Root.payload)
}


