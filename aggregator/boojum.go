package aggregator

func init() {
	// Initialize performs precomputation that are necessary before using boojum
	initialize()
}

// Boojum is struct wrapper for boojum aggregator
type Boojum struct {
	dir string
}

// NewBoojum is a boojum constructor
func NewBoojum() *Boojum {
	return &Boojum{dir: ""}
}

// WithDir attaches a directory to a boojum
func (boo *Boojum) WithDir(dir string) *Boojum {
	boo.dir = dir
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
func (boo *Boojum) AggregateTrees(left, right []byte) []byte {

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
	return res
}

// Verify returns a boolean indicating that a tree is valid
func (boo *Boojum) Verify(buff []byte) bool {
	tree := newTree(boo).SetDataFromBytes(buff)
	res := verify(
		tree.data,
	)
	return res
}
