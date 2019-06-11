package aggregator

import "C"

//Tree is a container for an aggregated snark
type Tree struct {
	boo  *Boojum
	data *byte
}

func newTree(boo *Boojum) *Tree {
	var data *byte
	data = nil
	return &Tree{
		boo:  boo,
		data: data,
	}
}

// Rm unallocate the memory associated with the Tree object
func (tree *Tree) Rm() {
	memFree(tree.data)
}

// Verify returns true if the tree is valid
func (tree *Tree) Verify() bool {
	return verify(tree.data)
}

// ToByte return a go slice of the datas
func (tree *Tree) ToByte() []byte {
	return toByte(tree.data)
}

// SetDataFromBytes sets the data from a slice
func (tree *Tree) SetDataFromBytes(data []byte) *Tree {
	tree.data = &data[0]
	return tree
}
