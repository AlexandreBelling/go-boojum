package boojum

// Boojum is struct wrapper for boojum aggregator
type Boojum struct {
	dir string
}

// New is a boojum constructor
func New() *Boojum {
	return &Boojum{dir: ""}
}

// WithDir attaches a directory to a boojum
func (boo *Boojum) WithDir(dir string) *Boojum {
	boo.dir = dir
	return boo
}

// Initialize performs precomputation that are necessary before using boojum
func (boo *Boojum) Initialize() *Boojum {
	initialize()
	return boo
}

// RunGenerators fetches the proving and verifications keys
func (boo *Boojum) RunGenerators() *Boojum {
	runGenerators(boo.dir)
	return boo
}

// MakeExample returns an example proof
func (boo *Boojum) MakeExample() *Tree {
	tree := newTree()
	makeExampleProof(&tree.data)
	return tree
}

// AggregateTrees returns the aggregated tree
func (boo *Boojum) AggregateTrees(left, right Tree) (output *Tree) {
	output = newTree()
	proveAggregation(
		left.data,
		right.data,
		&output.data,
	)
	return output
} 

// Verify returns a boolean indicating that a tree is valid
func (boo *Boojum) Verify(tree *Tree) bool {
	return verify(
		tree.data,
	)
}