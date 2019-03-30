package aggregator

import "C"

//Tree is a container
type Tree struct {
	boo *Boojum
	data *byte
}

func newTree(boo *Boojum) (*Tree) {
	var data *byte
	data = nil
	return &Tree{
		boo: boo,
		data : data,
	}
}

// Rm unallocate the memory associated with the Tree object
func (tree *Tree) Rm() {
	memFree(tree.data)
}

// Verify returns true if the tree is valid
func (tree *Tree) Verify() (bool) {
	return verify(tree.data)
}
