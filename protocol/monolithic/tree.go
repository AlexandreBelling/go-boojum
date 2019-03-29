package monolithic

import (
	"math"
	"github.com/AlexandreBelling/go-boojum/aggregator"
)

// Tree is a recursive helper that helps scheduling the aggregation
type Tree struct{
	left, right *Tree
	payloadChan chan aggregator.Tree
	payload *aggregator.Tree
	height int
}

// NewTree assigns a tree with depth 
func NewTree(height int) (t* Tree) {
	if height > 0 {

		t = &Tree{
			left: NewTree(height - 1),
			right: NewTree(height - 1),
			payloadChan: make(chan aggregator.Tree, 1),
			height: height,
		}

		return t
	}
	return nil
}

// GetLeaves returns the list of all leaves
func (t *Tree) GetLeaves() (leaves []Tree) {

	leaves = make ([]Tree, 0)
	t.getLeavesRecursive(&leaves)
	return leaves
}

func (t *Tree) getLeavesRecursive(leaves *[]Tree) {

	if t.left != nil && t.right != nil {
		t.left.getLeavesRecursive(leaves)
		t.right.getLeavesRecursive(leaves)
		return

	} else if t.left == nil && t.right == nil {
		*leaves = append(*leaves, *t)
		return
	}

	println("Found a node with a left but no right ")
	println(t.left)
	println(t.right)
}

//Schedule waits for operands to be aggregated then add to scheduler
func (t *Tree) Schedule(pendings chan Tree) {

	if t.left != nil && t.right != nil {
		go t.left.Schedule(pendings)
		go t.right.Schedule(pendings)
	
		// Waits for children to be completed before adding to pendings
		leftPayload := <- t.left.payloadChan
		rightPayload := <- t.right.payloadChan

		t.left.payload = &leftPayload
		t.right.payload = &rightPayload 

		pendings <- *t
	}
}

// Round keep tracks of aggregation rounds
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

// Worker is responsible to perform aggregation
type Worker struct {
	boo *aggregator.Boojum
}

// StartConsuming in the pendings queue until it receives done
func (w *Worker) StartConsuming(pendings chan Tree, done chan bool) {
	for {
		select {
		case <- done:
			return
		case job := <- pendings:

			// Scheduler ensures that whenever job is received
			// left and right are already assigned
			job.payloadChan <- *w.boo.AggregateTrees(
				*job.left.payload,
				*job.right.payload,
			)

			job.left.payload.memFree()
			job.right.payload.memFree()
			
		}
	}
}

