package protocol

// Node is an interface
type Node interface{}

// Tree is a recursive helper that helps scheduling the aggregation
type Tree struct {
	Parent   *Tree
	Children []*Tree
	Node     Node
}

// NewTree instantiate a tree with given depth and arity
func NewTree(height, arity int) *Tree {
	return newTree(height, arity, nil)
}

// Set as private to prevent user to mistakenly set a value to t
func newTree(height, arity int, parent *Tree) *Tree {
	t := &Tree{Parent: parent}

	if height > 1 {
		// Recursively create offspring
		t.Children = make([]*Tree, arity)
		for i := 0; i < arity; i++ {

			t.Children[i] = newTree(height-1, arity, t)
		}
	}

	return t
}

// TreeFunc is an alias for a closure taking a tree as argument and returning nothing
type TreeFunc func(t *Tree)

// Walk apply a function at each nodes of the subtree, starting with the root
// It's depth-first
func (t *Tree) Walk(f TreeFunc) {
	f(t)
	for _, child := range t.Children {
		child.Walk(f)
	}
}

// MakeLeavesAccumulator returns a TreeFunc and a closure to get the list of leaves
func MakeLeavesAccumulator() (TreeFunc, func() []*Tree) {
	accumulator := make([]*Tree, 0) // Closure local accumulator

	// Appends to accumulator if t is a leaf node
	accumulateIfLeaf := func(t *Tree) {
		if t.IsLeaf() {
			accumulator = append(accumulator, t)
		}
	}
	getAccumulator := func() []*Tree { return accumulator }

	return accumulateIfLeaf, getAccumulator
}

// IsLeaf returns true is this tree is a leaf
func (t *Tree) IsLeaf() bool {
	return t.Children == nil
}

// GetLeaves returns the list of all leaves of the subtree
func (t *Tree) GetLeaves() (leaves []*Tree) {

	// Initialize the accumulator
	accumulateIfLeaf, getAccumulator := MakeLeavesAccumulator()

	// Run the accumulation on leaves
	t.Walk(accumulateIfLeaf)
	return getAccumulator()
}
