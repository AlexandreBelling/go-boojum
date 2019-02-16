package boojum

type Boojum struct {
	dir string,
}

func New() *Boojum {
	return &Boojum{dir: ""}
}

func (boo *Boojum) WithDir(dir string) *Boojum {
	boo.dir = dir
	return boo
}

// Initialize performs precomputation that are necessary before using boojum
func (boo *Boojum) Initialize() *Boojum {
	initialize()
	return boo
}

// RunGenerator fetch the proving and verifications keys
func (boo *Boojum) RunGenerators() {
	runGenerators(boo.dir)
	return boo
}

func (boo *Boojum) MakeExampe() *Tree {
	tree := newTree()
	makeExampleProof(&tree.data)
}

func (boo *Boojum) aggregateTrees(left, right Tree) (output *Tree) {
	output = newTree()
	proveAggregation(
		left.data,
		right.data,
		&output.data
	)
	return output
} 

func (boo *Boojum) verify(tree *Tree) bool {
	return verify(
		tree.data
	)
}