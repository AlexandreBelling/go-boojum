package monolithic

import (
	"sync"
)

// Tree is a recursive helper that helps scheduling the aggregation
type Tree struct {
	left, right *Tree
	payloadChan chan []byte
	payload     []byte
	height      int
	mut         sync.Mutex
}

// NewTree assigns a tree with depth
func NewTree(height int) (t *Tree) {
	if height > 0 {

		t = &Tree{
			left:        NewTree(height - 1),
			right:       NewTree(height - 1),
			payloadChan: make(chan []byte, 1),
			height:      height,
		}

		return t
	}
	return nil
}

// GetLeaves returns the list of all leaves
func (t *Tree) GetLeaves() (leaves []Tree) {

	leaves = make([]Tree, 0)
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

// Schedule waits for operands to be aggregated then add to scheduler
func (t *Tree) Schedule(pendings chan Tree) {

	if t.left != nil && t.right != nil {
		go t.left.Schedule(pendings)
		go t.right.Schedule(pendings)

		// Waits for children to be completed before adding to pendings
		leftPayload := <-t.left.payloadChan
		rightPayload := <-t.right.payloadChan

		t.left.payload = leftPayload
		t.right.payload = rightPayload

		pendings <- *t
	}
}
