package aggregator

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
func (boo *Boojum) MakeExample() []byte {
	tree := newTree(boo)
	defer tree.Rm()
	makeExampleProof(&tree.data)
	return tree.ToByte()
}

// AggregateTrees returns the aggregated tree
func (boo *Boojum) AggregateTrees(left, right []byte) ([]byte) {
	
	// Initialize Tree object
	leftTree := newTree(boo).SetDataFromBytes(left)
	rightTree := newTree(boo).SetDataFromBytes(right)
	outputTree := newTree(boo)

	proveAggregation(
		leftTree.data,
		rightTree.data,
		&outputTree.data,
	)

	res := outputTree.ToByte()

	// Unallocate the memory
	//leftTree.Rm()
	//rightTree.Rm()
	//outputTree.Rm()

	return res
} 

// Verify returns a boolean indicating that a tree is valid
func (boo *Boojum) Verify(buff []byte) bool {
	tree := newTree(boo).SetDataFromBytes(buff)
	res := verify(
		tree.data,
	)
	//tree.Rm()
	return res
}